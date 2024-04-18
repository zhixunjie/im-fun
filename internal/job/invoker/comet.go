package invoker

import (
	"context"
	pb "github.com/zhixunjie/im-fun/api/pb"
	"github.com/zhixunjie/im-fun/internal/job/conf"
	"go.uber.org/atomic"
)

// Logic -> Job -> CometInvoker -> RPC To Comet -> Comet

// CometInvoker 调用comet的服务
type CometInvoker struct {
	ctx    context.Context
	Cancel context.CancelFunc

	serverId  string
	rpcClient pb.CometClient

	// counter
	counterToUser atomic.Uint64
	counterToRoom atomic.Uint64
	RoutineNum    uint64

	// send msg
	chUsers []chan *pb.SendToUsersReq // send msg: to some user
	chRoom  []chan *pb.SendToRoomReq  // send msg: to the room
	chAll   chan *pb.SendToAllReq     // send msg: to all user
}

func NewCometInvoker(serverId, grpcAddr string, conf *conf.CometInvoker) (*CometInvoker, error) {
	routineNum := conf.RoutineNum
	cmt := &CometInvoker{
		serverId:   serverId,
		chUsers:    make([]chan *pb.SendToUsersReq, routineNum),
		chRoom:     make([]chan *pb.SendToRoomReq, routineNum),
		chAll:      make(chan *pb.SendToAllReq, routineNum),
		RoutineNum: uint64(routineNum),
	}

	// create rpc client
	var err error
	if cmt.rpcClient, err = newCometClient(grpcAddr); err != nil {
		return nil, err
	}

	// creat channel and routine
	cmt.ctx, cmt.Cancel = context.WithCancel(context.Background())
	for i := 0; i < int(cmt.RoutineNum); i++ {
		cmt.chUsers[i] = make(chan *pb.SendToUsersReq, conf.ChanBufferSize)
		cmt.chRoom[i] = make(chan *pb.SendToRoomReq, conf.ChanBufferSize)
		go cmt.Run(i)
	}
	return cmt, nil
}
