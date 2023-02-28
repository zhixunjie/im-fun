package channel

import (
	"github.com/zhixunjie/im-fun/api/protocol"
	"github.com/zhixunjie/im-fun/internal/comet/errors"
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
	Writer bufio.Writer
	Reader bufio.Reader

	// User Info
	UserInfo *UserInfo

	// use to get proto，reduce GC
	ProtoAllocator Ring
}

// NewChannel new a channel.
func NewChannel(cli, svr int) *Channel {
	c := new(Channel)
	c.signal = make(chan *protocol.Proto, svr)
	return c
}

func (c *Channel) Push(p *protocol.Proto) (err error) {
	select {
	case c.signal <- p:
	default:
		err = errors.ErrSignalFullMsgDropped
	}
	return
}

func (c *Channel) Ready() *protocol.Proto {
	return <-c.signal
}

func (c *Channel) Signal() {
	c.signal <- protocol.ProtoReady
}

func (c *Channel) Close() {
	c.signal <- protocol.ProtoFinish
}
