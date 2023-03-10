package channel

import (
	"github.com/zhixunjie/im-fun/api/protocol"
	"github.com/zhixunjie/im-fun/internal/comet/conf"
	"github.com/zhixunjie/im-fun/pkg/buffer/bufio"
)

// Channel used by message pusher send msg to write goroutine.
type Channel struct {
	// The Room that Channel(User) is In
	Room   *Room
	signal chan *protocol.Proto
	Next   *Channel
	Prev   *Channel

	// use bufio to reuse buffer
	Writer *bufio.Writer
	Reader *bufio.Reader

	// User Info
	UserInfo *UserInfo

	// use to get proto，reduce GC
	ProtoAllocator Ring
}

// NewChannel new a channel.
func NewChannel(conf *conf.Config) *Channel {
	ch := new(Channel)
	ch.ProtoAllocator.Init(uint64(conf.Protocol.ClientProtoNum))
	ch.signal = make(chan *protocol.Proto, conf.Protocol.ServerProtoNum)
	ch.Reader = new(bufio.Reader)
	ch.Writer = new(bufio.Writer)
	ch.UserInfo = new(UserInfo)
	return ch
}
