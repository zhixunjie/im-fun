package dao

import (
	"github.com/go-redis/redis/v8"
	"github.com/zhixunjie/im-fun/internal/logic/conf"
	"github.com/zhixunjie/im-fun/pkg/gomysql"
	"github.com/zhixunjie/im-fun/pkg/goredis"
	"github.com/zhixunjie/im-fun/pkg/kafka"
	"github.com/zhixunjie/im-fun/pkg/logging"
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

	redisClient, err := goredis.CreatePool(redisConf.Addr, redisConf.Auth, 0)
	if err != nil {
		logging.Errorf("redisClient,err=%v", err)
		panic(err)
	}

	mysqlClient, err := gomysql.CreatePool(&gomysql.Config{
		Addr:     mysqlConf.Addr,
		Port:     mysqlConf.Port,
		UserName: mysqlConf.UserName,
		Password: mysqlConf.Password,
		Database: mysqlConf.DbName,
	})
	if err != nil {
		logging.Errorf("mysqlClient,err=%v", err)
		panic(err)
	}

	d := &Dao{
		conf:          c,
		KafkaProducer: kafkaProducer,
		RedisClient:   redisClient,
		MySQLClient:   mysqlClient,
		redisExpire:   int32(time.Duration(redisConf.Expire) / time.Second),
	}
	return d
}
