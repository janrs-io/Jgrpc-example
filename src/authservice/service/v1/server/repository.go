package serverV1

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"time"

	"authservice/config"
)

// Repository Repository
type Repository struct {
	redis *redis.Client
	conf  *config.Config
	trace *sdktrace.TracerProvider
}

// NewRepository New Repository
func NewRepository(
	redis *redis.Client,
	conf *config.Config,
	trace *sdktrace.TracerProvider,
) *Repository {
	return &Repository{
		redis: redis,
		conf:  conf,
		trace: trace,
	}
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

// IsWhiteListApi 判断请求的接口是否在接口白名单内
func (r *Repository) IsWhiteListApi(api any) bool {

	for _, v := range r.conf.WhiteList.Api {
		if api == v {
			return true
		}
	}
	return false

}
