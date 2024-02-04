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

	id1 := NewUserComponentId(1001)
	id2 := NewUserComponentId(1002)
	id3 := NewGroupComponentId(10)

	// test1
	msgId, err := MsgId(ctx, client, id1, id2)
	fmt.Printf("单聊,msgId=%v,err=%v\n", msgId, err)

	// test2
	msgId, err = MsgId(ctx, client, id2, id1)
	fmt.Printf("单聊,msgId=%v,err=%v\n", msgId, err)

	// test3
	msgId, err = MsgId(ctx, client, id1, id3)
	fmt.Printf("群聊,msgId=%v,err=%v\n", msgId, err)

	// test4
	msgId, err = MsgId(ctx, client, id3, id1)
	fmt.Printf("群聊,msgId=%v,err=%v\n", msgId, err)
}
