package tests

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/request"
	"github.com/zhixunjie/im-fun/pkg/gmodel"
	"github.com/zhixunjie/im-fun/pkg/utils"
	"log"
	"testing"
)

func TestUserLogin(t *testing.T) {
	ctx := context.Background()
	rsp, err := httpSrv.BzUser.Login(ctx, &request.LoginReq{
		AccountType: gmodel.AccountTypeDID,
		AccountID:   uuid.New().String(),
	})
	if err != nil {
		log.Fatalln(err)
	}
	buf, err := json.Marshal(&rsp)
	if err != nil {
		return
	}
	utils.PrettyJson(buf)
}
