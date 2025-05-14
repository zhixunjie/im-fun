package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/generate/model"
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/request"
	"github.com/zhixunjie/im-fun/pkg/gen_id"
	"github.com/zhixunjie/im-fun/pkg/utils"
	"log"
	"testing"
)

// 拉取消息：用户之间互相通信
func TestMessageFetchBetweenUser(t *testing.T) {
	ctx := context.Background()

	rsp, err := messageUseCase.Fetch(ctx, &request.MessageFetchReq{
		FetchType: model.FetchTypeForward,
		VersionId: 0,
		Owner:     gen_id.NewUserComponentId(1001),
		Peer:      gen_id.NewGroupComponentId(10001),
	})

	if err != nil {
		log.Fatal(err)
	}

	buf, err := json.Marshal(&rsp.Data)
	if err != nil {
		return
	}
	utils.PrettyJson(buf)
	fmt.Println(rsp.Data.MsgList)
}

// 拉取消息：用户与机器人之间互相通信
func TestFetchBetweenUserAndRobot(t *testing.T) {
	ctx := context.Background()

	rsp, err := messageUseCase.Fetch(ctx, &request.MessageFetchReq{
		//FetchType: model.FetchTypeBackward,
		//VersionID: 1705766012000002,
		//FetchType: model.FetchTypeForward,
		//VersionID: 1705766012000002,
		FetchType: model.FetchTypeForward,
		VersionId: 0,
		Owner:     gen_id.NewUserComponentId(1003),
		Peer:      gen_id.NewGroupComponentId(10003),
	})

	if err != nil {
		log.Fatal(err)
	}

	buf, err := json.Marshal(&rsp.Data)
	if err != nil {
		return
	}
	utils.PrettyJson(buf)
	fmt.Println(rsp.Data.MsgList)
}

// 发送消息：用户与群组之间的通信
func TestFetchBetweenUserAndGroup(t *testing.T) {
	ctx := context.Background()

	rsp, err := messageUseCase.Fetch(ctx, &request.MessageFetchReq{
		FetchType: model.FetchTypeForward,
		VersionId: 0,
		Owner:     gen_id.NewUserComponentId(1001),
		Peer:      gen_id.NewGroupComponentId(100000001),
	})

	if err != nil {
		log.Fatal(err)
	}

	buf, err := json.Marshal(&rsp.Data)
	if err != nil {
		return
	}
	utils.PrettyJson(buf)
	fmt.Println(len(rsp.Data.MsgList))
}
