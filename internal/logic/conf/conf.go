package conf

import (
	"fmt"
	"github.com/sirupsen/logrus"
	newtime "github.com/zhixunjie/im-fun/pkg/time"
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
	Node       *Node
}

type Discovery struct {
	Addr string
}

// HTTPServer is http server config.
type HTTPServer struct {
	Network      string
	Addr         string
	ReadTimeout  newtime.Duration
	WriteTimeout newtime.Duration
}

type RPC struct {
	Server *RPCServer
	Client *RPCClient
}

// RPCClient is RPC client config.
type RPCClient struct {
	Dial    newtime.Duration
	Timeout newtime.Duration
}

// RPCServer is RPC server config.
type RPCServer struct {
	Network           string
	Addr              string
	Timeout           newtime.Duration
	IdleTimeout       newtime.Duration
	MaxLifeTime       newtime.Duration
	ForceCloseWait    newtime.Duration
	KeepAliveInterval newtime.Duration
	KeepAliveTimeout  newtime.Duration
}

// Node node config.
type Node struct {
	DefaultDomain string
	HostDomain    string
	TCPPort       int
	WSPort        int
	WSSPort       int
	HeartbeatMax  int
	Heartbeat     newtime.Duration
	RegionWeight  float64
}
