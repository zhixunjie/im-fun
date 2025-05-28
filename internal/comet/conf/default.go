package conf

import (
	"github.com/zhixunjie/im-fun/pkg/buffer/bytes"
	"github.com/zhixunjie/im-fun/pkg/env"
	newtime "github.com/zhixunjie/im-fun/pkg/time"
	"os"
	"strconv"
	"strings"
	"time"
)

func defaultConfig() *Config {
	val := &Config{
		//Env:       DefaultEnv(),
		Env:       env.TypeLocal,
		Name:      "comet",
		Debug:     false,
		Discovery: &Discovery{},
		RPC:       DefaultRPC(),
		Connect:   DefaultConnect(),
		Protocol: &Protocol{
			TimerPool: &TimerPool{
				HashNum:        32,
				InitSizeInPool: 2048,
			},
			Proto: &Proto{
				ChannelSize:   64,
				AllocatorSize: 64,
			},
			HandshakeTimeout: newtime.Duration(time.Second * 5),
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
		HostName:  defHost,
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
		BufferOptions: &bytes.Options{
			ReadPool: bytes.PoolOptions{
				PoolNum:  10,
				BatchNum: 1024,
				BufSize:  8192,
			},
			WritePool: bytes.PoolOptions{
				PoolNum:  10,
				BatchNum: 1024,
				BufSize:  8192,
			},
		},
	}
}
