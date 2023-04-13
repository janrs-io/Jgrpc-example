package service

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

// Repository Repository
type Repository struct {
	redis *redis.Client
}

// NewRepository New Repository
func NewRepository(redis *redis.Client) *Repository {
	return &Repository{redis: redis}
}

// RegisterAuthentication 注册授权数据到 redis 缓存
func (r *Repository) RegisterAuthentication(ctx context.Context, accessToken string, duration int64) (err error) {
	return r.redis.Set(ctx, accessToken, "", time.Second*time.Duration(duration)).Err()
}

// GetAuthentication 获取授权数据
func (r *Repository) GetAuthentication(ctx context.Context, accessToken string, duration int64) error {

	exists := r.redis.Exists(ctx, accessToken)
	if exists.Err() != nil {
		return exists.Err()
	}
	if exists.Val() <= 0 {
		return errors.New("access token不存在")
	}

	return r.RefreshAccessTokenExpireTime(accessToken, duration)

}

// RefreshAccessTokenExpireTime 刷新授权数据过期时间
func (r *Repository) RefreshAccessTokenExpireTime(accessToken string, duration int64) error {

	expire := r.redis.Expire(context.Background(), accessToken, time.Second*time.Duration(duration))
	if expire.Err() != nil {
		return expire.Err()
	}
	if !expire.Val() {
		return errors.New("刷新授权缓存数据失败")
	}
	return nil

}

// DestroyAuthentication 销毁授权数据
func (r *Repository) DestroyAuthentication(ctx context.Context, accessToken string) (err error) {
	return r.redis.Del(ctx, accessToken).Err()
}
