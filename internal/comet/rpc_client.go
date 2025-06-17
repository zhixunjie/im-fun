package comet

import (
	"context"
	"github.com/zhixunjie/im-fun/api/pb"
	"github.com/zhixunjie/im-fun/internal/comet/conf"
	"github.com/zhixunjie/im-fun/pkg/grpc"
	"time"
)

func newLogicClient(conf *conf.RPCClient) pb.LogicClient {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// dial
	conn, err := grpc.Connection(ctx, "discovery:///logic", true)
	if err != nil {
		panic(err)
	}
	return pb.NewLogicClient(conn)
}
