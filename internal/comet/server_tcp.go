package comet

import (
	"context"
	"encoding/json"
	"github.com/sirupsen/logrus"
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
			logrus.Errorf("TCP ResolveTCPAddr(bind) err=%v", bind, err)
			return
		}
		if listener, err = net.ListenTCP("tcp", addr); err != nil {
			logrus.Errorf("TCP ListenTCP(bind) err=%v", bind, err)
			return
		}
		logrus.Infof("TCP服务器启动成功，正在监听：%s", bind)
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
			logrus.Errorf("listener.Accept(%s) error=%v", listener.Addr().String(), err)
			return
		}
		if err = conn.SetKeepAlive(server.conf.Connect.TCP.KeepAlive); err != nil {
			logrus.Errorf("conn.SetKeepAlive() error=%v", err)
			return
		}
		if err = conn.SetReadBuffer(server.conf.Connect.TCP.Rcvbuf); err != nil {
			logrus.Errorf("conn.SetReadBuffer() error=%v", err)
			return
		}
		if err = conn.SetWriteBuffer(server.conf.Connect.TCP.Sndbuf); err != nil {
			logrus.Errorf("conn.SetWriteBuffer() error=%v", err)
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
		logrus.Infof("connect success,lAddr=%v,rAddr=%v", lAddr, rAddr)
	}
	s.serveTCP(conn, rp, wp, tr)
}

// serveTCP serve a tcp connection.
func (s *Server) serveTCP(conn *net.TCPConn, readerPool, writerPool *buffer.Pool, timerPool *mytime.Timer) {
	var (
		err    error
		proto  *protocol.Proto
		bucket *Bucket
		//lastHb = time.Now()
	)
	var hb time.Duration
	var trd *mytime.TimerData
	var rb = readerPool.Get()
	var wb = writerPool.Get()
	var ch = channel.NewChannel(s.conf)
	ch.Reader.ResetBuffer(conn, rb.Bytes())
	ch.Writer.ResetBuffer(conn, wb.Bytes())
	ch.UserInfo.IP, _, _ = net.SplitHostPort(conn.RemoteAddr().String())
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// set timer
	var step = 0
	// TODO 暂时把timer关闭，感觉有点问题
	//trd = timerPool.Add(time.Duration(s.conf.Protocol.HandshakeTimeout), func() {
	//	conn.Close()
	//	logrus.Errorf("handshake timeout UserInfo=%+v,addr=%v,step=%v",
	//		ch.UserInfo, conn.RemoteAddr().String(), step)
	//})
	trd = timerPool.Add(100000*time.Second, func() {
		conn.Close()
		logrus.Errorf("handshake timeout UserInfo=%+v,addr=%v,step=%v",
			ch.UserInfo, conn.RemoteAddr().String(), step)
	})

	step = 1
	{
		release := func() {
			conn.Close()
			readerPool.Put(rb)
			writerPool.Put(wb)
			timerPool.Del(trd)
		}
		// get a proto to write
		proto, err = ch.ProtoAllocator.GetProtoCanWrite()
		if err != nil {
			release()
			logrus.Errorf("UserInfo=%v authTCP err=%v", ch.UserInfo, err)
			return
		}
		// auth（check token）
		hb, err = s.authTCP(ctx, ch, proto)
		if err != nil {
			release()
			logrus.Errorf("UserInfo=%v authTCP err=%v,hb=%v", ch.UserInfo, err, hb)
			return
		}
		// set bucket
		bucket = s.AllocBucket(ch.UserInfo.UserKey)
		err = bucket.Put(ch)
		if err != nil {
			release()
			logrus.Errorf("UserInfo=%v authTCP err=%v", ch.UserInfo, err)
			return
		}
	}
	step = 2
	// set new timer
	//trd.Key = ch.UserInfo.UserKey
	//timerPool.Set(trd, hb)

	// dispatch
	go s.dispatchTCP(conn, writerPool, wb, ch)

	// loop to read client msg
	//hbTime := s.RandHeartbeatTime()
	for {
		if proto, err = ch.ProtoAllocator.GetProtoCanWrite(); err != nil {
			break
		}
		// read msg from client（if there is no msg，it will block here）
		if err = proto.ReadTCP(ch.Reader); err != nil {
			break
		}

		if protocol.Operation(proto.Op) == protocol.OpHeartbeat {
			//timerPool.Set(trd, hb)
			proto.Op = int32(protocol.OpHeartbeatReply)
			proto.Body = nil
			//if now := time.Now(); now.Sub(lastHb) > hbTime {
			//	if err1 := s.Heartbeat(ctx, ch.UserInfo); err1 == nil {
			//		lastHb = now
			//	}
			//}
		} else {
			if err = s.Operate(ctx, proto, ch, bucket); err != nil {
				break
			}
		}

		ch.ProtoAllocator.AdvWritePointer()
		ch.SendReady()
	}
	if err != nil && err != io.EOF && !strings.Contains(err.Error(), "closed") {
		logrus.Errorf("UserInfo=%v sth has happened", ch.UserInfo, err)
	}
	// 回收相关资源
	{
		bucket.DelChannel(ch)
		timerPool.Del(trd)
		readerPool.Put(rb) // writePool will release buffer in Server.dispatchTCP()
		conn.Close()
		ch.Close()
		if err = s.Disconnect(ctx, ch); err != nil {
			logrus.Errorf("Disconnect UserInfo=%+v,err=%v", ch.UserInfo, err)
		}
	}
}

func (s *Server) authTCP(ctx context.Context, ch *channel.Channel, proto *protocol.Proto) (hb time.Duration, err error) {
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
			logrus.Errorf("tcp request operation(%d) not auth", proto.Op)
		}
	}

	var params struct {
		UserId   int64  `json:"user_id"`
		UserKey  string `json:"user_key"`
		RoomId   string `json:"room_id"`
		Platform string `json:"platform"`
		Token    string `json:"token"`
	}
	if err = json.Unmarshal(proto.Body, &params); err != nil {
		logrus.Errorf("Unmarshal body=%v,err=%v", proto.Body, err)
		return
	}

	// update channel
	ch.UserInfo.UserId = params.UserId
	ch.UserInfo.UserKey = params.UserKey
	ch.UserInfo.RoomId = params.RoomId
	ch.UserInfo.Platform = params.Platform
	if hb, err = s.Connect(ctx, ch, params.Token); err != nil {
		logrus.Errorf("Connect UserInfo=%v, err=%v", ch.UserInfo, err)
		return
	}

	// reply to client
	proto.Op = int32(protocol.OpAuthReply)
	proto.Body = nil
	if err = proto.WriteTCP(writer); err != nil {
		logrus.Errorf("WriteTCP UserInfo=%v, err=%v", ch.UserInfo, err)
		return
	}
	err = writer.Flush()
	return
}
