package conf

import (
	"errors"
	"github.com/zhixunjie/im-fun/pkg/buffer/bytes"
	"github.com/zhixunjie/im-fun/pkg/encoding/yaml"
	"github.com/zhixunjie/im-fun/pkg/env"
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

// Config is comet config.
type Config struct {
	Env       env.EnvType `yaml:"env"`       // 环境名
	Name      string      `yaml:"name"`      // 服务名
	Debug     bool        `yaml:"debug"`     // 是否开启debug
	Discovery *Discovery  `yaml:"discovery"` // etcd的配置
	Connect   *Connect    `yaml:"connect"`   // 长连接配置
	RPC       *RPC        `yaml:"rpc"`       // RPC配置
	Protocol  *Protocol   `yaml:"protocol"`  // 协议配置
	Bucket    *Bucket     `yaml:"bucket"`    // 桶配置
}

type Discovery struct {
	Addr string `yaml:"addr"`
}

type Connect struct {
	TCP           *TCP           `yaml:"tcp"`
	Websocket     *Websocket     `yaml:"websocket"`
	BufferOptions *bytes.Options `yaml:"bufferOptions"`
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
	KeepaliveInterval newtime.Duration `yaml:"keepaliveInterval"`
	KeepaliveTimeout  newtime.Duration `yaml:"keepaliveTimeout"`
}

// RPCClient is RPC client config.
type RPCClient struct {
	Dial    newtime.Duration `yaml:"dial"`
	Timeout newtime.Duration `yaml:"timeout"`
}

type TCP struct {
	Bind      []string `yaml:"bind"`      // 绑定的地址
	Sndbuf    int      `yaml:"sndbuf"`    // 内核缓冲区
	Rcvbuf    int      `yaml:"rcvbuf"`    // 内核缓冲区
	Keepalive bool     `yaml:"keepalive"` // 操作系统的Keepalive机制（自带的心跳机制）
}

type Websocket struct {
	Bind        []string `yaml:"bind"`        // 绑定的地址
	TLSOpen     bool     `yaml:"tlsOpen"`     // 是否开启TLS
	TLSBind     []string `yaml:"tlsBind"`     // TLS的绑定地址
	CertFile    string   `yaml:"certFile"`    // 证书文件
	PrivateFile string   `yaml:"privateFile"` // 私钥文件
}

type Protocol struct {
	TimerPool        *TimerPool       `yaml:"timerPool"`        // Timer池子的配置
	Proto            *Proto           `yaml:"proto"`            // proto相关的配置
	HandshakeTimeout newtime.Duration `yaml:"handshakeTimeout"` // TCP 握手超时
}

// TimerPool Timer池子的配置
type TimerPool struct {
	HashNum        int `yaml:"hashNum"`        // Hash切片数量：每个切片是一个Timer池子
	InitSizeInPool int `yaml:"initSizeInPool"` // 每个Timer池子中，拥有的Timer初始数量
}

type Proto struct {
	ChannelSize   int `yaml:"channelSize"`   // 每个TCP链接都有一个缓冲大小ChannelSize的Channel接收Proto
	AllocatorSize int `yaml:"allocatorSize"` // 一个Proto分配器的最大容量
}

type Bucket struct {
	HashNum            int `yaml:"hashNum"`            // Hash切片数量（每个切片是一个Bucket）
	InitSizeChannelMap int `yaml:"initSizeChannelMap"` // Channel Map的大小（初始大小）
	InitSizeRoomMap    int `yaml:"initSizeRoomMap"`    // Room Map的大小（初始大小）
	RoutineAmount      int `yaml:"routineHashNum"`     // Hash切片数量：每个切片是一个协程
	RoutineChannelSize int `yaml:"routineChannelSize"` // 每个协程拥有一个指定缓冲大小的Channel
}

// Env is env config.（暂无使用）
type Env struct {
	Region    string   `yaml:"region"`    // 地域
	Zone      string   `yaml:"zone"`      // 可用区
	DeployEnv string   `yaml:"deployEnv"` // 部署环境
	HostName  string   `yaml:"host"`      // 主机
	Weight    int64    `yaml:"weight"`    // 权重（负载均衡权重）
	Offline   bool     `yaml:"offline"`
	Addrs     []string `yaml:"addrs"`
}
