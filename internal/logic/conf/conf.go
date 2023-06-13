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
	bytes, err := ioutil.ReadFile("cmd/logic/logic.yaml")
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
	Debug      bool                 `yaml:"debug"`
	Discovery  *Discovery           `yaml:"discovery"`
	RPC        *RPC                 `yaml:"rpc"`
	HTTPServer *HTTPServer          `yaml:"http"`
	Kafka      []kafka.ProducerConf `yaml:"kafka"`
	Redis      []Redis              `yaml:"redis"`
	MySQL      []MySQL              `yaml:"mysql"`
	Node       *Node                `yaml:"node"`
}

type Discovery struct {
	Addr string `yaml:"addr"`
}

// HTTPServer is http server config.
type HTTPServer struct {
	Network      string           `yaml:"network"`
	Addr         string           `yaml:"addr"`
	ReadTimeout  newtime.Duration `yaml:"readTimeout"`
	WriteTimeout newtime.Duration `yaml:"writeTimeout"`
}

type RPC struct {
	Server *RPCServer `yaml:"server"`
	Client *RPCClient `yaml:"client"`
}

// RPCClient is RPC client config.
type RPCClient struct {
	Dial    newtime.Duration `yaml:"dial"`
	Timeout newtime.Duration `yaml:"timeout"`
}

// RPCServer is RPC server config.
type RPCServer struct {
	Network           string           `yaml:"network"`
	Addr              string           `yaml:"addr"`
	Timeout           newtime.Duration `yaml:"timeout"`
	IdleTimeout       newtime.Duration `yaml:"idleTimeout"`
	MaxLifeTime       newtime.Duration `yaml:"maxLifeTime"`
	ForceCloseWait    newtime.Duration `yaml:"forceCloseWait"`
	KeepAliveInterval newtime.Duration `yaml:"keepAliveInterval"`
	KeepAliveTimeout  newtime.Duration `yaml:"keepAliveTimeout"`
}

// Node node config.
type Node struct {
	DefaultDomain string           `yaml:"defaultDomain"`
	HostDomain    string           `yaml:"hostDomain"`
	TCPPort       int              `yaml:"tcpPort"`
	WSPort        int              `yaml:"wsPort"`
	WSSPort       int              `yaml:"wssPort"`
	HeartbeatMax  int              `yaml:"heartbeatMax"`
	Heartbeat     newtime.Duration `yaml:"heartbeat"`
	RegionWeight  float64          `yaml:"regionWeight"`
}
