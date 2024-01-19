package gen_id

import (
	"context"
	"fmt"
	"github.com/zhixunjie/im-fun/pkg/goredis"
	"github.com/zhixunjie/im-fun/pkg/logging"
	"testing"
)

func TestVersionIdContact(t *testing.T) {
	ctx := context.Background()
	mem, err := goredis.CreatePool("127.0.0.1:6379", "", 0)
	if err != nil {
		logging.Errorf("mem,err=%v", err)
		panic(err)
	}
	ownerId := uint64(1001)

	for i := 0; i < 10; i++ {
		fmt.Println(VersionId(ctx, &GenVersionParams{
			Mem:            mem,
			GenVersionType: GenVersionTypeContact,
			OwnerId:        ownerId,
		}))
	}
}

func TestVersionIdMsg(t *testing.T) {
	ctx := context.Background()
	mem, err := goredis.CreatePool("127.0.0.1:6379", "", 0)
	if err != nil {
		logging.Errorf("redisClient,err=%v", err)
		panic(err)
	}
	smallerId := uint64(1001)
	largerId := uint64(1002)

	for i := 0; i < 10; i++ {
		fmt.Println(VersionId(ctx, &GenVersionParams{
			Mem:            mem,
			GenVersionType: GenVersionTypeMsg,
			SmallerId:      smallerId,
			LargerId:       largerId,
		}))
	}
}
