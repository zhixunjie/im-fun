package invoker

import (
	"context"
	"github.com/zhixunjie/im-fun/api/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"time"
)

func newCometClient(addr string) (pb.CometClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// dial
	conn, err := grpc.DialContext(ctx, addr,
		[]grpc.DialOption{
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
		}...,
	)
	if err != nil {
		return nil, err
	}
	return pb.NewCometClient(conn), err
}
