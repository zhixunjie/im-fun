package comet

import (
	"github.com/zhixunjie/im-fun/internal/comet/conf"
	"github.com/zhixunjie/im-fun/pkg/buffer/bytes"
	"github.com/zhixunjie/im-fun/pkg/time"
)

// RoundOptions round options.
type RoundOptions struct {
	TimerPool *conf.TimerPool
}

// Round used for connection round-robin get a reader/writer/timer for split big lock.
type Round struct {
	Timers     []time.Timer
	BufferHash bytes.Hash

	options RoundOptions
}

// NewRound new a round struct.
func NewRound(c *conf.Config) (round *Round) {
	round = &Round{
		options: RoundOptions{
			TimerPool: c.Protocol.TimerPool,
		},
	}

	hashNum := c.Protocol.TimerPool.HashNum
	initSize := c.Protocol.TimerPool.InitSizeInPool

	// make timer
	round.Timers = make([]time.Timer, hashNum)
	for i := 0; i < hashNum; i++ {
		round.Timers[i].Init(initSize)
	}
	// make buffer pool
	round.BufferHash = bytes.NewHash(c.Connect.BufferOptions)
	return
}

// TimerPool get a timer.
func (r *Round) TimerPool(rn int) *time.Timer {
	return &(r.Timers[rn%r.options.TimerPool.HashNum])
}
