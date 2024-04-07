package gen_id

import (
	"fmt"
	"testing"
)

func TestIdSession(t *testing.T) {
	id1 := NewUserComponentId(1001)
	id2 := NewUserComponentId(1002)
	id3 := NewGroupComponentId(10)

	fmt.Println("单聊", SessionId(id1, id2))
	fmt.Println("单聊", SessionId(id2, id1))
	fmt.Println("群聊", SessionId(id1, id3))
	fmt.Println("群聊", SessionId(id3, id1))
	fmt.Println("群聊", SessionId(id2, id3))
	fmt.Println("群聊", SessionId(id3, id2))
}

func TestSort(t *testing.T) {
	id1 := &ComponentId{
		id:     1005,
		idType: 2,
	}
	id2 := &ComponentId{
		id:     1004,
		idType: 1,
	}
	fmt.Println(Sort(id1, id2))
}

func TestParseSessionId(t *testing.T) {
	sessionId := SessionId(NewUserComponentId(1001), NewGroupComponentId(100000000001))
	result := ParseSessionId(sessionId)
	fmt.Println(result, result.IdArr[0])

	sessionId = SessionId(NewUserComponentId(1001), NewUserComponentId(1002))
	result = ParseSessionId(sessionId)
	fmt.Println(result, result.IdArr[0], result.IdArr[1])

	sessionId = SessionId(NewUserComponentId(1001), NewRobotComponentId(111111))
	result = ParseSessionId(sessionId)
	fmt.Println(result, result.IdArr[0], result.IdArr[1])
}
