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
		lastCount1 int64
		//lastCount5  int64
		lastCount60 int64
		interval1   = int64(1)
		interval5   = int64(5)
		interval60  = int64(60)
		qps         int64
		qpm         int64
	)

	ticker1 := time.NewTicker(time.Duration(interval1) * time.Second)
	ticker5 := time.NewTicker(time.Duration(interval5) * time.Second)
	ticker60 := time.NewTicker(time.Duration(interval60) * time.Second)

	for {
		select {
		case <-ticker1.C:
			msgCount := atomic.LoadInt64(&mCount)
			if msgCount-lastCount1 > qps {
				qps = msgCount - lastCount1
			}
			lastCount1 = msgCount
		case <-ticker5.C:
			msgCount := atomic.LoadInt64(&mCount)
			aliveCount := atomic.LoadInt64(&aCount)
			fmt.Println(fmt.Sprintf("%s aliveCount=%d msgTotal=%d,qpm=%v,最近%v秒内最大的qps:%d",
				time.Now().Format("2006-01-02 15:04:05"),
				aliveCount, msgCount, qpm, interval60, qps))
			//lastCount5 = msgCount
		case <-ticker60.C:
			msgCount := atomic.LoadInt64(&mCount)
			qpm = msgCount - lastCount60
			lastCount60 = msgCount
			// reset qps
			qps = 0
		}
	}

}
