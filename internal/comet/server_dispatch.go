package comet

import (
	"errors"
	"github.com/zhixunjie/im-fun/api/protocol"
	"github.com/zhixunjie/im-fun/internal/comet/channel"
	"github.com/zhixunjie/im-fun/pkg/logging"
)

var (
	ErrNotAndProtoToRead = errors.New("not any proto to read")
	ErrTCPWriteError     = errors.New("write err")
)

// dispatch deal any proto send to signal channel（Just like a state machine）
// 可能出现的消息：SendReady（client message） or service job
func (s *Server) dispatch(logHead string, ch *channel.Channel) {
	logHead = logHead + "dispatch|,"
	var err error

	for {
		// wait any message from signal channel（if not, it will block here）
		var proto = ch.Waiting()
		switch protocol.Operation(proto.Op) {
		case protocol.OpProtoReady:
			// case1. read msg from client
			// 数据流：client -> comet -> read -> generate proto -> send protoReady -> protoReady
			if err = protoReady(logHead, ch); err != nil {
				logging.Errorf(logHead+"protoReady err=%v", err)
				goto fail
			}
		case protocol.OpBatchMsg:
			// case2. write msg to client
			if err = ch.ConnReaderWriter.WriteProto(proto); err != nil {
				logging.Errorf(logHead+"WriteProto err=%v", err)
				goto fail
			}
		case protocol.OpProtoFinish:
			// case3. close channel
			logging.Errorf(logHead+"get OpProtoFinish err=%v", err)
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
		logging.Errorf(logHead+"err=%v", err)
	}
	ch.CleanPath3()
}

func protoReady(logHead string, ch *channel.Channel) error {
	logHead = logHead + "protoReady|"
	var err error
	var online int32
	var proto *protocol.Proto
	for {
		// 1. read proto from client
		proto, err = ch.ProtoAllocator.GetProtoCanRead()
		if err != nil { // err != nil 说明没有东西可读了(not any proto to read)
			//logging.Errorf(logHead+"GetProtoCanRead err=%v", err)
			return nil
		}
		// 2. check proto
		switch protocol.Operation(proto.Op) {
		case protocol.OpHeartbeatReply:
			// 2.1 reply heartbeat
			if ch.Room != nil {
				online = ch.Room.OnlineNum()
			}
			if err = ch.ConnReaderWriter.WriteProtoHeart(proto, online); err != nil {
				logging.Errorf(logHead+"WriteTCPHeart err=%v", err)
				return ErrTCPWriteError
			}
		default:
			// 2.2 write msg to client directly
			if err = ch.ConnReaderWriter.WriteProto(proto); err != nil {
				logging.Errorf(logHead+"WriteTCP err=%v", err)
				return ErrTCPWriteError
			}
		}
		proto.Body = nil // avoid memory leak
		ch.ProtoAllocator.AdvReadPointer()
	}
}
