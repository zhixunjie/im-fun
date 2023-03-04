package job

import (
	"github.com/Shopify/sarama"
	"github.com/zhixunjie/im-fun/internal/job/conf"
	"github.com/zhixunjie/im-fun/pkg/kafka"
	"sync"
)

type Job struct {
	c            *conf.Config
	consumer     *kafka.ConsumerGroup
	cometServers map[string]*Comet

	//rooms      map[string]*Room
	roomsMutex sync.RWMutex
}

func New(c *conf.Config) *Job {
	job := &Job{
		c: c,
		//rooms:    make(map[string]*Room),
	}
	// make consumer
	fn := func(msg *sarama.ConsumerMessage) {
		job.Consume(msg)
	}
	tmp, err := kafka.NewConsumerGroup(&c.Kafka[0], fn)
	if err != nil {
		panic(err)
	}
	job.consumer = tmp
	return job
}
