package data

import (
	"context"
	"github.com/go-redis/redis/v8"
)

// SessionGetByUserKeys 获取多个userKey的string信息
func (d *Data) SessionGetByUserKeys(ctx context.Context, userKeys []string) (res []string, err error) {
	mem := d.RedisClient

	// exec command
	var cmds []redis.Cmder
	cmds, err = mem.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		for i := 0; i < len(userKeys); i++ {
			pipe.Get(ctx, keyStringUserKey(userKeys[i]))
		}
		return nil
	})
	if err != nil {
		return
	}

	// get command result
	res = make([]string, len(cmds))
	for i, cmd := range cmds {
		res[i] = cmd.(*redis.StringCmd).Val()
	}
	return res, nil
}

// SessionGetByUserIds 获取多个userId的Hash信息
func (d *Data) SessionGetByUserIds(ctx context.Context, userIds []int64) (res map[string]string, err error) {
	mem := d.RedisClient

	// exec command
	var cmds []redis.Cmder
	cmds, err = mem.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		for i := 0; i < len(userIds); i++ {
			pipe.HGetAll(ctx, keyHashUserId(userIds[i]))
		}
		return nil
	})
	if err != nil {
		return
	}
	// get command result
	res = make(map[string]string)
	for _, cmd := range cmds {
		tmpMap := cmd.(*redis.StringStringMapCmd).Val()
		for userKey, serverId := range tmpMap {
			res[userKey] = serverId
		}
	}
	return res, nil
}
