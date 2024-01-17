package grpc

import (
	"github.com/google/wire"
	"github.com/zhixunjie/im-fun/api/pb"
	"github.com/zhixunjie/im-fun/internal/logic/biz"
	"github.com/zhixunjie/im-fun/internal/logic/conf"
	"github.com/zhixunjie/im-fun/pkg/logging"
	google_grpc "google.golang.org/grpc"
	google_grpc_keepalive "google.golang.org/grpc/keepalive"
	"net"
	"time"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(NewServer)

// NewServer logic grpc server
func NewServer(conf *conf.Config, bz *biz.Biz, bzContact *biz.ContactUseCase, bzMessage *biz.MessageUseCase) *google_grpc.Server {
	rpcConfig := conf.RPC.Server
	params := google_grpc.KeepaliveParams(google_grpc_keepalive.ServerParameters{
		MaxConnectionIdle:     time.Duration(rpcConfig.IdleTimeout),
		MaxConnectionAgeGrace: time.Duration(rpcConfig.ForceCloseWait),
		Time:                  time.Duration(rpcConfig.KeepAliveInterval),
		Timeout:               time.Duration(rpcConfig.KeepAliveTimeout),
		MaxConnectionAge:      time.Duration(rpcConfig.MaxLifeTime),
	})

	// register
	srv := google_grpc.NewServer(params)
	pb.RegisterLogicServer(srv, &server{
		UnimplementedLogicServer: pb.UnimplementedLogicServer{},
		bz:                       bz,
		bzContact:                bzContact,
		bzMessage:                bzMessage,
	})

	// begin to listen
	listener, err := net.Listen(rpcConfig.Network, rpcConfig.Addr)
	if err != nil {
		panic(err)
	}
	logging.Infof("GRPC server is listening %v：%v", rpcConfig.Network, rpcConfig.Addr)
	go func() {
		if err = srv.Serve(listener); err != nil {
			panic(err)
		}
	}()
	return srv
}
