package comet

import (
	"context"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/zhixunjie/im-fun/api/protocol"
	"github.com/zhixunjie/im-fun/internal/comet/channel"
	"github.com/zhixunjie/im-fun/internal/comet/conf"
	"github.com/zhixunjie/im-fun/pkg/buffer/bytes"
	newtimer "github.com/zhixunjie/im-fun/pkg/time"
	"github.com/zhixunjie/im-fun/pkg/websocket"
	"io"
	"math"
	"net"
	"strings"
	"time"
)

func InitWs(server *Server, accept int) (listener *net.TCPListener, err error) {
	var addr *net.TCPAddr
	addrs := conf.Conf.Connect.Websocket.Bind
	for _, bind := range addrs {
		if addr, err = net.ResolveTCPAddr("tcp", bind); err != nil {
			logrus.Errorf("TCP ResolveTCPAddr(bind=%v) err=%v", bind, err)
			return
		}
		if listener, err = net.ListenTCP("tcp", addr); err != nil {
			logrus.Errorf("TCP ListenTCP(bind=%v) err=%v", bind, err)
			return
		}
		logrus.Infof("WebSocket server is listening：%s", bind)
		// 启动N个协程，每个协程开启后进行accept（需要使用REUSE解决惊群问题）
		for i := 0; i < accept; i++ {
			// TODO 使用协程池进行管理
			go acceptWebSocket(server, listener)
		}
	}
	return
}

