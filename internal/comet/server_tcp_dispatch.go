package comet

import (
	"errors"
	"github.com/zhixunjie/im-fun/api/protocol"
	"github.com/zhixunjie/im-fun/internal/comet/channel"
	"github.com/zhixunjie/im-fun/pkg/buffer/bufio"
	"github.com/zhixunjie/im-fun/pkg/buffer/bytes"
	"github.com/zhixunjie/im-fun/pkg/logging"
	"net"
)

var (
	ErrNotAndProtoToRead = errors.New("not any proto to read")
	ErrTCPWriteError     = errors.New("write err")
)

// dispatchTCP deal any proto send to signal channel（Just like a state machine）
// 可能出现的消息：SendReady（client message） or service job
func (s *Server) dispatchTCP(conn *net.TCPConn, writerPool *bytes.Pool, writeBuf *bytes.Buffer, ch *channel.Channel) {
	logHead := "dispatchTCP"
	var err error
	writer := ch.Writer

	for {
		// wait any message from signal channel（if not, it will block here）
		var proto = ch.Waiting()
		switch protocol.Operation(proto.Op) {
		case protocol.OpProtoReady:
			// 1. read msg from client
			if err = protoReady(ch, writer); errors.Is(err, ErrTCPWriteError) {
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
fail: // TODO 子协程的结束，需要通知到主协程（否则主协程不会结束）
	if err != nil {
		logging.Errorf(logHead+"UserInfo=%+v,err=%v", ch.UserInfo, err)
	}
	_ = conn.Close()
	writerPool.Put(writeBuf)
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
		if err != nil {
			logging.Errorf(logHead+"GetProtoCanRead err=%v", err)
			return ErrNotAndProtoToRead
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
