package job

import (
	"context"
	"fmt"
	"github.com/Shopify/sarama"
	pb "github.com/zhixunjie/im-fun/api/pb"
	"github.com/zhixunjie/im-fun/internal/job/conf"
	"github.com/zhixunjie/im-fun/internal/job/invoker"
	"github.com/zhixunjie/im-fun/pkg/kafka"
	"github.com/zhixunjie/im-fun/pkg/logging"
	"github.com/zhixunjie/im-fun/pkg/registry"
	"github.com/zhixunjie/im-fun/pkg/registry/endpoint"
	"google.golang.org/protobuf/proto"
	"sync"
	"time"
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
	b.Watch(conf)

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

func (b *Job) Watch(conf *conf.Config) {
	name := "comet"

	// get watcher
	watcher, err := registry.KratosEtcdRegistry.Watch(context.Background(), name)
	if err != nil {
		panic(err)
	}

	// start watcher
	go func() {
		for {
			var tErr error
			// 任意一个instance变化，都会触发Next（并返回当前有效的instance）
			instances, tErr := watcher.Next()
			if tErr != nil {
				logging.Errorf("watch Next err=%v", err)
				time.Sleep(time.Second)
				continue
			}

			// create invoker
			var cmt *invoker.CometInvoker
			var ept string
			oldComet := b.cometInvokers
			newComet := make(map[string]*invoker.CometInvoker, len(oldComet))
			for _, ins := range instances {
				serverId := ins.ID
				endpoints := ins.Endpoints
				logging.Infof("deal with instance,serverId=%v,endpoints=%v", serverId, endpoints)

				// parse endpoint
				ept, err = endpoint.ParseEndpoint(ins.Endpoints, endpoint.Scheme("grpc", false))
				if err != nil {
					logging.Errorf("parse endpoint err=%v,serverId=%v,endpoints=%v", err, serverId, endpoints)
					continue
				}
				if ept == "" {
					logging.Errorf("ept is empty,serverId=%v,endpoints=%v", serverId, endpoints)
					continue
				}

				// check if exists serverId, if not create new
				if _, ok := oldComet[serverId]; ok {
					newComet[serverId] = oldComet[serverId]
					logging.Infof("old comet,nothing change,serverId=%v", serverId)
					continue
				} else {
					cmt, tErr = invoker.NewCometInvoker(serverId, ept, conf.CometInvoker)
					if tErr != nil {
						logging.Errorf("NewCometInvoker err=%v", tErr)
						continue
					}
					logging.Infof("new commet success,serverId=%v", serverId)
					newComet[serverId] = cmt
				}
			}

			// 关闭无效的实例
			for serverId, old := range oldComet {
				if _, ok := newComet[serverId]; !ok {
					logging.Infof("remove comet,serverId=%v", serverId)
					old.Cancel()
				}
			}

			// 完整地覆盖当前的MAP
			b.cometInvokers = newComet
		}
	}()
}
