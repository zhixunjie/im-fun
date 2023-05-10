package job

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
	pb "github.com/zhixunjie/im-fun/api/logic"
	"google.golang.org/protobuf/proto"
)

// Consume messages, watch signals
func (job *Job) Consume(msg *sarama.ConsumerMessage) (err error) {
	logHead := "Consume|"

	// Unmarshal msg
	message := new(pb.PushMsg)
	if err = proto.Unmarshal(msg.Value, message); err != nil {
		logrus.Errorf(logHead+"err=%v", err)
		return
	}

	// deal msg
	switch message.Type {
	case pb.PushMsg_UserKeys:
		err = job.PushUserKeys(message.SubId, message.ServerId, message.UserKeys, message.Msg)
	case pb.PushMsg_UserRoom:
		room := job.getRoom(message.RoomId)
		err = room.PushToChannel(message.Msg)
	case pb.PushMsg_UserAll:
		//err = job.broadcast(message.SubId, message.Msg, message.Speed)
	default:
		err = fmt.Errorf("unknown push type: %s", message.Type)
	}
	if err != nil {
		logrus.Errorf(logHead+"err=%v", err)
	}
	return
}
