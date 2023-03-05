package job

import (
	"context"
	"github.com/zhixunjie/im-fun/api/comet"
	"github.com/zhixunjie/im-fun/internal/job/conf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"time"
)

// Logic -> Job -> Comet

type Comet struct {
	serverId     string
	rpcClient    comet.CometClient
	userKeysChan []chan *comet.PushUserKeysReq // push to some user
	userRoomChan []chan *comet.PushUserRoomReq // push to the room
	userAllChan  chan *comet.PushUserAllReq    // push to all user
	pushChanNum  uint64
	roomChanNum  uint64
	routineNum   uint64

	ctx    context.Context
	cancel context.CancelFunc
}

func NewComet(serverId string, c *conf.Comet) (*Comet, error) {
	routineNum := c.RoutineNum
	cmt := &Comet{
		serverId:     serverId,
		userKeysChan: make([]chan *comet.PushUserKeysReq, routineNum),
		userRoomChan: make([]chan *comet.PushUserRoomReq, routineNum),
		userAllChan:  make(chan *comet.PushUserAllReq, routineNum),
		routineNum:   uint64(routineNum),
	}

	// create rpc client
	var err error
	grpcAddr := "127.0.0.1:12570"
	if cmt.rpcClient, err = newCometClient(grpcAddr); err != nil {
		return nil, err
	}

	// creat channel and routine
	chanNum := c.ChanNum
	cmt.ctx, cmt.cancel = context.WithCancel(context.Background())
	for i := 0; i < routineNum; i++ {
		cmt.userKeysChan[i] = make(chan *comet.PushUserKeysReq, chanNum)
		cmt.userRoomChan[i] = make(chan *comet.PushUserRoomReq, chanNum)
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

func newCometClient(addr string) (comet.CometClient, error) {
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
	return comet.NewCometClient(conn), err
}
