package conf

import (
	"github.com/zhixunjie/im-fun/pkg/env"
	newtime "github.com/zhixunjie/im-fun/pkg/time"
	"time"
)

func defaultConfig() *Config {
	return &Config{
		Env:       env.EnvTypeLocal,
		Name:      "job",
		Debug:     false,
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
