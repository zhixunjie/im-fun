package gen_id

import (
	"fmt"
	"github.com/spf13/cast"
	"strings"
)

// SessionId 根据id的类型，生成sessionId
func SessionId(id1, id2 *ComponentId) (sessionId string) {
	switch {
	case id1.IsGroup(): // 群聊
		sessionId = groupSessionId(id1)
	case id2.IsGroup(): // 群聊
		sessionId = groupSessionId(id2)
	default: // 单聊
		sessionId = userSessionId(id1, id2)
	}

	return
}

// userSessionId 标识单聊timeline（使用双方id，小的id在前，大的id在后）
func userSessionId(id1, id2 *ComponentId) string {
	smallerId, largerId := Sort(id1, id2)

	// session_id的组成部分：[ smallerId ":" largerId]
	return fmt.Sprintf("%s:%s", smallerId.ToString(), largerId.ToString())
}

// groupSessionId 标识群聊timeline（使用群组id）
func groupSessionId(group *ComponentId) string {

	// session_id的组成部分：[ groupId ]
	return fmt.Sprintf("%s", group.ToString())
}

// ParseSessionId 解析SessionId
func ParseSessionId(sessionId string) (id1, id2 *ComponentId) {
	slice := strings.Split(sessionId, ":")
	if len(slice) == 1 { // 群聊
		val := strings.Split(slice[0], "_")
		id1 = NewComponentId(cast.ToUint64(val[1]), cast.ToUint32(val[0]))
	} else { // 单聊
		val := strings.Split(slice[0], "_")
		id1 = NewComponentId(cast.ToUint64(val[1]), cast.ToUint32(val[0]))

		val = strings.Split(slice[1], "_")
		id2 = NewComponentId(cast.ToUint64(val[1]), cast.ToUint32(val[0]))
	}

	return
}
