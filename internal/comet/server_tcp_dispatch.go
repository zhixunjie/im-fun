package comet

import (
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/zhixunjie/im-fun/api/protocol"
	"github.com/zhixunjie/im-fun/internal/comet/channel"
	"github.com/zhixunjie/im-fun/pkg/buffer/bufio"
	"github.com/zhixunjie/im-fun/pkg/buffer/bytes"
	"net"
)

var (
	ErrNotAndProtoToRead = errors.New("not any proto to read")
	ErrTCPWriteError     = errors.New("write err")
)

// dispatchTCP deal any proto send to signal channel（Just like a state machine）
// 可能出现的消息：SendReady（client message） or service job
func (s *Server) dispatchTCP(conn *net.TCPConn, writerPool *bytes.Pool, writeBuf *bytes.Buffer, ch *channel.Channel) {
	var err error
	writer := ch.Writer

	for {
		// wait any message from signal channel（if not, it will block here）
		var proto = ch.Waiting()
		switch proto {
		case protocol.ProtoFinish: // close channel
			goto fail
		case protocol.ProtoReady: // read msg from client
			if err = protoReady(ch, writer); errors.Is(err, ErrTCPWriteError) {
				goto fail
			}
		default: // write msg to client
			if err = proto.WriteTCP(writer); err != nil {
				goto fail
			}
		}
		if err = writer.Flush(); err != nil {
			break
		}
	}
fail: // TODO 子协程的结束，需要通知到主协程（否则主协程不会结束）
	if err != nil {
		logrus.Errorf("UserInfo=%+v,err=%v", ch.UserInfo, err)
	}
	_ = conn.Close()
	writerPool.Put(writeBuf)
}

func protoReady(ch *channel.Channel, writer *bufio.Writer) error {
	var err error
	var online int32
	var proto *protocol.Proto
	for {
		// read proto from client
		proto, err = ch.ProtoAllocator.GetProtoCanRead()
		if err != nil {
			logrus.Errorf("GetProtoCanRead err=%v", err)
			return ErrNotAndProtoToRead
		}
		// deal proto
		if protocol.Operation(proto.Op) == protocol.OpHeartbeatReply {
			if ch.Room != nil {
				online = ch.Room.OnlineNum()
			}
			if err = proto.WriteTCPHeart(writer, online); err != nil {
				logrus.Errorf("WriteTCPHeart err=%v", err)
				return ErrTCPWriteError
			}
		} else { // write msg to client
			if err = proto.WriteTCP(writer); err != nil {
				logrus.Errorf("WriteTCP err=%v", err)
				return ErrTCPWriteError
			}
		}
		proto.Body = nil // avoid memory leak
		ch.ProtoAllocator.AdvReadPointer()
	}
}
