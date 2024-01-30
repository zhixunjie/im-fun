package gen_id

import (
	"fmt"
	"github.com/spf13/cast"
	"strings"
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

// GenSessionId 根据id的类型，生成sessionId
func GenSessionId(id1, id2 *ComponentId) (sessionId string) {
	switch {
	case id1.Type() == uint32(ContactIdTypeGroup):
		sessionId = GroupSessionId(id1)
	case id2.Type() == uint32(ContactIdTypeGroup):
		sessionId = GroupSessionId(id2)
	default:
		sessionId = UserSessionId(id1, id2)
	}

	return
}

// ParseSessionId 解析SessionId
func ParseSessionId(sessionId string) (id1, id2 *ComponentId) {
	slice := strings.Split(sessionId, ":")
	if len(slice) == 1 {
		val := strings.Split(slice[0], "_")
		id1 = NewComponentId(cast.ToUint64(val[1]), cast.ToUint32(val[0]))
	} else {
		val := strings.Split(slice[0], "_")
		id1 = NewComponentId(cast.ToUint64(val[1]), cast.ToUint32(val[0]))

		val = strings.Split(slice[1], "_")
		id2 = NewComponentId(cast.ToUint64(val[1]), cast.ToUint32(val[0]))
	}

	return
}
