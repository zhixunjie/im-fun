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
		err = api.ErrChannelSignalFull
	}
	return
}

func (c *Channel) SendReady() {
	logging.Infof("[traceId=%v] SendReady", c.TraceId)
	c.signal <- protocol.ProtoReady
}

// SendFinish 通知Write协程结束
func (c *Channel) SendFinish(logHead string) {
	logging.Infof(logHead+"[traceId=%v] SendFinish", c.TraceId)
	c.signal <- protocol.ProtoFinish
}
