package gen_id

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cast"
	"github.com/zhixunjie/im-fun/pkg/gmodel"
	"strconv"
	"time"
)

// NewMsgId 根据id的类型，生成msgId
func NewMsgId(ctx context.Context, params *MsgIdParams) (msgId uint64, err error) {
	id1, id2 := params.Id1, params.Id2
	mem := params.Mem

	switch {
	case id1.IsGroup(): // 群聊
		msgId, err = genMsgId(ctx, mem, id1.GetId())
	case id2.IsGroup(): // 群聊
		msgId, err = genMsgId(ctx, mem, id2.GetId())
	default: // 单聊
		_, largerId := id1.Sort(id2)
		msgId, err = genMsgId(ctx, mem, largerId.GetId())
	}

	return
}

// genMsgId 生成msg_id
// 如果 msgId 使用int64，可以支持偏移28年。
// 如果 msgId 使用uint64，可以支持偏移58年。
func genMsgId(ctx context.Context, mem *redis.Client, slotId uint64) (msgId uint64, err error) {
	// redis：每秒一个key，在key上执行原子操作+1
	ts := time.Now().Unix()
	key := keyMsgId(ts)

	// incr
	afterIncr, err := incNumExpire(ctx, mem, key, 1, expireMsgKey)
	if err != nil {
		return
	}

	// msg_id的组成部分：[ 10位：相对时间戳 | 6位：自增id | 4位：槽id ]
	timeOffset := ts - baseTimeStampOffset
	idStr := fmt.Sprintf("%d%06d%04d", timeOffset, afterIncr%1000000, slotId%SlotBit)

	return cast.ToUint64E(idStr)
}

// genMsgId1 生成msg_id
// 如果 msgId 使用int64，可以支持偏移28年。
// 如果 msgId 使用uint64，可以支持偏移58年。
func genMsgId1(ctx context.Context, mem *redis.Client, slotId uint64, t time.Time, needType string) (msgId any, err error) {
	// redis：每秒一个key，在key上执行原子操作+1
	ts := t.Unix()
	key := keyMsgId(ts)

	// incr
	afterIncr, err := incNumExpire(ctx, mem, key, 1, expireMsgKey)
	if err != nil {
		return
	}

	// msg_id的组成部分：[ 10位：相对时间戳 | 6位：自增id | 4位：槽id ]
	timeOffset := ts - baseTimeStampOffset
	idStr := fmt.Sprintf("%d%06d%04d", timeOffset, afterIncr%1000000, slotId%SlotBit)

	switch needType {
	case "uint64":
		//msgId = cast.ToUint64(idStr) 有bug，相当于cast.ToInt64的效果
		msgId, _ = strconv.ParseUint(idStr, 10, 64)
	case "int64":
		msgId = cast.ToInt64(idStr)
		msgId, _ = strconv.ParseInt(idStr, 10, 64)
	case "uint32":
		msgId = cast.ToUint32(idStr)
	case "int32":
		msgId = cast.ToInt32(idStr)
	}
	return
}

// 解决：incr 和 expire 的原子性问题
// lua脚本: https://gitee.com/jasonzxj/LearnGo/blob/master/use/pkg/redis/goredis/lua/atomic/incry_expire.go
// pexpire: https://redis.io/docs/latest/commands/pexpire/
func incNumExpire(ctx context.Context, mem *redis.Client, key string, incr int64, expire time.Duration) (afterIncr int64, err error) {
	script := redis.NewScript(`
local incr = tonumber(ARGV[1])
local current = redis.call("INCRBY", KEYS[1], ARGV[1])
-- 只有第一次的incry操作才会设置过期时间
if current == incr then 
	redis.call("PEXPIRE", KEYS[1], ARGV[2])
end

return current
`)

	afterIncr, err = script.Run(ctx, mem, []string{key}, []any{incr, expire.Milliseconds()}).Int64()
	if err != nil {
		return
	}

	return
}

type MsgIdParams struct {
	Mem      *redis.Client
	Id1, Id2 *gmodel.ComponentId
}
