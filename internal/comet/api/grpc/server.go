package grpc

import (
	"context"
	"github.com/zhixunjie/im-fun/internal/comet"
	"github.com/zhixunjie/im-fun/pkg/logging"
	"net"
	"time"

	pb "github.com/zhixunjie/im-fun/api/pb"
	"github.com/zhixunjie/im-fun/internal/comet/conf"
	"github.com/zhixunjie/im-fun/internal/comet/errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

func New(s *comet.Server, conf *conf.RPCServer) *grpc.Server {
	srv := grpc.NewServer(grpc.KeepaliveParams(keepalive.ServerParameters{
		MaxConnectionIdle:     time.Duration(conf.IdleTimeout),
		MaxConnectionAgeGrace: time.Duration(conf.ForceCloseWait),
		Time:                  time.Duration(conf.KeepaliveInterval),
		Timeout:               time.Duration(conf.KeepaliveTimeout),
		MaxConnectionAge:      time.Duration(conf.MaxLifeTime),
	}))
	pb.RegisterCometServer(srv, &server{srv: s, UnimplementedCometServer: pb.UnimplementedCometServer{}})
	listener, err := net.Listen(conf.Network, conf.Addr)
	if err != nil {
		panic(err)
	}
	// begin to listen
	logging.Infof("GRPC server is listening %v：%v", conf.Network, conf.Addr)
	go func() {
		if err = srv.Serve(listener); err != nil {
			panic(err)
		}
	}()
	return srv
}

type server struct {
	srv *comet.Server
	pb.UnimplementedCometServer
}

var _ pb.CometServer = &server{}

func (s *server) SendToUsers(ctx context.Context, req *pb.SendToUsersReq) (reply *pb.SendToUsersReply, err error) {
	if len(req.TcpSessionIds) == 0 || req.Proto == nil {
		err = errors.ErrParamsNotAllow
		return
	}
	for _, key := range req.TcpSessionIds {
		bucket := s.srv.AllocBucket(key)
		if channel := bucket.GetChannel(key); channel != nil {
			if err = channel.Push(req.Proto); err != nil {
				return
			}
		}
	}
	return
}

// SendToRoom 发送消息到指定房间
func (s *server) SendToRoom(ctx context.Context, req *pb.SendToRoomReq) (reply *pb.SendToRoomReply, err error) {
	if req.Proto == nil || req.RoomId == "" {
		err = errors.ErrParamsNotAllow
		return
	}
	// 同一个房间ID，可能存在于多个Bucket中
	for _, bucket := range s.srv.Buckets() {
		bucket.BroadcastRoom(req)
	}
	return
}

// SendToAll 广播消息到所有的用户（所有bucket的所有channel）
func (s *server) SendToAll(ctx context.Context, req *pb.SendToAllReq) (reply *pb.SendToAllReply, err error) {
	if req.Proto == nil {
		err = errors.ErrParamsNotAllow
		return
	}
	comet.BroadcastToAllBucket(s.srv, req.GetProto(), int(req.Speed))
	return
}

// GetAllRoomId 获取所有在线人数大于0的房间
func (s *server) GetAllRoomId(ctx context.Context, req *pb.GetAllRoomIdReq) (reply *pb.GetAllRoomIdReply, err error) {
	var (
		roomIds = make(map[string]bool)
	)
	for _, bucket := range s.srv.Buckets() {
		for roomID := range bucket.GetRoomsOnline() {
			roomIds[roomID] = true
		}
	}
	reply.Rooms = roomIds
	return
}
