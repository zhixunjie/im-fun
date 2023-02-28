package service

import (
	"github.com/zhixunjie/im-fun/internal/logic/conf"
	"github.com/zhixunjie/im-fun/internal/logic/dao"
)

type Service struct {
	conf *conf.Config
	dao  *dao.Dao
}

func New(conf *conf.Config) *Service {
	return &Service{
		conf: conf,
		dao:  dao.New(conf),
	}
}
