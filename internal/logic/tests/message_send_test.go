package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/zhixunjie/im-fun/pkg/gmodel"
	"testing"
)

//func TestSendSimple(t *testing.T) {
//	ctx := context.Background()
//
//	rsp, err := messageUseCase.Send(ctx, &request.MessageSendReq{
//		SeqId:    uint64(gmodel.NewSeqId()),
//		Sender:   gmodel.NewUserComponentId(1001),
//		Receiver: gmodel.NewUserComponentId(1005),
//		MsgBody: &format.MsgBody{
//			MsgType: format.MsgTypeText,
//			MsgContent: &format.TextContent{
//				Text: "哈哈哈",
//			},
//		},
//	})
//	fmt.Printf("rsp=%+v,err=%v\n", rsp, err)
//}

// 发送消息：用户之间互相通信
func TestSendBetweenUser(t *testing.T) {
	ctx := context.Background()

	users1 := []uint64{1001}  // user
	users2 := []uint64{10001} // user

	for _, user1 := range users1 {
		for _, user2 := range users2 {
			for i := 1; i <= 5; i++ {
				var sender, receiver *gmodel.ComponentId
				if i%2 == 1 {
					sender = gmodel.NewUserComponentId(user1)
					receiver = gmodel.NewUserComponentId(user2)
				} else {
					sender = gmodel.NewUserComponentId(user2)
					receiver = gmodel.NewUserComponentId(user1)
				}

				// build data
				d := map[string]any{
					"type":    "level_up",
					"content": i,
				}
				JsonStr, _ := json.Marshal(d)

				rsp, err := messageUseCase.SendSimpleCustomMessage(ctx, sender, receiver, string(JsonStr))
				fmt.Printf("rsp=%+v,err=%v\n", rsp, err)
			}
		}
	}
}

// 发送消息：用户与机器人之间互相通信
func TestSendBetweenUserAndRobot(t *testing.T) {
	ctx := context.Background()

	users1 := []uint64{1001, 1002, 1003, 1004}     // user
	users2 := []uint64{10001, 10002, 10003, 10004} // robot

	for _, user1 := range users1 {
		for _, user2 := range users2 {

			for i := 1; i <= 5; i++ {
				var sender, receiver *gmodel.ComponentId
				if i%2 == 1 {
					sender = gmodel.NewUserComponentId(user1)
					receiver = gmodel.NewRobotComponentId(user2)
				} else {
					sender = gmodel.NewRobotComponentId(user2)
					receiver = gmodel.NewUserComponentId(user1)
				}

				// build data
				d := map[string]any{
					"type":    "level_up",
					"content": i,
				}
				JsonStr, _ := json.Marshal(d)

				rsp, err := messageUseCase.SendSimpleCustomMessage(ctx, sender, receiver, string(JsonStr))
				fmt.Printf("rsp=%+v,err=%v\n", rsp, err)
			}
		}
	}
}

// 发送消息：用户与群组之间的通信
func TestSendBetweenUserAndGroup(t *testing.T) {
	ctx := context.Background()

	senders := []uint64{1001, 1002, 1003, 1004, 1005, 1006, 1007, 1008, 1009}
	groups := []uint64{100000001, 100000002, 100000003}

	for _, groupId := range groups {
		for _, senderId := range senders {
			sender := gmodel.NewUserComponentId(senderId)
			receiver := gmodel.NewGroupComponentId(groupId)

			// build data
			d := map[string]any{
				"type":    "level_up",
				"content": 123456,
			}
			JsonStr, _ := json.Marshal(d)

			rsp, err := messageUseCase.SendSimpleCustomMessage(ctx, sender, receiver, string(JsonStr))
			fmt.Printf("rsp=%+v,err=%v\n", rsp, err)
		}
	}
}
