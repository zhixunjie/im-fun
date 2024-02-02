package gen_id

import (
	"context"
	"fmt"
	"testing"
)

func TestVersionIdContact(t *testing.T) {
	ctx := context.Background()

	for i := 0; i < 10; i++ {
		fmt.Println(ContactVersionId(ctx, &ContactVerParams{
			Mem:     client,
			OwnerId: NewUserComponentId(1001),
		}))
	}
}

func TestVersionIdMsg(t *testing.T) {
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
