package conf

import (
	"github.com/zhixunjie/im-fun/pkg/buffer"
	newtime "github.com/zhixunjie/im-fun/pkg/time"
	"time"
)

func DefaultRPC() *RPC {
	return &RPC{
		Server: &RPCServer{
			Network:           "tcp",
			Addr:              ":3109",
			Timeout:           newtime.Duration(time.Second),
			IdleTimeout:       newtime.Duration(time.Second * 60),
			MaxLifeTime:       newtime.Duration(time.Hour * 2),
			ForceCloseWait:    newtime.Duration(time.Second * 20),
			KeepAliveInterval: newtime.Duration(time.Second * 60),
			KeepAliveTimeout:  newtime.Duration(time.Second * 20),
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
			KeepAlive: false,
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
