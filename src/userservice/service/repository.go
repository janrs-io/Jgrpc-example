package service

import (
	authV1 "authservice/genproto/go/v1"
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"time"

	userV1 "userservice/genproto/go/v1"
	"userservice/migrate"
	"userservice/pkg"
)

// Repository Repository
type Repository struct {
	db         *gorm.DB
	redis      *redis.Client
	authClient authV1.AuthServiceClient
}

// NewRepository New Repository
func NewRepository(db *gorm.DB, redis *redis.Client, authClient authV1.AuthServiceClient) *Repository {
	return &Repository{
		db:         db,
		redis:      redis,
		authClient: authClient,
	}
}

// UserModel User model
func (r *Repository) UserModel() *gorm.DB {
	return r.db.Table("user")
}

// IsUsernameExists Check is username exists
func (r *Repository) IsUsernameExists(username string) (bool, error) {

	result := r.UserModel().Where("username = ?", username).Find(&migrate.User{}).Limit(1)

	if result.Error != nil {
		return false, result.Error
	}
	if result.RowsAffected > 0 {
		return true, nil
	}
	return false, nil

}

// Register Register a new user
func (r *Repository) Register(request *userV1.RegisterRequest) (bool, error) {

	encrypt := pkg.NewEncrypt()
	password, err := encrypt.EncryptPWD([]byte(request.Password))
	if err != nil {
		return false, err
	}
	user := &migrate.User{}
	user.Username = request.Username
	user.Password = password
	user.CreateTime = time.Now().Unix()
	user.UpdateTime = time.Now().Unix()

	result := r.UserModel().Create(&user)
	if result.Error != nil {
		return false, result.Error
	}
	return true, nil

}

// Login User login
func (r *Repository) Login(ctx context.Context, request *userV1.LoginRequest) (*migrate.User, error) {

	user := &migrate.User{}

	// 查询用户名是否存在
	result := r.UserModel().Where("username = ?", request.Username).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("账号或密码错误")
	}

	// 验证密码
	encrypt := pkg.NewEncrypt()
	if !encrypt.ValidatePWD(user.Password, []byte(request.Password)) {
		return nil, errors.New("账号或密码错误")
	}

	// 销毁旧的授权数据
	destroyAuthResp, err := r.authClient.DestroyAuth(ctx, &authV1.DestroyAuthRequest{
		AccessToken: user.AccessToken,
	})
	if err != nil || !destroyAuthResp.Success {
		return nil, errors.New("登录失败")
	}

	// 创建 accessToken 以及设置 token 过期时间
	accessToken, err := encrypt.UniqueStr()
	if err != nil {
		return nil, err
	}

	// 创建新的用户登录授权信息
	var duration int64 = 7 * 24 * 60 * 60
	user.AccessToken = accessToken
	user.AccessTokenExpireTime = time.Now().Unix() + duration
	user.UpdateTime = time.Now().Unix()

	// 注册授权到 auth service
	registerAuthResp, err := r.authClient.RegisterAuth(ctx, &authV1.RegisterAuthRequest{
		AccessToken: accessToken,
		Duration:    duration,
	})
	if err != nil {
		return nil, err
	}
	if !registerAuthResp.Success {
		return nil, errors.New("注册授权失败")
	}

	// 保存新的登录数据到数据库
	result = r.UserModel().Save(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil

}

// Logout user logout
func (r *Repository) Logout() (bool, error) {
	return true, nil
}
