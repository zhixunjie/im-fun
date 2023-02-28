package conf

import (
	newtime "github.com/zhixunjie/im-fun/pkg/time"
)

func DefaultRedis() []Redis {
	return []Redis{
		{
			Addr: "127.0.0.1:6379",
		},
	}
}

// Redis .
type Redis struct {
	Network      string
	Addr         string
	Auth         string
	Active       int
	Idle         int
	DialTimeout  newtime.Duration
	ReadTimeout  newtime.Duration
	WriteTimeout newtime.Duration
	IdleTimeout  newtime.Duration
	Expire       newtime.Duration
}
