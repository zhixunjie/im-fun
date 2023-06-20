package comet

import (
	"context"
	"github.com/zhixunjie/im-fun/internal/comet/channel"
	"github.com/zhixunjie/im-fun/internal/comet/conf"
	"github.com/zhixunjie/im-fun/pkg/logging"
	"github.com/zhixunjie/im-fun/pkg/websocket"
	"math"
	"net"
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
	s.serveTCP(conn, channel.ConnectionTypeWebSocket, rp, wp, tr)
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