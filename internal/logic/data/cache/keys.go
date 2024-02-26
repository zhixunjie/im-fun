package cache

import k "github.com/zhixunjie/im-fun/pkg/goredis/key"

// 待划分：划分出cache操作层（需要wire上Redis对象）

const Prefix = "im:logic"

// 分布式锁：保证version_id和数据库写入的时序一致性
const (
	TimelineMessageLock k.Key = Prefix + "timeline:message:lock{session_id}"
	TimelineContactLock k.Key = Prefix + "timeline:contact:lock{contact_id}"
)
