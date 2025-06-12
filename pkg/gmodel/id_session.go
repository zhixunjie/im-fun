package gmodel

import (
	"errors"
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
func (s SessionId) Parse() (result *ParseResult, err error) {
	result = new(ParseResult)

	// first time parse
	slice, err := s.parseOuter(string(s))
	if err != nil {
		return
	}
	result.Prefix = slice[0]

	// check case
	var v []string
	switch slice[0] {
	case PrefixPair: // 单聊
		if len(slice) != 3 {
			err = errors.New("invalid session id format(outer)")
			return
		}
		// parse inner
		v, err = s.parseInner(slice[1])
		if err != nil {
			return
		}
		result.Ids = append(result.Ids, NewComponentId(cast.ToUint64(v[1]), ContactIdType(cast.ToUint32(v[0]))))
		// parse inner
		v, err = s.parseInner(slice[2])
		if err != nil {
			return
		}
		result.Ids = append(result.Ids, NewComponentId(cast.ToUint64(v[1]), ContactIdType(cast.ToUint32(v[0]))))
		return
	case PrefixGroup: // 群聊
		if len(slice) != 2 {
			err = errors.New("invalid session id format(outer)")
			return
		}
		// parse inner
		v, err = s.parseInner(slice[1])
		if err != nil {
			return
		}
		result.Ids = append(result.Ids, NewComponentId(cast.ToUint64(v[1]), ContactIdType(cast.ToUint32(v[0]))))
		return
	}

	return
}

// 通过sender
func (s SessionId) ParseUserSessionId(sender *ComponentId) (receiver *ComponentId, err error) {

	return
}

func (s SessionId) parseOuter(src string) (dst []string, err error) {
	dst = strings.Split(src, ":")
	if len(dst) == 0 {
		err = errors.New("invalid session id format(outer)")
		return
	}
	return
}

func (s SessionId) parseInner(src string) (dst []string, err error) {
	dst = strings.Split(src, "_")
	if len(dst) != 2 {
		err = errors.New("invalid session id format(inner)")
		return
	}
	return
}
