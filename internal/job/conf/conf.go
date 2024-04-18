package conf

import (
	"github.com/zhixunjie/im-fun/pkg/encoding/yaml"
	"github.com/zhixunjie/im-fun/pkg/kafka"
	newtime "github.com/zhixunjie/im-fun/pkg/time"
)

var Conf *Config

func InitConfig(path string) (err error) {
	Conf = defaultConfig()
	return yaml.LoadConfig(path, Conf)
}

type Config struct {
	Name         string                    `yaml:"name"`         // 服务名
	Debug        bool                      `yaml:"debug"`        // 是否开启debug
	Discovery    *Discovery                `yaml:"discovery"`    // etcd的配置
	Kafka        []kafka.ConsumerGroupConf `yaml:"kafka"`        // Kafka的配置
	CometInvoker *CometInvoker             `yaml:"cometInvoker"` // Comet调用器配置
	Room         *Room                     `yaml:"room"`
}

type Discovery struct {
	Addr string `yaml:"addr"`
}

type Room struct {
	Batch    int              `yaml:"batch"`    // 每累计达到Batch数量的proto，发送一次批量消息到Room
	Interval newtime.Duration `yaml:"interval"` // 每间隔Interval的时间，发送一次批量消息到Room
}

type CometInvoker struct {
	RoutineNum     int `yaml:"routineNum"`     // 每个CometInvoker有RoutineNum个协程，用于消费Channel的消息
	ChanBufferSize int `yaml:"chanBufferSize"` // 每个协程对应1个Channel，1个Channel的缓冲区大小为ChanBufferSize
}
