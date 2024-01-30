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
	ownerId = gen_id.NewComponentId(1001, uint32(gen_id.ContactIdTypeUser))
	peerId = gen_id.NewComponentId(10001, uint32(gen_id.ContactIdTypeUser))

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
	ownerId = gen_id.NewComponentId(1003, uint32(gen_id.ContactIdTypeUser))
	peerId = gen_id.NewComponentId(10003, uint32(gen_id.ContactIdTypeRobot))

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
	ownerId = gen_id.NewComponentId(1001, uint32(gen_id.ContactIdTypeUser))
	peerId = gen_id.NewComponentId(100000001, uint32(gen_id.ContactIdTypeGroup))

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
