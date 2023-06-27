package service

import (
	"context"
	"github.com/zhixunjie/im-fun/api/pb"
	"github.com/zhixunjie/im-fun/api/protocol"
	"github.com/zhixunjie/im-fun/pkg/logging"
)

// Connect connected a conn.
func (svc *Service) Connect(ctx context.Context, req *pb.ConnectReq) (hb int64, err error) {
	logHead := "Connect|"

	// set return
	hb = int64(svc.conf.Node.Heartbeat) * int64(svc.conf.Node.HeartbeatMax)
	if err = svc.dao.SessionBinding(ctx, req.UserId, req.UserKey, req.ServerId); err != nil {
		logging.Errorf(logHead+"fail,SessionBinding error=%v,UserId=%v,UserKey=%v", err, req.UserId, req.UserKey)
		return
	}

	logging.Infof(logHead+"success,UserId=%v,UserKey=%v", req.UserId, req.UserKey)
	return
}

// Disconnect disconnect a conn.
func (svc *Service) Disconnect(c context.Context, req *pb.DisconnectReq) (has bool, err error) {
	logHead := "Disconnect|"
	if has, err = svc.dao.SessionDel(c, req.UserId, req.UserKey, req.ServerId); err != nil {
		logging.Errorf(logHead+"fail,SessionDel error=%v,UserId=%v,UserKey=%v", err, req.UserId, req.UserKey)
		return
	}
	logging.Infof(logHead+"success,UserId=%v,UserKey=%v", req.UserId, req.UserKey)
	return
}

// Heartbeat heartbeat a conn.
func (svc *Service) Heartbeat(c context.Context, userId int64, userKey, serverId string) (err error) {
	//has, err := svc.dao.ExpireMapping(c, userId, userKey)
	//if err != nil {
	//	logging.Errorf("l.dao.ExpireMapping(%d,%s,%s) error(%v)", userId, userKey, serverId, err)
	//	return
	//}
	//if !has {
	//	if err = svc.dao.SessionBinding(c, userId, userKey, serverId); err != nil {
	//		logging.Errorf("l.dao.SessionBinding(%d,%s,%s) error(%v)", userId, userKey, serverId, err)
	//		return
	//	}
	//}
	logging.Infof("conn heartbeat userKey:%s serverId:%s userId:%d", userKey, serverId, userId)
	return
}

// RenewOnline renew a server online.
func (svc *Service) RenewOnline(c context.Context, serverId string, roomCount map[string]int32) (map[string]int32, error) {
	//online := &model.Online{
	//	Server:    serverId,
	//	RoomCount: roomCount,
	//	Updated:   time.Now().Unix(),
	//}
	//if err := svc.dao.AddServerOnline(context.Background(), serverId, online); err != nil {
	//	return nil, err
	//}
	return map[string]int32{}, nil
}

// Receive receive a message.
func (svc *Service) Receive(c context.Context, userId int64, proto *protocol.Proto) (err error) {
	logging.Infof("receive userId:%d message:%+v", userId, proto)
	return
}
