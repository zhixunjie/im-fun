package biz

import (
	"context"
	"fmt"
	"github.com/zhixunjie/im-fun/api/pb"
	"github.com/zhixunjie/im-fun/api/protocol"
	"github.com/zhixunjie/im-fun/pkg/logging"
)

//type TcpUseCase struct {
//	contactRepo *data.TcpRepo
//}
//
//func NewTcpUseCase(contactRepo *data.TcpRepo) *TcpUseCase {
//	return &TcpUseCase{contactRepo: contactRepo}
//}

// Connect connected a conn.
func (bz *Biz) Connect(ctx context.Context, req *pb.ConnectReq) (hb int64, err error) {
	rr := req.Comm
	logHead := fmt.Sprintf("Connect,rr=%v|", rr)

	if err = bz.data.SessionBinding(ctx, rr.UserId, rr.TcpSessionId, rr.ServerId); err != nil {
		logging.Errorf(logHead+"SessionBinding fail,error=%v", err)
		return
	}

	// return hb
	hb = int64(bz.conf.Node.Heartbeat) * int64(bz.conf.Node.HeartbeatMax)

	logging.Infof(logHead + "success")
	return
}

// Disconnect disconnect a conn.
func (bz *Biz) Disconnect(c context.Context, req *pb.DisconnectReq) (has bool, err error) {
	rr := req.Comm
	logHead := fmt.Sprintf("Disconnect,rr=%v|", rr)

	if has, err = bz.data.SessionDel(c, rr.UserId, rr.TcpSessionId, rr.ServerId); err != nil {
		logging.Errorf(logHead+"SessionDel fail,error=%v", err, rr.UserId, rr.TcpSessionId)

		return
	}

	logging.Infof(logHead + "success")
	return
}

// Heartbeat heartbeat a conn.
func (bz *Biz) Heartbeat(c context.Context, req *pb.HeartbeatReq) (err error) {
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
	//logging.Infof("conn heartbeat userKey:%s serverId:%s userId:%d", userKey, serverId, userId)
	return
}

// RenewOnline renew a server online.
func (bz *Biz) RenewOnline(c context.Context, serverId string, roomCount map[string]int32) (map[string]int32, error) {
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
func (bz *Biz) Receive(c context.Context, userId int64, proto *protocol.Proto) (err error) {
	logging.Infof("receive userId:%d message:%+v", userId, proto)
	return
}
