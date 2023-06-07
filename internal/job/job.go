package job

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
	pb "github.com/zhixunjie/im-fun/api/logic"
	"github.com/zhixunjie/im-fun/internal/job/conf"
	"github.com/zhixunjie/im-fun/pkg/kafka"
	"google.golang.org/protobuf/proto"
	"os"
	"sync"
)

type Job struct {
	conf     *conf.Config
	consumer *kafka.ConsumerGroup
	allComet map[string]*Comet

	rooms   map[string]*Room
	rwMutex sync.RWMutex
}

func New(conf *conf.Config) *Job {
	job := &Job{
		conf:     conf,
		allComet: map[string]*Comet{},
		rooms:    make(map[string]*Room),
	}

	// make consumer
	fn := func(msg *sarama.ConsumerMessage) {
		job.Consume(msg)
	}
	tmp, err := kafka.NewConsumerGroup(&conf.Kafka[0], fn)
	if err != nil {
		panic(err)
	}
	job.consumer = tmp

	// make comet
	defHost, _ := os.Hostname()
	cm, err := NewComet(defHost, conf.Comet)
	if err != nil {
		panic(err)
	}
	job.allComet[defHost] = cm

	return job
}

// Consume messages, watch signals
func (job *Job) Consume(msg *sarama.ConsumerMessage) {
	logHead := "Consume|"
	var err error

	// Unmarshal msg
	message := new(pb.PushMsg)
	if err = proto.Unmarshal(msg.Value, message); err != nil {
		logrus.Errorf(logHead+"err=%v", err)
		return
	}

	// deal msg
	switch message.Type {
	case pb.PushMsg_UserKeys:
		err = job.SendToUsers(message.SubId, message.ServerId, message.UserKeys, message.Msg)
	case pb.PushMsg_UserRoom:
		err = job.CreateOrGetRoom(message.RoomId).Send(message.Msg)
	case pb.PushMsg_UserAll:
		//err = job.broadcast(message.SubId, message.Msg, message.Speed)
	default:
		err = fmt.Errorf("unknown push type: %s", message.Type)
	}
	if err != nil {
		logrus.Errorf(logHead+"err=%v", err)
	}
}
