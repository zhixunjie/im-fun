package comet

import (
	"context"
	"github.com/zhixunjie/im-fun/internal/comet/channel"
	"time"

	"github.com/golang/glog"
	"github.com/zhixunjie/im-fun/api/logic"
	"github.com/zhixunjie/im-fun/api/protocol"
)

// Connect connected a connection.
func (s *Server) Connect(c context.Context, p *protocol.Proto, cookie string) (mid int64, key, rid string, accepts []int32, heartbeat time.Duration, err error) {
	reply, err := s.rpcClient.Connect(c, &logic.ConnectReq{
		Server: s.serverID,
		Cookie: cookie,
		Token:  p.Body,
	})
	if err != nil {
		return
	}
	return reply.Mid, reply.Key, reply.RoomID, reply.Accepts, time.Duration(reply.Heartbeat), nil
}

// Disconnect disconnected a connection.
func (s *Server) Disconnect(c context.Context, mid int64, key string) (err error) {
	_, err = s.rpcClient.Disconnect(context.Background(), &logic.DisconnectReq{
		Server: s.serverID,
		Mid:    mid,
		Key:    key,
	})
	return
}

// Heartbeat heartbeat a connection session.
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
func (s *Server) Receive(ctx context.Context, mid int64, p *protocol.Proto) (err error) {
	_, err = s.rpcClient.Receive(ctx, &logic.ReceiveReq{Mid: mid, Proto: p})
	return
}

// Operate operate.
func (s *Server) Operate(ctx context.Context, p *protocol.Proto, ch *Channel, b *Bucket) error {
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
		// TODO
		//if err := s.Receive(ctx, ch.Mid, p); err != nil {
		//	glog.Errorf("s.Report(%d) op:%d error(%v)", ch.Mid, p.Op, err)
		//}
		//p.Body = nil
	}
	return nil
}
