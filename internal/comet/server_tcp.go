package comet

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/zhixunjie/im-fun/api/pb"
	"github.com/zhixunjie/im-fun/api/protocol"
	"github.com/zhixunjie/im-fun/internal/comet/channel"
	"github.com/zhixunjie/im-fun/internal/comet/conf"
	"github.com/zhixunjie/im-fun/pkg/gmodel"
	"github.com/zhixunjie/im-fun/pkg/logging"
	"github.com/zhixunjie/im-fun/pkg/websocket"
	"io"
	"math"
	"net"
	"strings"
	"time"
)

func InitTCP(server *TcpServer, numCPU int, connType channel.ConnType) (listener *net.TCPListener, err error) {
	var addr *net.TCPAddr
	var bindList []string

	// 同时支持TCP和WebSocket
	logHead := channel.LogHeadByConnType(connType)
	if connType == channel.ConnTypeWebSocket {
		bindList = conf.Conf.Connect.Websocket.Bind
	} else {
		bindList = conf.Conf.Connect.TCP.Bind
	}

	// bind address
	for _, bindItem := range bindList {
		if addr, err = net.ResolveTCPAddr("tcp", bindItem); err != nil {
			logging.Errorf(logHead+"ResolveTCPAddr(bind=%v) err=%v", bindItem, err)
			return
		}
		if listener, err = net.ListenTCP("tcp", addr); err != nil {
			logging.Errorf(logHead+"ListenTCP(bind=%v) err=%v", bindItem, err)
			return
		}
		logging.Infof(logHead+"server is listening：%s", bindItem)
		// 启动N个协程，每个协程开启后进行accept（需要使用REUSE解决惊群问题）
		for i := 0; i < numCPU; i++ {
			// TODO 使用协程池进行管理
			go accept(logHead, connType, server, listener)
		}
	}
	return
}

