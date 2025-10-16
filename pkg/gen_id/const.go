package gen_id

import (
	"fmt"
	"time"
)

const (
	RedisPrefix = "gen_id:"
)

// msg_id
const (
	// 计算：相对时间戳
	baseTimeStampOffset = 1677004307 // 2023-02-22 02:31:47
	// 生成 msg_id 时，redis key的有效期（5秒）
	// 需要考虑：时间回退的问题、过期key淘汰
	expireMsgKey = 5 * time.Second

	SlotBit = 10000
)

// version_id
const (
	shiftVersionKey  = 5                                      // 生成 version_id 时，每隔2^5(32)秒对应一个redis key
	expireVersionKey = (1<<shiftVersionKey + 3) * time.Second // 生成 version_id 时，redis key的有效期
)

func keyMsgId(timestamp int64) string {
	return fmt.Sprintf(RedisPrefix+"mid:{%v}", timestamp)
}

func keyContactVersion(ownerUniId string, verIdTimeKey int64) string {
	return fmt.Sprintf(RedisPrefix+"cvid:{%v}:%v", ownerUniId, verIdTimeKey)
}

func keyMsgVersion(smallerUniId, largerUniId string, verIdTimeKey int64) string {
	return fmt.Sprintf(RedisPrefix+"mvid:{%v:%v}:%v", smallerUniId, largerUniId, verIdTimeKey)
}

func keyMsgGroupVersion(groupUniId string, verIdTimeKey int64) string {
	return fmt.Sprintf(RedisPrefix+"mvid:g:{%v}:%v", groupUniId, verIdTimeKey)
}
