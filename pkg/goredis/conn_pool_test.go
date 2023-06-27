package goredis

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestCreatePool(t *testing.T) {
	ctx := context.Background()
	client, err := CreatePool("127.0.0.1:6379", "", 0)
	if err != nil {
		panic(err)
	}
	_, _ = client.Set(ctx, "a", 100, time.Duration(100*time.Second)).Result()

	for i := 1; i <= 10; i++ {
		go func(index int) {
			val, err := client.Get(ctx, "a").Result()
			fmt.Printf("index=%d,get=%v,err=%v\n", index, val, err)
		}(i)
	}
	time.Sleep(5 * time.Second)
}