func accept(logHead string, connType channel.ConnType, server *TcpServer, listener *net.TCPListener) {
	var conn *net.TCPConn
	var err error
	var r int
	var traceId = time.Now().UnixNano()
	logHead = fmt.Sprintf("[traceId=%v] ", traceId) + logHead

	for {
		if conn, err = listener.AcceptTCP(); err != nil {
			// if listener close then return
			logging.Errorf(logHead+"listener.Accept=%s error=%v", listener.Addr().String(), err)
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
		go func(s *TcpServer, conn *net.TCPConn, r int) {
			var tr = s.round.TimerPool(r)
			var rp = s.round.BufferHash.ReaderPool(r)
			var wp = s.round.BufferHash.WriterPool(r)
			var ch = channel.NewChannel(s.conf, conn, traceId, connType, rp, wp, tr)

			// let get to server
			logging.Infof(logHead+"connect success,addr_local=%v,addr_remote=%v", conn.LocalAddr().String(), conn.RemoteAddr().String())
			s.serveTCP(logHead, ch, connType)
		}(server, conn, r)
	}
}

// 清除资源
func (s *TcpServer) cleanAfterFn(ctx context.Context, logHead string, cleanPath channel.CleanPath, ch *channel.Channel, bucket *Bucket) {
	var err error
	logging.Infof(logHead + "run cleanAfterFn")

	switch cleanPath {
	case channel.CleanPath1:
		// clean
		ch.CleanPath1()
	case channel.CleanPath2:
		// delete channel
		if bucket != nil {
			bucket.DelChannel(ch)
		}
		// clean
		ch.CleanPath2()
		// rpc: remove redis binding info
		if err = s.Disconnect(ctx, ch); err != nil {
			logging.Errorf(logHead+"Disconnect,err=%v", err)
		}
	case channel.CleanPath3:
		// clean
		ch.CleanPath3()
	}
}

func (s *TcpServer) readLoopFail(ctx context.Context, logHead string, ch *channel.Channel, bucket *Bucket, err error) {
	// check error
	if err != nil {
		switch {
		case err == io.EOF, strings.Contains(err.Error(), "closed") == true:
			logging.Infof(logHead+"fail: err=%v (client close or server close by dispatch)", err)
		default:
			logging.Errorf(logHead+"fail: sth has happened,err=%v", err)
		}
	} else {
		logging.Errorf(logHead + "fail: sth has happened")
	}

	// clean
	s.cleanAfterFn(ctx, logHead, channel.CleanPath2, ch, bucket)
}

// loop to read client msg
// 数据流：client -> [comet] -> read -> send protoReady -> dispatch
func (s *TcpServer) readLoop(ctx context.Context, logHead string, ch *channel.Channel, bucket *Bucket) {
	var err error
	var proto *protocol.Proto

	for {
		if proto, err = ch.ProtoAllocator.GetProtoForWrite(); err != nil {
			logging.Errorf(logHead+"GetProtoForWrite,err=%v", err)
			goto fail
		}
		// read msg from client
		// note：if there is no msg，it will block here！！！！
		//logging.Infof(logHead + "waiting proto from client...")
		if err = ch.ConnReadWriter.ReadProto(proto); err != nil {
			//logging.Errorf(logHead+"ReadProto err=%v", err)
			goto fail
		}

		// handle msg
		if err = s.handleClientMsg(ctx, logHead, proto, ch, bucket); err != nil {
			logging.Errorf(logHead+"handleClientMsg err=%v", err)
			goto fail
		}

		// dispatch msg
		ch.ProtoAllocator.AdvWritePointer()
		ch.SendReady()
	}
fail:
	s.readLoopFail(ctx, logHead, ch, bucket, err)
}

// serveTCP serve a tcp connection.
func (s *TcpServer) serveTCP(logHead string, ch *channel.Channel, connType channel.ConnType) {
	logHead = logHead + "serveTCP|"

	var (
		err error

		bucket *Bucket
	)
	var hbCfg *pb.HbCfg
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fn1 := func() {
		defer func() {
			if err != nil {
				s.cleanAfterFn(ctx, logHead, channel.CleanPath1, ch, nil)
			}
		}()

		// set timer
		s.SetTimerHandshake(logHead, ch)

		// upgrade to websocket
		if connType == channel.ConnTypeWebSocket {
			if err = s.upgradeToWebSocket(ctx, logHead, ch); err != nil {
				logging.Errorf(logHead+"upgradeToWebSocket err=%v", err)
				return
			}
		}

		// auth（check token）
		hbCfg, bucket, err = s.auth(ctx, logHead, ch)
		if err != nil {
			logging.Errorf(logHead+"auth err=%v", err)
			return
		}
	}
	fn1()

	fn2 := func() {
		defer func() {
			if err != nil {
				s.cleanAfterFn(ctx, logHead, channel.CleanPath2, ch, bucket)
			}
		}()

		// allocate bucket
		bucket = s.AllocBucket(ch.UserInfo.TcpSessionId.ToString())
		err = bucket.Put(ch)
		if err != nil {
			logging.Errorf(logHead+"AllocBucket err=%v", err)
			return
		}
		// recreate: timer
		s.SetTimerHeartbeat(ctx, logHead, ch, hbCfg, bucket)
	}
	fn2()

	// routine new: dispatch
	go s.dispatch(ctx, logHead, ch)

	// routine main: read in loop
	s.readLoop(ctx, logHead, ch, bucket)
}

// 一直读取，直到读取到的Proto操作类型为：protocol.OpAuth
func (s *TcpServer) getAuthProto(logHead string, ch *channel.Channel, proto *protocol.Proto) (err error) {
	for {
		if err = ch.ConnReadWriter.ReadProto(proto); err != nil {
			logging.Infof(logHead+"ReadProto err=%v", err)
			return
		}
		if protocol.Operation(proto.Op) == protocol.OpAuth {
			return
		}
		logging.Infof(logHead+"tcp request op=%d, but not auth", proto.Op)
	}
}

func (s *TcpServer) auth(ctx context.Context, logHead string, ch *channel.Channel) (hbCfg *pb.HbCfg, bucket *Bucket, err error) {
	logHead += "auth|"

	// 获取：proto（用于写入）
	proto, err := ch.ProtoAllocator.GetProtoForWrite()
	if err != nil {
		logging.Errorf(logHead+"GetProtoForWrite err=%v", err)
		return
	}

	// 读取：auth信息
	err = s.getAuthProto(logHead, ch, proto)
	if err != nil {
		return
	}

	// 解析：授权信息
	authParams := new(channel.AuthParams)
	if err = json.Unmarshal(proto.Body, authParams); err != nil {
		logging.Errorf(logHead+"Unmarshal body=%s,err=%v", proto.Body, err)
		return
	}
	authParams.UserInfo.IP = ch.UserInfo.IP

	// RPC：建立绑定关系
	if hbCfg, err = s.Connect(ctx, authParams); err != nil {
		logging.Errorf(logHead+"Connect err=%v,params=%+v", err, authParams)
		return
	}
	ch.UserInfo = authParams.UserInfo
	logging.Infof(logHead+"update user info after Connect[%v]", ch.UserInfo)

	// TCP响应：下发TCP消息给给客户端（授权结果）
	proto.Op = int32(protocol.OpAuthResp)
	proto.Seq = gmodel.NewSeqId32()
	proto.Body = nil
	if err = ch.ConnReadWriter.WriteProto(proto); err != nil {
		logging.Errorf(logHead+"WriteTCP err=%v,UserInfo=%v", err, ch.UserInfo)
		return
	}
	err = ch.ConnReadWriter.Flush()

	return
}

func (s *TcpServer) upgradeToWebSocket(ctx context.Context, logHead string, ch *channel.Channel) (err error) {
	conn := ch.Conn
	logHead += fmt.Sprintf("upgradeToWebSocket,UserInfo=%+v,addr=%v|", ch.UserInfo, conn.RemoteAddr().String())

	// read request line && upgrade（websocket独有）
	var req *websocket.Request
	if req, err = websocket.ReadRequest(ch.Reader); err != nil {
		logging.Errorf(logHead+"ReadRequest err=%v", err)
		return
	}
	var wsConn *websocket.Conn
	if wsConn, err = websocket.Upgrade(conn, ch.Reader, ch.Writer, req); err != nil {
		logging.Errorf(logHead+"Upgrade err=%v", err)
		return
	}
	ch.SetWebSocketConnReaderWriter(wsConn)

	return
}

// SetTimerHandshake set timer: for handshake
func (s *TcpServer) SetTimerHandshake(logHead string, ch *channel.Channel) {
	conn := ch.Conn
	duration := time.Duration(s.conf.Protocol.HandshakeTimeout)

	ch.Trd = ch.TimerPool.Add(duration, func() {
		err := conn.Close()
		logging.Errorf(logHead+"TCP handshake timeout,err=%v", err)
	})

	return
}

// SetTimerHeartbeat set timer: for heartbeat
func (s *TcpServer) SetTimerHeartbeat(ctx context.Context, logHead string, ch *channel.Channel, hbCfg *pb.HbCfg, bucket *Bucket) {
	logHead += fmt.Sprintf("SetTimerHeartbeat,hbCfg=%+v|", hbCfg)
	hbInterval := time.Duration(hbCfg.Interval) * time.Second
	hbExpire := hbInterval * time.Duration(hbCfg.FailCount)

	if hbExpire.Seconds() == 0 {
		logging.Errorf(logHead + "hbDuration not allow")
		return
	}

	ch.TimerPool.Del(ch.Trd)
	ch.Trd = ch.TimerPool.Add(hbExpire, func() {
		logging.Errorf(logHead + "trigger timer(hbExpire)")
		s.cleanAfterFn(ctx, logHead, channel.CleanPath2, ch, bucket)
	})
	ch.LastHb = time.Now()
	ch.HbExpire = hbExpire
	ch.HbInterval = hbInterval
	logging.Infof(logHead+"timer set success,params(LastHb=%v,HbExpire=%v)", ch.LastHb, ch.HbExpire)
}

// ResetTimerHeartbeat reset timer: for heartbeat
func (s *TcpServer) ResetTimerHeartbeat(ctx context.Context, logHead string, ch *channel.Channel) {
	if ch.HbExpire.Seconds() == 0 {
		logging.Errorf(logHead + "hbDuration not allow")
		return
	}

	ch.TimerPool.Set(ch.Trd, ch.HbExpire)
}
