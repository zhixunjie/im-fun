package distrib_lock

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/cast"
	"github.com/zhixunjie/im-fun/pkg/logging"
	"math/rand"
	"time"
)

var (
	ErrorLockExists         = errors.New("lock exists")
	ErrorResponseNotAllowed = errors.New("response is not allowed")
	ErrorReplyNotAllowed    = errors.New("reply is not allowed")
)

const (
	randomLen = 16

	// 默认的锁过期时间，单位：milliseconds
	// 防止：
	// 1. 使用者不设置锁过期时间，导致锁永不过期
	// 2. 使用者设置的锁时间太短（如：1ms），导致锁马上就过期
	tolerance = 500 * time.Millisecond
)

// lua 脚本
const (
	// 加锁脚本
	// 注意：PX使用的单位是毫秒
	lockScript = `if redis.call("GET", KEYS[1]) == ARGV[1] then
    redis.call("SET", KEYS[1], ARGV[1], "PX", ARGV[2])
    return "OK"
else
    return redis.call("SET", KEYS[1], ARGV[1], "NX", "PX", ARGV[2])
end`

	// 解锁脚本
	delScript = `if redis.call("GET", KEYS[1]) == ARGV[1] then
    return redis.call("DEL", KEYS[1])
else
    return 0
end`
)

// A RedisLock is a redis lock.
type RedisLock struct {
	ctx    context.Context
	store  *redis.Client // redis客户端
	expire time.Duration // 超时时间
	key    string        // 锁的redis key
	id     string        // 锁的唯一id，用于标识锁的唯一性（防止：被别人解锁、被别人续约等）

	logHead string
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

// NewLock returns a RedisLock.
func NewLock(store *redis.Client, key string, expire time.Duration) *RedisLock {
	if expire.Seconds() < 1 {
		panic("NewRedisLockWithExpire|expire not allow")
	}

	return &RedisLock{
		ctx:     context.TODO(),
		store:   store,
		expire:  expire,
		key:     key,
		id:      Randn(randomLen),
		logHead: "RedisLock|",
	}
}

// Acquire acquires the lock.
func (rl *RedisLock) Acquire() error {
	logHead := rl.logHead

	// 设置过期时间(在原来的基础上加上tolerance)
	lockTime := cast.ToString(rl.expire.Milliseconds() + tolerance.Milliseconds())
	//logging.Infof(logHead+"ready to acquire lock,lockTime=%vms", lockTime)

	// 执行lua脚本
	resp, err := rl.store.Eval(rl.ctx, lockScript, []string{rl.key}, []string{rl.id, lockTime}).Result()
	if err != nil {
		if err == redis.Nil { // why nil？https://redis.io/commands/set/
			err = ErrorLockExists
		}
		logging.Errorf(logHead+"Error on acquiring lock for key=%s,err=%s", rl.key, err.Error())
		return err
	}

	// type conversion for reply
	reply, ok := resp.(string)
	if !ok {
		err = ErrorResponseNotAllowed
		logging.Errorf(logHead+"Error on acquiring lock for key=%s,err=%s,resp=%v", rl.key, err.Error(), resp)
		return err
	}

	// check reply
	if reply == "OK" {
		logging.Infof(logHead+"Success on acquiring lock for key=%s", rl.key)
		return nil
	} else {
		err = ErrorReplyNotAllowed
		logging.Errorf(logHead+"Error on acquiring lock for key=%s,err=%s,reply=%v", rl.key, err.Error(), reply)
		return err
	}
}

// Release releases the lock.
// 释放锁
func (rl *RedisLock) Release() error {
	logHead := rl.logHead

	resp, err := rl.store.Eval(rl.ctx, delScript, []string{rl.key}, []string{rl.id}).Result()
	if err != nil {
		logging.Errorf(logHead+"Error on releasing lock for key=%s,err=%s", rl.key, err.Error())
		return err
	}

	// type conversion for reply
	reply, ok := resp.(int64)
	if !ok {
		err = ErrorResponseNotAllowed
		logging.Errorf(logHead+"Error on releasing lock for key=%s,err=%s,resp=%v", rl.key, err.Error(), resp)
		return err
	}

	// check reply
	if reply == 1 {
		logging.Infof(logHead+"Success on releasing lock for key=%s", rl.key)
		return nil
	} else {
		err = ErrorReplyNotAllowed
		logging.Errorf(logHead+"Error on releasing lock for key=%s,err=%s,reply=%v", rl.key, err.Error(), reply)
		return err
	}
}
