package biz

import (
	"context"
	"fmt"
	"github.com/zhixunjie/im-fun/api/pb"
	"github.com/zhixunjie/im-fun/api/protocol"
	"github.com/zhixunjie/im-fun/pkg/logging"
	"time"
)

//type TcpUseCase struct {
//	contactRepo *data.TcpRepo
//}
//
//func NewTcpUseCase(contactRepo *data.TcpRepo) *TcpUseCase {
//	return &TcpUseCase{contactRepo: contactRepo}
//}

// Connect connected a conn.
func (bz *Biz) Connect(ctx context.Context, req *pb.ConnectReq) (reply *pb.ConnectReply, err error) {
	reply = new(pb.ConnectReply)
	rr := req.Comm
	expire := bz.GetHeartbeatExpire()
	logHead := fmt.Sprintf("Connect,req.Comm=%v,expire=%v|", rr, expire)

	if err = bz.data.SessionBinding(ctx, logHead, rr, expire); err != nil {
		logging.Errorf(logHead+"SessionBinding fail,error=%v", err)
		return
	}

	// return hb
	hbConf := bz.conf.Node.Heartbeat
	interval := int64(hbConf.Interval)
	reply.HbCfg = &pb.HbCfg{
		Interval:  interval,
		FailCount: hbConf.FailCount,
	}
	logging.Infof(logHead + "success")

	return
}

// Disconnect disconnect a conn.
func (bz *Biz) Disconnect(ctx context.Context, req *pb.DisconnectReq) (reply *pb.DisconnectReply, err error) {
	reply = new(pb.DisconnectReply)
	rr := req.Comm
	logHead := fmt.Sprintf("Disconnect,req.Comm=%v|", rr)

	reply.Has, err = bz.data.SessionDel(ctx, logHead, rr)
	if err != nil {
		logging.Errorf(logHead+"SessionDel fail,error=%v", err, rr.UserId, rr.TcpSessionId)

		return
	}

	logging.Infof(logHead + "success")
	return
}

// Heartbeat heartbeat a conn.
func (bz *Biz) Heartbeat(ctx context.Context, req *pb.HeartbeatReq) (reply *pb.HeartbeatReply, err error) {
	reply = new(pb.HeartbeatReply)
	rr := req.Comm
	logHead := fmt.Sprintf("Heartbeat,req.Comm=%v|", rr)

	// 续约KEY
	expire := bz.GetHeartbeatExpire()
	reply.Has, err = bz.data.SessionLease(ctx, logHead, rr, expire)
	if err != nil {
		logging.Errorf(logHead+"SessionLease fail,error=%v", err)
		return
	}
	// 重新建立绑定关系
	if !reply.Has {
		if err = bz.data.SessionBinding(ctx, logHead, rr, expire); err != nil {
			logging.Errorf(logHead+"SessionBinding fail,error=%v", err)
			return
		}
	}
	logging.Infof(logHead + "success")

	return
}

// RenewOnline renew a server online.
func (bz *Biz) RenewOnline(ctx context.Context, serverId string, roomCount map[string]int32) (reply *pb.OnlineReply, err error) {
	reply = new(pb.OnlineReply)
	//online := &model.Online{
	//	Server:    serverId,
	//	RoomCount: roomCount,
	//	Updated:   time.Now().Unix(),
	//}
	//if err := svc.dao.AddServerOnline(context.Background(), serverId, online); err != nil {
	//	return nil, err
	//}

	reply.AllRoomCount = map[string]int32{}

	return
}

// Receive receive a message.
func (bz *Biz) Receive(ctx context.Context, userId int64, proto *protocol.Proto) (err error) {
	logHead := fmt.Sprintf("Receive,userId=%v|", userId)

	logging.Infof(logHead+"get message:%+v", userId, proto)
	return
}

// Nodes 获取节点信息
// - 对于comet来说，用于建立TCP服务器、WS服务器、WSS服务器
// - 对于客户端来说，用于建立TCP客户端、WS客户端、WSS客户端
// - 可以配合一定的负载均衡算法，实现节点的负载（根据每个节点的连接数进行负载分发）
func (bz *Biz) Nodes(ctx context.Context, req *pb.NodesReq) (reply *pb.NodesReply, err error) {
	nodeConf := bz.conf.Node
	backoffConf := bz.conf.Backoff
	reply = &pb.NodesReply{
		Domain:  nodeConf.DefaultDomain,
		TcpPort: nodeConf.TCPPort,
		WsPort:  nodeConf.WSPort,
		WssPort: nodeConf.WSSPort,
		HbCfg: &pb.HbCfg{
			Interval:  int64(time.Duration(nodeConf.Heartbeat.Interval).Seconds()),
			FailCount: nodeConf.Heartbeat.FailCount,
		},
		Backoff: &pb.Backoff{
			BaseDelay:  backoffConf.BaseDelay,
			Multiplier: backoffConf.Multiplier,
			Jitter:     backoffConf.Jitter,
			MaxDelay:   backoffConf.MaxDelay,
		},
	}

	// TODO 获取Nodes配置
	switch req.Platform {
	case pb.Platform_Platform_Web: // Web平台
		reply.Nodes = []string{}
	default:
		reply.Nodes = []string{}
	}
	if len(reply.Nodes) == 0 {
		reply.Nodes = []string{nodeConf.DefaultDomain}
	}

	return
}

func (bz *Biz) GetHeartbeatExpire() (result time.Duration) {
	nodeConf := bz.conf.Node
	result = time.Duration(nodeConf.Heartbeat.Interval) * time.Duration(nodeConf.Heartbeat.FailCount)

	return
}
