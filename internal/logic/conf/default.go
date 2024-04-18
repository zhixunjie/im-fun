package conf

import (
	newtime "github.com/zhixunjie/im-fun/pkg/time"
	"time"
)

func defaultConfig() *Config {
	return &Config{
		Name:  "logic",
		Debug: false,
		Discovery: &Discovery{
			Addr: "127.0.0.1:7171",
		},
		RPC: DefaultRPC(),
		HTTPServer: &HTTPServer{
			Network:      "tcp",
			Addr:         ":8080",
			ReadTimeout:  newtime.Duration(time.Second),
			WriteTimeout: newtime.Duration(time.Second),
		},
		Kafka: DefaultKafka(),
		MySQL: DefaultMySQL(),
		Redis: DefaultRedis(),
	}
}

func DefaultRPC() *RPC {
	return &RPC{
		Server: &RPCServer{
			Network:           "tcp",
			Addr:              ":12670",
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
