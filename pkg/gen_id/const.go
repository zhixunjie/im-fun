package gen_id

const (
	RedisPrefix = "m0_"
)

// msg
const (
	// 计算：相对时间戳
	baseTimeStampOffset = 1677004307 // 2023-02-22 02:31:47
)

// version
const (
	TimeStampKeyShift  = 7                        // 生成versionId时，2的7次幂秒生成一个redis key
	TimeStampKeyExpire = 1<<TimeStampKeyShift + 3 // 生成versionId时，redis key的有效期
)
