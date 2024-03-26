package data

import (
	"context"
	"fmt"
	"github.com/zhixunjie/im-fun/api/pb"
	"github.com/zhixunjie/im-fun/pkg/logging"
	"time"
)

const (
	KeyExpire = 3600
)

// Hash：userId -> [ tcpSessionId : serverId ]
func keyHashUserId(userId uint64) string {
	return fmt.Sprintf("session_hash_%d", userId)
}

// String：tcpSessionId -> serverId
func keyStringTcpSessionId(tcpSessionId string) string {
	return fmt.Sprintf("session_string_%s", tcpSessionId)
}

// server -> online
func keyServerOnline(key string) string {
	return fmt.Sprintf("online_%s", key)
}

// SessionBinding KEY绑定
func (d *Data) SessionBinding(ctx context.Context, logHead string, rr *pb.ConnectCommon) (err error) {
	logHead += "SessionBinding|"
	mem := d.RedisClient
	expire := KeyExpire * time.Second
	serverId := rr.ServerId
	userId := rr.UserId
	tcpSessionId := rr.TcpSessionId

	// set hash
	if userId > 0 {
		key := keyHashUserId(userId)
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
		key := keyStringTcpSessionId(tcpSessionId)
		if err = mem.SetEX(ctx, key, serverId, expire).Err(); err != nil {
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
		key := keyHashUserId(userId)
		if err = mem.HDel(ctx, key, tcpSessionId).Err(); err != nil {
			logging.Errorf(logHead+"HDel error=%v,key=%v", err, key)
			return
		}
		logging.Infof(logHead+"HDel success,key=%v", key)
	}
	// delete string
	key := keyStringTcpSessionId(tcpSessionId)
	if err = mem.Del(ctx, key).Err(); err != nil {
		logging.Errorf(logHead+"Del error=%v,key=%v", err, key)
		return
	}
	logging.Infof(logHead+"Del success,key=%v", key)

	return
}

// SessionLease KEY续约
func (d *Data) SessionLease(ctx context.Context, logHead string, rr *pb.ConnectCommon) (has bool, err error) {
	logHead += "SessionLease|"

	mem := d.RedisClient
	expire := KeyExpire * time.Second
	//serverId := rr.ServerId
	userId := rr.UserId
	tcpSessionId := rr.TcpSessionId

	// expire 1
	key := keyHashUserId(userId)
	if err = mem.Expire(ctx, key, expire).Err(); err != nil {
		logging.Errorf(logHead+"Expire(1) error=%v,key=%v", err, key)
		return
	}
	logging.Infof(logHead+"Expire(1) success,key=%v", key)

	// expire 2
	key = keyStringTcpSessionId(tcpSessionId)
	if err = mem.Expire(ctx, key, expire).Err(); err != nil {
		logging.Errorf(logHead+"Expire(2) error=%v,key=%v", err, key)
		return
	}
	logging.Infof(logHead+"Expire(2) success,key=%v", key)

	return
}
