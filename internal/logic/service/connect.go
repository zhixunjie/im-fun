package service

import (
	"context"
	"encoding/json"
	"github.com/golang/glog"
	"github.com/google/uuid"
	"github.com/zhixunjie/im-fun/api/protocol"
)

type param struct {
	Mid      int64   `json:"mid"`
	Key      string  `json:"key"`
	RoomID   string  `json:"room_id"`
	Platform string  `json:"platform"`
	Accepts  []int32 `json:"accepts"`
}

// Connect connected a conn.
func (svc *Service) Connect(c context.Context, server, cookie string, token []byte) (mid int64, key, roomID string, accepts []int32, hb int64, err error) {
	var params param
	if err = json.Unmarshal(token, &params); err != nil {
		glog.Errorf("json.Unmarshal(%s) error(%v)", token, err)
		return
	}
	mid = params.Mid
	roomID = params.RoomID
	accepts = params.Accepts
	//hb = int64(l.c.Node.Heartbeat) * int64(svc.conf.Node.HeartbeatMax)
	if key = params.Key; key == "" {
		key = uuid.New().String()
	}
	//if err = svc.dao.AddMapping(c, mid, key, server); err != nil {
	//	glog.Errorf("l.dao.AddMapping(%d,%s,%s) error(%v)", mid, key, server, err)
	//}
	//glog.Infof("conn connected key:%s server:%s mid:%d token:%s", key, server, mid, token)
	return
}

// Disconnect disconnect a conn.
func (svc *Service) Disconnect(c context.Context, mid int64, key, server string) (has bool, err error) {
	//if has, err = svc.dao.DelMapping(c, mid, key, server); err != nil {
	//	glog.Errorf("l.dao.DelMapping(%d,%s) error(%v)", mid, key, server)
	//	return
	//}
	glog.Infof("conn disconnected key:%s server:%s mid:%d", key, server, mid)
	return
}

// Heartbeat heartbeat a conn.
func (svc *Service) Heartbeat(c context.Context, mid int64, key, server string) (err error) {
	//has, err := svc.dao.ExpireMapping(c, mid, key)
	//if err != nil {
	//	glog.Errorf("l.dao.ExpireMapping(%d,%s,%s) error(%v)", mid, key, server, err)
	//	return
	//}
	//if !has {
	//	if err = svc.dao.AddMapping(c, mid, key, server); err != nil {
	//		glog.Errorf("l.dao.AddMapping(%d,%s,%s) error(%v)", mid, key, server, err)
	//		return
	//	}
	//}
	glog.Infof("conn heartbeat key:%s server:%s mid:%d", key, server, mid)
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
	glog.Infof("receive mid:%d message:%+v", mid, proto)
	return
}
