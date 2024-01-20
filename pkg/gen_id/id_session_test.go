package gen_id

import (
	"fmt"
	"testing"
)

func TestIdSession(t *testing.T) {
	id1 := &ComponentId{
		id:     1001,
		idType: 2,
	}
	id2 := &ComponentId{
		id:     1002,
		idType: 1,
	}
	fmt.Println(UserSessionId(id1, id2))
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
