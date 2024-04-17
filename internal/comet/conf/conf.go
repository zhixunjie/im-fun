package conf

import (
	"github.com/zhixunjie/im-fun/pkg/buffer/bytes"
	"github.com/zhixunjie/im-fun/pkg/encoding/yaml"
	newtime "github.com/zhixunjie/im-fun/pkg/time"
	"time"
)

var Conf *Config

func InitConfig(path string) (err error) {
	Conf = defaultConfig()
	return yaml.LoadConfig(path, Conf)
}

func defaultConfig() *Config {
	val := &Config{
		Debug: false,
		Env:   DefaultEnv(),
		Discovery: &Discovery{
			Addr: "127.0.0.1:7171",
		},
		RPC:     DefaultRPC(),
		Connect: DefaultConnect(),
		Protocol: &Protocol{
			TimerHashNum:       32,
			InitSizeTimerPool:  2048,
			ProtoAllocatorSize: 64,
			ProtoChannelSize:   64,
			HandshakeTimeout:   newtime.Duration(time.Second * 5),
		},
		Bucket: &Bucket{
			HashNum:            32,
			InitSizeChannelMap: 1024,
			InitSizeRoomMap:    1024,
			RoutineAmount:      32,
			RoutineChannelSize: 1024,
		},
	}

	return val
}

// Config is comet config.
type Config struct {
	Name      string     `yaml:"name"`
	Debug     bool       `yaml:"debug"`
	Env       *Env       `yaml:"env"`
	Discovery *Discovery `yaml:"discovery"`
	Connect   *Connect   `yaml:"connect"`
	RPC       *RPC       `yaml:"rpc"`
	Protocol  *Protocol  `yaml:"protocol"`
	Bucket    *Bucket    `yaml:"bucket"`
}

type Discovery struct {
	Addr string `yaml:"addr"`
}

// Env is env config.
type Env struct {
	Region    string   `yaml:"region"`    // 地域
	Zone      string   `yaml:"zone"`      // 可用区
	DeployEnv string   `yaml:"deployEnv"` // 部署环境
	HostName  string   `yaml:"host"`      // 主机
	Weight    int64    `yaml:"weight"`    // 权重（负载均衡权重）
	Offline   bool     `yaml:"offline"`
	Addrs     []string `yaml:"addrs"`
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
	Network           string           `yaml:"network"`
	Addr              string           `yaml:"addr"`
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
	Bind      []string `yaml:"bind"`
	Sndbuf    int      `yaml:"sndbuf"`    // 内核缓冲区
	Rcvbuf    int      `yaml:"rcvbuf"`    // 内核缓冲区
	Keepalive bool     `yaml:"keepalive"` // 操作系统的Keepalive机制（自带的心跳机制）
}

type Websocket struct {
	Bind        []string `yaml:"bind"`
	TLSOpen     bool     `yaml:"tlsOpen"`
	TLSBind     []string `yaml:"tlsBind"`
	CertFile    string   `yaml:"certFile"`
	PrivateFile string   `yaml:"privateFile"`
}

type Protocol struct {
	TimerHashNum       int              `yaml:"timerHashNum"`       // Hash切片数量（每个切片是一个Timer池子）
	InitSizeTimerPool  int              `yaml:"initSizeTimerPool"`  // 每个Timer池子拥有的Timer数量（初始数量）
	ProtoChannelSize   int              `yaml:"protoChannelSize"`   // 接收Proto的Channel的大小
	ProtoAllocatorSize int              `yaml:"protoAllocatorSize"` // Proto分配器的大小（本质是一个Ring）
	HandshakeTimeout   newtime.Duration `yaml:"handshakeTimeout"`   // TCP 握手超时
}

type Bucket struct {
	HashNum            int `yaml:"hashNum"`            // Hash切片数量（每个切片是一个Bucket）
	InitSizeChannelMap int `yaml:"initSizeChannelMap"` // Channel Map的大小（初始大小）
	InitSizeRoomMap    int `yaml:"initSizeRoomMap"`    // Room Map的大小（初始大小）
	RoutineAmount      int `yaml:"routineHashNum"`     // Hash切片数量（每个切片是一个协程）
	RoutineChannelSize int `yaml:"routineChannelSize"` // 每个协程拥有一个指定缓冲大小的Channel
}
