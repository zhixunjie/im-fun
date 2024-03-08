package job

import (
	"fmt"
	"github.com/Shopify/sarama"
	pb "github.com/zhixunjie/im-fun/api/pb"
	"github.com/zhixunjie/im-fun/internal/job/conf"
	"github.com/zhixunjie/im-fun/internal/job/invoker"
	"github.com/zhixunjie/im-fun/pkg/kafka"
	"github.com/zhixunjie/im-fun/pkg/logging"
	"google.golang.org/protobuf/proto"
	"os"
	"sync"
)

// Logic -> Kafka -> Job -> invoker.CometInvoker -> RPC To Comet -> Comet

// Job 任务（消费KAFKA，执行指定行为）
type Job struct {
	conf          *conf.Config
	consumerGroup *kafka.ConsumerGroup             // KAFKA：消息者组
	cometInvokers map[string]*invoker.CometInvoker // 记录：所有的CometInvoker

	roomJobs map[string]*RoomJob
	rwMutex  sync.RWMutex
}

func NewJob(conf *conf.Config) *Job {
	b := &Job{
		conf:          conf,
		cometInvokers: map[string]*invoker.CometInvoker{},
		roomJobs:      make(map[string]*RoomJob),
	}

	// 1. make consumer group
	var err error
	b.consumerGroup, err = kafka.NewConsumerGroup(&conf.Kafka[0], func(msg *sarama.ConsumerMessage) {
		b.Consume(msg)
	})
	if err != nil {
		panic(err)
	}

	// 2. make comet invoker
	defHost, _ := os.Hostname()
	cmt, err := invoker.NewCometInvoker(defHost, conf.CometInvoker)
	if err != nil {
		panic(err)
	}
	b.cometInvokers[defHost] = cmt

	return b
}

// Consume messages, watch signals
func (b *Job) Consume(msg *sarama.ConsumerMessage) {
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
	case pb.KafkaSendMsg_ToUsers:
		err = b.SendToUser(message.SubId, message.ServerId, message.TcpSessionIds, message.Msg)
	case pb.KafkaSendMsg_ToRoom:
		err = b.CreateOrGetRoom(message.RoomId).SendToCh(message.Msg)
	case pb.KafkaSendMsg_ToAll:
		err = b.SendToAll(message.SubId, message.Speed, message.Msg)
	default:
		err = fmt.Errorf("unknown send type: %s", message.Type)
	}
	if err != nil {
		logging.Errorf(logHead+"err=%v", err)
	}
}
