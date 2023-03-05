package job

import (
	"github.com/Shopify/sarama"
	"github.com/zhixunjie/im-fun/internal/job/conf"
	"github.com/zhixunjie/im-fun/pkg/kafka"
	"os"
	"sync"
)

type Job struct {
	conf     *conf.Config
	consumer *kafka.ConsumerGroup
	allComet map[string]*Comet

	//rooms      map[string]*Room
	rwMutex sync.RWMutex
}

func New(conf *conf.Config) *Job {
	job := &Job{
		conf:     conf,
		allComet: map[string]*Comet{},
		//rooms:    make(map[string]*Room),
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
