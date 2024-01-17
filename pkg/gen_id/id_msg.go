package gen_id

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/cast"
)

// MsgId 生成msg_id
// - msg_id要求全局唯一
// - msg_id跟largerId的后4位是相同的（slotId其实就是largerId）
func MsgId(ctx context.Context, mem *redis.Client, slotId uint64, currTimestamp int64) (uint64, error) {
	// 每秒一个Key，在Key上面进行+1操作
	key := keyMsgId(currTimestamp)
	expireSec := 2
	incr, err := incNum(ctx, mem, key, expireSec)
	if err != nil {
		return 0, err
	}

	// msg_id的组成部分：[ 10位：相对时间戳 | 6位：自增id | 4位：槽id ]
	// 槽id的作用：使用msg_id也能定位到对应的数据库和数据表
	timeOffset := currTimestamp - baseTimeStampOffset
	idStr := fmt.Sprintf("%d%06d%04d", timeOffset, incr%1000000, slotId%10000)

	return cast.ToUint64(idStr), nil
}

// incNum 每秒一个Key，进行累加
func incNum(ctx context.Context, mem *redis.Client, key string, expireSec int) (value int64, err error) {
	value, err = mem.IncrBy(ctx, key, 1).Result()
	if err != nil {
		return
	}
	// 这里的命令可能会失败
	// 解决办法：
	// - 使用lua脚本：https://gitee.com/jasonzxj/LearnGo/blob/master/use/pkg/redis/goredis/lua/atomic/incry_expire.go
	if value == 1 {
		_, err = mem.Do(ctx, "EXPIRE", key, expireSec).Result()
		if err != nil {
			return
		}
	}
	return
}

func keyMsgId(timestamp int64) string {
	return fmt.Sprintf(RedisPrefix+"mid_%v", timestamp)
}
