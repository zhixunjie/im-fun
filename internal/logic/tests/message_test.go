package tests

import (
	"context"
	"fmt"
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/request"
	"github.com/zhixunjie/im-fun/pkg/gmodel"
	"log"
	"testing"
)

func TestClearHistory(t *testing.T) {
	ctx := context.Background()
	rsp, err := messageUseCase.ClearHistory(ctx, &request.MessageClearHistoryReq{
		MsgID: 726942620000030001,
		Owner: gmodel.NewUserComponentId(1001),
		Peer:  gmodel.NewUserComponentId(10001),
	})
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(rsp)
}
