package gen_id

import (
	"context"
	"fmt"
	"github.com/spf13/cast"
)

// GetMsgVersionId 获取"消息表"的version_id
// 注意：version_id不需要全局唯一，只要在同一个会话中唯一即可
func GetMsgVersionId(ctx context.Context, currTimestamp int64, smallerId uint64, largerId uint64) (uint64, error) {
	// 每小时一个Key，在Key上面进行+1操作
	key := getMsgVersionIncrNum(currTimestamp, smallerId, largerId)
	incr, err := incNum(ctx, key, 86400+120)
	if err != nil {
		return 0, err
	}

	// version_id的组成部分：[ 10位：相对时间戳 | 6位：自增id ]
	timeOffset := currTimestamp - baseTimeStampOffset
	idStr := fmt.Sprintf("%d%06d", timeOffset, incr%1000000)
	return cast.ToUint64(idStr), nil
}

// GetContactVersionId 获取"会话表"的version_id
// 注意：version_id不需要全局唯一，只要在同一个用户中唯一即可
func GetContactVersionId(ctx context.Context, currTimestamp int64, ownerId uint64) (uint64, error) {
	// 每小时一个Key，在Key上面进行+1操作
	key := getContactVersionIncrNum(currTimestamp, ownerId)
	incr, err := incNum(ctx, key, 86400+120)
	if err != nil {
		return 0, err
	}

	// version_id的组成部分：[ 10位：相对时间戳 | 6位：自增id ]
	timeOffset := currTimestamp - baseTimeStampOffset
	idStr := fmt.Sprintf("%d%06d", timeOffset, incr%1000000)
	return cast.ToUint64(idStr), nil
}

func getMsgVersionIncrNum(currTimestamp int64, smallerId uint64, largerId uint64) string {
	offset := currTimestamp % 86400
	return fmt.Sprintf("msg_version_num_%v_%v_%v", smallerId, largerId, offset)
}

func getContactVersionIncrNum(currTimestamp int64, ownerId uint64) string {
	offset := currTimestamp % 86400
	return fmt.Sprintf("contact_version_num_%v_%v", ownerId, offset)
}
