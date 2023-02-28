package conf

import (
	"fmt"
	"github.com/sirupsen/logrus"
	mytime "github.com/zhixunjie/im-fun/pkg/time"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

var Conf *Config

func InitConfig() (err error) {
	Conf = defaultConfig()
	bytes, err := ioutil.ReadFile("cmd/logic/logic.yaml")
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
	Debug      bool
	Discovery  *Discovery
	RPC        *RPC
	HTTPServer *HTTPServer
	Kafka      []Kafka
	Redis      []Redis
	MySQL      []MySQL
}

type Discovery struct {
	Addr string
}

// HTTPServer is http server config.
type HTTPServer struct {
	Network      string
	Addr         string
	ReadTimeout  mytime.Duration
	WriteTimeout mytime.Duration
}

type RPC struct {
	Server *RPCServer
	Client *RPCClient
}

// RPCClient is RPC client config.
type RPCClient struct {
	Dial    mytime.Duration
	Timeout mytime.Duration
}

// RPCServer is RPC server config.
type RPCServer struct {
	Network           string
	Addr              string
	Timeout           mytime.Duration
	IdleTimeout       mytime.Duration
	MaxLifeTime       mytime.Duration
	ForceCloseWait    mytime.Duration
	KeepAliveInterval mytime.Duration
	KeepAliveTimeout  mytime.Duration
}
