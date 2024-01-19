package gen_id

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/cast"
	"time"
)

// MsgId 生成msg_id
// - msg_id要求全局唯一
// - msg_id跟largerId的后4位是相同的（slotId其实就是largerId）
func MsgId(ctx context.Context, mem *redis.Client, slotId uint64) (msgId uint64, err error) {
	// redis：每秒一个key，在key上执行原子操作+1
	ts := time.Now().Unix()
	key := keyMsgId(ts)

	// incr
	afterIncr, err := incNum(ctx, mem, key, msgKeyExpire)
	if err != nil {
		return
	}

	// msg_id的组成部分：[ 10位：相对时间戳 | 6位：自增id | 4位：槽id ]
	// 槽id的作用：使用msg_id也能定位到对应的数据库和数据表
	timeOffset := ts - baseTimeStampOffset
	idStr := fmt.Sprintf("%d%06d%04d", timeOffset, afterIncr%1000000, slotId%10000)
	msgId = cast.ToUint64(idStr)

	return
}

// incNum 每秒一个Key，进行累加
func incNum(ctx context.Context, mem *redis.Client, key string, expireSec int) (value int64, err error) {
	value, err = mem.IncrBy(ctx, key, 1).Result()
	if err != nil {
		return
	}
	// 这里的命令可能会失败
	// 解决办法：lua脚本：https://gitee.com/jasonzxj/LearnGo/blob/master/use/pkg/redis/goredis/lua/atomic/incry_expire.go
	if value == 1 {
		_, err = mem.Do(ctx, "EXPIRE", key, expireSec).Result()
		if err != nil {
			return
		}
	}
	return
}
