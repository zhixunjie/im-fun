package comet

import (
	"context"
	"fmt"
	"github.com/zhixunjie/im-fun/internal/comet/channel"
	"github.com/zhixunjie/im-fun/pkg/gmodel"
	"github.com/zhixunjie/im-fun/pkg/logging"
	"time"

	"github.com/zhixunjie/im-fun/api/pb"
	"github.com/zhixunjie/im-fun/api/protocol"
)

// handleClientMsg 专门处理客户端上行的TCP消息
func (s *TcpServer) handleClientMsg(ctx context.Context, logHead string, proto *protocol.Proto, ch *channel.Channel, bucket *Bucket) (err error) {
	logHead += "handleClientMsg|"

	switch protocol.Operation(proto.Op) {
	case protocol.OpHeartbeatReq: // 心跳上报
		proto.Op = int32(protocol.OpHeartbeatResp)
		proto.Ver = protocol.ProtoVersion
		proto.Seq = gmodel.NewSeqId32()
		proto.Body = nil
		logging.Infof(logHead + "Heartbeat Proto is generated")
		// reset timer
		s.ResetTimerHeartbeat(ctx, logHead, ch)
		// rpc: lease
		// 节流: 即使客户端上报心跳过来，也不一定要调用RPC接口进行续约
		if now := time.Now(); now.Sub(ch.LastHb) > ch.HbInterval {
			tErr := s.Heartbeat(ctx, ch)
			if tErr != nil {
				logging.Errorf(logHead+"Heartbeat lease fail,err=%v", tErr)
				return
			}
			ch.LastHb = now
			logging.Infof(logHead + "Heartbeat lease success")
		}
	case protocol.OpChangeRoomReq: // 房间切换
		err = bucket.ChangeRoom(string(proto.Body), ch)
		if err != nil {
			logging.Errorf(logHead+"bucket.ChangeRoom(%s) error(%v)", proto.Body, err)
			return
		}
		proto.Op = int32(protocol.OpChangeRoomResp)
	case protocol.OpSubReq: // 订阅
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

func (s *TcpServer) Connect(ctx context.Context, params *pb.AuthParams) (resp *pb.ConnectResp, err error) {
	if params == nil {
		err = fmt.Errorf("params is nil")
		return
	}
	resp, err = s.rpcToLogic.Connect(ctx, &pb.ConnectReq{
		AuthParams: params,
		ServerId:   s.serverId,
	})
	if err != nil {
		err = fmt.Errorf("RPC Connect,err=%w", err)
		return
	}
	if resp.SessionId == "" {
		err = fmt.Errorf("sessionId is empty")
		return
	}
	if resp.HbCfg == nil {
		err = fmt.Errorf("hbCfg is empty")
		return
	}
	logging.Infof("RPC Connect success,resp=%v", resp)

	return
}

func (s *TcpServer) Disconnect(ctx context.Context, ch *channel.Channel) (err error) {
	if ch.UserInfo == nil {
		err = fmt.Errorf("userInfo is empty")
		return
	}
	_, err = s.rpcToLogic.Disconnect(ctx, &pb.DisconnectReq{Connect: ch.UserInfo.Connect})
	if err != nil {
		err = fmt.Errorf("RPC Disconnect,err=%v", err)
		return
	}
	logging.Infof("RPC Disconnect success")

	return
}

func (s *TcpServer) Heartbeat(ctx context.Context, ch *channel.Channel) (err error) {
	if ch.UserInfo == nil {
		err = fmt.Errorf("userInfo is empty")
		return
	}
	_, err = s.rpcToLogic.Heartbeat(ctx, &pb.HeartbeatReq{
		Connect:    ch.UserInfo.Connect,
		BindExpire: ch.UserInfo.HbCfg.BindExpire,
	})
	if err != nil {
		err = fmt.Errorf("RPC Heartbeat,err=%v", err)
		return
	}
	logging.Infof("RPC Heartbeat success")
	return
}

// RenewOnline renew room online.
//func (s *TcpServer) RenewOnline(ctx context.Context, serverID string, roomCount map[string]int32) (allRoom map[string]int32, err error) {
//	resp, err := s.rpcToLogic.RenewOnline(ctx, &logic.OnlineReq{
//		TcpServer:    s.serverID,
//		RoomCount: roomCount,
//	}, grpc.UseCompressor(gzip.Name))
//	if err != nil {
//		return
//	}
//	return resp.AllRoomCount, nil
//}

// Receive 接收到消息
func (s *TcpServer) Receive(ctx context.Context, ch *channel.Channel, p *protocol.Proto) (err error) {
	_, err = s.rpcToLogic.Receive(ctx, &pb.ReceiveReq{
		UniId: ch.UserInfo.Connect.UniId,
		Proto: p,
	})
	logging.Infof("RPC Receive,err=%v", err)

	return
}
