package apigrpc

import (
	"context"
	"github.com/zhixunjie/im-fun/internal/logic/service"
	"net"
	"time"

	pb "github.com/zhixunjie/im-fun/api/logic"
	"github.com/zhixunjie/im-fun/internal/logic/conf"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	// use gzip decoder
	_ "google.golang.org/grpc/encoding/gzip"
)

// New logic grpc server
func New(conf *conf.RPCServer, svc *service.Service) *grpc.Server {
	params := grpc.KeepaliveParams(keepalive.ServerParameters{
		MaxConnectionIdle:     time.Duration(conf.IdleTimeout),
		MaxConnectionAgeGrace: time.Duration(conf.ForceCloseWait),
		Time:                  time.Duration(conf.KeepAliveInterval),
		Timeout:               time.Duration(conf.KeepAliveTimeout),
		MaxConnectionAge:      time.Duration(conf.MaxLifeTime),
	})
	srv := grpc.NewServer(params)
	pb.RegisterLogicServer(srv, &server{svc: svc})
	listener, err := net.Listen(conf.Network, conf.Addr)
	if err != nil {
		panic(err)
	}
	// begin to listen
	go func() {
		if err = srv.Serve(listener); err != nil {
			panic(err)
		}
	}()
	return srv
}

type server struct {
	pb.UnimplementedLogicServer
	svc *service.Service
}

var _ pb.LogicServer = &server{}

// Connect connect a conn.
func (s *server) Connect(ctx context.Context, req *pb.ConnectReq) (*pb.ConnectReply, error) {
	//mid, key, room, accepts, hb, err := s.svc.Connect(ctx, req.Server, req.Cookie, req.Token)
	//if err != nil {
	//	return &pb.ConnectReply{}, err
	//}
	return &pb.ConnectReply{}, nil
}

// Disconnect disconnect a conn.
func (s *server) Disconnect(ctx context.Context, req *pb.DisconnectReq) (*pb.DisconnectReply, error) {
	//has, err := s.svc.Disconnect(ctx, req.Mid, req.Key, req.Server)
	//if err != nil {
	//	return &pb.DisconnectReply{}, err
	//}
	return &pb.DisconnectReply{}, nil
}

// Heartbeat beartbeat a conn.
func (s *server) Heartbeat(ctx context.Context, req *pb.HeartbeatReq) (*pb.HeartbeatReply, error) {
	//if err := s.svc.Heartbeat(ctx, req.Mid, req.Key, req.Server); err != nil {
	//	return &pb.HeartbeatReply{}, err
	//}
	return &pb.HeartbeatReply{}, nil
}

// RenewOnline renew server online.
func (s *server) RenewOnline(ctx context.Context, req *pb.OnlineReq) (*pb.OnlineReply, error) {
	allRoomCount, err := s.svc.RenewOnline(ctx, req.ServerId, req.RoomCount)
	if err != nil {
		return &pb.OnlineReply{}, err
	}
	return &pb.OnlineReply{AllRoomCount: allRoomCount}, nil
}

// Receive receive a message.
func (s *server) Receive(ctx context.Context, req *pb.ReceiveReq) (*pb.ReceiveReply, error) {
	//if err := s.svc.Receive(ctx, req.Mid, req.Proto); err != nil {
	//	return &pb.ReceiveReply{}, err
	//}
	return &pb.ReceiveReply{}, nil
}

// nodes return nodes.
func (s *server) Nodes(ctx context.Context, req *pb.NodesReq) (*pb.NodesReply, error) {
	//return s.svc.NodesWeighted(ctx, req.Platform, req.ClientIP), nil
	return &pb.NodesReply{}, nil
}
