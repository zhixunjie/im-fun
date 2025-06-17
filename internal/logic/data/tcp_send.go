package data

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

// GetServerIds 获取TcpSessionId对应的ServerId
func (d *Data) GetServerIds(ctx context.Context, tcpSessionIds []string) (res []string, err error) {
	mem := d.RedisClient

	// exec command
	var cmds []redis.Cmder
	cmds, err = mem.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		for i := 0; i < len(tcpSessionIds); i++ {
			pipe.Get(ctx, fmt.Sprintf(TcpSessionToSrv, tcpSessionIds[i]))
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

// GetSessionByUniIds 获取用户ID对应的Session信息
func (d *Data) GetSessionByUniIds(ctx context.Context, uniIds []string) (res map[string]string, err error) {
	mem := d.RedisClient

	// exec command
	var cmds []redis.Cmder
	cmds, err = mem.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		for i := 0; i < len(uniIds); i++ {
			pipe.HGetAll(ctx, fmt.Sprintf(TcpUserAllSession, uniIds[i]))
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
