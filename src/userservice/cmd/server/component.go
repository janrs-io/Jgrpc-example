package server

import (
	"context"
	"os"
	"strconv"

	"github.com/go-kit/log"
	"github.com/oklog/run"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"userservice/config"
)

// NewRedis 实例化 redis 组件
func NewRedis(conf *config.Config) *redis.Client {

	rdb := redis.NewClient(&redis.Options{
		Addr:         conf.Redis.Host + conf.Redis.Port,
		Username:     conf.Redis.Username,
		Password:     conf.Redis.Password,
		DB:           conf.Redis.Database,
		DialTimeout:  conf.Redis.DialTimeout,
		ReadTimeout:  conf.Redis.ReadTimeout,
		WriteTimeout: conf.Redis.WriteTimeout,
		PoolSize:     conf.Redis.PoolSize,
		PoolTimeout:  conf.Redis.PoolTimeout,
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		panic("redis connect failed [ERROR]=> " + err.Error())
	}
	return rdb

}

// NewLogger 实例化 go-kit/log 组件
func NewLogger() log.Logger {
	return log.NewLogfmtLogger(os.Stderr)
}

// NewRunGroup 实例化 run.Group 组件
func NewRunGroup() *run.Group {
	return &run.Group{}
}

// NewDB 实例化数据库组件
func NewDB(conf *config.Config) *gorm.DB {

	switch conf.Database.Driver {
	case "mysql":
		return newMysqlDB(conf)
	default:
		return newMysqlDB(conf)
	}

}

// newMysqlDB 初始化 mysql 连接
func newMysqlDB(conf *config.Config) *gorm.DB {

	// Database configuration
	dbConf := conf.Database
	if dbConf.Database == "" {
		panic("database config is empty.")
	}

	// Database connection dsn
	dsn := dbConf.UserName + ":" +
		dbConf.Password +
		"@tcp(" + dbConf.Host + ":" + strconv.Itoa(dbConf.Port) + ")/" + dbConf.Database +
		"?charset=" + dbConf.Charset + "&parseTime=True&loc=Local"

	// Set Config
	mysqlConfig := mysql.Config{
		DSN:                       dsn,   // DSN data source name
		DefaultStringSize:         255,   // Default length of the string type field
		DisableDatetimePrecision:  true,  // Disable datetime precision, not supported on databases prior to MySQL 5.6
		DontSupportRenameIndex:    true,  // Renaming indexes is done by deleting and creating new ones.
		DontSupportRenameColumn:   true,  // Rename columns with `change`, renaming columns is not supported in databases prior to MySQL 8 and MariaDB
		SkipInitializeWithVersion: false, // Automatic configuration based on version
	}

	// New mysql with config
	newMysql := mysql.New(mysqlConfig)

	// Connect mysql
	conn, err := gorm.Open(
		newMysql,
		&gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true, // Disable automatic foreign key creation constraints
			SkipDefaultTransaction:                   true, // Close global open transactions
		})

	if err != nil {
		panic("mysql connect failed [ERROR]=> " + err.Error())
	}

	sqlDB, _ := conn.DB()
	sqlDB.SetMaxIdleConns(dbConf.MaxIdleCons)
	sqlDB.SetMaxOpenConns(dbConf.MaxOpenCons)

	return conn

}
