package dao

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"runtime"
)

var (
	RedisClient *redis.Client
	MySQLClient *gorm.DB
)

func InitDao() {
	RedisClient = CreateClient("127.0.0.1:6379", "")
	MySQLClient = CreateMySqlConn("127.0.0.1", 3306, "root", "", "im")
}

func CreateMySqlConn(ip string, port int, userName string, password string, database string) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=true", userName, password, ip, port, database)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	return db
}

// CreateClient 创建Redis客户端
func CreateClient(ip, password string) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     ip,
		Password: password,
		DB:       0,
	})

	// ping pong
	pong, err := client.Ping(context.Background()).Result()
	if err != nil {
		fmt.Printf("Redis Ping fail,err=%v", err)
		runtime.Goexit()
	}

	// fmt.Println("PING Result：", pong, err) // Output: PONG <nil>
	if pong != "PONG" {
		fmt.Printf("NewClient res=%v,err=%v\n", pong, err)
	}

	return client
}
