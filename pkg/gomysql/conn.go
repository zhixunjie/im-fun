package gomysql

import (
	"fmt"
	"github.com/zhixunjie/im-fun/pkg/logging"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

// https://gorm.io/zh_CN/docs/logger.html
func setLogger(config *gorm.Config) {
	// 1. 使用默认的logger: 日志级别为Info
	//config.Logger = logger.Default.LogMode(logger.Info)

	// 2. 使用新的logger:
	config.Logger = logger.New(
		// set writer: 将Stdout作为Writer
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		// set config
		logger.Config{
			SlowThreshold:             1 * time.Second, // 设定慢查询时间阈值为1秒
			LogLevel:                  logger.Info,     // 设置日志级别（小于这个值的日志都会被输出）
			IgnoreRecordNotFoundError: true,            // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  true,            // 彩色打印
		},
	)
}

type Config struct {
	Addr     string
	Port     string
	UserName string
	Password string
	Database string
}

func CreatePool(cfg *Config) (*gorm.DB, error) {
	// set config
	var config = &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	}
	setLogger(config)

	// open connection
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true&loc=Local", cfg.UserName, cfg.Password, cfg.Addr, cfg.Port, cfg.Database)
	db, err := gorm.Open(mysql.Open(dsn), config)
	if err != nil {
		logging.Error("failed to connect database")
		return nil, err
	}

	return db, nil
}

//func CreateSQLiteConn() *gorm.DB {
//	db, err := gorm.Open(sqlite.Open("test.Db"), &gorm.Config{})
//	if err != nil {
//		panic("failed to connect database")
//	}
//
//	return db
//}
