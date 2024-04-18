package comet

import (
	"context"
	"errors"
	"github.com/zhixunjie/im-fun/api/protocol"
	"github.com/zhixunjie/im-fun/internal/comet/channel"
	"github.com/zhixunjie/im-fun/pkg/logging"
)

var (
	ErrTCPWriteError = errors.New("write err")
)

// dispatch 处理发送到Channel的Proto（Just like a state machine）
func (s *TcpServer) dispatch(ctx context.Context, logHead string, ch *channel.Channel) {
	logHead = logHead + "dispatch|"
	var err error
	var finish bool

	for {
		// waiting for messages from channel（otherwise, it will block here）
		//logging.Infof(logHead + "waiting get proto...")
		proto := ch.Waiting()
		switch protocol.Operation(proto.Op) {
		case protocol.OpProtoReady: // 处理TCP客户端的消息
			// 数据流：client -> [comet] -> read -> send ProtoReady -> ch.Waiting()
			if err = ProtoReady(logHead, ch); err != nil {
				logging.Errorf(logHead+"ProtoReady err=%v", err)
				goto fail
			}
		case protocol.OpBatchMsg: // 下发消息到TCP客户端（批量消息）
			// 数据流：client -> [logic] -> [job] -> [comet] -> ch.Waiting()
			if err = ch.ConnReadWriter.WriteProto(proto); err != nil {
				logging.Errorf(logHead+"WriteProto err=%v", err)
				goto fail
			}
		case protocol.OpProtoFinish: // 结束dispatch的流程
			// 数据流：send -> OpProtoFinish -> channel
			finish = true
			logging.Errorf(logHead+"get OpProtoFinish err=%v", err)
			goto fail
		default:
			logging.Errorf(logHead + "unknown proto")
			goto fail
		}
		// TCP响应：下发TCP消息给给客户端
		if err = ch.ConnReadWriter.Flush(); err != nil {
			goto fail
		}
	}
fail:
	s.dispatchFail(ctx, logHead, ch, err, finish)
}

func (s *TcpServer) dispatchFail(ctx context.Context, logHead string, ch *channel.Channel, err error, finish bool) {
	// check error
	if err != nil {
		logging.Errorf(logHead+"err=%v", err)
	} else {
		logging.Infof(logHead + "fail: sth has happened")
	}

	// clean
	s.cleanAfterFn(ctx, logHead, channel.CleanPath3, ch, nil)

	// 把signal这个channel的消息全部消费掉
	// 防止：其他地方往signal发送东西时，由于signal没有消费方，而导致进入无限的阻塞。
	for !finish {
		finish = ch.Waiting() == protocol.ProtoFinish
	}
	logging.Infof(logHead + "finally ended")
}

// ProtoReady 处理TCP客户端的消息
// 使用for循环处理：如果客户端发送频率非常高，可以通过for循环把连续的proto读取出来了
func ProtoReady(logHead string, ch *channel.Channel) (err error) {
	logHead = logHead + "ProtoReady|"
	var online int32
	var proto *protocol.Proto
	var tErr error

	for {
		// 1. read proto from client
		proto, tErr = ch.ProtoAllocator.GetProtoForRead()
		if tErr != nil {
			if tErr == channel.ErrRingEmpty { // 说明：没有东西可读（此错误无需报错）
				break
			} else {
				err = tErr
			}
		}
		// 2. check proto
		switch protocol.Operation(proto.Op) {
		case protocol.OpHeartbeatReply: // 下发心跳响应消息给客户端
			if ch.Room != nil {
				online = ch.Room.OnlineNum()
			}
			if err = ch.ConnReadWriter.WriteProtoHeart(proto, online); err != nil {
				logging.Errorf(logHead+"WriteTCPHeart err=%v", err)
				return
			}
		default: // 直接把消息下发给客户端（未知的类型）
			//if err = ch.ConnReadWriter.WriteProto(proto); err != nil {
			//	logging.Errorf(logHead+"WriteTCP err=%v", err)
			//	return ErrTCPWriteError
			//}
			//err = ErrTCPWriteError
			logging.Errorf(logHead+"unknown proto=%+v", proto)
			return
		}
		proto.Body = nil // avoid memory leak
		ch.ProtoAllocator.AdvReadPointer()
	}

	return
}
