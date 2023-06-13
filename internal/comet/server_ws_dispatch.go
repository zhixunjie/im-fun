package comet

import (
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/zhixunjie/im-fun/api/protocol"
	"github.com/zhixunjie/im-fun/internal/comet/channel"
	"github.com/zhixunjie/im-fun/pkg/buffer/bytes"
	"github.com/zhixunjie/im-fun/pkg/websocket"
)

// dispatchWebSocket deal any proto send to signal channel（Just like a state machine）
// 可能出现的消息：SendReady（client message） or service job
func (s *Server) dispatchWebSocket(wsConn *websocket.Conn, writerPool *bytes.Pool, writeBuf *bytes.Buffer, ch *channel.Channel) {
	var err error
	for {
		// wait message from signal channel（if not, it will block here）
		var proto = ch.Waiting()
		switch protocol.Operation(proto.Op) {
		case protocol.OpProtoFinish:
			// 1. close channel
			goto fail
		case protocol.OpProtoReady:
			// 2. read msg from client
			// 链路：client -> server -> read -> proto -> send protoReady
			if err = protoReadyWebsocket(ch, wsConn); errors.Is(err, ErrTCPWriteError) {
				goto fail
			}
		default: // write msg to client
			if err = proto.WriteWs(wsConn); err != nil {
				goto fail
			}
		}
		if err = wsConn.Flush(); err != nil {
			break
		}
	}
fail: // TODO 子协程的结束，需要通知到主协程（否则主协程不会结束）
	if err != nil {
		logrus.Errorf("UserInfo=%+v,err=%v", ch.UserInfo, err)
	}
	_ = wsConn.Close()
	writerPool.Put(writeBuf)
}

func protoReadyWebsocket(ch *channel.Channel, wsConn *websocket.Conn) error {
	var err error
	var online int32
	var proto *protocol.Proto
	for {
		// 1. read proto from client
		proto, err = ch.ProtoAllocator.GetProtoCanRead()
		if err != nil {
			logrus.Errorf("GetProtoCanRead err=%v", err)
			return ErrNotAndProtoToRead
		}
		// 2. deal proto
		switch protocol.Operation(proto.Op) {
		case protocol.OpHeartbeatReply:
			if ch.Room != nil {
				online = ch.Room.OnlineNum()
			}
			if err = proto.WriteWsHeart(wsConn, online); err != nil {
				logrus.Errorf("WriteTCPHeart err=%v", err)
				return ErrTCPWriteError
			}
		default:
			// write msg to client
			if err = proto.WriteWs(wsConn); err != nil {
				logrus.Errorf("WriteTCP err=%v", err)
				return ErrTCPWriteError
			}
		}
		proto.Body = nil // avoid memory leak
		ch.ProtoAllocator.AdvReadPointer()
	}
}
