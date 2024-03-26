package comet

import (
	"github.com/zhixunjie/im-fun/internal/comet/conf"
	"github.com/zhixunjie/im-fun/pkg/buffer/bytes"
	"github.com/zhixunjie/im-fun/pkg/time"
)

// RoundOptions round options.
type RoundOptions struct {
	TimerHashNum      int
	InitSizeTimerPool int
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
			TimerHashNum:      c.Protocol.TimerHashNum,
			InitSizeTimerPool: c.Protocol.InitSizeTimerPool,
		}}

	// make timer
	round.Timers = make([]time.Timer, round.options.TimerHashNum)
	for i := 0; i < round.options.TimerHashNum; i++ {
		round.Timers[i].Init(round.options.InitSizeTimerPool)
	}
	// make buffer pool
	round.BufferHash = bytes.NewHash(c.Connect.BufferOptions)
	return
}

// TimerPool get a timer.
func (r *Round) TimerPool(rn int) *time.Timer {
	return &(r.Timers[rn%r.options.TimerHashNum])
}
