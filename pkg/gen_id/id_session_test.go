package gen_id

import (
	"fmt"
	"testing"
)

func TestIdSession(t *testing.T) {
	id1 := &ComponentId{
		id:     1001,
		idType: uint32(ContactIdTypeUser),
	}
	id2 := &ComponentId{
		id:     1002,
		idType: uint32(ContactIdTypeUser),
	}
	id3 := &ComponentId{
		id:     10,
		idType: uint32(ContactIdTypeGroup),
	}
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
