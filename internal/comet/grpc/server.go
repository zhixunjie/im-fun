package grpc

import (
	"context"
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
		Time:                  time.Duration(conf.KeepAliveInterval),
		Timeout:               time.Duration(conf.KeepAliveTimeout),
		MaxConnectionAge:      time.Duration(conf.MaxLifeTime),
	}))
	pb.RegisterCometServer(srv, &server{srv: s, UnimplementedCometServer: pb.UnimplementedCometServer{}})
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
	srv *comet.Server
	pb.UnimplementedCometServer
}

var _ pb.CometServer = &server{}

// PushMsg 发送消息到 UserKey 数组
func (s *server) PushMsg(ctx context.Context, req *pb.PushMsgReq) (reply *pb.PushMsgReply, err error) {
	if len(req.UserKeys) == 0 || req.Proto == nil {
		return nil, errors.ErrPushMsgArg
	}
	for _, key := range req.UserKeys {
		bucket := s.srv.AllocBucket(key)
		if channel := bucket.GetChannelByUserKey(key); channel != nil {
			if err = channel.Push(req.Proto); err != nil {
				return
			}
		}
	}
	return &pb.PushMsgReply{}, nil
}

// Broadcast 广播消息到所有的用户（所有bucket的所有channel）
func (s *server) Broadcast(ctx context.Context, req *pb.BroadcastReq) (*pb.BroadcastReply, error) {
	if req.Proto == nil {
		return nil, errors.ErrBroadCastArg
	}
	comet.BroadcastToAllBucket(s.srv, req.GetProto(), int(req.Speed))
	return &pb.BroadcastReply{}, nil
}

// BroadcastRoom 发送消息到指定房间
func (s *server) BroadcastRoom(ctx context.Context, req *pb.BroadcastRoomReq) (*pb.BroadcastRoomReply, error) {
	if req.Proto == nil || req.RoomId == "" {
		return nil, errors.ErrBroadCastRoomArg
	}
	// 同一个房间ID，可能存在于多个Bucket中
	for _, bucket := range s.srv.Buckets() {
		bucket.SendToTheRoom(req)
	}
	return &pb.BroadcastRoomReply{}, nil
}

// Rooms 获取所有在线人数大于0的房间
func (s *server) Rooms(ctx context.Context, req *pb.RoomsReq) (*pb.RoomsReply, error) {
	var (
		roomIds = make(map[string]bool)
	)
	for _, bucket := range s.srv.Buckets() {
		for roomID := range bucket.GetRoomsOnline() {
			roomIds[roomID] = true
		}
	}
	return &pb.RoomsReply{Rooms: roomIds}, nil
}
