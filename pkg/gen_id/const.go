package gen_id

const (
	RedisPrefix = "m0_"
)

// msg_id
const (
	// 计算：相对时间戳
	baseTimeStampOffset = 1677004307 // 2023-02-22 02:31:47
	msgKeyExpire        = 2          // 生成 msg_id 时，redis key的有效期（2秒）
)

// version_id
const (
	versionKeyShift  = 7                      // 生成 version_id 时，每隔2^7次方秒对应一个redis key
	versionKeyExpire = 1<<versionKeyShift + 3 // 生成 version_id 时，redis key的有效期
)
