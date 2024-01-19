package conf

import "github.com/zhixunjie/im-fun/pkg/kafka"

func DefaultKafka() []kafka.ConsumerGroupConf {
	return []kafka.ConsumerGroupConf{
		{
			Topic:   "im_push",
			Brokers: []string{"127.0.0.1:9092"},
			GroupId: "im_push_group",
		},
	}
}
