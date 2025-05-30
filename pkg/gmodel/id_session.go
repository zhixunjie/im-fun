package gmodel

import (
	"fmt"
	"github.com/spf13/cast"
	"strings"
)

const (
	PrefixPair  = "pair"
	PrefixGroup = "group"
)

type SessionId string

// NewSessionId 根据id的类型，生成sessionId
func NewSessionId(id1, id2 *ComponentId) (sessionId SessionId) {
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
func userSessionId(id1, id2 *ComponentId) SessionId {
	smallerId, largerId := id1.Sort(id2)

	// session_id的组成部分：[ smallerId ":" largerId]
	return SessionId(fmt.Sprintf(PrefixPair+":%s:%s", smallerId.ToString(), largerId.ToString()))
}

// groupSessionId 标识群聊timeline（使用群组id）
func groupSessionId(group *ComponentId) SessionId {

	// session_id的组成部分：[ groupId ]
	return SessionId(fmt.Sprintf(PrefixGroup+":%s", group.ToString()))
}

type ParseResult struct {
	Prefix string
	Ids    []*ComponentId
}

// Parse 解析SessionId
func (s SessionId) Parse() (result *ParseResult) {
	slice := strings.Split(string(s), ":")
	result = new(ParseResult)

	if len(slice) > 0 {
		result.Prefix = slice[0]
		switch slice[0] {
		case PrefixPair: // 单聊
			v := strings.Split(slice[1], "_")
			result.Ids = append(result.Ids, NewComponentId(cast.ToUint64(v[1]), ContactIdType(cast.ToUint32(v[0]))))

			v = strings.Split(slice[2], "_")
			result.Ids = append(result.Ids, NewComponentId(cast.ToUint64(v[1]), ContactIdType(cast.ToUint32(v[0]))))
		case PrefixGroup: // 群聊
			val := strings.Split(slice[1], "_")
			result.Ids = append(result.Ids, NewComponentId(cast.ToUint64(val[1]), ContactIdType(cast.ToUint32(val[0]))))
		}
	}

	return
}
