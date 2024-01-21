package gen_id

import (
	"fmt"
	"time"
)

const (
	RedisPrefix = "m0_"
)

// msg_id
const (
	// 计算：相对时间戳
	baseTimeStampOffset = 1677004307      // 2023-02-22 02:31:47
	expireMsgKey        = 2 * time.Second // 生成 msg_id 时，redis key的有效期（2秒）
)

// version_id
const (
	shiftVersionKey  = 7                                      // 生成 version_id 时，每隔2^7次方秒对应一个redis key
	expireVersionKey = (1<<shiftVersionKey + 3) * time.Second // 生成 version_id 时，redis key的有效期
)

func keyMsgId(timestamp int64) string {
	return fmt.Sprintf(RedisPrefix+"mid_%v", timestamp)
}

func keyContactVersion(ownerUniId string, verIdTimeKey int64) string {
	return fmt.Sprintf(RedisPrefix+"cvid_%v_%v", ownerUniId, verIdTimeKey)
}

func keyMsgVersion(smallerUniId, largerUniId string, verIdTimeKey int64) string {
	return fmt.Sprintf(RedisPrefix+"mvid_%v_%v_%v", smallerUniId, largerUniId, verIdTimeKey)
}
