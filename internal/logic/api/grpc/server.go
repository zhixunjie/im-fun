package grpc

import (
	"context"
	"errors"
	pb "github.com/zhixunjie/im-fun/api/pb"
	"github.com/zhixunjie/im-fun/internal/logic/biz"
	"github.com/zhixunjie/im-fun/pkg/logging"

	// use gzip decoder
	_ "google.golang.org/grpc/encoding/gzip"
)

type server struct {
	pb.UnimplementedLogicServer
	bz        *biz.Biz
	bzContact *biz.ContactUseCase
	bzMessage *biz.MessageUseCase
}

var _ pb.LogicServer = &server{}

func (s *server) Connect(ctx context.Context, req *pb.ConnectReq) (resp *pb.ConnectResp, err error) {
	defer func() {
		if err != nil {
			logging.Errorf("err=%v,req=%+v", err, req)
			return
		}
	}()
	if req.Comm.UserId == 0 {
		err = errors.New("UserId not allow")
		return
	}
	if req.Comm.TcpSessionId == "" {
		err = errors.New("TcpSessionId not allow")
		return
	}
	if req.Token == "" {
		err = errors.New("token not allow")
		return
	}

	// invoke svc
	resp, err = s.bz.Connect(ctx, req)
	if err != nil {
		return
	}

	return
}

func (s *server) Disconnect(ctx context.Context, req *pb.DisconnectReq) (resp *pb.DisconnectResp, err error) {
	// invoke svc
	resp, err = s.bz.Disconnect(ctx, req)
	if err != nil {
		return
	}

	return
}

func (s *server) Heartbeat(ctx context.Context, req *pb.HeartbeatReq) (resp *pb.HeartbeatResp, err error) {
	// invoke svc
	resp, err = s.bz.Heartbeat(ctx, req)
	if err != nil {
		return
	}

	return
}

func (s *server) RenewOnline(ctx context.Context, req *pb.OnlineReq) (resp *pb.OnlineResp, err error) {
	// invoke svc
	resp, err = s.bz.RenewOnline(ctx, req.ServerId, req.RoomCount)
	if err != nil {
		return
	}

	return
}

func (s *server) Receive(ctx context.Context, req *pb.ReceiveReq) (resp *pb.ReceiveResp, err error) {
	resp = new(pb.ReceiveResp)

	//if err := s.svc.Receive(ctx, req.Mid, req.Proto); err != nil {
	//	return &pb.ReceiveResp{}, err
	//}
	return
}

func (s *server) Nodes(ctx context.Context, req *pb.NodesReq) (resp *pb.NodesResp, err error) {
	resp = new(pb.NodesResp)

	// invoke svc
	resp, err = s.bz.Nodes(ctx, req)
	if err != nil {
		return
	}
	//return s.svc.NodesWeighted(ctx, req.Platform, req.ClientIP), nil
	return
}
