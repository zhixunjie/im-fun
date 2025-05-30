package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/request"
	"github.com/zhixunjie/im-fun/pkg/gmodel"
	"github.com/zhixunjie/im-fun/pkg/utils"
	"log"
	"testing"
)

func TestQueryContactLogic(t *testing.T) {
	//res, err := wire.GetContactRepo(conf.Conf).InfoWithCache(1001, 1002)
	//if err != nil {
	//	log.Fatal()
	//}
	//fmt.Printf("%+v\n", res)
}

func TestContactFetch(t *testing.T) {
	ctx := context.Background()

	rsp, err := contactUseCase.Fetch(ctx, &request.ContactFetchReq{
		VersionId: 0,
		Owner:     gmodel.NewUserComponentId(1005),
	})

	if err != nil {
		log.Fatal(err)
	}

	buf, err := json.Marshal(&rsp.Data)
	if err != nil {
		return
	}
	utils.PrettyJson(buf)
	fmt.Println(rsp.Data.ContactList, rsp.Data.NextVersionId)
}
