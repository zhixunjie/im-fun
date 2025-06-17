package comet

import (
	"context"
	"github.com/zhixunjie/im-fun/api/pb"
	"github.com/zhixunjie/im-fun/internal/comet/conf"
	"github.com/zhixunjie/im-fun/pkg/grpc"
)

func newLogicClient(conf *conf.RPCClient) pb.LogicClient {
	// dial
	conn, err := grpc.Connection(context.Background(), "discovery:///logic", true)
	if err != nil {
		panic(err)
	}
	return pb.NewLogicClient(conn)
}
