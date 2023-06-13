package conf

import (
	"github.com/zhixunjie/im-fun/pkg/kafka"
	"github.com/zhixunjie/im-fun/pkg/logging"
	newtime "github.com/zhixunjie/im-fun/pkg/time"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

var Conf *Config

func InitConfig() (err error) {
	Conf = defaultConfig()
	bytes, err := ioutil.ReadFile("cmd/job/job.yaml")
	if err != nil {
		logging.Errorf("err=%v", err)
		return err
	}

	// begin to unmarshal
	err = yaml.Unmarshal(bytes, &Conf)
	if err != nil {
		logging.Errorf("err=%v", err)
		return err
	}
	logging.Infof("config=%+v", Conf)

	return nil
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
