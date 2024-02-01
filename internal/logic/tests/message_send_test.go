package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/zhixunjie/im-fun/pkg/gen_id"
	"testing"
)

// 发送消息：用户之间互相通信
func TestSendBetweenUser(t *testing.T) {
	ctx := context.Background()

	users1 := []uint64{1001, 1002}   // user
	users2 := []uint64{10001, 10002} // user

	for _, user1 := range users1 {
		for _, user2 := range users2 {

			for i := 1; i <= 5; i++ {
				var senderId, receiverId *gen_id.ComponentId
				if i%2 == 1 {
					senderId = gen_id.NewUserComponentId(user1)
					receiverId = gen_id.NewUserComponentId(user2)
				} else {
					senderId = gen_id.NewUserComponentId(user2)
					receiverId = gen_id.NewUserComponentId(user1)
				}

				// build data
				d := map[string]interface{}{
					"type":    111111,
					"content": i,
				}
				JsonStr, _ := json.Marshal(d)

				rsp, err := messageUseCase.SendSimpleCustomMessage(ctx, senderId, receiverId, string(JsonStr))
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
				var senderId, receiverId *gen_id.ComponentId
				if i%2 == 1 {
					senderId = gen_id.NewUserComponentId(user1)
					receiverId = gen_id.NewRobotComponentId(user2)
				} else {
					senderId = gen_id.NewRobotComponentId(user2)
					receiverId = gen_id.NewUserComponentId(user1)
				}

				// build data
				d := map[string]interface{}{
					"type":    111111,
					"content": i,
				}
				JsonStr, _ := json.Marshal(d)

				rsp, err := messageUseCase.SendSimpleCustomMessage(ctx, senderId, receiverId, string(JsonStr))
				fmt.Printf("rsp=%+v,err=%v\n", rsp, err)
			}
		}
	}
}

// 发送消息：用户与群组之间的通信
func TestSendBetweenUserAndGroup(t *testing.T) {
	ctx := context.Background()

	senderIds := []uint64{1001, 1002, 1003, 1004, 1005, 1006, 1007, 1008, 1009}
	groupIds := []uint64{100000001, 100000002, 100000003}

	for _, groupId := range groupIds {
		for _, senderId := range senderIds {
			sender := gen_id.NewUserComponentId(senderId)
			receiver := gen_id.NewGroupComponentId(groupId)

			// build data
			d := map[string]interface{}{
				"type":    111111,
				"content": 123456,
			}
			JsonStr, _ := json.Marshal(d)

			rsp, err := messageUseCase.SendSimpleCustomMessage(ctx, sender, receiver, string(JsonStr))
			fmt.Printf("rsp=%+v,err=%v\n", rsp, err)
		}
	}
}
