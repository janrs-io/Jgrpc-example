package serverV1

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	Jgrpc_otelspan "github.com/janrs-io/Jgrpc-otel-span"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	orderPBV1 "orderservice/genproto/go/v1"
	productPBV1 "productservice/genproto/go/v1"
	"time"
	"userservice/config"
	userPBV1 "userservice/genproto/go/v1"
	"userservice/service/model"
)

// Repository Repository
type Repository struct {
	mysqlDB       *gorm.DB
	redis         *redis.Client
	orderClient   orderPBV1.OrderServiceClient
	productClient productPBV1.ProductServiceClient
	userClient    userPBV1.UserServiceClient
	span          *Jgrpc_otelspan.OtelSpan
	conf          *config.Config
}

// NewRepository New Repository
func NewRepository(
	mysqlDB *gorm.DB,
	redis *redis.Client,
	orderClient orderPBV1.OrderServiceClient,
	productClient productPBV1.ProductServiceClient,
	userClient userPBV1.UserServiceClient,
	span *Jgrpc_otelspan.OtelSpan,
	conf *config.Config,
) *Repository {
	return &Repository{
		mysqlDB:       mysqlDB,
		redis:         redis,
		orderClient:   orderClient,
		productClient: productClient,
		userClient:    userClient,
		span:          span,
		conf:          conf,
	}
}

// UserModel User model
func (r *Repository) UserModel() *gorm.DB {
	return r.mysqlDB.Table("user")
}

// IsUsernameExists 查询用户名是否存在
func (r *Repository) IsUsernameExists(ctx context.Context, username string) (bool, error) {

	_, span := r.span.Record(ctx, r.conf.Trace.TracerName)
	defer span.End()

	result := r.UserModel().Where("username = ?", username).Find(&model.User{}).Limit(1)

	if result.Error != nil {
		return false, r.span.Error(span, result.Error.Error())
	}
	if result.RowsAffected > 0 {
		return true, nil
	}
	return false, nil

}

// Register 注册一个新用户
func (r *Repository) Register(ctx context.Context, request *userPBV1.RegisterRequest) (bool, error) {

	_, span := r.span.Record(ctx, r.conf.Trace.TracerName)
	defer span.End()

	password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return false, r.span.Error(span, err.Error())
	}
	user := &model.User{}
	user.Username = request.Username
	user.Password = string(password)
	user.CreateTime = time.Now().Unix()
	user.UpdateTime = time.Now().Unix()

	result := r.UserModel().Create(&user)
	if result.Error != nil {
		return false, r.span.Error(span, result.Error.Error())
	}
	return true, nil

}

// Login 用户登录
func (r *Repository) Login(ctx context.Context, request *userPBV1.LoginRequest) (*model.User, error) {

	_, span := r.span.Record(ctx, r.conf.Trace.TracerName)
	defer span.End()

	user := &model.User{}

	// 查询用户名是否存在
	result := r.UserModel().Where("username = ?", request.Username).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, r.span.Error(span, "账号或密码不存在")
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		return nil, r.span.Error(span, "账号或密码错误")
	}

	// 创建 accessToken 以及设置 token 过期时间
	accessToken := uuid.New().String()

	// 创建新的用户登录授权信息
	var duration int64 = 7 * 24 * 60 * 60
	user.AccessToken = accessToken
	user.AccessTokenExpireTime = time.Now().Unix() + duration
	user.UpdateTime = time.Now().Unix()

	// 保存新的登录数据到数据库
	result = r.UserModel().Save(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil

}

// Logout 用户退出登录
// 用户退出登录后调用 auth 服务删除授权数据并且删除数据库 access token
func (r *Repository) Logout(ctx context.Context, accessToken string) (bool, error) {

	_, span := r.span.Record(ctx, r.conf.Trace.TracerName)
	defer span.End()

	// 删除数据库 access token
	user := &model.User{}
	result := r.UserModel().Where("access_token = ?", accessToken).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, r.span.Error(span, result.Error.Error())
	}
	user.AccessToken = ""
	user.AccessTokenExpireTime = 0
	result = r.UserModel().Save(&user)
	if result.Error != nil {
		return false, r.span.Error(span, result.Error.Error())
	}
	return true, nil

}

// Info 获取用户信息
func (r *Repository) Info(ctx context.Context, accessToken string) (*model.User, error) {

	_, span := r.span.Record(ctx, r.conf.Trace.TracerName)
	defer span.End()

	user := &model.User{}
	result := r.UserModel().Where("access_token = ?", accessToken).First(&user)
	if result.Error != nil {
		return nil, r.span.Error(span, result.Error.Error())
	}
	return user, nil

}

// UserInfo 获取用户详情
func (r *Repository) UserInfo(ctx context.Context) (*userPBV1.UserDetail_Detail, error) {

	_, span := r.span.Record(ctx, r.conf.Trace.TracerName)
	defer span.End()

	user := &model.User{}

	accessToken, err := auth.AuthFromMD(ctx, "Bearer")
	if err != nil {
		return nil, r.span.Error(span, err.Error())
	}

	result := r.UserModel().Where("access_token = ?", accessToken).First(&user)
	if result.Error != nil {
		return nil, r.span.Error(span, result.Error.Error())
	}

	return &userPBV1.UserDetail_Detail{
		Id:                    user.ID,
		Username:              user.Username,
		Sex:                   user.Sex,
		IdNumber:              user.IDNumber,
		Email:                 user.Email,
		Phone:                 user.Phone,
		IsDisable:             user.IsDisable,
		AccessToken:           user.AccessToken,
		AccessTokenExpireTime: user.AccessTokenExpireTime,
		NickName:              user.NickName,
		RealName:              user.RealName,
		CreateTime:            user.CreateTime,
		UpdateTime:            user.UpdateTime,
	}, nil

}

// ProductInfo 获取产品详情
func (r *Repository) ProductInfo(ctx context.Context, request *userPBV1.OrderInfoRequest) (*productPBV1.Response, error) {

	ctx, span := r.span.Record(ctx, r.conf.Trace.TracerName)
	defer span.End()

	// 获取产品信息
	return r.productClient.Detail(ctx, &productPBV1.DetailRequest{Id: request.ProductId})

}

// OrderInfo 获取订单详情
func (r *Repository) OrderInfo(ctx context.Context, request *userPBV1.OrderInfoRequest) (*orderPBV1.Response, error) {

	ctx, span := r.span.Record(ctx, r.conf.Trace.TracerName)
	defer span.End()

	// 获取订单信息
	return r.orderClient.Detail(ctx, &orderPBV1.DetailRequest{Id: request.OrderId})

}
