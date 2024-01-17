package gen_id

import "time"

func SeqId() int64 {
	return time.Now().UnixMilli() // 毫秒
}
