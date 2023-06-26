package dao

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/zhixunjie/im-fun/internal/logic/conf"
	"github.com/zhixunjie/im-fun/pkg/kafka"
	"github.com/zhixunjie/im-fun/pkg/logging"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

type Dao struct {
	conf          *conf.Config
	RedisClient   *redis.Client
	MySQLClient   *gorm.DB
	KafkaProducer kafka.SyncProducer
	redisExpire   int32
}

func New(c *conf.Config) *Dao {
	redisConf := c.Redis[0]
	mysqlConf := c.MySQL[0]

	kafkaProducer, err := kafka.NewSyncProducer(&c.Kafka[0])
	if err != nil {
		logging.Errorf("kafka.NewSyncProducer,err=%v", err)
		panic(err)
	}

	d := &Dao{
		conf:          c,
		KafkaProducer: kafkaProducer,
		RedisClient:   CreateRedisPool(redisConf.Addr, redisConf.Auth),
		MySQLClient:   CreateMySqlClient(mysqlConf.Addr, mysqlConf.UserName, mysqlConf.Password, mysqlConf.DbName),
		redisExpire:   int32(time.Duration(redisConf.Expire) / time.Second),
	}
	return d
}

var (
	RedisClient *redis.Client
	MySQLClient *gorm.DB
)

//func InitDao() {
//	RedisClient = CreateRedisClient("127.0.0.1:6379", "")
//	MySQLClient = CreateMySqlClient("127.0.0.1:3306", "root", "", "im")
//}

func CreateMySqlClient(addr string, userName string, password string, database string) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true", userName, password, addr, database)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logging.Errorf("gorm.Open,err=%v", err)
		panic("failed to connect database")
	}

	return db
}

func CreateRedisPool(addr, password string) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,

		// 连接池容量及闲置连接数量
		// 连接池最大socket连接数，默认为4倍CPU数， 4 * runtime.NumCPU
		PoolSize: 16,
		// 在启动阶段创建指定数量的Idle连接，并长期维持idle状态的连接数不少于指定数量。
		MinIdleConns: 10,

		// 闲置连接检查包括IdleTimeout，MaxConnAge
		// IdleCheckFrequency: 60 * time.Second, // 闲置连接检查的周期，默认为1分钟，-1表示不做周期性检查，只在客户端获取连接时对闲置连接进行处理。
		// IdleTimeout:        5 * time.Minute,  // 闲置超时，默认5分钟，-1表示取消闲置超时检查
		// MaxConnAge:         0 * time.Second,  // 连接存活时长，从创建开始计时，超过指定时长则关闭连接，默认为0，即不关闭存活时长较长的连接

		// 命令执行失败时的重试策略
		// MaxRetries:      0,                      // 命令执行失败时，最多重试多少次，默认为0即不重试
		// MinRetryBackoff: 8 * time.Millisecond,   // 每次计算重试间隔时间的下限，默认8毫秒，-1表示取消间隔
		// MaxRetryBackoff: 512 * time.Millisecond, // 每次计算重试间隔时间的上限，默认512毫秒，-1表示取消间隔

		// 可自定义连接函数
		// Dialer: func() (net.Conn, error) {
		// 	netDialer := &net.Dialer{
		// 		Timeout:   5 * time.Second,
		// 		KeepAlive: 5 * time.Minute,
		// 	}
		// 	return netDialer.Dial("tcp", "127.0.0.1:6379")
		// },

		// 钩子函数
		// 仅当客户端执行命令时需要从连接池获取连接时，如果连接池需要新建连接时则会调用此钩子函数
		OnConnect: func(ctx context.Context, conn *redis.Conn) error {
			logging.Infof("conn=%v", conn)
			return nil
		},
	})

	// ping pong
	pong, err := client.Ping(context.Background()).Result()
	logging.Infof("PING Result：", pong, err) // Output: PONG <nil>
	if pong != "PONG" {
		logging.Errorf("NewClient res=%v,err=%v", pong, err)
		return nil
	}

	return client
}
