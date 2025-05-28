package distrib_lock

import (
	"github.com/redis/go-redis/v9"
	"time"
)

const (
	// 默认的自旋间隔
	defaultSpinInterval = 50 * time.Millisecond
	// 默认的自旋超时时间
	defaultSpinTimout = 3 * time.Second
	// 默认的自旋次数
	defaultSpinTimes = 10
)

// A RedisSpinLock is a redis spin lock.
type RedisSpinLock struct {
	l       *RedisLock
	options *SpinOption
}

type SpinOption struct {
	Interval time.Duration // 自旋的间隔，默认为: defaultSpinInterval
	Timeout  time.Duration // 自旋的超时时间
	Times    uint32        // 自旋的次数
}

// NewSpinLock returns a RedisLock.
// expireSec: 单位: 秒
func NewSpinLock(store *redis.Client, key string, expire time.Duration, options *SpinOption) *RedisSpinLock {
	// deal options
	if options.Interval == 0 {
		options.Interval = defaultSpinInterval
	}
	if options.Timeout == 0 {
		options.Timeout = defaultSpinTimout
	}
	if options.Times == 0 {
		options.Times = defaultSpinTimes
	}

	return &RedisSpinLock{
		l:       NewLock(store, key, expire),
		options: options,
	}
}

// AcquireWithTimeout acquires lock in specified time.
func (rl *RedisSpinLock) AcquireWithTimeout() (err error) {
	timer := time.NewTimer(rl.options.Timeout)
	defer timer.Stop()

	// begin
	for {
		select {
		case <-timer.C:
			return
		default:
			if err = rl.l.Acquire(); err != nil {
				if err != ErrorLockExists {
					return
				}
				time.Sleep(rl.options.Interval)
				continue
			}
			return
		}
	}
}

// AcquireWithTimes acquires lock in specified times.
func (rl *RedisSpinLock) AcquireWithTimes() (err error) {

	// begin
	for i := uint32(0); i < rl.options.Times; i++ {
		if err = rl.l.Acquire(); err != nil {
			if err != ErrorLockExists {
				return
			}
			time.Sleep(rl.options.Interval)
			continue
		}
		return
	}
	return
}

func (rl *RedisSpinLock) Release() error {
	return rl.l.Release()
}
