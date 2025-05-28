package conf

import (
	"errors"
	gconf "github.com/zhixunjie/im-fun/pkg/conf"
	"github.com/zhixunjie/im-fun/pkg/encoding/yaml"
	"github.com/zhixunjie/im-fun/pkg/env"
	"github.com/zhixunjie/im-fun/pkg/gomysql"
	"github.com/zhixunjie/im-fun/pkg/kafka"
	newtime "github.com/zhixunjie/im-fun/pkg/time"
)

var Conf *Config

func InitConfig(path string) (err error) {
	Conf = defaultConfig()
	err = yaml.LoadConfig(path, Conf)
	if err != nil {
		return
	}
	env.InitEnv(Conf.Env, Conf.Name)
	if !env.IsLocal() && !env.IsTest() && !env.IsProd() {
		err = errors.New("env is invalid")
		return
	}
	return
}

type Config struct {
	Env          env.Type               `yaml:"env"`           // 环境名
	Name         string                 `yaml:"name"`          // 服务名
	Debug        bool                   `yaml:"debug"`         // 是否开启debug
	Discovery    *Discovery             `yaml:"discovery"`     // etcd的配置
	RPC          *RPC                   `yaml:"rpc"`           // RPC配置
	HTTPServer   *HTTPServer            `yaml:"http"`          // HTTP配置
	Kafka        []kafka.ProducerConf   `yaml:"kafka"`         // Kafka的配置
	Redis        []gconf.Redis          `yaml:"redis"`         // Redis的配置
	MySQL        []gconf.MySQL          `yaml:"mysql"`         // MySQL的配置
	MySQLCluster []gomysql.MySQLCluster `yaml:"mysql_cluster"` // MySQL的集群配置
	Node         *Node                  `yaml:"node"`          // 节点配置（客户端获取节点配置，然后根据配置进行TCP连接）
	Backoff      *Backoff               `yaml:"backoff"`       // GRPC Client用到的配置
	Regions      map[string][]string    `yaml:"regions"`       // 省份与Region的映射关系
}

type Discovery struct {
	Addr string `yaml:"addr"`
}

// HTTPServer is http server config.
type HTTPServer struct {
	Network      string           `yaml:"network"` //
	Addr         string           `yaml:"addr"`
	ReadTimeout  newtime.Duration `yaml:"readTimeout"`
	WriteTimeout newtime.Duration `yaml:"writeTimeout"`
}

type RPC struct {
	Server *RPCServer `yaml:"server"`
	Client *RPCClient `yaml:"client"`
}

// RPCServer is RPC server config.
type RPCServer struct {
	Network           string           `yaml:"network"` // 使用协议，如：tcp
	Addr              string           `yaml:"addr"`    // 服务器地址
	Timeout           newtime.Duration `yaml:"timeout"`
	IdleTimeout       newtime.Duration `yaml:"idleTimeout"`
	MaxLifeTime       newtime.Duration `yaml:"maxLifeTime"`
	ForceCloseWait    newtime.Duration `yaml:"forceCloseWait"`
	KeepAliveInterval newtime.Duration `yaml:"keepAliveInterval"`
	KeepAliveTimeout  newtime.Duration `yaml:"keepAliveTimeout"`
}

// RPCClient is RPC client config.
type RPCClient struct {
	Dial    newtime.Duration `yaml:"dial"`
	Timeout newtime.Duration `yaml:"timeout"`
}

// Node node config.
type Node struct {
	DefaultDomain string    `yaml:"defaultDomain"` // 默认域名
	HostDomain    string    `yaml:"hostDomain"`    // 主机域名
	TCPPort       int32     `yaml:"tcpPort"`       // 如果采用TCP连接，采用哪个端口？
	WSPort        int32     `yaml:"wsPort"`        // 如果采用WS连接，采用哪个端口？
	WSSPort       int32     `yaml:"wssPort"`       // 如果采用WSS连接，采用哪个端口？
	Heartbeat     Heartbeat `yaml:"heartbeat"`     // 心跳配置
	RegionWeight  float64   `yaml:"regionWeight"`  // 权重扩大比例（如果节点地址所属的Region与当前客户端Region一致，则增加该节点的权重值）
}

type Heartbeat struct {
	Interval  newtime.Duration `yaml:"interval"`  // 每次心跳发送的间隔
	FailCount int64            `yaml:"failCount"` // 允许心跳失败的次数（超过次数才算是心跳失败）
}

// Backoff GRPC Client用到的配置
type Backoff struct {
	BaseDelay  int32   `yaml:"baseDelay"`
	Multiplier float32 `yaml:"multiplier"`
	Jitter     float32 `yaml:"jitter"`
	MaxDelay   int32   `yaml:"maxDelay"`
}
