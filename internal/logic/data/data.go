package data

import (
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	"github.com/zhixunjie/im-fun/internal/logic/conf"
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/generate/query"
	"github.com/zhixunjie/im-fun/pkg/gomysql"
	"github.com/zhixunjie/im-fun/pkg/goredis"
	"github.com/zhixunjie/im-fun/pkg/kafka"
	"github.com/zhixunjie/im-fun/pkg/logging"
	"gorm.io/gorm"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewContactRepo, NewMessageRepo, NewData)

type Data struct {
	conf          *conf.Config
	RedisClient   *redis.Client
	MySQLClient   *gorm.DB
	Db            *query.Query
	KafkaProducer kafka.SyncProducer
}

func NewData(c *conf.Config) *Data {
	redisConf := c.Redis[0]
	mysqlConf := c.MySQL[0]

	// kafka producer
	kafkaProducer, err := kafka.NewSyncProducer(&c.Kafka[0])
	if err != nil {
		logging.Errorf("kafka.NewSyncProducer,err=%v", err)
		panic(err)
	}

	// redis
	redisClient, err := goredis.CreatePool(redisConf.Addr, redisConf.Auth, 0)
	if err != nil {
		logging.Errorf("redisClient,err=%v", err)
		panic(err)
	}

	// mysql
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
	query.SetDefault(mysqlClient)

	return &Data{
		conf:          c,
		KafkaProducer: kafkaProducer,
		RedisClient:   redisClient,
		MySQLClient:   mysqlClient,
		Db:            query.Q,
	}
}
