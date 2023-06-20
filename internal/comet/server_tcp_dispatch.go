package comet

import (
	"errors"
	"github.com/zhixunjie/im-fun/api/protocol"
	"github.com/zhixunjie/im-fun/internal/comet/channel"
	"github.com/zhixunjie/im-fun/pkg/buffer/bufio"
	"github.com/zhixunjie/im-fun/pkg/logging"
)

var (
	ErrNotAndProtoToRead = errors.New("not any proto to read")
	ErrTCPWriteError     = errors.New("write err")
)

// dispatchTCP deal any proto send to signal channel（Just like a state machine）
// 可能出现的消息：SendReady（client message） or service job
func (s *Server) dispatchTCP(ch *channel.Channel) {
	logHead := "dispatchTCP"
	var err error
	writer := ch.Writer

	for {
		// wait any message from signal channel（if not, it will block here）
		var proto = ch.Waiting()
		switch protocol.Operation(proto.Op) {
		case protocol.OpProtoReady:
			// 1. read msg from client
			if err = protoReady(ch, writer); err != nil {
				goto fail
			}
		case protocol.OpBatchMsg:
			// 2. write msg to client
			if err = proto.WriteTCP(writer); err != nil {
				goto fail
			}
		case protocol.OpProtoFinish:
			// 3. close channel
			goto fail
		default:
			logging.Errorf(logHead + "unknown proto")
			goto fail
		}
		if err = writer.Flush(); err != nil {
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
func protoReady(ch *channel.Channel, writer *bufio.Writer) error {
	logHead := "protoReady"
	var err error
	var online int32
	var proto *protocol.Proto
	for {
		// 1. read proto from client（）
		proto, err = ch.ProtoAllocator.GetProtoCanRead()
		if err != nil { // err != nil 说明没有东西可读了(not any proto to read)
			//logging.Errorf(logHead+"GetProtoCanRead err=%v", err)
			return nil
		}
		// 2. deal proto
		switch protocol.Operation(proto.Op) {
		case protocol.OpHeartbeatReply:
			if ch.Room != nil {
				online = ch.Room.OnlineNum()
			}
			if err = proto.WriteTCPHeart(writer, online); err != nil {
				logging.Errorf(logHead+"WriteTCPHeart err=%v", err)
				return ErrTCPWriteError
			}
		default:
			// 3. write msg to client
			if err = proto.WriteTCP(writer); err != nil {
				logging.Errorf(logHead+"WriteTCP err=%v", err)
				return ErrTCPWriteError
			}
		}
		proto.Body = nil // avoid memory leak
		ch.ProtoAllocator.AdvReadPointer()
	}
}
