package job

import (
	"context"
	pb "github.com/zhixunjie/im-fun/api/pb"
	"github.com/zhixunjie/im-fun/internal/job/conf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"time"
)

// Logic -> Job -> CometInvoker -> Comet

type CometInvoker struct {
	ctx    context.Context
	cancel context.CancelFunc

	serverId  string
	rpcClient pb.CometClient

	// number
	pushChanNum uint64
	roomChanNum uint64
	routineNum  uint64

	// send msg
	chUserKeys []chan *pb.SendToUserKeysReq // send msg: to some user
	chRoom     []chan *pb.SendToRoomReq     // send msg: to the room
	chAll      chan *pb.SendToAllReq        // send msg: to all user
}

func NewCometInvoker(serverId string, c *conf.CometInvoker) (*CometInvoker, error) {
	routineNum := c.RoutineNum
	cmt := &CometInvoker{
		serverId:   serverId,
		chUserKeys: make([]chan *pb.SendToUserKeysReq, routineNum),
		chRoom:     make([]chan *pb.SendToRoomReq, routineNum),
		chAll:      make(chan *pb.SendToAllReq, routineNum),
		routineNum: uint64(routineNum),
	}

	// create rpc client
	var err error
	grpcAddr := "127.0.0.1:12570"
	if cmt.rpcClient, err = newCometClient(grpcAddr); err != nil {
		return nil, err
	}

	// creat channel and routine
	cmt.ctx, cmt.cancel = context.WithCancel(context.Background())
	for i := 0; i < routineNum; i++ {
		cmt.chUserKeys[i] = make(chan *pb.SendToUserKeysReq, c.ChanNum)
		cmt.chRoom[i] = make(chan *pb.SendToRoomReq, c.ChanNum)
		go cmt.Process(i)
	}
	return cmt, nil
}

var (
	// grpc options
	grpcKeepAliveTime    = time.Duration(10) * time.Second
	grpcKeepAliveTimeout = time.Duration(3) * time.Second
	grpcBackoffMaxDelay  = time.Duration(3) * time.Second
	grpcMaxSendMsgSize   = 1 << 24
	grpcMaxCallMsgSize   = 1 << 24
)

const (
	// grpc options
	grpcInitialWindowSize     = 1 << 24
	grpcInitialConnWindowSize = 1 << 24
)

func newCometClient(addr string) (pb.CometClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, addr,
		[]grpc.DialOption{
			grpc.WithInsecure(),
			grpc.WithInitialWindowSize(grpcInitialWindowSize),
			grpc.WithInitialConnWindowSize(grpcInitialConnWindowSize),
			grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(grpcMaxCallMsgSize)),
			grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(grpcMaxSendMsgSize)),
			grpc.WithBackoffMaxDelay(grpcBackoffMaxDelay),
			grpc.WithKeepaliveParams(keepalive.ClientParameters{
				Time:                grpcKeepAliveTime,
				Timeout:             grpcKeepAliveTimeout,
				PermitWithoutStream: true,
			}),
		}...,
	)
	if err != nil {
		return nil, err
	}
	return pb.NewCometClient(conn), err
}
