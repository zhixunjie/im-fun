package gomysql

import (
	"github.com/zhixunjie/im-fun/pkg/logging"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MySQLCluster struct {
	Name      string   `yaml:"name"`
	Master    string   `yaml:"master"`
	Slaves    []string `yaml:"slaves"`
	IsDefault bool     `yaml:"is_default"`
}

var dbMap map[string]map[dbType]*gorm.DB

type dbType string

const (
	master dbType = "master"
	slave  dbType = "slave"
)

func InitMysqlCluster(arr []MySQLCluster) (defaultDb *gorm.DB, err error) {
	dbMap = make(map[string]map[dbType]*gorm.DB)
	for _, item := range arr {
		var db *gorm.DB
		// create master
		if len(item.Master) > 0 {
			db, err = createPoolWithDsn(item.Master)
			if err != nil {
				return
			}
			if _, ok := dbMap[item.Name]; !ok {
				dbMap[item.Name] = make(map[dbType]*gorm.DB)
			}
			dbMap[item.Name][master] = db
		}
		// create slave
		if len(item.Slaves) > 0 {
			db, err = createPoolWithDsn(item.Slaves[0])
			if err != nil {
				return
			}
			if _, ok := dbMap[item.Name]; !ok {
				dbMap[item.Name] = make(map[dbType]*gorm.DB)
			}
			dbMap[item.Name][slave] = db
		}
		if item.IsDefault {
			defaultDb = db
		}
	}
	return
}

func Master(name string) *gorm.DB {
	if _, ok := dbMap[name]; !ok {
		return dbMap[name][master]
	}

	return nil
}

func Slave(name string) *gorm.DB {
	if _, ok := dbMap[name]; !ok {
		return dbMap[name][slave]
	}

	return nil
}

func createPoolWithDsn(dsn string) (db *gorm.DB, err error) {
	// set config
	var config = &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	}
	setLogger(config)

	// open connection
	db, err = gorm.Open(mysql.Open(dsn), config)
	if err != nil {
		logging.Error("failed to connect database")
		return
	}

	return
}

//func createPoolWithDsn(item MySQLCluster) (db *gorm.DB, err error) {
//	// set config
//	var config = &gorm.Config{
//		DisableForeignKeyConstraintWhenMigrating: true,
//	}
//	setLogger(config)
//
//	// open connection
//	db, err = gorm.Open(mysql.Open(item.Master), config)
//	if err != nil {
//		logging.Error("failed to connect database")
//		return
//	}
//	// open slave
//	var slaves []gorm.Dialector
//	for _, slave := range item.Slaves {
//		slaves = append(slaves, mysql.Open(slave))
//	}
//
//	err = db.Use(dbresolver.Register(dbresolver.Config{
//		Sources: []gorm.Dialector{
//			mysql.Open(item.Master),
//		},
//		Replicas: slaves,
//		// sources/replicas load balancing policy
//		Policy: dbresolver.RandomPolicy{},
//		// print sources/replicas mode in logger
//		TraceResolverMode: true,
//	}))
//
//	return
//}
