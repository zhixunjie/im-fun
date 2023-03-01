package comet

import (
	"context"
	"github.com/zhixunjie/im-fun/internal/comet/channel"
	"time"

	"github.com/golang/glog"
	"github.com/zhixunjie/im-fun/api/logic"
	"github.com/zhixunjie/im-fun/api/protocol"
)

func (s *Server) Connect(ctx context.Context, ch *channel.Channel, token string) (heartbeat time.Duration, err error) {
	reply, err := s.rpcClient.Connect(ctx, &logic.ConnectReq{
		ServerId: s.serverId,
		UserId:   ch.UserInfo.UserId,
		UserKey:  ch.UserInfo.UserKey,
		RoomId:   ch.UserInfo.RoomId,
		Platform: ch.UserInfo.Platform,
		Token:    token,
	})
	if err != nil {
		return
	}
	return time.Duration(reply.Heartbeat), nil
}

func (s *Server) Disconnect(ctx context.Context, ch *channel.Channel) (err error) {
	_, err = s.rpcClient.Disconnect(ctx, &logic.DisconnectReq{
		ServerId: s.serverId,
		UserId:   ch.UserInfo.UserId,
		UserKey:  ch.UserInfo.UserKey,
	})
	return
}

func (s *Server) Heartbeat(ctx context.Context, userInfo *channel.UserInfo) (err error) {
	_, err = s.rpcClient.Heartbeat(ctx, &logic.HeartbeatReq{
		UserId:  userInfo.UserId,
		UserKey: userInfo.UserKey,
	})
	return
}

// RenewOnline renew room online.
//func (s *Server) RenewOnline(ctx context.Context, serverID string, roomCount map[string]int32) (allRoom map[string]int32, err error) {
//	reply, err := s.rpcClient.RenewOnline(ctx, &logic.OnlineReq{
//		Server:    s.serverID,
//		RoomCount: roomCount,
//	}, grpc.UseCompressor(gzip.Name))
//	if err != nil {
//		return
//	}
//	return reply.AllRoomCount, nil
//}

// Receive receive a message.
func (s *Server) Receive(ctx context.Context, ch *channel.Channel, p *protocol.Proto) (err error) {
	_, err = s.rpcClient.Receive(ctx, &logic.ReceiveReq{UserId: ch.UserInfo.UserId, Proto: p})
	return
}

func (s *Server) Operate(ctx context.Context, p *protocol.Proto, ch *channel.Channel, b *Bucket) error {
	switch protocol.Operation(p.Op) {
	case protocol.OpChangeRoom:
		if err := b.ChangeRoom(string(p.Body), ch); err != nil {
			glog.Errorf("b.ChangeRoom(%s) error(%v)", p.Body, err)
		}
		p.Op = int32(protocol.OpChangeRoomReply)
	case protocol.OpSub:
		// TBD
	case protocol.OpUnsub:
		// TBD
	default:
		// TBD
		if err := s.Receive(ctx, ch, p); err != nil {
			glog.Errorf("UserInfo=%+v,op=%v,err=%v", ch.UserInfo, p.Op, err)
		}
		p.Body = nil
	}
	return nil
}
