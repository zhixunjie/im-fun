package gen_id

import (
	"context"
	"fmt"
	"github.com/zhixunjie/im-fun/pkg/gmodel"
	"testing"
)

func TestVersionIdContact(t *testing.T) {
	ctx := context.Background()

	for i := 0; i < 10; i++ {
		fmt.Println(ContactVersionId(ctx, &ContactVerParams{
			Mem:   client,
			Owner: gmodel.NewUserComponentId(1001),
		}))
	}
}

func TestVersionIdMsg(t *testing.T) {
	ctx := context.Background()
	id1 := gmodel.NewUserComponentId(1001)
	id2 := gmodel.NewUserComponentId(1002)
	id3 := gmodel.NewGroupComponentId(10)

	fmt.Println(MsgVersionId(ctx, &MsgVerParams{
		Mem: client,
		Id1: id1,
		Id2: id2,
	}))

	fmt.Println(MsgVersionId(ctx, &MsgVerParams{
		Mem: client,
		Id1: id2,
		Id2: id1,
	}))

	fmt.Println(MsgVersionId(ctx, &MsgVerParams{
		Mem: client,
		Id1: id1,
		Id2: id3,
	}))

	fmt.Println(MsgVersionId(ctx, &MsgVerParams{
		Mem: client,
		Id1: id3,
		Id2: id1,
	}))
	fmt.Println(MsgVersionId(ctx, &MsgVerParams{
		Mem: client,
		Id1: id2,
		Id2: id3,
	}))

	fmt.Println(MsgVersionId(ctx, &MsgVerParams{
		Mem: client,
		Id1: id3,
		Id2: id2,
	}))
}
