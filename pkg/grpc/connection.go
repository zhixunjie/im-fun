package grpc

import (
	"context"
	"fmt"
	kratos_discovery "github.com/go-kratos/kratos/v2/transport/grpc/resolver/discovery"
	"github.com/zhixunjie/im-fun/pkg/registry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"time"
)

func Connection(ctx context.Context, addr string, useKratosResolver bool) (conn *grpc.ClientConn, err error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	// set options
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithInitialWindowSize(1 << 24),
		grpc.WithInitialConnWindowSize(1 << 24),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1 << 24)),
		grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(1 << 24)),
		// https://pkg.go.dev/google.golang.org/grpc/backoff
		grpc.WithConnectParams(grpc.ConnectParams{
			Backoff: backoff.Config{
				BaseDelay:  1 * time.Second,
				Multiplier: 1.6,
				Jitter:     0.2,
				MaxDelay:   3 * time.Second,
				//MaxDelay:   120 * time.Second,
			},
			MinConnectTimeout: 20 * time.Second,
		}),
		// 设置 keepalive 参数
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                30 * time.Second, // 推荐 30s ~ 60s
			Timeout:             10 * time.Second, // 等待响应的超时时间
			PermitWithoutStream: false,            // 空闲连接不发送 PING，避免被判为“过于活跃”
		}),
	}

	// GRPC的服务发现机制
	// 使用Kratos的注册中心和解析器
	if useKratosResolver {
		opts = append(opts,
			grpc.WithResolvers(kratos_discovery.NewBuilder(registry.KratosEtcdRegistry, kratos_discovery.WithInsecure(true))),
		)
	}

	// dial
	conn, err = grpc.DialContext(ctx, addr, opts...)
	if err != nil {
		err = fmt.Errorf("failed to dial grpc connection: %w", err)
		return
	}
	return
}
