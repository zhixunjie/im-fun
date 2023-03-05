package conf

import (
	"github.com/zhixunjie/im-fun/pkg/buffer"
	newtime "github.com/zhixunjie/im-fun/pkg/time"
	"os"
	"strconv"
	"strings"
	"time"
)

func DefaultEnv() *Env {
	defHost, _ := os.Hostname()
	defWeight, _ := strconv.ParseInt(os.Getenv("WEIGHT"), 10, 32)
	defAddrs := os.Getenv("ADDRS")
	defOffline, _ := strconv.ParseBool(os.Getenv("OFFLINE"))
	//defDebug, _ := strconv.ParseBool(os.Getenv("DEBUG"))

	return &Env{
		Region:    os.Getenv("REGION"),
		Zone:      os.Getenv("ZONE"),
		DeployEnv: os.Getenv("DEPLOY_ENV"),
		Host:      defHost,
		Weight:    defWeight,
		Addrs:     strings.Split(defAddrs, ","),
		Offline:   defOffline,
	}
}

func DefaultRPC() *RPC {
	return &RPC{
		Server: &RPCServer{
			Network:           "tcp",
			Addr:              ":3109",
			Timeout:           newtime.Duration(time.Second),
			IdleTimeout:       newtime.Duration(time.Second * 60),
			MaxLifeTime:       newtime.Duration(time.Hour * 2),
			ForceCloseWait:    newtime.Duration(time.Second * 20),
			KeepaliveInterval: newtime.Duration(time.Second * 60),
			KeepaliveTimeout:  newtime.Duration(time.Second * 20),
		},
		Client: &RPCClient{
			Dial:    newtime.Duration(time.Second),
			Timeout: newtime.Duration(time.Second),
		},
	}
}

func DefaultConnect() *Connect {
	return &Connect{
		TCP: &TCP{
			Bind:      []string{":3101"},
			Sndbuf:    4096,
			Rcvbuf:    4096,
			Keepalive: false,
		},
		Websocket: &Websocket{
			Bind: []string{":3102"},
		},
		BufferOptions: &buffer.Options{
			ReadPoolOption: buffer.PoolOptions{
				PoolNum:  10,
				BatchNum: 1024,
				BufSize:  8192,
			},
			WritePoolOption: buffer.PoolOptions{
				PoolNum:  10,
				BatchNum: 1024,
				BufSize:  8192,
			},
		},
	}
}
