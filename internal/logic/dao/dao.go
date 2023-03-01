package dao

import (
	"context"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"github.com/zhixunjie/im-fun/internal/logic/conf"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"runtime"
)

type Dao struct {
	conf          *conf.Config
	RedisClient   *redis.Client
	MySQLClient   *gorm.DB
	KafkaProducer sarama.SyncProducer
}

func New(c *conf.Config) *Dao {
	redisConf := c.Redis[0]
	mysqlConf := c.MySQL[0]

	kafkaProducer, err := newKafkaProducer(&c.Kafka[0])
	if err != nil {
		panic(err)
	}

	d := &Dao{
		conf:          c,
		KafkaProducer: kafkaProducer,
		RedisClient:   CreateRedisClient(redisConf.Addr, redisConf.Auth),
		MySQLClient:   CreateMySqlClient(mysqlConf.Addr, mysqlConf.UserName, mysqlConf.Password, mysqlConf.DbName),
	}
	return d
}

var (
	RedisClient *redis.Client
	MySQLClient *gorm.DB
)

func InitDao() {
	RedisClient = CreateRedisClient("127.0.0.1:6379", "")
	MySQLClient = CreateMySqlClient("127.0.0.1:3306", "root", "", "im")
}

func CreateMySqlClient(addr string, userName string, password string, database string) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true", userName, password, addr, database)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	return db
}

// CreateRedisClient 创建Redis客户端
func CreateRedisClient(addr, password string) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})

	// ping pong
	pong, err := client.Ping(context.Background()).Result()
	if err != nil {
		logrus.Errorf("Redis Ping fail,err=%v", err)
		runtime.Goexit()
	}

	// fmt.Println("PING Result：", pong, err) // Output: PONG <nil>
	if pong != "PONG" {
		logrus.Errorf("NewClient res=%v,err=%v\n", pong, err)
	}

	return client
}
