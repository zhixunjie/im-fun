package data

import (
	"context"
	"fmt"
	"github.com/zhixunjie/im-fun/api/pb"
	"github.com/zhixunjie/im-fun/pkg/logging"
	"time"
)

// SessionBinding KEY绑定
func (d *Data) SessionBinding(ctx context.Context, logHead string, rr *pb.ConnectCommon, expire time.Duration) (err error) {
	logHead += fmt.Sprintf("SessionBinding,expire=%v|", expire)
	mem := d.RedisClient
	serverId := rr.ServerId
	userId := rr.UserId
	tcpSessionId := rr.TcpSessionId

	// set hash
	if userId > 0 {
		key := fmt.Sprintf(TcpUserAllSession, userId)
		// HSet
		if err = mem.HSet(ctx, key, tcpSessionId, serverId).Err(); err != nil {
			logging.Errorf(logHead+"HSet error=%v,key=%v", key)
			return
		}
		logging.Infof(logHead+"HSet success,key=%v", key)
		// Expire
		if err = mem.Expire(ctx, key, expire).Err(); err != nil {
			logging.Errorf(logHead+"Expire error=%v,key=%v", key)
			return
		}
	}
	// set string
	{
		key := fmt.Sprintf(TcpSessionToSrv, tcpSessionId)
		if err = mem.SetEx(ctx, key, serverId, expire).Err(); err != nil {
			logging.Errorf(logHead+"SetEX error=%v,key=%v", key)
			return
		}
		logging.Infof(logHead+"SetEX success,key=%v", key)
	}

	return
}

// SessionDel KEY删除
func (d *Data) SessionDel(ctx context.Context, logHead string, rr *pb.ConnectCommon) (has bool, err error) {
	logHead += "SessionDel|"
	mem := d.RedisClient
	//serverId := rr.ServerId
	userId := rr.UserId
	tcpSessionId := rr.TcpSessionId

	// delete hash
	if userId > 0 {
		// HDel
		key := fmt.Sprintf(TcpUserAllSession, userId)
		if err = mem.HDel(ctx, key, tcpSessionId).Err(); err != nil {
			logging.Errorf(logHead+"HDel error=%v,key=%v", err, key)
			return
		}
		logging.Infof(logHead+"HDel success,key=%v", key)
	}
	// delete string
	key := fmt.Sprintf(TcpSessionToSrv, tcpSessionId)
	if err = mem.Del(ctx, key).Err(); err != nil {
		logging.Errorf(logHead+"Del error=%v,key=%v", err, key)
		return
	}
	logging.Infof(logHead+"Del success,key=%v", key)

	return
}

// SessionLease KEY续约
func (d *Data) SessionLease(ctx context.Context, logHead string, rr *pb.ConnectCommon, expire time.Duration) (has bool, err error) {
	logHead += "SessionLease|"

	mem := d.RedisClient
	//serverId := rr.ServerId
	userId := rr.UserId
	tcpSessionId := rr.TcpSessionId

	// expire 1（续约 Hash KEY）
	key := fmt.Sprintf(TcpUserAllSession, userId)
	has, err = mem.Expire(ctx, key, expire).Result()
	if err != nil {
		logging.Errorf(logHead+"Expire(1) error=%v,key=%v", err, key)
		return
	}
	logging.Infof(logHead+"Expire(1) success,key=%v", key)

	// expire 2（续约 String KEY）
	key = fmt.Sprintf(TcpSessionToSrv, tcpSessionId)
	has, err = mem.Expire(ctx, key, expire).Result()
	if err != nil {
		logging.Errorf(logHead+"Expire(2) error=%v,key=%v", err, key)
		return
	}
	logging.Infof(logHead+"Expire(2) success,key=%v", key)

	return
}
