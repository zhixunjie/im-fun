package grpc

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/zhixunjie/im-fun/internal/comet"
	"net"
	"time"

	pb "github.com/zhixunjie/im-fun/api/comet"
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
	logrus.Infof("GRPC server is listening %v：%v\n", conf.Network, conf.Addr)
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

// PushUserKeys 发送消息到 UserKey 数组
func (s *server) PushUserKeys(ctx context.Context, req *pb.PushUserKeysReq) (reply *pb.PushUserKeysReply, err error) {
	if len(req.UserKeys) == 0 || req.Proto == nil {
		return nil, errors.ErrParamsNotAllow
	}
	for _, key := range req.UserKeys {
		bucket := s.srv.AllocBucket(key)
		if channel := bucket.GetChannelByUserKey(key); channel != nil {
			if err = channel.Push(req.Proto); err != nil {
				return
			}
		}
	}
	return &pb.PushUserKeysReply{}, nil
}

// PushUserAll 广播消息到所有的用户（所有bucket的所有channel）
func (s *server) PushUserAll(ctx context.Context, req *pb.PushUserAllReq) (*pb.PushUserAllReply, error) {
	if req.Proto == nil {
		return nil, errors.ErrParamsNotAllow
	}
	comet.BroadcastToAllBucket(s.srv, req.GetProto(), int(req.Speed))
	return &pb.PushUserAllReply{}, nil
}

// PushUserRoom 发送消息到指定房间
func (s *server) PushUserRoom(ctx context.Context, req *pb.PushUserRoomReq) (*pb.PushUserRoomReply, error) {
	if req.Proto == nil || req.RoomId == "" {
		return nil, errors.ErrParamsNotAllow
	}
	// 同一个房间ID，可能存在于多个Bucket中
	for _, bucket := range s.srv.Buckets() {
		bucket.BroadcastRoom(req)
	}
	return &pb.PushUserRoomReply{}, nil
}

// GetAllRoomId 获取所有在线人数大于0的房间
func (s *server) GetAllRoomId(ctx context.Context, req *pb.GetAllRoomIdReq) (*pb.GetAllRoomIdReply, error) {
	var (
		roomIds = make(map[string]bool)
	)
	for _, bucket := range s.srv.Buckets() {
		for roomID := range bucket.GetRoomsOnline() {
			roomIds[roomID] = true
		}
	}
	return &pb.GetAllRoomIdReply{Rooms: roomIds}, nil
}
