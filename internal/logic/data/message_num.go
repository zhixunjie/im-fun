package data

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cast"
	"github.com/zhixunjie/im-fun/pkg/gen_id"
	"github.com/zhixunjie/im-fun/pkg/logging"
	"time"
)

// 维护：消息未读数（会话未读数、总未读数）

const (
	KeySessionUnreadExpire = 86400 * 15
	KeyTotalUnreadExpire   = 86400 * 15
)

// Hash：owner -> [ peer : serverId ]
func keyHashSessionUnread(u *gen_id.ComponentId) string {
	return fmt.Sprintf("hash_session_unread_%v", u.ToString())
}

// String：owner -> 1000
func keyStringTotalUnread(u *gen_id.ComponentId) string {
	return fmt.Sprintf("string_total_unread_%v", u.ToString())
}

// IncrUnreadAfterSend 发送消息后，增加未读数
func (repo *MessageRepo) IncrUnreadAfterSend(ctx context.Context, logHead string, receiverId, senderId *gen_id.ComponentId, incr int64) (err error) {
	// clean before add
	err = repo.checkBeforeIncrSessionUnread(ctx, logHead, receiverId, senderId)
	if err != nil {
		return
	}

	// 增加：会话未读数
	_, err = repo.incrSessionUnread(ctx, logHead, receiverId, senderId, incr)
	if err != nil {
		return
	}

	// clean before add
	err = repo.checkBeforeIncrTotalUnread(ctx, logHead, receiverId)
	if err != nil {
		return
	}

	// 增加：总未读数（全部会话）
	_, err = repo.incrTotalUnread(ctx, logHead, receiverId, incr)
	if err != nil {
		return
	}

	return
}

// DecrUnreadAfterFetch 获取消息后，清空未读数
func (repo *MessageRepo) DecrUnreadAfterFetch(ctx context.Context, logHead string, owner, peer *gen_id.ComponentId, decr int64) (err error) {
	// 减少：会话未读数
	_, err = repo.incrSessionUnread(ctx, logHead, owner, peer, -decr)
	if err != nil {
		return
	}

	// 减少：总未读数
	_, err = repo.incrTotalUnread(ctx, logHead, owner, -decr)
	if err != nil {
		return
	}

	return
}

// clean before add
func (repo *MessageRepo) checkBeforeIncrSessionUnread(ctx context.Context, logHead string, receiverId, senderId *gen_id.ComponentId) (err error) {
	retMap, err := repo.MGetSessionUnread(ctx, logHead, receiverId, []*gen_id.ComponentId{senderId})
	if err != nil {
		return
	}

	// check result map
	srcVal, ok := retMap[senderId.ToString()]
	if ok && srcVal < 0 {
		err = repo.cleanSessionUnread(ctx, logHead, receiverId, senderId)
		if err != nil {
			return
		}
	}

	return
}

// clean before add
func (repo *MessageRepo) checkBeforeIncrTotalUnread(ctx context.Context, logHead string, receiverId *gen_id.ComponentId) (err error) {
	srcVal, err := repo.GetTotalUnread(ctx, logHead, receiverId)
	if err != nil {
		return
	}

	// 兼容错误：当遇到错误的数据时，把未读数据进行重置
	if srcVal < 0 {
		err = repo.cleanTotalUnread(ctx, logHead, receiverId)
		if err != nil {
			return
		}
	}

	return
}

////////////////////// 会话未读数

// incrSessionUnread 增减未读数（会话未读数）
func (repo *MessageRepo) incrSessionUnread(ctx context.Context, logHead string, id1, id2 *gen_id.ComponentId, incr int64) (afterIncr int64, err error) {
	mem := repo.RedisClient
	key := keyHashSessionUnread(id1)
	expire := KeySessionUnreadExpire * time.Second
	logHead += fmt.Sprintf("incrSessionUnread,key=%v|", key)

	// HIncrBy
	res := mem.HIncrBy(ctx, key, id2.ToString(), incr)
	if err = res.Err(); err != nil {
		logging.Errorf(logHead+"HIncrBy error=%v", err)
		return
	}
	afterIncr = res.Val()
	logging.Infof(logHead+"HIncrBy success,afterIncr=%v", afterIncr)

	// Expire
	if err = mem.Expire(ctx, key, expire).Err(); err != nil {
		logging.Errorf(logHead+"Expire error=%v", err)
		return
	}

	return
}

