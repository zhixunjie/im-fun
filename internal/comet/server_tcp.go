package comet

import (
	"context"
	"github.com/golang/glog"
	"github.com/zhixunjie/im-fun/api/protocol"
	"github.com/zhixunjie/im-fun/internal/comet/channel"
	"github.com/zhixunjie/im-fun/internal/comet/conf"
	"github.com/zhixunjie/im-fun/pkg/buffer"
	mytime "github.com/zhixunjie/im-fun/pkg/time"
	"io"
	"math"
	"net"
	"strings"
	"time"
)

// InitTCP listen all tcp.bind and start accept connections.
func InitTCP(server *Server, addrs []string, accept int) (err error) {
	var (
		bind     string
		listener *net.TCPListener
		addr     *net.TCPAddr
	)
	for _, bind = range addrs {
		if addr, err = net.ResolveTCPAddr("tcp", bind); err != nil {
			glog.Errorf("TCP ResolveTCPAddr(bind) err=%v", bind, err)
			return
		}
		if listener, err = net.ListenTCP("tcp", addr); err != nil {
			glog.Errorf("TCP ListenTCP(bind) err=%v", bind, err)
			return
		}
		glog.Infof("TCP服务器启动成功，正在监听：%s", bind)
		// 启动N个协程，每个协程开启后进行accept（需要使用REUSE解决惊群问题）
		for i := 0; i < accept; i++ {
			// TODO 使用协程池进行管理
			go acceptTCP(server, listener)
		}
	}
	return
}

func acceptTCP(server *Server, listener *net.TCPListener) {
	var conn *net.TCPConn
	var err error
	var r int

	for {
		if conn, err = listener.AcceptTCP(); err != nil {
			// if listener close then return
			glog.Errorf("listener.Accept(%s) error=%v", listener.Addr().String(), err)
			return
		}
		if err = conn.SetKeepAlive(server.conf.Connect.TCP.KeepAlive); err != nil {
			glog.Errorf("conn.SetKeepAlive() error=%v", err)
			return
		}
		if err = conn.SetReadBuffer(server.conf.Connect.TCP.Rcvbuf); err != nil {
			glog.Errorf("conn.SetReadBuffer() error=%v", err)
			return
		}
		if err = conn.SetWriteBuffer(server.conf.Connect.TCP.Sndbuf); err != nil {
			glog.Errorf("conn.SetWriteBuffer() error=%v", err)
			return
		}
		go serveTCPInit(server, conn, r)
		if r++; r == math.MaxInt {
			r = 0
		}
	}
}

func serveTCPInit(s *Server, conn *net.TCPConn, r int) {
	var (
		// timer
		tr = s.round.TimerPool(r)
		rp = s.round.BufferPool.ReaderPool(r)
		wp = s.round.BufferPool.WriterPool(r)
		// ip addr
		lAddr = conn.LocalAddr().String()
		rAddr = conn.RemoteAddr().String()
	)
	if conf.Conf.Debug {
		glog.Infof("connect success,lAddr=%v,rAddr=%v", lAddr, rAddr)
	}
	s.serveTCP(conn, rp, wp, tr)
}

// serveTCP serve a tcp connection.
func (s *Server) serveTCP(conn *net.TCPConn, readerPool, writerPool *buffer.Pool, timerPool *mytime.Timer) {
	var (
		err    error
		roomId string
		hb     time.Duration
		proto  *protocol.Proto
		bucket *Bucket
		trd    *mytime.TimerData
		lastHb = time.Now()
		rb     = readerPool.Get()
		wb     = writerPool.Get()

		ch = channel.NewChannel(s.conf)
		rr = ch.Reader
		wr = ch.Writer
	)
	ch.Reader.ResetBuffer(conn, rb.Bytes())
	ch.Writer.ResetBuffer(conn, wb.Bytes())
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// handshake
	trd = timerPool.Add(time.Duration(s.conf.Protocol.HandshakeTimeout), func() {
		conn.Close()
		glog.Errorf("key: %s remoteIP: %s step: %d tcp handshake timeout", ch.UserInfo.UserKey, conn.RemoteAddr().String(), step)
	})
	ch.UserInfo.IP, _, _ = net.SplitHostPort(conn.RemoteAddr().String())

	if proto, err = ch.ProtoAllocator.GetProtoCanWrite(); err == nil {
		ch.UserInfo.UserKey, roomId, hb, err = s.authTCP(ctx, rr, wr, proto)
		if err == nil {
			bucket = s.AllocBucket(ch.UserInfo.UserKey)
			err = bucket.Put(roomId, ch)
		}
	}

	if err != nil {
		conn.Close()
		readerPool.Put(rb)
		writerPool.Put(wb)
		timerPool.Del(trd)
		glog.Errorf("key: %s handshake failed error(%v)", ch.UserInfo.UserKey, err)
		return
	}
	trd.Key = ch.UserInfo.UserKey
	timerPool.Set(trd, hb)

	go s.dispatchTCP(conn, wr, writerPool, wb, ch)
	hbTime := s.RandHeartbeatTime()
	for {
		if proto, err = ch.ProtoAllocator.GetProtoCanWrite(); err != nil {
			break
		}
		if err = proto.ReadTCP(rr); err != nil {
			break
		}

		if protocol.Operation(proto.Op) == protocol.OpHeartbeat {
			timerPool.Set(trd, hb)
			proto.Op = int32(protocol.OpHeartbeatReply)
			proto.Body = nil
			// NOTE: send server heartbeat for a long time
			if now := time.Now(); now.Sub(lastHb) > hbTime {
				if err1 := s.Heartbeat(ctx, ch.UserInfo); err1 == nil {
					lastHb = now
				}
			}
			if conf.Conf.Debug {
				glog.Infof("tcp heartbeat receive key:%s, mid:%d", ch.UserInfo.UserKey, ch.Mid)
			}
		} else {
			if err = s.Operate(ctx, proto, ch, bucket); err != nil {
				break
			}
		}

		ch.ProtoAllocator.AdvWritePointer()
		ch.Ready()
	}
	if err != nil && err != io.EOF && !strings.Contains(err.Error(), "closed") {
		glog.Errorf("key: %s server tcp failed error(%v)", ch.UserInfo.UserKey, err)
	}
	bucket.DelChannel(ch)
	timerPool.Del(trd)
	readerPool.Put(rb)
	conn.Close()
	ch.Close()
	if err = s.Disconnect(ctx, ch.Mid, ch.UserInfo.UserKey); err != nil {
		glog.Errorf("key: %s mid: %d operator do disconnect error(%v)", ch.UserInfo.UserKey, ch.Mid, err)
	}
}

func (s *Server) authTCP(ctx context.Context, ch *channel.Channel, proto *protocol.Proto) (UserKey, rid string, hb time.Duration, err error) {
	reader := ch.Reader
	writer := ch.Writer
	// 一直读取，直到读取到的Proto的操作类型为protocol.OpAuth
	for {
		if err = proto.ReadTCP(reader); err != nil {
			return
		}
		if protocol.Operation(proto.Op) == protocol.OpAuth {
			break
		} else {
			glog.Errorf("tcp request operation(%d) not auth", proto.Op)
		}
	}
	if UserKey, rid, hb, err = s.Connect(ctx, proto); err != nil {
		glog.Errorf("Connect UserKey=%v, err=%v", UserKey, err)
		return
	}
	proto.Op = int32(protocol.OpAuthReply)
	proto.Body = nil
	if err = proto.WriteTCP(writer); err != nil {
		glog.Errorf("WriteTCP UserKey=%v, err=%v", UserKey, err)
		return
	}
	err = writer.Flush()
	return
}
