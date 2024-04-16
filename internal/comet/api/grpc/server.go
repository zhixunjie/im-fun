package grpc

import (
	"context"
	pb "github.com/zhixunjie/im-fun/api/pb"
	"github.com/zhixunjie/im-fun/internal/comet"
	"github.com/zhixunjie/im-fun/internal/comet/api"
)

type server struct {
	srv *comet.Server
	pb.UnimplementedCometServer
}

var _ pb.CometServer = &server{}

func (s *server) SendToUsers(ctx context.Context, req *pb.SendToUsersReq) (reply *pb.SendToUsersReply, err error) {
	reply = new(pb.SendToUsersReply)
	if len(req.TcpSessionIds) == 0 || req.Proto == nil {
		err = api.ErrParamsNotAllow
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
	reply = new(pb.SendToRoomReply)
	if req.Proto == nil || req.RoomId == "" {
		err = api.ErrParamsNotAllow
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
	reply = new(pb.SendToAllReply)
	if req.Proto == nil {
		err = api.ErrParamsNotAllow
		return
	}
	comet.BroadcastToAllBucket(s.srv, req.GetProto(), int(req.Speed))
	return
}

// GetAllRoomId 获取所有在线人数大于0的房间
func (s *server) GetAllRoomId(ctx context.Context, req *pb.GetAllRoomIdReq) (reply *pb.GetAllRoomIdReply, err error) {
	reply = new(pb.GetAllRoomIdReply)
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
