package invoker

import (
	"context"
	"github.com/zhixunjie/im-fun/api/pb"
	"github.com/zhixunjie/im-fun/pkg/grpc"
)

func newCometClient(addr string) (pb.CometClient, error) {
	// dial
	conn, err := grpc.Connection(context.Background(), addr, false)
	if err != nil {
		return nil, err
	}
	return pb.NewCometClient(conn), err
}
