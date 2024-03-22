package channel

import (
	"github.com/zhixunjie/im-fun/api/protocol"
	"github.com/zhixunjie/im-fun/internal/comet/api"
	"github.com/zhixunjie/im-fun/pkg/logging"
)

func (c *Channel) Waiting() *protocol.Proto {
	return <-c.signal
}

func (c *Channel) Push(p *protocol.Proto) (err error) {
	select {
	case c.signal <- p:
	default:
		err = api.ErrSignalFullMsgDropped
	}
	return
}

func (c *Channel) SendReady() {
	c.signal <- protocol.ProtoReady
}

// SendFinish 通知Write协程结束
func (c *Channel) SendFinish() {
	logging.Infof("traceId=%v] SendFinish|", c.TraceId)
	c.signal <- protocol.ProtoFinish
}
