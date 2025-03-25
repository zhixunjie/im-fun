package data

import (
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	"github.com/zhixunjie/im-fun/internal/logic/conf"
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/generate/model"
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/generate/query"
	"github.com/zhixunjie/im-fun/pkg/gomysql"
	"github.com/zhixunjie/im-fun/pkg/goredis"
	"github.com/zhixunjie/im-fun/pkg/kafka"
	"github.com/zhixunjie/im-fun/pkg/logging"
	"github.com/zhixunjie/im-fun/pkg/utils"
	"gorm.io/gorm"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewContactRepo, NewMessageRepo, NewData)

type Data struct {
	conf          *conf.Config
	RedisClient   *redis.Client      // Redis数据库对象
	MySQLDb       *query.Query       // MySQL的默认数据库对象
	KafkaProducer kafka.SyncProducer // Kafka生产者对象
}

func NewData(c *conf.Config) *Data {
	// kafka producer
	kafkaProducer, err := kafka.NewSyncProducer(&c.Kafka[0])
	if err != nil {
		logging.Errorf("kafka.NewSyncProducer,err=%v", err)
		panic(err)
	}

	// redis
	var redisClient *redis.Client
	if len(c.Redis) > 0 {
		redisConf := c.Redis[0]
		redisClient, err = goredis.CreatePool(redisConf.Addr, redisConf.Auth, 0)
		if err != nil {
			logging.Errorf("goredis.CreatePool,err=%v", err)
			panic(err)
		}
	}

	// mysql cluster
	if len(c.MySQLCluster) > 0 {
		var defaultDb *gorm.DB
		defaultDb, err = gomysql.InitMysqlCluster(c.MySQLCluster)
		if err != nil {
			logging.Errorf("gomysql.InitMysqlCluster,err=%v", err)
			panic(err)
		}
		if defaultDb == nil {
			logging.Errorf("gomysql.InitMysqlCluster，need to set default")
			panic(err)
		}
		// set gorm gen default db
		query.SetDefault(defaultDb)
	}
	d := &Data{
		conf:          c,
		KafkaProducer: kafkaProducer,
		RedisClient:   redisClient,
		MySQLDb:       query.Q,
	}
	// create or drop table
	d.CreateOrDrop()

	return d
}

func (d *Data) Master(dbName string) *query.Query {
	return query.Use(gomysql.Master(dbName))
}
func (d *Data) Slave(dbName string) *query.Query {
	return query.Use(gomysql.Slave(dbName))
}

func (d *Data) CreateOrDrop() {
	dbNames := d.conf.MySQLCluster
	for _, item := range dbNames {
		if item.IsDefault {
			continue
		}
		utils.CreateOrDrop(gomysql.Master(item.Name), "create", model.TableNameMessage, int64(model.DBNum()))
		utils.CreateOrDrop(gomysql.Master(item.Name), "create", model.TableNameContact, int64(model.DBNum()))
	}
}
