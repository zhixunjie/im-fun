//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package wire

import (
	"context"
	"github.com/google/wire"
	"github.com/zhixunjie/im-fun/internal/logic/api/grpc"
	"github.com/zhixunjie/im-fun/internal/logic/api/http"
	"github.com/zhixunjie/im-fun/internal/logic/biz"
	"github.com/zhixunjie/im-fun/internal/logic/conf"
	"github.com/zhixunjie/im-fun/internal/logic/data"
	google_grpc "google.golang.org/grpc"
)

func InitGrpc(ctx context.Context, c *conf.Config) (*google_grpc.Server, func(), error) {
	wire.Build(grpc.ProviderSet, biz.ProviderSet, data.ProviderSet)

	return &google_grpc.Server{}, nil, nil
}

func InitHttp(c *conf.Config) *http.Server {
	wire.Build(http.ProviderSet, data.ProviderSet, biz.ProviderSet)

	return nil
}

// 用于测试用例的各种对象

func GetMessageRepo(c *conf.Config) *data.MessageRepo {
	wire.Build(data.ProviderSet)

	return nil
}

func GetContactRepo(c *conf.Config) *data.ContactRepo {
	wire.Build(data.ProviderSet)

	return nil
}

func GetMessageUseCase(c *conf.Config) *biz.MessageUseCase {
	wire.Build(biz.ProviderSet, data.ProviderSet)

	return nil
}

func GetContactUseCase(c *conf.Config) *biz.ContactUseCase {
	wire.Build(biz.ProviderSet, data.ProviderSet)

	return nil
}
