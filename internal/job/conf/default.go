package conf

import (
	newtime "github.com/zhixunjie/im-fun/pkg/time"
	"time"
)

func defaultConfig() *Config {
	return &Config{
		Debug: false,
		Discovery: &Discovery{
			Addr: "127.0.0.1:7171",
		},
		Kafka: DefaultKafka(),
		Comet: &CometInvoker{ChanNum: 1024, RoutineNum: 32},
		Room: &Room{
			Batch:    20,
			Duration: newtime.Duration(time.Millisecond * 500),
		},
	}
}
