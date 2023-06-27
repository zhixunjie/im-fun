package operation

import (
	"fmt"
	"sync/atomic"
	"time"
)

var (
	mCount int64
	aCount int64
)

func DashBoard() {
	var (
		lastCount int64
		interval  = int64(5)
	)
	for {
		msgCount := atomic.LoadInt64(&mCount)
		aliveCount := atomic.LoadInt64(&aCount)

		// print data
		fmt.Println(fmt.Sprintf("%s interval=%vs,aliveCount=%d msgCount=%d(incr/s:%d,lastCount=%v)", time.Now().Format("2006-01-02 15:04:05"),
			interval, aliveCount, msgCount, (msgCount-lastCount)/interval, lastCount))
		lastCount = msgCount

		time.Sleep(time.Second * time.Duration(interval))
	}
}
