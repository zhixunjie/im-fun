package gen_id

import "time"

// UTC时间：2024-04-17 00:00:00
var baseTime = time.Date(2024, time.April, 17, 0, 0, 0, 0, time.UTC)

func SeqId() int64 {
	return time.Now().UnixMilli() // 毫秒
}

func SeqId32() int32 {
	return int32(time.Now().Sub(baseTime).Seconds())
}
