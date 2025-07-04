package biz

import (
	"github.com/google/wire"
	"github.com/zhixunjie/im-fun/internal/logic/conf"
	"github.com/zhixunjie/im-fun/internal/logic/data"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(
	NewBiz,
	NewMessageFilterUseCase,
	NewMessageUseCase, NewGroupMessageUseCase, NewContactUseCase,
	NewUserUseCase, NewUserGroupUseCase,
)

// Biz 通用的对象
type Biz struct {
	conf        *conf.Config
	data        *data.Data
	lb          *LoadBalance // 暂无使用（比较复杂）
	userUseCase *UserUseCase
}

func NewBiz(conf *conf.Config, userUseCase *UserUseCase) *Biz {
	// TODO：负载均衡机制（结合Job的watch机制+配置中心：https://go-kratos.dev/docs/component/config）
	return &Biz{
		conf:        conf,
		data:        data.NewData(conf),
		lb:          NewLoadBalance(conf),
		userUseCase: userUseCase,
	}
}

func (bz *Biz) GetData() *data.Data {
	return bz.data
}

func (bz *Biz) Ping() {

}
