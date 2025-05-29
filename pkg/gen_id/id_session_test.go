package gen_id

import (
	"fmt"
	"github.com/zhixunjie/im-fun/pkg/gmodel"
	"testing"
)

func TestIdSession(t *testing.T) {
	id1 := gmodel.NewUserComponentId(1001)
	id2 := gmodel.NewUserComponentId(1002)
	id3 := gmodel.NewGroupComponentId(10)

	fmt.Println("单聊", SessionId(id1, id2))
	fmt.Println("单聊", SessionId(id2, id1))
	fmt.Println("群聊", SessionId(id1, id3))
	fmt.Println("群聊", SessionId(id3, id1))
	fmt.Println("群聊", SessionId(id2, id3))
	fmt.Println("群聊", SessionId(id3, id2))
}

func TestSort(t *testing.T) {
	id1 := &gmodel.ComponentId{
		id:     1005,
		idType: 2,
	}
	id2 := &gmodel.ComponentId{
		id:     1004,
		idType: 1,
	}
	fmt.Println(id1.Sort(id2))
}

func TestParseSessionId(t *testing.T) {
	sessionId := SessionId(gmodel.NewUserComponentId(1001), gmodel.NewGroupComponentId(100000000001))
	result := ParseSessionId(sessionId)
	fmt.Println(result, result.IdArr[0])

	sessionId = SessionId(gmodel.NewUserComponentId(1001), gmodel.NewUserComponentId(1002))
	result = ParseSessionId(sessionId)
	fmt.Println(result, result.IdArr[0], result.IdArr[1])

	sessionId = SessionId(gmodel.NewUserComponentId(1001), gmodel.NewRobotComponentId(111111))
	result = ParseSessionId(sessionId)
	fmt.Println(result, result.IdArr[0], result.IdArr[1])
}
