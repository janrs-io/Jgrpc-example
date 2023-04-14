package server

import (
	"context"
	"github.com/oklog/run"
	"os"

	"github.com/go-kit/log"
	"github.com/redis/go-redis/v9"

	"authservice/config"
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

// NewLogger 实例化 logger 组件
func NewLogger() log.Logger {
	return log.NewLogfmtLogger(os.Stderr)
}

// NewRunGroup 实例化 run.Group 组件
func NewRunGroup() *run.Group {
	return &run.Group{}
}
