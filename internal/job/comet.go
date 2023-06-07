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
	ctx    context.Context
	cancel context.CancelFunc

	serverId  string
	rpcClient comet.CometClient

	pushChanNum uint64
	roomChanNum uint64
	routineNum  uint64

	// send msg
	chUserKeys []chan *comet.SendToUserKeysReq // send msg to some user
	chRoom     []chan *comet.SendToRoomReq     // send msg to the room
	chAll      chan *comet.SendToAllReq        // send msg to all user
}

func NewComet(serverId string, c *conf.Comet) (*Comet, error) {
	routineNum := c.RoutineNum
	cmt := &Comet{
		serverId:   serverId,
		chUserKeys: make([]chan *comet.SendToUserKeysReq, routineNum),
		chRoom:     make([]chan *comet.SendToRoomReq, routineNum),
		chAll:      make(chan *comet.SendToAllReq, routineNum),
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
		cmt.chUserKeys[i] = make(chan *comet.SendToUserKeysReq, c.ChanNum)
		cmt.chRoom[i] = make(chan *comet.SendToRoomReq, c.ChanNum)
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
