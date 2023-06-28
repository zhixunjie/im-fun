package comet

import (
	"errors"
	"github.com/zhixunjie/im-fun/api/protocol"
	"github.com/zhixunjie/im-fun/internal/comet/channel"
	"github.com/zhixunjie/im-fun/pkg/logging"
)

var (
	ErrTCPWriteError = errors.New("write err")
)

// dispatch deal any proto send to signal channel（Just like a state machine）
// 可能出现的消息：SendReady（client message） or service job
func (s *Server) dispatch(logHead string, ch *channel.Channel) {
	logHead = logHead + "dispatch|"
	var err error
	var finish bool

	for {
		// blocking here !!!
		// wait any message from signal channel（if not, it will block here）
		//logging.Infof(logHead + "waiting get proto...")
		var proto = ch.Waiting()
		switch protocol.Operation(proto.Op) {
		case protocol.OpProtoReady:
			// case1: read msg from client
			// 数据流：client -> [comet] -> read -> send ProtoReady -> ch.Waiting()
			if err = ProtoReady(logHead, ch); err != nil {
				logging.Errorf(logHead+"ProtoReady err=%v", err)
				goto fail
			}
		case protocol.OpBatchMsg:
			// case2: write msg to client
			// 数据流：client -> [logic] -> [job] -> [comet] -> ch.Waiting()
			if err = ch.ConnReaderWriter.WriteProto(proto); err != nil {
				logging.Errorf(logHead+"WriteProto err=%v", err)
				goto fail
			}
		case protocol.OpProtoFinish:
			// case3: OpProtoFinish is used to end the process of dispatch
			finish = true
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
	// 把signal这个channel的消息全部消费掉
	// 防止：其他地方往signal发送东西时，由于signal没有消费方，而被导致进入无限的阻塞当中。
	for !finish {
		tmp := ch.Waiting() == protocol.ProtoFinish
		finish = tmp
	}
	logging.Infof(logHead + "finally ended")
}

func ProtoReady(logHead string, ch *channel.Channel) error {
	logHead = logHead + "ProtoReady|"
	var err error
	var online int32
	var proto *protocol.Proto

	// 使用for循环：
	// 如果客户端发送消息的频率非常高，那么就可以通过for循环把连续的proto读取出来了
	for {
		// 1. read proto from client
		proto, err = ch.ProtoAllocator.GetProtoCanRead()
		if err != nil {
			// 说明：没有东西可读了（not any proto to read，因为外层使用了for循环，所以这类的情况还是会出现的）
			//logging.Infof(logHead+"GetProtoCanRead err=%v", err)
			break
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
			//if err = ch.ConnReaderWriter.WriteProto(proto); err != nil {
			//	logging.Errorf(logHead+"WriteTCP err=%v", err)
			//	return ErrTCPWriteError
			//}
			logging.Errorf(logHead+"unknown proto=%+v", proto)
			return ErrTCPWriteError
		}
		proto.Body = nil // avoid memory leak
		ch.ProtoAllocator.AdvReadPointer()
	}

	return nil
}
