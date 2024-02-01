package distrib_lock

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// 测试：自旋锁（通过超时时间获取）
func TestAcquireWithTimeout(t *testing.T) {
	lockKey := "spin_lock"

	// routine 1: 加锁成功，3秒后解锁返回
	go func() {
		logHead := "routine[1]|"
		err := TrySpinLockTimeout(logHead, lockKey, 3*time.Second)
		assert.NoError(t, err)
	}()
	time.Sleep(1 * time.Second)

	// routine 2: 加锁失败，自旋尝试获取锁（重试3次后获锁成功）
	go func() {
		logHead := "routine[2]|"
		err := TrySpinLockTimeout(logHead, lockKey, 3*time.Second)
		assert.NoError(t, err)
	}()

	time.Sleep(5 * time.Second)
}

// 测试：自旋锁（通过超时次数获取）
func TestAcquireWithTimes(t *testing.T) {
	lockKey := "spin_lock"

	// routine 1: 加锁成功，3秒后解锁返回
	go func() {
		logHead := "routine[1]|"
		err := TrySpinLockTimes(logHead, lockKey, 3*time.Second)
		assert.NoError(t, err)
	}()
	time.Sleep(1 * time.Second)

	// routine 2: 加锁失败，自旋尝试获取锁（重试3次后获锁成功）
	go func() {
		logHead := "routine[2]|"
		err := TrySpinLockTimes(logHead, lockKey, 3*time.Second)
		assert.NoError(t, err)
	}()

	time.Sleep(5 * time.Second)
}

func TrySpinLockTimeout(logHead, lockKey string, sleepAfterSuccess time.Duration) (err error) {
	redisSpinLock := NewSpinLock(client, lockKey, 5*time.Second, &SpinOption{
		Interval: 1000 * time.Millisecond,
		Timeout:  15 * time.Second,
	})
	if err = redisSpinLock.AcquireWithTimeout(); err != nil {
		fmt.Printf(logHead+"acquire fail,lockKey=%v,err=%v\n", lockKey, err)
		return
	}
	defer redisSpinLock.Release()
	fmt.Printf(logHead+"acquire success,lockKey=%v\n", lockKey)

	// sleep
	time.Sleep(sleepAfterSuccess)

	return
}

func TrySpinLockTimes(logHead, lockKey string, sleepAfterSuccess time.Duration) (err error) {
	options := &SpinOption{
		Interval: 1000 * time.Millisecond,
		Times:    10,
	}
	redisSpinLock := NewSpinLock(client, lockKey, 5*time.Second, options)
	if err = redisSpinLock.AcquireWithTimes(); err != nil {
		fmt.Printf(logHead+"acquire fail,lockKey=%v,err=%v\n", lockKey, err)
		return
	}
	defer redisSpinLock.Release()
	fmt.Printf(logHead+"acquire success,lockKey=%v\n", lockKey)

	// sleep
	time.Sleep(sleepAfterSuccess)

	return
}
