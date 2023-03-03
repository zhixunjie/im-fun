package comet

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/zhixunjie/im-fun/internal/comet/channel"
	"time"

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

func (s *Server) Operate(ctx context.Context, proto *protocol.Proto, ch *channel.Channel, bucket *Bucket) error {
	switch protocol.Operation(proto.Op) {
	case protocol.OpHeartbeat: // 客户端-心跳上报
		proto.Op = int32(protocol.OpHeartbeatReply)
		proto.Body = nil
		//timerPool.Set(trd, hb)
		//if now := time.Now(); now.Sub(lastHb) > hbTime {
		//	if err1 := s.Heartbeat(ctx, ch.UserInfo); err1 == nil {
		//		lastHb = now
		//	}
		//}
	case protocol.OpChangeRoom: // 客户端房间切换
		if err := bucket.ChangeRoom(string(proto.Body), ch); err != nil {
			logrus.Errorf("bucket.ChangeRoom(%s) error(%v)", proto.Body, err)
		}
		proto.Op = int32(protocol.OpChangeRoomReply)
	case protocol.OpSub: // 客户端-添加订阅消息
		// TBD
	case protocol.OpUnsub: // 客户端-取消订阅消息
		// TBD
	default: // 客户端-收到其他消息（直接转到logic进行处理）
		// TBD
		if err := s.Receive(ctx, ch, proto); err != nil {
			logrus.Errorf("UserInfo=%+v,op=%v,err=%v", ch.UserInfo, proto.Op, err)
		}
		proto.Body = nil
	}
	return nil
}