func acceptWebSocket(server *Server, listener *net.TCPListener) {
	var conn *net.TCPConn
	var err error
	var r int

	for {
		if conn, err = listener.AcceptTCP(); err != nil {
			// if listener close then return
			logrus.Errorf("listener.Accept(%s) error=%v", listener.Addr().String(), err)
			return
		}
		if err = conn.SetKeepAlive(server.conf.Connect.TCP.Keepalive); err != nil {
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
		go serveWebSocketInit(server, conn, r)
		if r++; r == math.MaxInt {
			r = 0
		}
	}
}

func serveWebSocketInit(s *Server, conn *net.TCPConn, r int) {
	var (
		// timer
		tr = s.round.TimerPool(r)
		rp = s.round.BufferPool.ReaderPool(r)
		wp = s.round.BufferPool.WriterPool(r)
	)
	logrus.Infof("connect success,LocalAddr=%v,RemoteAddr=%v",
		conn.LocalAddr().String(), conn.RemoteAddr().String())
	s.serveWebSocket(conn, rp, wp, tr)
}

// serveTCP serve a tcp connection.
func (s *Server) serveWebSocket(conn *net.TCPConn, readerPool, writerPool *bytes.Pool, timerPool *newtimer.Timer) {
	var (
		err    error
		proto  *protocol.Proto
		bucket *Bucket
		//lastHb = time.Now()
	)
	var hb time.Duration
	var trd *newtimer.TimerData
	var ch = channel.NewChannel(s.conf)
	// set reader
	var rb = readerPool.Get()
	ch.Reader.ResetBuffer(conn, rb.Bytes())
	// set writer
	var wb = writerPool.Get()
	ch.Writer.ResetBuffer(conn, wb.Bytes())
	// set user ip
	ch.UserInfo.IP, _, _ = net.SplitHostPort(conn.RemoteAddr().String())
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// set timer
	var step = 0
	// TODO 暂时把timer关闭，感觉有点问题
	//trd = timerPool.Add(time.Duration(s.conf.Protocol.HandshakeTimeout), func() {
	//	conn.Close()
	//	logrus.Errorf("TCP handshake timeout UserInfo=%+v,addr=%v,step=%v,hb=%v",
	//		ch.UserInfo, conn.RemoteAddr().String(), step, hb)
	//})

	releaseOld := func() {
		readerPool.Put(rb)
		writerPool.Put(wb)
		timerPool.Del(trd)
		_ = conn.Close()
	}

	// read request line && upgrade（websocket独有）
	var wsConn *websocket.Conn // websocket
	var req *websocket.Request
	if req, err = websocket.ReadRequest(ch.Reader); err != nil {
		releaseOld()
		logrus.Errorf("websocket.ReadRequest err=%v,UserInfo=%+v,addr=%v",
			err, ch.UserInfo, conn.RemoteAddr().String())
		return
	}
	if wsConn, err = websocket.Upgrade(conn, ch.Reader, ch.Writer, req); err != nil {
		releaseOld()
		logrus.Errorf("websocket.Upgrade err=%v,UserInfo=%+v,addr=%v",
			err, ch.UserInfo, conn.RemoteAddr().String())
		return
	}

	step = 1
	{
		releaseNew := func() {
			readerPool.Put(rb)
			writerPool.Put(wb)
			timerPool.Del(trd)
			_ = wsConn.Close()
		}
		// get a proto to write
		proto, err = ch.ProtoAllocator.GetProtoCanWrite()
		if err != nil {
			releaseNew()
			logrus.Errorf("GetProtoCanWrite err=%v,UserInfo=%v,step=%v,hb=%v", err, ch.UserInfo, step, hb)
			return
		}
		// auth（check token）
		hb, err = s.authWebsocket(ctx, wsConn, ch, proto)
		if err != nil {
			releaseNew()
			logrus.Errorf("authWebsocket err=%v,UserInfo=%v", err, ch.UserInfo)
			return
		}
		// set bucket
		bucket = s.AllocBucket(ch.UserInfo.UserKey)
		err = bucket.Put(ch)
		if err != nil {
			releaseNew()
			logrus.Errorf("AllocBucket err=%v,UserInfo=%v", err, ch.UserInfo)
			return
		}
	}
	step = 2
	// set new timer
	//trd.Key = ch.UserInfo.UserKey
	//timerPool.Set(trd, hb)

	// dispatch
	go s.dispatchWebSocket(wsConn, writerPool, wb, ch)

	// loop to read client msg
	//hbTime := s.RandHeartbeatTime()
	for {
		if proto, err = ch.ProtoAllocator.GetProtoCanWrite(); err != nil {
			goto fail
		}
		// read msg from client
		// note：if there is no msg，it will block here
		if err = proto.ReadWs(wsConn); err != nil {
			goto fail
		}

		// deal with the msg
		if err = s.Operate(ctx, proto, ch, bucket); err != nil {
			goto fail
		}

		// dispatch msg
		ch.ProtoAllocator.AdvWritePointer()
		ch.SendReady()
	}
fail:
	if err != nil && err != io.EOF && !strings.Contains(err.Error(), "closed") {
		logrus.Errorf("UserInfo=%v sth has happened,err=%v", ch.UserInfo, err)
	}
	// 回收相关资源
	bucket.DelChannel(ch)
	timerPool.Del(trd)
	readerPool.Put(rb) // writePool's buffer will be released  in Server.dispatchTCP()
	ch.Close()
	_ = wsConn.Close()
	if err = s.Disconnect(ctx, ch); err != nil {
		logrus.Errorf("Disconnect UserInfo=%+v,err=%v", ch.UserInfo, err)
	}
}

func (s *Server) authWebsocket(ctx context.Context, wsConn *websocket.Conn, ch *channel.Channel, proto *protocol.Proto) (hb time.Duration, err error) {
	logHead := "authWebsocket|"
	// 一直读取，直到读取到的Proto的操作类型为protocol.OpAuth
	for {
		if err = proto.ReadWs(wsConn); err != nil {
			return
		}
		if protocol.Operation(proto.Op) == protocol.OpAuth {
			break
		} else {
			logrus.Errorf(logHead+"tcp request operation(%d) not auth", proto.Op)
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
		logrus.Errorf(logHead+"Unmarshal body=%s,err=%v", proto.Body, err)
		return
	}

	// update channel
	ch.UserInfo.UserId = params.UserId
	ch.UserInfo.UserKey = params.UserKey
	ch.UserInfo.RoomId = params.RoomId
	ch.UserInfo.Platform = params.Platform
	if hb, err = s.Connect(ctx, ch, params.Token); err != nil {
		logrus.Errorf(logHead+"Connect UserInfo=%v, err=%v", ch.UserInfo, err)
		return
	}

	// reply to client
	proto.Op = int32(protocol.OpAuthReply)
	proto.Body = nil
	if err = proto.WriteWs(wsConn); err != nil {
		logrus.Errorf(logHead+"WriteTCP UserInfo=%v, err=%v", ch.UserInfo, err)
		return
	}
	err = wsConn.Flush()
	return
}
