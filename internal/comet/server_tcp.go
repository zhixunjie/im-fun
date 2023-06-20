package comet

import (
	"context"
	"encoding/json"
	"github.com/zhixunjie/im-fun/api/protocol"
	"github.com/zhixunjie/im-fun/internal/comet/channel"
	"github.com/zhixunjie/im-fun/internal/comet/conf"
	"github.com/zhixunjie/im-fun/pkg/buffer/bytes"
	"github.com/zhixunjie/im-fun/pkg/logging"
	newtimer "github.com/zhixunjie/im-fun/pkg/time"
	"github.com/zhixunjie/im-fun/pkg/websocket"
	"io"
	"math"
	"net"
	"strings"
	"time"
)

func InitTCP(server *Server, numCPU, connType int) (listener *net.TCPListener, err error) {
	var addr *net.TCPAddr
	var addrS []string
	var logHead string

	// 同时支持TCP和WebSocket
	if connType == channel.ConnTypeWebSocket {
		addrS = conf.Conf.Connect.TCP.Bind
		logHead = channel.GetLogHeadByConnType(connType)
	} else {
		addrS = conf.Conf.Connect.Websocket.Bind
		logHead = channel.GetLogHeadByConnType(connType)
	}

	// bind address
	for _, bind := range addrS {
		if addr, err = net.ResolveTCPAddr("tcp", bind); err != nil {
			logging.Errorf(logHead+"ResolveTCPAddr(bind=%v) err=%v", bind, err)
			return
		}
		if listener, err = net.ListenTCP("tcp", addr); err != nil {
			logging.Errorf(logHead+"ListenTCP(bind=%v) err=%v", bind, err)
			return
		}
		logging.Infof(logHead+"server is listening：%s", bind)
		// 启动N个协程，每个协程开启后进行accept（需要使用REUSE解决惊群问题）
		for i := 0; i < numCPU; i++ {
			// TODO 使用协程池进行管理
			go acceptTCP(logHead, connType, server, listener)
		}
	}
	return
}

func acceptTCP(logHead string, connType int, server *Server, listener *net.TCPListener) {
	var conn *net.TCPConn
	var err error
	var r int

	for {
		if conn, err = listener.AcceptTCP(); err != nil {
			// if listener close then return
			logging.Errorf(logHead+"listener.Accept(%s) error=%v", listener.Addr().String(), err)
			return
		}
		// set params for the connection
		if err = conn.SetKeepAlive(server.conf.Connect.TCP.Keepalive); err != nil {
			logging.Errorf(logHead+"conn.SetKeepAlive() error=%v", err)
			return
		}
		if err = conn.SetReadBuffer(server.conf.Connect.TCP.Rcvbuf); err != nil {
			logging.Errorf(logHead+"conn.SetReadBuffer() error=%v", err)
			return
		}
		if err = conn.SetWriteBuffer(server.conf.Connect.TCP.Sndbuf); err != nil {
			logging.Errorf(logHead+"conn.SetWriteBuffer() error=%v", err)
			return
		}
		if r++; r == math.MaxInt {
			r = 0
		}

		// begin to serve
		go func(s *Server, conn *net.TCPConn, r int) {
			var tr = s.round.TimerPool(r)
			var rp = s.round.BufferPool.ReaderPool(r)
			var wp = s.round.BufferPool.WriterPool(r)
			logging.Infof("connect success,LocalAddr=%v,RemoteAddr=%v", conn.LocalAddr().String(), conn.RemoteAddr().String())
			s.serveTCP(logHead, conn, connType, rp, wp, tr)
		}(server, conn, r)
	}
}