// MGetSessionUnread 获取未读数（会话未读数）
func (repo *MessageRepo) MGetSessionUnread(ctx context.Context, logHead string, id1 *gen_id.ComponentId, id2Arr []*gen_id.ComponentId) (retMap map[string]int64, err error) {
	mem := repo.RedisClient
	key := keyHashSessionUnread(id1)
	logHead += fmt.Sprintf("MGetSessionUnread,key=%v|", key)

	retMap = make(map[string]int64, len(id2Arr))
	for _, id2 := range id2Arr {
		// HGet
		res, tErr := mem.HGet(ctx, key, id2.ToString()).Result()
		if tErr != nil && !errors.Is(tErr, redis.Nil) {
			err = tErr
			logging.Errorf(logHead+"HGet error=%v", err)
			return
		}
		retMap[id2.ToString()] = cast.ToInt64(res)
	}

	return
}

// cleanSessionUnread 清空所有的未读数（会话未读数）
func (repo *MessageRepo) cleanSessionUnread(ctx context.Context, logHead string, OwnerId, PeerId *gen_id.ComponentId) (err error) {
	mem := repo.RedisClient
	key := keyHashSessionUnread(OwnerId)
	logHead += fmt.Sprintf("cleanSessionUnread,key=%v|", key)

	// HDel
	if err = mem.HDel(ctx, key, PeerId.ToString()).Err(); err != nil {
		logging.Errorf(logHead+"HDel error=%v", err)
		return
	}
	logging.Infof(logHead + "HDel success")

	return
}

////////////////////// 总未读数

// incrTotalUnread 增减未读数（总未读数）
func (repo *MessageRepo) incrTotalUnread(ctx context.Context, logHead string, id *gen_id.ComponentId, incr int64) (afterIncr int64, err error) {
	mem := repo.RedisClient
	key := keyStringTotalUnread(id)
	expire := KeyTotalUnreadExpire * time.Second
	logHead += fmt.Sprintf("incrTotalUnread,key=%v|", key)

	// IncrBy()
	res := mem.IncrBy(ctx, key, incr)
	if err = res.Err(); err != nil {
		logging.Errorf(logHead+"IncrBy error=%v", err)
		return
	}
	afterIncr = res.Val()
	logging.Infof(logHead+"IncrBy success,afterIncr=%v", afterIncr)

	// Expire
	if err = mem.Expire(ctx, key, expire).Err(); err != nil {
		logging.Errorf(logHead+"Expire error=%v", err)
		return
	}

	return
}

// GetTotalUnread 获取未读数（总未读数）
func (repo *MessageRepo) GetTotalUnread(ctx context.Context, logHead string, id *gen_id.ComponentId) (val int64, err error) {
	mem := repo.RedisClient
	key := keyStringTotalUnread(id)
	logHead += fmt.Sprintf("GetTotalUnread,key=%v|", key)

	// HGet
	res := mem.Get(ctx, key)
	if tErr := res.Err(); tErr != nil && tErr != redis.Nil {
		err = tErr
		logging.Errorf(logHead+"Get error=%v", err)
		return
	}
	val = cast.ToInt64(res.Val())

	return
}

// cleanTotalUnread 清空所有的未读数（总未读数）
func (repo *MessageRepo) cleanTotalUnread(ctx context.Context, logHead string, id *gen_id.ComponentId) (err error) {
	mem := repo.RedisClient
	key := keyStringTotalUnread(id)
	logHead += fmt.Sprintf("cleanTotalUnread,key=%v|", key)

	// HDel
	if err = mem.Del(ctx, key).Err(); err != nil {
		logging.Errorf(logHead+"Del error=%v", err)
		return
	}
	logging.Infof(logHead + "Del success")

	return
}
