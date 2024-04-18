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

// OpFromClient 专门处理客户端上行的TCP消息
func (s *TcpServer) OpFromClient(ctx context.Context, logHead string, proto *protocol.Proto, ch *channel.Channel, bucket *Bucket) (err error) {
	logHead += "OpFromClient|"

	switch protocol.Operation(proto.Op) {
	case protocol.OpHeartbeat: // 心跳上报
		proto.Op = int32(protocol.OpHeartbeatReply)
		proto.Ver = protocol.ProtoVersion
		proto.Seq = gen_id.SeqId32()
		proto.Body = nil
		logging.Infof(logHead + "Heartbeat Proto is generated")
		// reset timer
		s.ResetTimerHeartbeat(ctx, logHead, ch)
		// rpc: lease
		// 节流: 即使客户端上报心跳过来，也不一定要调用RPC接口进行续约
		if now := time.Now(); now.Sub(ch.LastHb) > ch.HbLeaseDuration {
			tErr := s.Heartbeat(ctx, ch.UserInfo)
			if tErr != nil {
				logging.Errorf(logHead+"Heartbeat lease fail,err=%v", tErr)
				return
			}
			ch.LastHb = now
			logging.Infof(logHead + "Heartbeat lease success")
		}
	case protocol.OpChangeRoom: // 房间切换
		err = bucket.ChangeRoom(string(proto.Body), ch)
		if err != nil {
			logging.Errorf(logHead+"bucket.ChangeRoom(%s) error(%v)", proto.Body, err)
			return
		}
		proto.Op = int32(protocol.OpChangeRoomReply)
	case protocol.OpSub: // 订阅
		// TBD
	case protocol.OpUnsub: // 取消订阅
		// TBD
	default: // 其他类型的消息（直接转到logic进行处理）
		// TBD
		err = s.Receive(ctx, ch, proto)
		if err != nil {
			logging.Errorf(logHead+"UserInfo=%+v,op=%v,err=%v", ch.UserInfo, proto.Op, err)
			return
		}
		//proto.Body = nil
	}
	return
}

func (s *TcpServer) Connect(ctx context.Context, params *channel.AuthParams) (heartbeat time.Duration, err error) {
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
	heartbeat = time.Duration(reply.Heartbeat)
	logging.Infof("RPC Connect,err=%v", err)

	return
}

func (s *TcpServer) Disconnect(ctx context.Context, ch *channel.Channel) (err error) {
	_, err = s.rpcToLogic.Disconnect(ctx, &pb.DisconnectReq{
		Comm: &pb.ConnectCommon{
			ServerId:     s.serverId,
			UserId:       ch.UserInfo.TcpSessionId.UserId,
			TcpSessionId: ch.UserInfo.TcpSessionId.ToString(),
		},
	})
	logging.Infof("RPC Disconnect,err=%v", err)

	return
}

func (s *TcpServer) Heartbeat(ctx context.Context, userInfo *channel.UserInfo) (err error) {
	_, err = s.rpcToLogic.Heartbeat(ctx, &pb.HeartbeatReq{
		Comm: &pb.ConnectCommon{
			ServerId:     s.serverId,
			UserId:       userInfo.TcpSessionId.UserId,
			TcpSessionId: userInfo.TcpSessionId.ToString(),
		},
	})
	logging.Infof("RPC Heartbeat,err=%v", err)
	return
}

// RenewOnline renew room online.
//func (s *TcpServer) RenewOnline(ctx context.Context, serverID string, roomCount map[string]int32) (allRoom map[string]int32, err error) {
//	reply, err := s.rpcToLogic.RenewOnline(ctx, &logic.OnlineReq{
//		TcpServer:    s.serverID,
//		RoomCount: roomCount,
//	}, grpc.UseCompressor(gzip.Name))
//	if err != nil {
//		return
//	}
//	return reply.AllRoomCount, nil
//}

// Receive receive a message.
func (s *TcpServer) Receive(ctx context.Context, ch *channel.Channel, p *protocol.Proto) (err error) {
	_, err = s.rpcToLogic.Receive(ctx, &pb.ReceiveReq{UserId: ch.UserInfo.TcpSessionId.UserId, Proto: p})
	logging.Infof("RPC Receive,err=%v", err)

	return
}
