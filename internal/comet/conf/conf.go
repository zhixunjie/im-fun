package conf

import (
	"fmt"
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
	bytes, err := ioutil.ReadFile("file/example.yaml")
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

// defaultConfig new a config with specified defualt value.
func defaultConfig() *Config {
	val := &Config{
		Debug: false,
		Discovery: &Discovery{
			Addr: "127.0.0.1:7171",
		},
		RPC:     DefaultRPC(),
		Connect: DefaultConnect(),
		Protocol: &Protocol{
			Timer:            32,
			TimerSize:        2048,
			CliProto:         5,
			SvrProto:         10,
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
	Debug     bool
	Discovery *Discovery
	Connect   *Connect
	RPC       *RPC
	Protocol  *Protocol
	Bucket    *Bucket
}

type Discovery struct {
	Addr string
}
type Connect struct {
	TCP           *TCP
	Websocket     *Websocket
	BufferOptions *buffer.Options
}
type RPC struct {
	Server *RPCServer
	Client *RPCClient
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

// RPCClient is RPC client config.
type RPCClient struct {
	Dial    newtime.Duration
	Timeout newtime.Duration
}

type TCP struct {
	Bind      []string
	Sndbuf    int
	Rcvbuf    int
	KeepAlive bool
}

type Websocket struct {
	Bind        []string
	TLSOpen     bool
	TLSBind     []string
	CertFile    string
	PrivateFile string
}

type Protocol struct {
	Timer            int
	TimerSize        int
	SvrProto         int
	CliProto         int
	HandshakeTimeout newtime.Duration
}

type Bucket struct {
	Size          int
	Channel       int
	Room          int
	RoutineAmount uint64
	RoutineSize   int
}
