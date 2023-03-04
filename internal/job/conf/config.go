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
	Debug     bool
	Discovery *Discovery
	Kafka     []kafka.ConsumerGroupConf
	Comet     *Comet
	Room      *Room
}

type Discovery struct {
	Addr string
}

type Room struct {
	Batch  int
	Signal newtime.Duration
	Idle   newtime.Duration
}

type Comet struct {
	ChanNum    int
	RoutineNum int
}
