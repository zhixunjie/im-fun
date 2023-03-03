package comet

import (
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/zhixunjie/im-fun/api/protocol"
	"github.com/zhixunjie/im-fun/internal/comet/channel"
	"github.com/zhixunjie/im-fun/pkg/buffer"
	"github.com/zhixunjie/im-fun/pkg/buffer/bufio"
	"net"
)

var (
	ErrNotAndProtoToRead = errors.New("not any proto to read")
	ErrTCPWriteError     = errors.New("write err")
)

// dispatchTCP deal any proto send to signal channel（Just like a state machine）
// possibility：SendReady（client message） or service job
func (s *Server) dispatchTCP(conn *net.TCPConn, writerPool *buffer.Pool, writeBuf *buffer.Buffer, ch *channel.Channel) {
	var err error
	//var finish bool
	writer := ch.Writer

	for {
		// waiting any message from signal channel（if not, it will block here）
		var waiting = ch.Waiting()
		switch waiting {
		case protocol.ProtoFinish: // close channel
			//finish = true
			goto failed
		case protocol.ProtoReady: // read msg from client
			if err = protoReady(ch, writer); errors.Is(err, ErrTCPWriteError) {
				goto failed
			}
		default: // write msg to client
			if err = waiting.WriteTCP(writer); err != nil {
				goto failed
			}
		}
		if err = writer.Flush(); err != nil {
			break
		}
	}
failed: // TODO 子协程的结束，需要通知到主协程（否则主协程不会结束）
	if err != nil {
		logrus.Errorf("UserInfo=%+v,err=%v", ch.UserInfo, err)
	}
	conn.Close()
	writerPool.Put(writeBuf)
	// must ensure all channel message discard, for reader won't be blocking Signal
	//for !finish {
	//	finish = ch.Waiting() == protocol.ProtoFinish
	//}
}

func protoReady(ch *channel.Channel, writer *bufio.Writer) error {
	var err error
	var online int32
	var proto *protocol.Proto
	for {
		proto, err = ch.ProtoAllocator.GetProtoCanRead()
		if err != nil {
			logrus.Errorf("GetProtoCanRead err=%v", err)
			return ErrNotAndProtoToRead
		}
		if protocol.Operation(proto.Op) == protocol.OpHeartbeatReply {
			if ch.Room != nil {
				online = ch.Room.OnlineNum()
			}
			if err = proto.WriteTCPHeart(writer, online); err != nil {
				logrus.Errorf("WriteTCPHeart err=%v", err)
				return ErrTCPWriteError
			}
		} else {
			// write msg to client
			if err = proto.WriteTCP(writer); err != nil {
				logrus.Errorf("WriteTCP err=%v", err)
				return ErrTCPWriteError
			}
		}
		proto.Body = nil // avoid memory leak
		ch.ProtoAllocator.AdvReadPointer()
	}

	return nil
}
