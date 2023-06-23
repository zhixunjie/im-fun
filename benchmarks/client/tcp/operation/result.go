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
		lastCount = msgCount
		
		fmt.Println(fmt.Sprintf("%s aliveCount=%d msgCount=%d(incr/s:%d)", time.Now().Format("2006-01-02 15:04:05"),
			aliveCount, msgCount, (msgCount-lastCount)/interval))

		time.Sleep(time.Second * time.Duration(interval))
	}
}
