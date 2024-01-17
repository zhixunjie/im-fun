package gen_id

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/cast"
)

// ContactVersionId 获取"会话表"的version_id
// 注意：version_id不需要全局唯一，只要在同一个用户中唯一即可
func ContactVersionId(ctx context.Context, mem *redis.Client, currTimestamp int64, ownerId uint64) (id uint64, err error) {
	// key:
	// - ownerId：contact's owner
	// - verIdTimeKey = timeStamp / 128
	//   - 每隔128，verIdTimeKey的值增加1，所以KEY随着随着时间的过去，会不断增大（我们就是需要它能不断增大的）
	//   - 128秒内，使用同一KEY进行累加，如果128秒的请求数超出100w（同一个用户下），那么version_id的值就有问题了
	verIdTimeKey := currTimestamp >> TimeStampKeyShift
	key := keyContactVersion(ownerId, verIdTimeKey)

	// IncrBy
	value, err := mem.IncrBy(ctx, key, 1).Result()
	if err != nil {
		return
	}
	if value == 1 {
		_, err = mem.Do(ctx, "EXPIRE", key, TimeStampKeyExpire).Result()
		if err != nil {
			mem.Do(ctx, "EXPIRE", key, TimeStampKeyExpire)
			return
		}
	}
	// version_id的组成部分：[ 10位：当前时间戳 | 6位：自增id ]
	idStr := fmt.Sprintf("%d%06d", currTimestamp, value%1000000)
	id = cast.ToUint64(idStr)

	return
}

// MsgVersionId 获取"消息表"的version_id
// 注意：version_id不需要全局唯一，只要在同一个会话中唯一即可
func MsgVersionId(ctx context.Context, mem *redis.Client, currTimestamp int64, smallerId, largerId uint64) (id uint64, err error) {
	// key:
	// - smallerId、largerId：people that in chatting
	// - verIdTimeKey = timeStamp / 128
	//   - 每隔128，verIdTimeKey的值增加1，所以KEY随着随着时间的过去，会不断增大（我们就是需要它能不断增大的）
	//   - 128秒内，使用同一KEY进行累加，如果128秒的请求数超出100w（同一个会话下），那么version_id的值就有问题了
	verIdTimeKey := currTimestamp >> TimeStampKeyShift
	key := keyMsgVersion(smallerId, largerId, verIdTimeKey)

	// IncrBy
	value, err := mem.IncrBy(ctx, key, 1).Result()
	if err != nil {
		return
	}
	if value == 1 {
		_, err = mem.Do(ctx, "EXPIRE", key, TimeStampKeyExpire).Result()
		if err != nil {
			mem.Do(ctx, "EXPIRE", key, TimeStampKeyExpire)
			return
		}
	}
	// version_id的组成部分：[ 10位：当前时间戳 | 6位：自增id ]
	idStr := fmt.Sprintf("%d%06d", currTimestamp, value%1000000)
	id = cast.ToUint64(idStr)

	return
}

func keyContactVersion(ownerId uint64, verIdTimeKey int64) string {
	return fmt.Sprintf(RedisPrefix+"cvid_%v_%v", ownerId, verIdTimeKey)
}

func keyMsgVersion(smallerId, largerId uint64, verIdTimeKey int64) string {
	return fmt.Sprintf(RedisPrefix+"mvid_%v_%v_%v", smallerId, largerId, verIdTimeKey)
}
