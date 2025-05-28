package distrib_lock

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

var client = redis.NewClient(&redis.Options{
	Addr:     "127.0.0.1:6379",
	Password: "",
	DB:       0,
})

// 测试：加锁 & 解锁
func TestAcquireAndRelease(t *testing.T) {
	lockKey := "acquire:aaa"

	var wg sync.WaitGroup
	wg.Add(2)
	expire := 60 * time.Second

	// routine 1
	go func() {
		defer wg.Done()
		logHead := "routine[1]|"

		err := TryLock(logHead, lockKey, expire, 2*time.Second)
		assert.NoError(t, err)
	}()

	// routine 2
	go func() {
		defer wg.Done()
		time.Sleep(1 * time.Second)
		logHead := "routine[2]|"

		err := TryLock(logHead, lockKey, expire, 0*time.Second)
		assert.Error(t, err)
	}()

	//time.Sleep(time.Second * 20)
	wg.Wait()
}

// 测试：锁续约
func TestLockLease(t *testing.T) {
	logHead := "TestLockLease|"
	lockKey := "lease:aaa"

	// set time
	expire := 10 * time.Second
	leaseTime := expire / 2

	// get lock
	redisLock, err := GetLock(logHead, lockKey, expire)
	assert.NoError(t, err)
	defer redisLock.Release()

	// lease lock
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		ticker := time.NewTicker(leaseTime)
		// begin
		for {
			select {
			case <-ticker.C:
				if err := redisLock.Acquire(); err != nil {
					fmt.Printf(logHead+"lock lease fail,lockKey=%v,err=%v\n", lockKey, err)
					return
				}
				fmt.Printf(logHead+"lock lease success,lockKey=%v\n", lockKey)
			}
		}
	}()
	wg.Wait()
}

func TryLock(logHead string, lockKey string, expire, sleepAfterSuccess time.Duration) (err error) {
	redisLock := NewLock(client, lockKey, expire)
	if err = redisLock.Acquire(); err != nil {
		fmt.Printf(logHead+"acquire fail,lockKey=%v,err=%v\n", lockKey, err)
		return
	}
	defer redisLock.Release()
	fmt.Printf(logHead+"acquire success,lockKey=%v\n", lockKey)

	// sleep
	time.Sleep(sleepAfterSuccess)

	return
}

func GetLock(logHead string, lockKey string, expire time.Duration) (redisLock *RedisLock, err error) {
	redisLock = NewLock(client, lockKey, expire)
	if err = redisLock.Acquire(); err != nil {
		fmt.Printf(logHead+"acquire fail,lockKey=%v,err=%v\n", lockKey, err)
		return
	}
	fmt.Printf(logHead+"acquire success,lockKey=%v\n", lockKey)

	return
}
