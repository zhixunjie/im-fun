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

func (s *server) Connect(ctx context.Context, req *pb.ConnectReq) (*pb.ConnectReply, error) {
	if req.UserId == 0 {
		logging.Errorf("UserId not allow,token=%+v", req.GetToken())
		return &pb.ConnectReply{}, errors.New("req.UserId not allow")
	}
	if req.UserKey == "" {
		logging.Errorf("UserKey not allow,token=%+v", req.GetToken())
		return &pb.ConnectReply{}, errors.New("req.UserKey not allow")
	}
	if req.Token == "" {
		logging.Errorf("Token not allow,token=%+v", req.GetToken())
		return &pb.ConnectReply{}, errors.New("req.Token not allow")
	}

	hb, err := s.bz.Connect(ctx, req)
	if err != nil {
		return &pb.ConnectReply{}, err
	}
	return &pb.ConnectReply{
		Heartbeat: hb,
	}, nil
}

func (s *server) Disconnect(ctx context.Context, req *pb.DisconnectReq) (*pb.DisconnectReply, error) {
	has, err := s.bz.Disconnect(ctx, req)
	if err != nil {
		return &pb.DisconnectReply{}, err
	}
	return &pb.DisconnectReply{
		Has: has,
	}, nil
}

func (s *server) Heartbeat(ctx context.Context, req *pb.HeartbeatReq) (*pb.HeartbeatReply, error) {
	//if err := s.svc.Heartbeat(ctx, req.Mid, req.Key, req.Server); err != nil {
	//	return &pb.HeartbeatReply{}, err
	//}
	return &pb.HeartbeatReply{}, nil
}

func (s *server) RenewOnline(ctx context.Context, req *pb.OnlineReq) (*pb.OnlineReply, error) {
	allRoomCount, err := s.bz.RenewOnline(ctx, req.ServerId, req.RoomCount)
	if err != nil {
		return &pb.OnlineReply{}, err
	}
	return &pb.OnlineReply{AllRoomCount: allRoomCount}, nil
}

func (s *server) Receive(ctx context.Context, req *pb.ReceiveReq) (*pb.ReceiveReply, error) {
	//if err := s.svc.Receive(ctx, req.Mid, req.Proto); err != nil {
	//	return &pb.ReceiveReply{}, err
	//}
	return &pb.ReceiveReply{}, nil
}

func (s *server) Nodes(ctx context.Context, req *pb.NodesReq) (*pb.NodesReply, error) {
	//return s.svc.NodesWeighted(ctx, req.Platform, req.ClientIP), nil
	return &pb.NodesReply{}, nil
}
