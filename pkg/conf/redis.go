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
	Network      string           `yaml:"network"`
	Addr         string           `yaml:"addr"`
	Auth         string           `yaml:"auth"`
	Active       int              `yaml:"active"`
	Idle         int              `yaml:"idle"`
	DialTimeout  newtime.Duration `yaml:"dialTimeout"`
	ReadTimeout  newtime.Duration `yaml:"readTimeout"`
	WriteTimeout newtime.Duration `yaml:"writeTimeout"`
	IdleTimeout  newtime.Duration `yaml:"idleTimeout"`
	KeyExpire    newtime.Duration `yaml:"keyExpire"`
}
