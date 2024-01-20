package gen_id

import (
	"fmt"
)

// UserSessionId 获取会话id，小的id在前，大的id在后
// 用户与用户之间，使用的timeline标识
func UserSessionId(id1, id2 *ComponentId) string {
	smallerId, largerId := Sort(id1, id2)

	// session_id的组成部分：[ smallerId ":" largerId]
	return fmt.Sprintf("%s:%s", smallerId.ToString(), largerId.ToString())
}

// GroupSessionId 只有一个id
// 用户与群组之间，群组使用的timeline标识
func GroupSessionId(group *ComponentId) string {

	// session_id的组成部分：[ groupId ]
	return fmt.Sprintf("%s", group.ToString())
}
