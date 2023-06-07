package comet

import (
	"github.com/zhixunjie/im-fun/internal/comet/conf"
	"github.com/zhixunjie/im-fun/pkg/buffer"
	"github.com/zhixunjie/im-fun/pkg/time"
)

// RoundOptions round options.
type RoundOptions struct {
	Timer     int
	TimerSize int
}

// Round used for connection round-robin get a reader/writer/timer for split big lock.
type Round struct {
	Timers     []time.Timer
	BufferPool buffer.Hash

	options RoundOptions
}

// NewRound new a round struct.
func NewRound(c *conf.Config) (round *Round) {
	round = &Round{
		options: RoundOptions{
			Timer:     c.Protocol.Timer,
			TimerSize: c.Protocol.TimerSize,
		}}

	// make timer
	round.Timers = make([]time.Timer, round.options.Timer)
	for i := 0; i < round.options.Timer; i++ {
		round.Timers[i].Init(round.options.TimerSize)
	}
	// make buffer pool
	round.BufferPool = buffer.NewHash(c.Connect.BufferOptions)
	return
}

// TimerPool get a timer.
func (r *Round) TimerPool(rn int) *time.Timer {
	return &(r.Timers[rn%r.options.Timer])
}
