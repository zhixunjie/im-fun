package gen_id

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/cast"
	"time"
)

type GenVersionType int

const (
	GenVersionTypeMsg GenVersionType = iota
	GenVersionTypeContact
)

type GenVersionParams struct {
	Mem                 *redis.Client
	GenVersionType      GenVersionType
	OwnerId             *ComponentId
	SmallerId, LargerId *ComponentId
}

// VersionId 获取"消息表/会话表"的version_id
func VersionId(ctx context.Context, params *GenVersionParams) (versionId uint64, err error) {
	// 每隔128秒，verIdTimeKey的值增加1（随着时间过去，KEY会不断增大）
	ts := time.Now().Unix()
	verIdTimeKey := ts >> shiftVersionKey

	var key string
	switch params.GenVersionType {
	case GenVersionTypeMsg:
		// smallerId、largerId：people that in chatting
		// 不需要全局唯一，只要在「同一个会话」中唯一即可
		key = keyMsgVersion(params.SmallerId.ToString(), params.LargerId.ToString(), verIdTimeKey)
	case GenVersionTypeContact:
		// ownerId：contact's owner
		// 不需要全局唯一，只要在「同一个用户」中唯一即可
		key = keyContactVersion(params.OwnerId.ToString(), verIdTimeKey)
	}

	// incr
	afterIncr, err := incNum(ctx, params.Mem, key, expireVersionKey)
	if err != nil {
		return
	}

	// version_id的组成部分：[ 10位：当前时间戳 | 6位：自增id ]
	// 这里的总位数只有16位，不像 msg_id 那样需要20位那么紧张，所以不需要用到相对时间戳
	versionId = cast.ToUint64(fmt.Sprintf("%d%06d", ts, afterIncr%1000000))

	return
}
