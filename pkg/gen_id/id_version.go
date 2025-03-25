package gen_id

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/cast"
	"time"
)

type MsgVerParams struct {
	Mem      *redis.Client
	Id1, Id2 *ComponentId
}

type ContactVerParams struct {
	Mem   *redis.Client
	Owner *ComponentId
}

// MsgVersionId 专门为message表生成version_id
func MsgVersionId(ctx context.Context, params *MsgVerParams) (versionId uint64, err error) {
	// 每隔128秒，verIdTimeKey的值增加1（随着时间过去，KEY会不断增大）
	ts := time.Now().Unix()
	verIdTimeKey := ts >> shiftVersionKey

	var key string
	switch {
	case params.Id1.IsGroup(): // 群聊（群组id固定放在后面）
		key = keyMsgGroupVersion(params.Id1.ToString(), verIdTimeKey)
	case params.Id2.IsGroup(): // 群聊（群组id固定放在后面）
		key = keyMsgGroupVersion(params.Id2.ToString(), verIdTimeKey)
	default: // 单聊
		smallerId, largerId := Sort(params.Id1, params.Id2)
		key = keyMsgVersion(smallerId.ToString(), largerId.ToString(), verIdTimeKey)
	}

	// incr
	afterIncr, err := incNum(ctx, params.Mem, key, expireVersionKey)
	if err != nil {
		return
	}

	// version_id的组成部分：[ 10位：当前时间戳 | 6位：自增id ]
	versionId = cast.ToUint64(fmt.Sprintf("%d%06d", ts, afterIncr%1000000))

	return
}

// ContactVersionId 专门为contact表生成version_id
func ContactVersionId(ctx context.Context, params *ContactVerParams) (versionId uint64, err error) {
	// 每隔128秒，verIdTimeKey的值增加1（随着时间过去，KEY会不断增大）
	ts := time.Now().Unix()
	verIdTimeKey := ts >> shiftVersionKey

	key := keyContactVersion(params.Owner.ToString(), verIdTimeKey)

	// incr
	afterIncr, err := incNum(ctx, params.Mem, key, expireVersionKey)
	if err != nil {
		return
	}

	// version_id的组成部分：[ 10位：当前时间戳 | 6位：自增id ]
	versionId = cast.ToUint64(fmt.Sprintf("%d%06d", ts, afterIncr%1000000))

	return
}
