package biz

import (
	"github.com/google/wire"
	"github.com/zhixunjie/im-fun/internal/logic/conf"
	"github.com/zhixunjie/im-fun/internal/logic/data"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(NewBiz, NewMessageUseCase, NewContactUseCase)

// Biz 通用的对象
type Biz struct {
	conf *conf.Config
	data *data.Data
}

func NewBiz(conf *conf.Config) *Biz {
	return &Biz{
		conf: conf,
		data: data.NewData(conf),
	}
}

func (bz *Biz) GetData() *data.Data {
	return bz.data
}

func (bz *Biz) Ping() {

}
