package grpc

import (
	"context"
	"github.com/google/wire"
	"github.com/zhixunjie/im-fun/api/pb"
	"github.com/zhixunjie/im-fun/internal/logic/biz"
	"github.com/zhixunjie/im-fun/internal/logic/conf"
	"github.com/zhixunjie/im-fun/pkg/logging"
	"github.com/zhixunjie/im-fun/pkg/registry"
	google_grpc "google.golang.org/grpc"
	google_grpc_keepalive "google.golang.org/grpc/keepalive"
	"time"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(NewServer)

// NewServer logic grpc server
func NewServer(ctx context.Context, conf *conf.Config, bz *biz.Biz, bzContact *biz.ContactUseCase, bzMessage *biz.MessageUseCase) (*google_grpc.Server, func(), error) {
	rpcConfig := conf.RPC.Server
	params := google_grpc.KeepaliveParams(google_grpc_keepalive.ServerParameters{
		MaxConnectionIdle:     time.Duration(rpcConfig.IdleTimeout),
		MaxConnectionAgeGrace: time.Duration(rpcConfig.ForceCloseWait),
		Time:                  time.Duration(rpcConfig.KeepAliveInterval),
		Timeout:               time.Duration(rpcConfig.KeepAliveTimeout),
		MaxConnectionAge:      time.Duration(rpcConfig.MaxLifeTime),
	})

	// build instance
	instance, err := registry.BuildServiceInstance(conf.Name, rpcConfig.Network, rpcConfig.Addr)
	if err != nil {
		panic(err)
	}

	// register to grpc.Server
	srv := google_grpc.NewServer(params)
	pb.RegisterLogicServer(srv, &server{
		UnimplementedLogicServer: pb.UnimplementedLogicServer{},
		bz:                       bz,
		bzContact:                bzContact,
		bzMessage:                bzMessage,
	})

	// begin to listen
	logging.Infof("GRPC server is listening %v：%v", rpcConfig.Network, rpcConfig.Addr)
	go func() {
		if err = srv.Serve(instance.GrpcServerListener); err != nil {
			panic(err)
		}
	}()

	// register micro
	deRegisterFn, err := registry.Register(ctx, instance)
	if err != nil {
		panic(err)
	}

	return srv, deRegisterFn, nil
}
