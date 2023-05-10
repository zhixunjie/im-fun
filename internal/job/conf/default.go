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
		Comet: &Comet{ChanNum: 1024, RoutineNum: 32},
		Room: &Room{
			Batch:         20,
			TimerDuration: newtime.Duration(time.Second),
			Idle:          newtime.Duration(time.Minute * 15),
		},
	}
}
