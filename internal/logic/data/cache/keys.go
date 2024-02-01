package cache

import k "github.com/zhixunjie/im-fun/pkg/goredis/key"

// 待划分：划分出cache操作层（需要wire上Redis对象）

const Prefix = "im:logic"

const (
	TimelineMessageLock k.Key = Prefix + "timeline:message:lock{session_id}"
	TimelineContactLock k.Key = Prefix + "timeline:contact:lock{uid}"
)
