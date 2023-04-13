package pkg

import (
	"context"
	"github.com/redis/go-redis/v9"

	"userservice/config"
)

// Redis Service's redis component
type Redis struct {
}

// NewRedis Initial service's Redis
func NewRedis(c *config.Config) *redis.Client {

	rdb := redis.NewClient(&redis.Options{
		Addr:         c.Redis.Host + c.Redis.Port,
		Username:     c.Redis.Username,
		Password:     c.Redis.Password,
		DB:           c.Redis.Database,
		DialTimeout:  c.Redis.DialTimeout,
		ReadTimeout:  c.Redis.ReadTimeout,
		WriteTimeout: c.Redis.WriteTimeout,
		PoolSize:     c.Redis.PoolSize,
		PoolTimeout:  c.Redis.PoolTimeout,
	})

	err := rdb.Ping(context.Background()).Err()
	if err != nil {
		panic("redis connect failed [ERROR]=> " + err.Error())
	}
	return rdb

}
