package gen_id

import "time"

func GenerateSeqId() int64 {
	return time.Now().UnixMilli() // 毫秒
}