// serveTCP serve a tcp connection.
func (s *Server) serveTCP(logHead string, conn *net.TCPConn, connType int, readerPool, writerPool *bytes.Pool, timerPool *newtimer.Timer) {
	logHead = logHead + "serveTCP|"

	var (
		err    error
		proto  *protocol.Proto
		bucket *Bucket
		//lastHb = time.Now()
	)
	var hb time.Duration
	//var trd *newtimer.TimerData
	var ch = channel.NewChannel(s.conf, conn, connType, readerPool, writerPool, timerPool)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var step = 0
	{
		// set timer
		// TODO 暂时把timer关闭，感觉有点问题
		//trd = timerPool.Add(time.Duration(s.conf.Protocol.HandshakeTimeout), func() {
		//	conn.Close()
		//	logging.Errorf("TCP handshake timeout UserInfo=%+v,addr=%v,step=%v,hb=%v",
		//		ch.UserInfo, conn.RemoteAddr().String(), step, hb)
		//})

		// upgrade to websocket
		if connType == channel.ConnTypeWebSocket {
			if err = s.upgradeToWebSocket(ctx, logHead, ch); err != nil {
				logging.Errorf(logHead+"upgradeToWebSocket err=%v,UserInfo=%v", err, ch.UserInfo)
				ch.CleanPath1()
				return
			}
		}

	}
	step = 1
	{
		// auth（check token）
		hb, err = s.auth(ctx, logHead, ch, proto, step)
		if err != nil {
			logging.Errorf(logHead+"auth err=%v,UserInfo=%v,hb=%v", err, ch.UserInfo, hb)
			ch.CleanPath1()
			return
		}
		// set bucket
		bucket = s.AllocBucket(ch.UserInfo.UserKey)
		err = bucket.Put(ch)
		if err != nil {
			logging.Errorf("AllocBucket err=%v,UserInfo=%v", err, ch.UserInfo)
			ch.CleanPath1()
			return
		}
	}
	step = 2
	// set new timer
	//trd.Key = ch.UserInfo.UserKey
	//timerPool.Set(trd, hb)

	// dispatch
	go s.dispatch(logHead, ch)

	// loop to read client msg
	// 数据流：client -> comet -> read -> generate proto -> send protoReady(dispatch proto)
	//hbTime := s.RandHeartbeatTime()
	for {
		if proto, err = ch.ProtoAllocator.GetProtoCanWrite(); err != nil {
			goto fail
		}
		// read msg from client
		// note：if there is no msg，it will block here
		if err = ch.ConnReaderWriter.ReadProto(proto); err != nil {
			goto fail
		}

		// deal with the msg
		if err = s.Operate(ctx, logHead, proto, ch, bucket); err != nil {
			goto fail
		}

		// dispatch msg
		ch.ProtoAllocator.AdvWritePointer()
		ch.SendReady()
	}
fail:
	if err != nil && err != io.EOF && !strings.Contains(err.Error(), "closed") {
		logging.Errorf("UserInfo=%v sth has happened,err=%v", ch.UserInfo, err)
	}
	// 回收相关资源
	bucket.DelChannel(ch)
	ch.CleanPath2()
	if err = s.Disconnect(ctx, ch); err != nil {
		logging.Errorf("Disconnect UserInfo=%+v,err=%v", ch.UserInfo, err)
	}
}

func (s *Server) auth(ctx context.Context, logHead string, ch *channel.Channel, proto *protocol.Proto, step int) (hb time.Duration, err error) {
	logHead = logHead + "auth|"

	// get a proto to write
	proto, err = ch.ProtoAllocator.GetProtoCanWrite()
	if err != nil {
		logging.Errorf(logHead+"GetProtoCanWrite err=%v,UserInfo=%v,step=%v,hb=%v", err, ch.UserInfo, step, hb)
		return
	}
	// 一直读取，直到读取到的Proto的操作类型为protocol.OpAuth
	for {
		if err = ch.ConnReaderWriter.ReadProto(proto); err != nil {
			return
		}
		if protocol.Operation(proto.Op) == protocol.OpAuth {
			break
		} else {
			logging.Errorf(logHead+"tcp request operation(%d) not auth", proto.Op)
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
		logging.Errorf(logHead+"Unmarshal body=%s,err=%v", proto.Body, err)
		return
	}

	// update channel
	ch.UserInfo.UserId = params.UserId
	ch.UserInfo.UserKey = params.UserKey
	ch.UserInfo.RoomId = params.RoomId
	ch.UserInfo.Platform = params.Platform
	if hb, err = s.Connect(ctx, ch, params.Token); err != nil {
		logging.Errorf(logHead+"Connect UserInfo=%v, err=%v", ch.UserInfo, err)
		return
	}

	// reply to client
	proto.Op = int32(protocol.OpAuthReply)
	proto.Body = nil
	if err = ch.ConnReaderWriter.WriteProto(proto); err != nil {
		logging.Errorf(logHead+"WriteTCP UserInfo=%v, err=%v", ch.UserInfo, err)
		return
	}
	err = ch.ConnReaderWriter.Flush()
	return
}

func (s *Server) upgradeToWebSocket(ctx context.Context, logHead string, ch *channel.Channel) (err error) {
	logHead = logHead + "upgradeToWebSocket|"
	conn := ch.Conn

	// read request line && upgrade（websocket独有）
	var req *websocket.Request
	if req, err = websocket.ReadRequest(ch.Reader); err != nil {
		logging.Errorf(logHead+"websocket.ReadRequest err=%v,UserInfo=%+v,addr=%v", err, ch.UserInfo, conn.RemoteAddr().String())
		return
	}
	var wsConn *websocket.Conn
	if wsConn, err = websocket.Upgrade(conn, ch.Reader, ch.Writer, req); err != nil {
		logging.Errorf(logHead+"websocket.Upgrade err=%v,UserInfo=%+v,addr=%v", err, ch.UserInfo, conn.RemoteAddr().String())
		return
	}
	ch.SetWebSocketConnReaderWriter(wsConn)

	return
}
