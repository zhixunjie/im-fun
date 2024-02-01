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

	var ownerId, peerId *gen_id.ComponentId
	ownerId = gen_id.NewUserComponentId(1001)
	peerId = gen_id.NewUserComponentId(10001)

	rsp, err := messageUseCase.Fetch(ctx, &request.MessageFetchReq{
		FetchType: model.FetchTypeForward,
		VersionId: 0,
		OwnerId:   ownerId.Id(),
		OwnerType: gen_id.ContactIdType(ownerId.Type()),
		PeerId:    peerId.Id(),
		PeerType:  gen_id.ContactIdType(peerId.Type()),
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

	var ownerId, peerId *gen_id.ComponentId
	ownerId = gen_id.NewUserComponentId(1003)
	peerId = gen_id.NewRobotComponentId(10003)

	rsp, err := messageUseCase.Fetch(ctx, &request.MessageFetchReq{
		//FetchType: model.FetchTypeBackward,
		//VersionId: 1705766012000002,
		//FetchType: model.FetchTypeForward,
		//VersionId: 1705766012000002,
		FetchType: model.FetchTypeForward,
		VersionId: 0,
		OwnerId:   ownerId.Id(),
		OwnerType: gen_id.ContactIdType(ownerId.Type()),
		PeerId:    peerId.Id(),
		PeerType:  gen_id.ContactIdType(peerId.Type()),
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

	var ownerId, peerId *gen_id.ComponentId
	ownerId = gen_id.NewUserComponentId(1001)
	peerId = gen_id.NewGroupComponentId(100000001)

	rsp, err := messageUseCase.Fetch(ctx, &request.MessageFetchReq{
		FetchType: model.FetchTypeForward,
		VersionId: 0,
		OwnerId:   ownerId.Id(),
		OwnerType: gen_id.ContactIdType(ownerId.Type()),
		PeerId:    peerId.Id(),
		PeerType:  gen_id.ContactIdType(peerId.Type()),
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
