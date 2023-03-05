package dao

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"time"
)

const (
	KeyExpire = 3600
)

// Hash：userId -> userKey => serverId
func keyHashUserId(userId int64) string {
	return fmt.Sprintf("session_hash_%d", userId)
}

// String：userKey => server
func keyStringUserKey(userKey string) string {
	return fmt.Sprintf("session_string_%s", userKey)
}

// server => online
func keyServerOnline(key string) string {
	return fmt.Sprintf("online_%s", key)
}

// SessionBinding add relationship
func (d *Dao) SessionBinding(ctx context.Context, userId int64, userKey, serverId string) (err error) {
	mem := d.RedisClient

	// set hash
	if userId > 0 {
		k1 := keyHashUserId(userId)
		if err = mem.HSet(ctx, k1, userKey, serverId).Err(); err != nil {
			logrus.Errorf("mem.HSet(%d,%s,%s) error(%v)", userId, userKey, serverId, err)
			return
		}
		if err = mem.Expire(ctx, k1, KeyExpire*time.Second).Err(); err != nil {
			logrus.Errorf("mem.Expire(%d,%s,%s) error(%v)", userId, userKey, serverId, err)
			return
		}
	}
	// set string
	{
		if err = mem.SetEX(ctx, keyStringUserKey(userKey), serverId, KeyExpire*time.Second).Err(); err != nil {
			logrus.Errorf("mem.SetEX(%d,%s,%s) error(%v)", userId, serverId, userKey, err)
			return
		}
	}

	return
}

func (d *Dao) SessionDel(ctx context.Context, userId int64, userKey, serverId string) (has bool, err error) {
	mem := d.RedisClient

	// delete hash
	if userId > 0 {
		if err = mem.HDel(ctx, keyHashUserId(userId), userKey).Err(); err != nil {
			logrus.Errorf("mem.HDel(%d,%s,%s) error(%v)", userId, serverId, userKey, err)
			return
		}
	}
	// delete string
	if err = mem.Del(ctx, keyStringUserKey(userKey)).Err(); err != nil {
		logrus.Errorf("mem.Del(%d,%s,%s) error(%v)", userId, serverId, userKey, err)
		return
	}

	return
}
