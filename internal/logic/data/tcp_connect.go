package data

import (
	"context"
	"fmt"
	"github.com/zhixunjie/im-fun/pkg/logging"
	"time"
)

const (
	KeyExpire = 3600
)

// Hash：userId [ tcpSessionId => serverId ]
func keyHashUserId(userId uint64) string {
	return fmt.Sprintf("session_hash_%d", userId)
}

// String：tcpSessionId => serverId
func keyStringUserTcpSessionId(tcpSessionId string) string {
	return fmt.Sprintf("session_string_%s", tcpSessionId)
}

// server => online
func keyServerOnline(key string) string {
	return fmt.Sprintf("online_%s", key)
}

// SessionBinding add relationship
func (d *Data) SessionBinding(ctx context.Context, userId uint64, tcpSessionId, serverId string) (err error) {
	mem := d.RedisClient

	// set hash
	if userId > 0 {
		key := keyHashUserId(userId)
		if err = mem.HSet(ctx, key, tcpSessionId, serverId).Err(); err != nil {
			logging.Errorf("mem.HSet(%d,%s,%s) error(%v)", userId, tcpSessionId, serverId, err)
			return
		}
		if err = mem.Expire(ctx, key, KeyExpire*time.Second).Err(); err != nil {
			logging.Errorf("mem.Expire(%d,%s,%s) error(%v)", userId, tcpSessionId, serverId, err)
			return
		}
	}
	// set string
	{
		if err = mem.SetEX(ctx, keyStringUserTcpSessionId(tcpSessionId), serverId, KeyExpire*time.Second).Err(); err != nil {
			logging.Errorf("mem.SetEX(%d,%s,%s) error(%v)", userId, serverId, tcpSessionId, err)
			return
		}
	}

	return
}

func (d *Data) SessionDel(ctx context.Context, userId uint64, tcpSessionId, serverId string) (has bool, err error) {
	mem := d.RedisClient

	// delete hash
	if userId > 0 {
		if err = mem.HDel(ctx, keyHashUserId(userId), tcpSessionId).Err(); err != nil {
			logging.Errorf("mem.HDel(%d,%s,%s) error(%v)", userId, serverId, tcpSessionId, err)
			return
		}
	}
	// delete string
	if err = mem.Del(ctx, keyStringUserTcpSessionId(tcpSessionId)).Err(); err != nil {
		logging.Errorf("mem.Del(%d,%s,%s) error(%v)", userId, serverId, tcpSessionId, err)
		return
	}

	return
}
