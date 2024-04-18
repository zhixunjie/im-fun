package conf

import (
	newtime "github.com/zhixunjie/im-fun/pkg/time"
	"time"
)

func defaultConfig() *Config {
	return &Config{
		Debug:     false,
		Name:      "job",
		Discovery: &Discovery{},
		Kafka:     DefaultKafka(),
		CometInvoker: &CometInvoker{
			RoutineNum:     32,
			ChanBufferSize: 1024,
		},
		Room: &Room{
			Batch:    20,
			Interval: newtime.Duration(time.Millisecond * 500),
		},
	}
}
