package job

import (
	"fmt"
	"github.com/Shopify/sarama"
	pb "github.com/zhixunjie/im-fun/api/pb"
	"github.com/zhixunjie/im-fun/internal/job/conf"
	"github.com/zhixunjie/im-fun/pkg/kafka"
	"github.com/zhixunjie/im-fun/pkg/logging"
	"google.golang.org/protobuf/proto"
	"os"
	"sync"
)

type Job struct {
	conf            *conf.Config
	consumer        *kafka.ConsumerGroup
	allCometInvoker map[string]*CometInvoker

	rooms   map[string]*Room
	rwMutex sync.RWMutex
}

func NewJob(conf *conf.Config) *Job {
	job := &Job{
		conf:            conf,
		allCometInvoker: map[string]*CometInvoker{},
		rooms:           make(map[string]*Room),
	}

	// make consumer
	tmp, err := kafka.NewConsumerGroup(&conf.Kafka[0], func(msg *sarama.ConsumerMessage) {
		job.Consume(msg)
	})
	if err != nil {
		panic(err)
	}
	job.consumer = tmp

	// make comet
	defHost, _ := os.Hostname()
	cm, err := NewCometInvoker(defHost, conf.CometInvoker)
	if err != nil {
		panic(err)
	}
	job.allCometInvoker[defHost] = cm

	return job
}

// Consume messages, watch signals
func (job *Job) Consume(msg *sarama.ConsumerMessage) {
	logHead := "Consume|"
	var err error

	// Unmarshal msg
	message := new(pb.KafkaSendMsg)
	if err = proto.Unmarshal(msg.Value, message); err != nil {
		logging.Errorf(logHead+"err=%v", err)
		return
	}

	// deal msg
	switch message.Type {
	case pb.KafkaSendMsg_UserKeys:
		err = job.SendToUserKeys(message.SubId, message.ServerId, message.UserKeys, message.Msg)
	case pb.KafkaSendMsg_UserRoom:
		err = job.CreateOrGetRoom(message.RoomId).SendToCh(message.Msg)
	case pb.KafkaSendMsg_UserAll:
		err = job.SendToAll(message.SubId, message.Speed, message.Msg)
	default:
		err = fmt.Errorf("unknown send type: %s", message.Type)
	}
	if err != nil {
		logging.Errorf(logHead+"err=%v", err)
	}
}
