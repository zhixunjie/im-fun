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
	Debug     bool                      `yaml:"debug"`
	Discovery *Discovery                `yaml:"discovery"`
	Kafka     []kafka.ConsumerGroupConf `yaml:"kafka"`
	Comet     *Comet                    `yaml:"comet"`
	Room      *Room                     `yaml:"room"`
}

type Discovery struct {
	Addr string `yaml:"addr"`
}

type Room struct {
	Batch    int              `yaml:"batch"`
	Duration newtime.Duration `yaml:"duration"`
}

type Comet struct {
	ChanNum    int `yaml:"chanNum"`    // 每个协程对应多个Channel，这里设置每个Channel的缓冲区大小
	RoutineNum int `yaml:"routineNum"` // 协程数目，用于消费Channel的消息
}
