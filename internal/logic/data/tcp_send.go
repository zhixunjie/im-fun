package data

import (
	"context"
	"github.com/redis/go-redis/v9"
)

// GetServerIds 获取TcpSessionId对应的ServerId
func (d *Data) GetServerIds(ctx context.Context, tcpSessionIds []string) (res []string, err error) {
	mem := d.RedisClient

	// exec command
	var cmds []redis.Cmder
	cmds, err = mem.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		for i := 0; i < len(tcpSessionIds); i++ {
			pipe.Get(ctx, keyStringTcpSessionId(tcpSessionIds[i]))
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

// GetSessionByUserIds 获取用户ID对应的Session信息
func (d *Data) GetSessionByUserIds(ctx context.Context, userIds []uint64) (res map[string]string, err error) {
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
		tmpMap := cmd.(*redis.MapStringStringCmd).Val()
		for tcpSessionId, serverId := range tmpMap {
			res[tcpSessionId] = serverId
		}
	}
	return res, nil
}
