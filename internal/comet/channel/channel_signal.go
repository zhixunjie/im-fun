package channel

import (
	"github.com/zhixunjie/im-fun/api/protocol"
	"github.com/zhixunjie/im-fun/internal/comet/errors"
)

func (c *Channel) Waiting() *protocol.Proto {
	return <-c.signal
}

func (c *Channel) Push(p *protocol.Proto) (err error) {
	select {
	case c.signal <- p:
	default:
		err = errors.ErrSignalFullMsgDropped
	}
	return
}

func (c *Channel) SendReady() {
	c.signal <- protocol.ProtoReady
}

func (c *Channel) Close() {
	c.signal <- protocol.ProtoFinish
}
