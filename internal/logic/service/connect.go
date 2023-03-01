package service

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	pb "github.com/zhixunjie/im-fun/api/logic"
	"github.com/zhixunjie/im-fun/api/protocol"
)

// Connect connected a conn.
func (svc *Service) Connect(ctx context.Context, proto *pb.ConnectReq) (hb int64, err error) {
	if proto.UserId == 0 {
		logrus.Errorf("UserId not allow,token=%+v", proto.GetToken())
		return hb, errors.New("UserId not allow")
	}
	if proto.UserKey == "" {
		logrus.Errorf("UserId not allow,token=%+v", proto.GetToken())
		return hb, errors.New("UserId not allow")
	}
	if proto.Token == "" {
		logrus.Errorf("UserId not allow,token=%+v", proto.GetToken())
		return hb, errors.New("UserId not allow")
	}

	// set return
	hb = int64(svc.conf.Node.Heartbeat) * int64(svc.conf.Node.HeartbeatMax)
	if err = svc.dao.AddMapping(ctx, proto.UserId, proto.UserKey, proto.ServerId); err != nil {
		logrus.Errorf("AddMapping error=%v,UserId=%v,UserKey=%v", err, proto.UserId, proto.UserKey)
	}
	logrus.Infof("Connect success error=%v,UserId=%v,UserKey=%v", err, proto.UserId, proto.UserKey)
	return
}

// Disconnect disconnect a conn.
func (svc *Service) Disconnect(c context.Context, proto *pb.DisconnectReq) (has bool, err error) {
	if has, err = svc.dao.DelMapping(c, proto.UserId, proto.UserKey, proto.ServerId); err != nil {
		logrus.Errorf("DelMapping error=%v,UserId=%v,UserKey=%v", err, proto.UserId, proto.UserKey)
		return
	}
	logrus.Infof("Disconnect success error=%v,UserId=%v,UserKey=%v", err, proto.UserId, proto.UserKey)
	return
}

// Heartbeat heartbeat a conn.
func (svc *Service) Heartbeat(c context.Context, mid int64, key, server string) (err error) {
	//has, err := svc.dao.ExpireMapping(c, mid, key)
	//if err != nil {
	//	logrus.Errorf("l.dao.ExpireMapping(%d,%s,%s) error(%v)", mid, key, server, err)
	//	return
	//}
	//if !has {
	//	if err = svc.dao.AddMapping(c, mid, key, server); err != nil {
	//		logrus.Errorf("l.dao.AddMapping(%d,%s,%s) error(%v)", mid, key, server, err)
	//		return
	//	}
	//}
	logrus.Infof("conn heartbeat key:%s server:%s mid:%d", key, server, mid)
	return
}

// RenewOnline renew a server online.
func (svc *Service) RenewOnline(c context.Context, server string, roomCount map[string]int32) (map[string]int32, error) {
	//online := &model.Online{
	//	Server:    server,
	//	RoomCount: roomCount,
	//	Updated:   time.Now().Unix(),
	//}
	//if err := svc.dao.AddServerOnline(context.Background(), server, online); err != nil {
	//	return nil, err
	//}
	return map[string]int32{}, nil
}

// Receive receive a message.
func (svc *Service) Receive(c context.Context, mid int64, proto *protocol.Proto) (err error) {
	logrus.Infof("receive mid:%d message:%+v", mid, proto)
	return
}
