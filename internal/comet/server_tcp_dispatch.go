package comet

import (
	"errors"
	"github.com/golang/glog"
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

func (s *Server) dispatchTCP(conn *net.TCPConn, writer *bufio.Writer, writerPool *buffer.Pool, writeBuf *buffer.Buffer, ch *channel.Channel) {
	var err error
	var finish bool

	for {
		var waiting = ch.Waiting()
		switch waiting {
		case protocol.ProtoFinish:
			finish = true
			goto failed
		case protocol.ProtoReady:
			err = dealReady(ch, writer)
			if errors.Is(err, ErrTCPWriteError) {
				goto failed
			}
		default:
			if err = waiting.WriteTCP(writer); err != nil {
				goto failed
			}
		}
		if err = writer.Flush(); err != nil {
			break
		}
	}
failed:
	if err != nil {
		glog.Errorf("key=%s, err=%v", ch.UserInfo.UserKey, err)
	}
	conn.Close()
	writerPool.Put(writeBuf)
	// must ensure all channel message discard, for reader won't be blocking Signal
	for !finish {
		finish = ch.Waiting() == protocol.ProtoFinish
	}
}

func dealReady(ch *channel.Channel, writer *bufio.Writer) error {
	var err error
	var online int32
	var proto *protocol.Proto
	for {
		proto, err = ch.ProtoAllocator.GetProtoCanRead()
		if err != nil {
			glog.Errorf("GetProtoCanRead err=%v", err)
			break
		}
		if protocol.Operation(proto.Op) == protocol.OpHeartbeatReply {
			if ch.Room != nil {
				online = ch.Room.OnlineNum()
			}
			if err = proto.WriteTCPHeart(writer, online); err != nil {
				return ErrTCPWriteError
			}
		} else {
			if err = proto.WriteTCP(writer); err != nil {
				return ErrTCPWriteError
			}
		}
		proto.Body = nil // avoid memory leak
		ch.ProtoAllocator.AdvReadPointer()
	}

	return nil
}
