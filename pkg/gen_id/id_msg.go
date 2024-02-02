package gen_id

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/cast"
	"time"
)

// MsgId 根据id的类型，生成msgId
func MsgId(ctx context.Context, mem *redis.Client, id1, id2 *ComponentId) (msgId uint64, err error) {
	switch {
	case id1.IsGroup(): // 群聊
		msgId, err = genMsgId(ctx, mem, id1.Id())
	case id2.IsGroup(): // 群聊
		msgId, err = genMsgId(ctx, mem, id2.Id())
	default: // 单聊
		_, largerId := Sort(id1, id2)
		msgId, err = genMsgId(ctx, mem, largerId.Id())
	}

	return
}

// genMsgId 生成msg_id
func genMsgId(ctx context.Context, mem *redis.Client, slotId uint64) (msgId uint64, err error) {
	// redis：每秒一个key，在key上执行原子操作+1
	ts := time.Now().Unix()
	key := keyMsgId(ts)

	// incr
	afterIncr, err := incNum(ctx, mem, key, expireMsgKey)
	if err != nil {
		return
	}

	// msg_id的组成部分：[ 10位：相对时间戳 | 6位：自增id | 4位：槽id ]
	timeOffset := ts - baseTimeStampOffset
	idStr := fmt.Sprintf("%d%06d%04d", timeOffset, afterIncr%1000000, slotId%10000)
	msgId = cast.ToUint64(idStr)

	return
}

// incNum 每秒一个Key，进行累加
func incNum(ctx context.Context, mem *redis.Client, key string, expire time.Duration) (value int64, err error) {
	value, err = mem.IncrBy(ctx, key, 1).Result()
	if err != nil {
		return
	}
	// 这里的命令可能会失败
	// 解决办法：lua脚本：https://gitee.com/jasonzxj/LearnGo/blob/master/use/pkg/redis/goredis/lua/atomic/incry_expire.go
	if value == 1 {
		_, err = mem.Expire(ctx, key, expire).Result()
		if err != nil {
			return
		}
	}
	return
}
