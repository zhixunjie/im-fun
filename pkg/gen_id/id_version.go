package gen_id

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cast"
	"github.com/zhixunjie/im-fun/pkg/gmodel"
	"time"
)

// NewMsgVersionId 专门为message表生成version_id
func NewMsgVersionId(ctx context.Context, params *MsgVerParams) (versionId uint64, err error) {
	// 每隔32秒，verIdTimeKey的值增加1（随着时间过去，KEY会不断增大）
	ts := time.Now().Unix()
	verIdTimeKey := ts >> shiftVersionKey

	var key string
	switch {
	case params.Id1.IsGroup(): // 群聊｜群组维度的递增
		key = keyMsgGroupVersion(params.Id1.ToString(), verIdTimeKey)
	case params.Id2.IsGroup(): // 群聊｜群组维度的递增
		key = keyMsgGroupVersion(params.Id2.ToString(), verIdTimeKey)
	default: // 单聊｜会话维度的递增
		smallerId, largerId := params.Id1.Sort(params.Id2)
		key = keyMsgVersion(smallerId.ToString(), largerId.ToString(), verIdTimeKey)
	}

	// incr
	afterIncr, err := incNumExpire(ctx, params.Mem, key, 1, expireVersionKey)
	if err != nil {
		return
	}

	// version_id的组成部分：[ 10位：当前时间戳 | 6位：自增id ]
	versionId = cast.ToUint64(fmt.Sprintf("%d%06d", ts, afterIncr%1000000))

	return
}

// NewContactVersionId 专门为contact表生成version_id
func NewContactVersionId(ctx context.Context, params *ContactVerParams) (versionId uint64, err error) {
	// 每隔32秒，verIdTimeKey的值增加1（随着时间过去，KEY会不断增大）
	ts := time.Now().Unix()
	verIdTimeKey := ts >> shiftVersionKey

	// UID 维度的递增
	key := keyContactVersion(params.Owner.ToString(), verIdTimeKey)

	// incr
	afterIncr, err := incNumExpire(ctx, params.Mem, key, 1, expireVersionKey)
	if err != nil {
		return
	}

	// version_id的组成部分：[ 10位：当前时间戳 | 6位：自增id ]
	versionId = cast.ToUint64(fmt.Sprintf("%d%06d", ts, afterIncr%1000000))

	return
}

type MsgVerParams struct {
	Mem      *redis.Client
	Id1, Id2 *gmodel.ComponentId
}

type ContactVerParams struct {
	Mem   *redis.Client
	Owner *gmodel.ComponentId
}
