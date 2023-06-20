package comet

import (
	"context"
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

func InitWs(server *Server, accept int) (listener *net.TCPListener, err error) {
	var addr *net.TCPAddr
	addrs := conf.Conf.Connect.Websocket.Bind
	for _, bind := range addrs {
		if addr, err = net.ResolveTCPAddr("tcp", bind); err != nil {
			logging.Errorf("TCP ResolveTCPAddr(bind=%v) err=%v", bind, err)
			return
		}
		if listener, err = net.ListenTCP("tcp", addr); err != nil {
			logging.Errorf("TCP ListenTCP(bind=%v) err=%v", bind, err)
			return
		}
		logging.Infof("WebSocket server is listening：%s", bind)
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
			logging.Errorf("listener.Accept(%s) error=%v", listener.Addr().String(), err)
			return
		}
		if err = conn.SetKeepAlive(server.conf.Connect.TCP.Keepalive); err != nil {
			logging.Errorf("conn.SetKeepAlive() error=%v", err)
			return
		}
		if err = conn.SetReadBuffer(server.conf.Connect.TCP.Rcvbuf); err != nil {
			logging.Errorf("conn.SetReadBuffer() error=%v", err)
			return
		}
		if err = conn.SetWriteBuffer(server.conf.Connect.TCP.Sndbuf); err != nil {
			logging.Errorf("conn.SetWriteBuffer() error=%v", err)
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
	logging.Infof("connect success,LocalAddr=%v,RemoteAddr=%v",
		conn.LocalAddr().String(), conn.RemoteAddr().String())
	s.serveWebSocket(conn, rp, wp, tr)
}

func (s *Server) upgradeToWebSocket(ctx context.Context, ch *channel.Channel) (err error) {
	conn := ch.Conn

	// read request line && upgrade（websocket独有）
	var req *websocket.Request
	if req, err = websocket.ReadRequest(ch.Reader); err != nil {
		logging.Errorf("websocket.ReadRequest err=%v,UserInfo=%+v,addr=%v", err, ch.UserInfo, conn.RemoteAddr().String())
		return
	}
	var wsConn *websocket.Conn
	if wsConn, err = websocket.Upgrade(conn, ch.Reader, ch.Writer, req); err != nil {
		logging.Errorf("websocket.Upgrade err=%v,UserInfo=%+v,addr=%v", err, ch.UserInfo, conn.RemoteAddr().String())
		return
	}
	ch.SetWebSocketConnReaderWriter(wsConn)

	return
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
	//var trd *newtimer.TimerData
	var ch = channel.NewChannel(s.conf, conn, channel.ConnectionTypeWebSocket, readerPool, writerPool, timerPool)

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
		if err = s.upgradeToWebSocket(ctx, ch); err != nil {
			logging.Errorf("upgradeToWebSocket err=%v,UserInfo=%v", err, ch.UserInfo)
			ch.CleanPath1()
			return
		}
	}

	step = 1
	{
		// auth（check token）
		hb, err = s.auth(ctx, ch, proto, step)
		if err != nil {
			logging.Errorf("auth err=%v,UserInfo=%v,hb=%v", err, ch.UserInfo, hb)
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
	go s.dispatchWebSocket(ch)

	// loop to read client msg
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
		if err = s.Operate(ctx, proto, ch, bucket); err != nil {
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
