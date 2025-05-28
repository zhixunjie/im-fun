package goredis

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/zhixunjie/im-fun/pkg/logging"
	"time"
)

// CreatePool 创建Redis连接池
func CreatePool(addr, password string, db int) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,

		// 连接池容量及闲置连接数量
		// 连接池最大socket连接数，默认为4倍CPU数， 4 * runtime.NumCPU
		PoolSize: 500,
		// 在启动阶段创建指定数量的Idle连接，并长期维持idle状态的连接数不少于指定数量。
		MinIdleConns: 100,

		// 闲置连接检查包括IdleTimeout，MaxConnAge
		// IdleCheckFrequency: 60 * time.Second, // 闲置连接检查的周期，默认为1分钟，-1表示不做周期性检查，只在客户端获取连接时对闲置连接进行处理。
		// IdleTimeout:        5 * time.Minute,  // 闲置超时，默认5分钟，-1表示取消闲置超时检查
		// MaxConnAge:         0 * time.Second,  // 连接存活时长，从创建开始计时，超过指定时长则关闭连接，默认为0，即不关闭存活时长较长的连接
		ConnMaxIdleTime: 5 * time.Minute, // 连接最大空闲时间，达到此时间则将连接关闭（等同于v8的 IdleTimeout）
		ConnMaxLifetime: 0 * time.Second, // 连接存活时长，从创建开始计时，超过指定时长则关闭连接，默认为0，即不关闭存活时长较长的连接（等同于v8的 MaxConnAge）

		// 重试策略：命令执行失败时的重试策略
		// MaxRetries:      0,                      // 命令执行失败时，最多重试多少次，默认为0即不重试
		// MinRetryBackoff: 8 * time.Millisecond,   // 每次计算重试间隔时间的下限，默认8毫秒，-1表示取消间隔
		// MaxRetryBackoff: 512 * time.Millisecond, // 每次计算重试间隔时间的上限，默认512毫秒，-1表示取消间隔

		// 可自定义连接函数
		//Dialer: func(ctx context.Context, network, addr string) (net.Conn, error) {
		//	netDialer := &net.Dialer{
		//		Timeout:   5 * time.Second,
		//		KeepAlive: 5 * time.Minute,
		//	}
		//	return netDialer.Dial(network, addr)
		//},

		// 钩子函数
		// 仅当客户端执行命令时需要从连接池获取连接时，如果连接池需要新建连接时则会调用此钩子函数
		OnConnect: func(ctx context.Context, cn *redis.Conn) error {
			logging.Infof("conn success,%v\n", cn)
			return nil
		},
	})

	// ping pong
	pong, err := client.Ping(context.Background()).Result()
	logging.Infof("PING Result：pong=%v, err=%v", pong, err) // Output: PONG <nil>
	if pong != "PONG" {
		logging.Errorf("NewClient res=%v,err=%v\n", pong, err)
		return nil, errors.New("no response for the PING")
	}

	return client, nil
}

func printRedisPool(stats *redis.PoolStats) {
	fmt.Printf("Hits=%d Misses=%d Timeouts=%d TotalConns=%d IdleConns=%d StaleConns=%d\n",
		stats.Hits, stats.Misses, stats.Timeouts, stats.TotalConns, stats.IdleConns, stats.StaleConns)
}

func printRedisOption(opt *redis.Options) {
	fmt.Printf("Network=%v\n", opt.Network)
	fmt.Printf("Addr=%v\n", opt.Addr)
	fmt.Printf("Password=%v\n", opt.Password)
	fmt.Printf("DB=%v\n", opt.DB)
	fmt.Printf("MaxRetries=%v\n", opt.MaxRetries)
	fmt.Printf("MinRetryBackoff=%v\n", opt.MinRetryBackoff)
	fmt.Printf("MaxRetryBackoff=%v\n", opt.MaxRetryBackoff)
	fmt.Printf("DialTimeout=%v\n", opt.DialTimeout)
	fmt.Printf("ReadTimeout=%v\n", opt.ReadTimeout)
	fmt.Printf("WriteTimeout=%v\n", opt.WriteTimeout)
	fmt.Printf("PoolSize=%v\n", opt.PoolSize)
	fmt.Printf("MinIdleConns=%v\n", opt.MinIdleConns)
	//fmt.Printf("MaxConnAge=%v\n", opt.MaxConnAge)
	fmt.Printf("PoolTimeout=%v\n", opt.PoolTimeout)
	//fmt.Printf("IdleTimeout=%v\n", opt.IdleTimeout)
	//fmt.Printf("IdleCheckFrequency=%v\n", opt.IdleCheckFrequency)
	fmt.Printf("TLSConfig=%v\n", opt.TLSConfig)
}
