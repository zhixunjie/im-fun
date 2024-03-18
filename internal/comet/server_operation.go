package comet

import (
	"context"
	"github.com/zhixunjie/im-fun/internal/comet/channel"
	"github.com/zhixunjie/im-fun/pkg/gen_id"
	"github.com/zhixunjie/im-fun/pkg/logging"
	"time"

	"github.com/zhixunjie/im-fun/api/pb"
	"github.com/zhixunjie/im-fun/api/protocol"
)

func (s *Server) Connect(ctx context.Context, params *channel.AuthParams) (heartbeat time.Duration, err error) {
	userInfo := params.UserInfo
	reply, err := s.rpcToLogic.Connect(ctx, &pb.ConnectReq{
		Comm: &pb.ConnectCommon{
			ServerId:     s.serverId,
			UserId:       userInfo.TcpSessionId.UserId,
			TcpSessionId: userInfo.TcpSessionId.ToString(),
		},
		RoomId:   userInfo.RoomId,
		Token:    params.Token,
		Platform: userInfo.Platform,
	})
	if err != nil {
		return
	}
	return time.Duration(reply.Heartbeat), nil
}

func (s *Server) Disconnect(ctx context.Context, ch *channel.Channel) (err error) {
	_, err = s.rpcToLogic.Disconnect(ctx, &pb.DisconnectReq{
		Comm: &pb.ConnectCommon{
			ServerId:     s.serverId,
			UserId:       ch.UserInfo.TcpSessionId.UserId,
			TcpSessionId: ch.UserInfo.TcpSessionId.ToString(),
		},
	})
	return
}

func (s *Server) Heartbeat(ctx context.Context, userInfo *channel.UserInfo) (err error) {
	_, err = s.rpcToLogic.Heartbeat(ctx, &pb.HeartbeatReq{
		Comm: &pb.ConnectCommon{
			ServerId:     s.serverId,
			UserId:       userInfo.TcpSessionId.UserId,
			TcpSessionId: userInfo.TcpSessionId.ToString(),
		},
	})
	return
}

// RenewOnline renew room online.
//func (s *Server) RenewOnline(ctx context.Context, serverID string, roomCount map[string]int32) (allRoom map[string]int32, err error) {
//	reply, err := s.rpcToLogic.RenewOnline(ctx, &logic.OnlineReq{
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
	_, err = s.rpcToLogic.Receive(ctx, &pb.ReceiveReq{UserId: ch.UserInfo.TcpSessionId.UserId, Proto: p})
	return
}

func (s *Server) Operate(ctx context.Context, logHead string, proto *protocol.Proto, ch *channel.Channel, bucket *Bucket) error {
	logHead = logHead + "Operate|"

	switch protocol.Operation(proto.Op) {
	case protocol.OpHeartbeat:
		// 1. 客户端-心跳上报
		proto.Op = int32(protocol.OpHeartbeatReply)
		proto.Ver = protocol.ProtoVersion
		proto.Seq = int32(gen_id.SeqId())
		proto.Body = nil
		//logging.Infof(logHead + "OpHeartbeat generate")
		//timerPool.Set(trd, hb)
		//if now := time.Now(); now.Sub(lastHb) > hbTime {
		//	if err1 := s.Heartbeat(ctx, ch.UserInfo); err1 == nil {
		//		lastHb = now
		//	}
		//}
	case protocol.OpChangeRoom:
		// 2. 客户端房间切换
		if err := bucket.ChangeRoom(string(proto.Body), ch); err != nil {
			logging.Errorf(logHead+"bucket.ChangeRoom(%s) error(%v)", proto.Body, err)
		}
		proto.Op = int32(protocol.OpChangeRoomReply)
	case protocol.OpSub:
		// 客户端-添加订阅消息
		// TBD
	case protocol.OpUnsub:
		// 客户端-取消订阅消息
		// TBD
	default: // 客户端-收到其他消息（直接转到logic进行处理）
		// TBD
		if err := s.Receive(ctx, ch, proto); err != nil {
			logging.Errorf(logHead+"UserInfo=%+v,op=%v,err=%v", ch.UserInfo, proto.Op, err)
		}
		proto.Body = nil
	}
	return nil
}
