//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package wire

import (
	"github.com/google/wire"
	"github.com/zhixunjie/im-fun/internal/logic/api/grpc"
	"github.com/zhixunjie/im-fun/internal/logic/api/http"
	"github.com/zhixunjie/im-fun/internal/logic/biz"
	"github.com/zhixunjie/im-fun/internal/logic/conf"
	"github.com/zhixunjie/im-fun/internal/logic/data"
	google_grpc "google.golang.org/grpc"
)

func InitGrpc(c *conf.Config) *google_grpc.Server {
	wire.Build(grpc.ProviderSet, biz.ProviderSet, data.ProviderSet)

	return nil
}

func InitHttp(c *conf.Config) *http.Server {
	wire.Build(http.ProviderSet, data.ProviderSet, biz.ProviderSet)

	return nil
}
