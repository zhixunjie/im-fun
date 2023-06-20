package comet

import (
	"github.com/zhixunjie/im-fun/api/protocol"
	"github.com/zhixunjie/im-fun/internal/comet/channel"
	"github.com/zhixunjie/im-fun/pkg/logging"
)

// dispatchWebSocket deal any proto send to signal channel（Just like a state machine）
// 可能出现的消息：SendReady（client message） or service job
func (s *Server) dispatchWebSocket(ch *channel.Channel) {
	logHead := "dispatchWebSocket|"
	var err error
	for {
		// wait message from signal channel（if not, it will block here）
		var proto = ch.Waiting()
		switch protocol.Operation(proto.Op) {
		case protocol.OpProtoReady:
			// 1. read msg from client
			if err = protoReadyWebsocket(ch); err != nil {
				goto fail
			}
		case protocol.OpBatchMsg:
			// 2. write msg to client
			if err = ch.ConnReaderWriter.WriteProto(proto); err != nil {
				goto fail
			}
		case protocol.OpProtoFinish:
			// 3. close channel
			goto fail
		default:
			logging.Errorf(logHead + "unknown proto")
			goto fail
		}
		if err = ch.ConnReaderWriter.Flush(); err != nil {
			goto fail
		}
	}
fail:
	if err != nil {
		logging.Errorf(logHead+"UserInfo=%+v,err=%v", ch.UserInfo, err)
	}
	ch.CleanPath3()
}

// 数据流：client -> comet -> read -> generate proto -> send protoReady(dispatch proto) -> deal protoReady
func protoReadyWebsocket(ch *channel.Channel) error {
	var err error
	var online int32
	var proto *protocol.Proto
	for {
		// 1. read proto from client
		proto, err = ch.ProtoAllocator.GetProtoCanRead()
		if err != nil { // err != nil 说明没有东西可读了(not any proto to read)
			//logging.Errorf("GetProtoCanRead err=%v", err)
			return nil
		}
		// 2. deal proto
		switch protocol.Operation(proto.Op) {
		case protocol.OpHeartbeatReply:
			if ch.Room != nil {
				online = ch.Room.OnlineNum()
			}
			if err = ch.ConnReaderWriter.WriteProtoHeart(proto, online); err != nil {
				logging.Errorf("WriteTCPHeart err=%v", err)
				return ErrTCPWriteError
			}
		default:
			// 3. write msg to client
			if err = ch.ConnReaderWriter.WriteProto(proto); err != nil {
				logging.Errorf("WriteTCP err=%v", err)
				return ErrTCPWriteError
			}
		}
		proto.Body = nil // avoid memory leak
		ch.ProtoAllocator.AdvReadPointer()
	}
}
