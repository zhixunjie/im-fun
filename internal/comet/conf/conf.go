package conf

import (
	"github.com/sirupsen/logrus"
	"github.com/zhixunjie/im-fun/pkg/buffer"
	newtime "github.com/zhixunjie/im-fun/pkg/time"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"time"
)

var Conf *Config

func InitConfig() (err error) {
	Conf = defaultConfig()
	bytes, err := ioutil.ReadFile("cmd/comet/comet.yaml")
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
	logrus.Infof("config=%+v\n", Conf)
	
	return nil
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
			Timer:            32,
			TimerSize:        2048,
			ClientProtoNum:   64,
			ServerProtoNum:   64,
			HandshakeTimeout: newtime.Duration(time.Second * 5),
		},
		Bucket: &Bucket{
			Size:          32,
			Channel:       1024,
			Room:          1024,
			RoutineAmount: 32,
			RoutineSize:   1024,
		},
	}

	return val
}

// Config is comet config.
type Config struct {
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
	Host      string   `yaml:"host"`      // 主机
	Weight    int64    `yaml:"weight"`    // 权重（负载均衡权重）
	Offline   bool     `yaml:"offline"`
	Addrs     []string `yaml:"addrs"`
}

type Connect struct {
	TCP           *TCP            `yaml:"tcp"`
	Websocket     *Websocket      `yaml:"websocket"`
	BufferOptions *buffer.Options `yaml:"bufferOptions"`
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
	Sndbuf    int      `yaml:"sndbuf"`
	Rcvbuf    int      `yaml:"rcvbuf"`
	Keepalive bool     `yaml:"keepalive"`
}

type Websocket struct {
	Bind        []string `yaml:"bind"`
	TLSOpen     bool     `yaml:"tlsOpen"`
	TLSBind     []string `yaml:"tlsBind"`
	CertFile    string   `yaml:"certFile"`
	PrivateFile string   `yaml:"privateFile"`
}

type Protocol struct {
	Timer            int              `yaml:"timer"`
	TimerSize        int              `yaml:"timerSize"`
	ServerProtoNum   int              `yaml:"serverProtoNum"`
	ClientProtoNum   int              `yaml:"clientProtoNum"`
	HandshakeTimeout newtime.Duration `yaml:"handshakeTimeout"`
}

type Bucket struct {
	Size          int    `yaml:"size"`
	Channel       int    `yaml:"channel"`
	Room          int    `yaml:"room"`
	RoutineAmount uint64 `yaml:"routineAmount"`
	RoutineSize   int    `yaml:"routineSize"`
}
