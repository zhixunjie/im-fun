package conf

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/zhixunjie/im-fun/pkg/kafka"
	newtime "github.com/zhixunjie/im-fun/pkg/time"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

var Conf *Config

func InitConfig() (err error) {
	Conf = defaultConfig()
	bytes, err := ioutil.ReadFile("cmd/job/job.yaml")
	if err != nil {
		logrus.Errorf("err=%v", err)
		return err
	}

	// begin to unmarshal
	err = yaml.Unmarshal(bytes, &Conf)
	if err != nil {
		logrus.Errorf("err=%v", err)
		return err
	}
	fmt.Printf("config=%+v\n", Conf)
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
	Batch  int              `yaml:"batch"`
	Signal newtime.Duration `yaml:"signal"`
	Idle   newtime.Duration `yaml:"idle"`
}

type Comet struct {
	ChanNum    int `yaml:"chanNum"`
	RoutineNum int `yaml:"routineNum"`
}
