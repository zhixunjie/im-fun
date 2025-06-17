package biz

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/spf13/cast"
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
func (bz *Biz) Connect(ctx context.Context, req *pb.ConnectReq) (resp *pb.ConnectResp, err error) {
	resp = new(pb.ConnectResp)
	authParams := req.AuthParams
	expire := bz.GetHeartbeatExpire()
	logHead := fmt.Sprintf("Connect,authParams.UniId=%+v,expire=%v|", authParams.UniId, expire)

	// check token
	claims, err := bz.userUseCase.checkToken(authParams.Token)
	if err != nil {
		err = fmt.Errorf("token check failed: %w", err)
		return
	}
	if cast.ToString(claims.Uid) != authParams.UniId {
		err = fmt.Errorf("token not allows(%v,%v)", authParams.UniId, claims.Uid)
		return
	}

	sessionId := uuid.NewString()
	serverId := req.ServerId
	uniId := req.AuthParams.UniId
	if err = bz.data.SessionBinding(ctx, logHead, uniId, sessionId, serverId, expire); err != nil {
		err = fmt.Errorf("session binding failed: %w", err)
		return
	}

	// return hb
	hbCfg := bz.conf.Node.Heartbeat
	resp = &pb.ConnectResp{
		HbCfg: &pb.HbCfg{
			Interval:   int64(time.Duration(hbCfg.Interval).Seconds()),
			FailCount:  hbCfg.FailCount,
			BindExpire: int64(expire.Seconds()),
		},
		SessionId: sessionId,
	}
	logging.Infof(logHead+"success,hbCfg=%+v", hbCfg)

	return
}

// Disconnect disconnect a conn.
func (bz *Biz) Disconnect(ctx context.Context, req *pb.DisconnectReq) (resp *pb.DisconnectResp, err error) {
	resp = new(pb.DisconnectResp)
	connect := req.Connect
	logHead := fmt.Sprintf("Disconnect,connect=%+v|", connect)

	resp.Has, err = bz.data.SessionDel(ctx, logHead, connect)
	if err != nil {
		logging.Errorf(logHead+"SessionDel fail", err)

		return
	}

	logging.Infof(logHead + "success")
	return
}

// Heartbeat heartbeat a conn.
func (bz *Biz) Heartbeat(ctx context.Context, req *pb.HeartbeatReq) (resp *pb.HeartbeatResp, err error) {
	resp = new(pb.HeartbeatResp)
	connect := req.Connect
	logHead := fmt.Sprintf("Heartbeat,req=%+v|", req)

	// 续约KEY
	expire := time.Duration(req.BindExpire) * time.Second
	resp.Has, err = bz.data.SessionLease(ctx, logHead, connect, expire)
	if err != nil {
		logging.Errorf(logHead+"SessionLease fail,error=%v", err)
		return
	}
	// 重新建立绑定关系
	if !resp.Has {
		uniId := connect.UniId
		sessionId := connect.SessionId
		serverId := connect.ServerId
		if err = bz.data.SessionBinding(ctx, logHead, uniId, sessionId, serverId, expire); err != nil {
			logging.Errorf(logHead+"SessionBinding fail,error=%v", err)
			return
		}
	}
	logging.Infof(logHead + "success")

	return
}

// RenewOnline renew a server online.
func (bz *Biz) RenewOnline(ctx context.Context, serverId string, roomCount map[string]int32) (resp *pb.OnlineResp, err error) {
	resp = new(pb.OnlineResp)
	//online := &model.Online{
	//	Server:    serverId,
	//	RoomCount: roomCount,
	//	Updated:   time.Now().Unix(),
	//}
	//if err := svc.dao.AddServerOnline(context.Background(), serverId, online); err != nil {
	//	return nil, err
	//}

	resp.AllRoomCount = map[string]int32{}

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
func (bz *Biz) Nodes(ctx context.Context, req *pb.NodesReq) (resp *pb.NodesResp, err error) {
	nodeConf := bz.conf.Node
	backoffConf := bz.conf.Backoff
	resp = &pb.NodesResp{
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
		resp.Nodes = []string{}
	default:
		resp.Nodes = []string{}
	}
	if len(resp.Nodes) == 0 {
		resp.Nodes = []string{nodeConf.DefaultDomain}
	}

	return
}

// GetHeartbeatExpire 单位：秒
func (bz *Biz) GetHeartbeatExpire() (result time.Duration) {
	nodeConf := bz.conf.Node
	result = time.Duration(nodeConf.Heartbeat.Interval)*time.Duration(nodeConf.Heartbeat.FailCount) + 10*time.Second

	return
}
