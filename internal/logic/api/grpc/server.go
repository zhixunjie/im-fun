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

func (s *server) Connect(ctx context.Context, req *pb.ConnectReq) (reply *pb.ConnectReply, err error) {
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
	reply, err = s.bz.Connect(ctx, req)
	if err != nil {
		return
	}

	return
}

func (s *server) Disconnect(ctx context.Context, req *pb.DisconnectReq) (reply *pb.DisconnectReply, err error) {
	// invoke svc
	reply, err = s.bz.Disconnect(ctx, req)
	if err != nil {
		return
	}

	return
}

func (s *server) Heartbeat(ctx context.Context, req *pb.HeartbeatReq) (reply *pb.HeartbeatReply, err error) {
	// invoke svc
	reply, err = s.bz.Heartbeat(ctx, req)
	if err != nil {
		return
	}

	return
}

func (s *server) RenewOnline(ctx context.Context, req *pb.OnlineReq) (reply *pb.OnlineReply, err error) {
	// invoke svc
	reply, err = s.bz.RenewOnline(ctx, req.ServerId, req.RoomCount)
	if err != nil {
		return
	}

	return
}

func (s *server) Receive(ctx context.Context, req *pb.ReceiveReq) (reply *pb.ReceiveReply, err error) {
	reply = new(pb.ReceiveReply)

	//if err := s.svc.Receive(ctx, req.Mid, req.Proto); err != nil {
	//	return &pb.ReceiveReply{}, err
	//}
	return
}

func (s *server) Nodes(ctx context.Context, req *pb.NodesReq) (reply *pb.NodesReply, err error) {
	reply = new(pb.NodesReply)

	// invoke svc
	reply, err = s.bz.Nodes(ctx, req)
	if err != nil {
		return
	}
	//return s.svc.NodesWeighted(ctx, req.Platform, req.ClientIP), nil
	return
}
