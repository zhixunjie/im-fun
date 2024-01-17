package comet

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/zhixunjie/im-fun/api/protocol"
	"github.com/zhixunjie/im-fun/internal/comet/channel"
	"github.com/zhixunjie/im-fun/internal/comet/conf"
	"github.com/zhixunjie/im-fun/pkg/gen_id"
	"github.com/zhixunjie/im-fun/pkg/logging"
	"github.com/zhixunjie/im-fun/pkg/utils"
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

	// 同时支持TCP和WebSocket
	logHead := channel.GetLogHeadByConnType(connType)
	if connType == channel.ConnTypeWebSocket {
		addrS = conf.Conf.Connect.Websocket.Bind
	} else {
		addrS = conf.Conf.Connect.TCP.Bind
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
			go accept(logHead, connType, server, listener)
		}
	}
	return
}

func accept(logHead string, connType int, server *Server, listener *net.TCPListener) {
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
		go func(s *Server, conn *net.TCPConn, r int) {
			var tr = s.round.TimerPool(r)
			var rp = s.round.BufferPool.ReaderPool(r)
			var wp = s.round.BufferPool.WriterPool(r)
			var ch = channel.NewChannel(s.conf, conn, traceId, connType, rp, wp, tr)

			// let get to server
			logging.Infof(logHead+"connect success,addr_local=%v,addr_remote=%v",
				conn.LocalAddr().String(), conn.RemoteAddr().String())
			s.serveTCP(logHead, ch, connType)
		}(server, conn, r)
	}
}

// serveTCP serve a tcp connection.
func (s *Server) serveTCP(logHead string, ch *channel.Channel, connType int) {
	logHead = logHead + "serveTCP|"

	var (
		err    error
		proto  *protocol.Proto
		bucket *Bucket
		//lastHb = time.Now()
	)
	var hb time.Duration
	//var trd *newtimer.TimerData

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
				logging.Errorf(logHead+"upgradeToWebSocket err=%v", err)
				ch.CleanPath1()
				return
			}
		}

	}
	step = 1
	{
		// auth（check token）
		hb, err = s.auth(ctx, logHead, ch, step)
		if err != nil {
			logging.Errorf(logHead+"auth err=%v,hb=%v", err, hb)
			ch.CleanPath1()
			return
		}

		// set bucket
		bucket = s.AllocBucket(ch.UserInfo.UserKey)
		err = bucket.Put(ch)
		if err != nil {
			logging.Errorf(logHead+"AllocBucket err=%v", err)
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
	// 数据流：client -> [comet] -> read -> send protoReady -> dispatch
	//hbTime := s.RandHeartbeatTime()
	for {
		if proto, err = ch.ProtoAllocator.GetProtoCanWrite(); err != nil {
			logging.Errorf(logHead+"GetProtoCanWrite,err=%v", err)
			goto fail
		}
		// blocking here !!!
		// read msg from client
		// note：if there is no msg，it will block here
		//logging.Infof(logHead + "waiting proto from client...")
		if err = ch.ConnReaderWriter.ReadProto(proto); err != nil {
			//logging.Errorf(logHead+"ReadProto err=%v", err)
			goto fail
		}

		// deal with the msg
		if err = s.Operate(ctx, logHead, proto, ch, bucket); err != nil {
			logging.Errorf(logHead+"Operate err=%v", err)
			goto fail
		}

		// dispatch msg
		ch.ProtoAllocator.AdvWritePointer()
		ch.SendReady()
	}
fail:
	// check error
	if err != nil {
		switch {
		case err == io.EOF, strings.Contains(err.Error(), "closed") == true:
			logging.Infof(logHead+"fail: err=%v (client close or server close by dispatch)", err)
		default:
			logging.Errorf(logHead+"fail: sth has happened,err=%v", err)
		}
	} else {
		logging.Infof(logHead + "fail: sth has happened")
	}
	// 回收相关资源
	bucket.DelChannel(ch)
	ch.CleanPath2()
	if err = s.Disconnect(ctx, ch); err != nil {
		logging.Errorf(logHead+"Disconnect,err=%v", err)
	}
}

func (s *Server) auth(ctx context.Context, logHead string, ch *channel.Channel, step int) (hb time.Duration, err error) {
	logHead = logHead + "auth|"

	// get a proto to write
	var proto *protocol.Proto
	proto, err = ch.ProtoAllocator.GetProtoCanWrite()
	if err != nil {
		logging.Errorf(logHead+"GetProtoCanWrite err=%v,step=%v,hb=%v", err, step, hb)
		return
	}
	// 一直读取，直到读取到的Proto的操作类型为protocol.OpAuth
	for {
		if err = ch.ConnReaderWriter.ReadProto(proto); err != nil {
			logging.Errorf(logHead+"ReadProto err=%v", err)
			return
		}
		if protocol.Operation(proto.Op) == protocol.OpAuth {
			break
		}
		logging.Errorf(logHead+"tcp request op=%d,but not auth", proto.Op)
	}

	params := new(channel.AuthParams)
	if err = json.Unmarshal(proto.Body, params); err != nil {
		logging.Errorf(logHead+"Unmarshal body=%s,err=%v", proto.Body, err)
		return
	}

	// update channel
	newUserKey := utils.GetMergeUserKey(params.UserId, params.UserKey)
	params.UserKey = newUserKey
	params.IP = ch.UserInfo.IP
	if hb, err = s.Connect(ctx, params); err != nil {
		logging.Errorf(logHead+"Connect err=%v,params=%+v", err, params)
		return
	}
	ch.UserInfo = &params.UserInfo
	logging.Infof(logHead+"update user info after Connect,[%v]", ch.UserInfo)

	// reply to client
	proto.Op = int32(protocol.OpAuthReply)
	proto.Seq = int32(gen_id.SeqId())
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
