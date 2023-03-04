package conf

import "github.com/zhixunjie/im-fun/pkg/kafka"

func DefaultKafka() []kafka.ProducerConf {
	return []kafka.ProducerConf{
		{
			Topic:   "im_push",
			Brokers: []string{"127.0.0.1:9092"},
		},
	}
}

// Kafka .
type Kafka struct {
	Topic   string   `yaml:"topic"`
	Brokers []string `yaml:"brokers"`
}
