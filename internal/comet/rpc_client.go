package comet

import (
	"context"
	kratos_discovery "github.com/go-kratos/kratos/v2/transport/grpc/resolver/discovery"
	"github.com/zhixunjie/im-fun/api/pb"
	"github.com/zhixunjie/im-fun/internal/comet/conf"
	"github.com/zhixunjie/im-fun/pkg/registry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"time"
)

// grpc options
const (
	grpcInitialWindowSize     = 1 << 24
	grpcInitialConnWindowSize = 1 << 24
	grpcMaxSendMsgSize        = 1 << 24
	grpcMaxCallMsgSize        = 1 << 24

	grpcKeepAliveTime    = time.Second * 10
	grpcKeepAliveTimeout = time.Second * 3
	grpcBackoffMaxDelay  = time.Second * 3
)

func newLogicClient(conf *conf.RPCClient) pb.LogicClient {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// dial
	conn, err := grpc.DialContext(ctx, "discovery:///logic",
		[]grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithInitialWindowSize(grpcInitialWindowSize),
			grpc.WithInitialConnWindowSize(grpcInitialConnWindowSize),
			grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(grpcMaxCallMsgSize)),
			grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(grpcMaxSendMsgSize)),
			// https://pkg.go.dev/google.golang.org/grpc/backoff
			grpc.WithConnectParams(grpc.ConnectParams{
				Backoff: backoff.Config{
					BaseDelay:  1.0 * time.Second,
					Multiplier: 1.6,
					Jitter:     0.2,
					//MaxDelay:   120 * time.Second,
					MaxDelay: grpcBackoffMaxDelay,
				},
				MinConnectTimeout: 20 * time.Second,
			}),
			grpc.WithKeepaliveParams(keepalive.ClientParameters{
				Time:                grpcKeepAliveTime,
				Timeout:             grpcKeepAliveTimeout,
				PermitWithoutStream: true,
			}),
			// GRPC服务发现机制（使用Kratos的注册中心和解析器）
			grpc.WithResolvers(kratos_discovery.NewBuilder(registry.KratosEtcdRegistry, kratos_discovery.WithInsecure(true))),
		}...)
	if err != nil {
		panic(err)
	}
	return pb.NewLogicClient(conn)
}
