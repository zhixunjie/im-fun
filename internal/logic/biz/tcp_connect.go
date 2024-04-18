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
	logHead := fmt.Sprintf("Connect,req.Comm=%v|", rr)

	if err = bz.data.SessionBinding(ctx, logHead, rr); err != nil {
		logging.Errorf(logHead+"SessionBinding fail,error=%v", err)
		return
	}

	// return hb
	hb = int64(bz.conf.Node.Heartbeat) * int64(bz.conf.Node.HeartbeatMax)
	logging.Infof(logHead + "success")

	return
}

// Disconnect disconnect a conn.
func (bz *Biz) Disconnect(ctx context.Context, req *pb.DisconnectReq) (has bool, err error) {
	rr := req.Comm
	logHead := fmt.Sprintf("Disconnect,req.Comm=%v|", rr)

	has, err = bz.data.SessionDel(ctx, logHead, rr)
	if err != nil {
		logging.Errorf(logHead+"SessionDel fail,error=%v", err, rr.UserId, rr.TcpSessionId)

		return
	}

	logging.Infof(logHead + "success")
	return
}

// Heartbeat heartbeat a conn.
func (bz *Biz) Heartbeat(ctx context.Context, req *pb.HeartbeatReq) (err error) {
	rr := req.Comm
	logHead := fmt.Sprintf("Heartbeat,req.Comm=%v|", rr)

	// 如果KEY不存在，就不要再去续约Redis的KEY了
	_, err = bz.data.SessionLease(ctx, logHead, rr)
	if err != nil {
		logging.Errorf(logHead+"SessionLease fail,error=%v", err)
		return
	}
	// 重新建立绑定关系
	//if !has {
	//	if err = bz.data.SessionBinding(ctx, logHead, rr); err != nil {
	//		logging.Errorf(logHead+"SessionBinding fail,error=%v", err)
	//		return
	//	}
	//}
	logging.Infof(logHead + "success")

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
	logHead := fmt.Sprintf("Receive,userId=%v|", userId)

	logging.Infof(logHead+"get message:%+v", userId, proto)
	return
}
