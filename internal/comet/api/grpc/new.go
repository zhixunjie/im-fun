package grpc

import (
	"context"
	"github.com/zhixunjie/im-fun/api/pb"
	"github.com/zhixunjie/im-fun/internal/comet"
	"github.com/zhixunjie/im-fun/internal/comet/conf"
	"github.com/zhixunjie/im-fun/pkg/logging"
	"github.com/zhixunjie/im-fun/pkg/micro_registry"
	google_grpc "google.golang.org/grpc"
	google_grpc_keepalive "google.golang.org/grpc/keepalive"
	"net"
	"time"
)

func NewServer(ctx context.Context, s *comet.Server, conf *conf.Config) (*google_grpc.Server, func(), error) {
	rpcConfig := conf.RPC.Server
	params := google_grpc.KeepaliveParams(google_grpc_keepalive.ServerParameters{
		MaxConnectionIdle:     time.Duration(rpcConfig.IdleTimeout),
		MaxConnectionAgeGrace: time.Duration(rpcConfig.ForceCloseWait),
		Time:                  time.Duration(rpcConfig.KeepaliveInterval),
		Timeout:               time.Duration(rpcConfig.KeepaliveTimeout),
		MaxConnectionAge:      time.Duration(rpcConfig.MaxLifeTime),
	})

	// register
	srv := google_grpc.NewServer(params)
	pb.RegisterCometServer(srv, &server{srv: s, UnimplementedCometServer: pb.UnimplementedCometServer{}})
	listener, err := net.Listen(rpcConfig.Network, rpcConfig.Addr)
	if err != nil {
		panic(err)
	}

	// begin to listen
	logging.Infof("GRPC server is listening %v：%v", rpcConfig.Network, rpcConfig.Addr)
	go func() {
		if err = srv.Serve(listener); err != nil {
			panic(err)
		}
	}()

	// register micro
	deRegisterFn, err := micro_registry.Register(ctx, conf.Name, rpcConfig.Addr, listener)
	if err != nil {
		panic(err)
	}

	return srv, deRegisterFn, nil
}
