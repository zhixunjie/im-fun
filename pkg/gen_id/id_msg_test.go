package gen_id

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"testing"
)

var client = redis.NewClient(&redis.Options{
	Addr:     "127.0.0.1:6379",
	Password: "",
	DB:       0,
})

func TestIdMsg(t *testing.T) {
	ctx := context.Background()

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

	// test1
	smallerId, largerId := Sort(id1, id2)
	msgId, err := GenMsgId(ctx, client, smallerId, largerId)
	fmt.Printf("单聊,msgId=%v,err=%v\n", msgId, err)

	// test2
	smallerId, largerId = Sort(id2, id1)
	msgId, err = GenMsgId(ctx, client, smallerId, largerId)
	fmt.Printf("单聊,msgId=%v,err=%v\n", msgId, err)

	// test3
	smallerId, largerId = Sort(id1, id3)
	msgId, err = GenMsgId(ctx, client, smallerId, largerId)
	fmt.Printf("群聊,msgId=%v,err=%v\n", msgId, err)

	// test4
	smallerId, largerId = Sort(id3, id1)
	msgId, err = GenMsgId(ctx, client, smallerId, largerId)
	fmt.Printf("群聊,msgId=%v,err=%v\n", msgId, err)
}