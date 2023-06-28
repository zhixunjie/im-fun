package gen_id

import "time"

func GenerateSeqId() int64 {
	return time.Now().UnixNano()
}
