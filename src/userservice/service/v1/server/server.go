package serverV1

import (
	"context"
	"errors"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"

	orderPBV1 "orderservice/genproto/go/v1"
	productPBV1 "productservice/genproto/go/v1"
	userPBV1 "userservice/genproto/go/v1"
)

// Server Server struct
type Server struct {
	userPBV1.UnimplementedUserServiceServer
	userClient    userPBV1.UserServiceClient
	orderClient   orderPBV1.OrderServiceClient
	productClient productPBV1.ProductServiceClient
	repo          *Repository
	logger        log.Logger
}

// NewServer New service grpc server
func NewServer(
	repo *Repository,
	logger log.Logger,
	userClient userPBV1.UserServiceClient,
	orderClient orderPBV1.OrderServiceClient,
	productClient productPBV1.ProductServiceClient,
) userPBV1.UserServiceServer {
	return &Server{
		repo:          repo,
		userClient:    userClient,
		logger:        logger,
		orderClient:   orderClient,
		productClient: productClient,
	}
}

// Register 用户注册
func (s *Server) Register(ctx context.Context, req *userPBV1.RegisterRequest) (*userPBV1.Response, error) {

	isExists, err := s.repo.IsUsernameExists(ctx, req.Username)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if isExists {
		return nil, status.Error(codes.FailedPrecondition, "用户名已存在")
	}
	_, err = s.repo.Register(ctx, req)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &userPBV1.Response{}, nil

}

// Login 用户登录
func (s *Server) Login(ctx context.Context, req *userPBV1.LoginRequest) (*userPBV1.Response, error) {

	result, err := s.repo.Login(ctx, req)
	loginResp := &userPBV1.LoginResponse{}

	if err != nil {
		_ = level.Error(s.logger).Log("msg", "用户登录失败，错误[1]："+err.Error())
		return nil, status.Error(codes.FailedPrecondition, "账号或密码错误")
	}

	// 返回数据
	loginResp.AccessToken = result.AccessToken
	loginResp.Username = req.Username
	loginResp.ExpireIn = result.AccessTokenExpireTime

	resp := &userPBV1.Response{}
	returnAnyData, err := anypb.New(loginResp)
	if err != nil {
		_ = level.Error(s.logger).Log("msg", "用户登录失败，错误[2]："+err.Error())
		return nil, status.Error(codes.Internal, "系统错误")
	}
	resp.ProtoAnyData = returnAnyData

	return resp, nil
}

// Logout 用户退出登录
func (s *Server) Logout(ctx context.Context, _ *emptypb.Empty) (*userPBV1.Response, error) {

	accessToken, err := auth.AuthFromMD(ctx, "Bearer")
	if err != nil {
		_ = level.Error(s.logger).Log("msg", "退出登录失败，错误[1]："+err.Error())
		return nil, status.Error(codes.FailedPrecondition, "退出登录失败")
	}
	result, err := s.repo.Logout(ctx, accessToken)
	if err != nil {
		return nil, err
	}
	if !result {
		return nil, status.Error(codes.FailedPrecondition, "退出登录失败")
	}
	return &userPBV1.Response{}, nil

}

// Info 获取用户信息
func (s *Server) Info(ctx context.Context, _ *emptypb.Empty) (*userPBV1.Response, error) {

	resp := &userPBV1.Response{}
	accessToken, err := auth.AuthFromMD(ctx, "Bearer")
	info, err := s.repo.Info(ctx, accessToken)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return resp, nil
		}
		_ = level.Error(s.logger).Log("msg", "获取用户信息失败，错误[2]："+err.Error())
		return nil, status.Error(codes.FailedPrecondition, "获取用户信息失败")
	}
	detail := &userPBV1.UserDetail_Detail{
		Id:                    info.ID,
		Username:              info.Username,
		Sex:                   info.Sex,
		IdNumber:              info.IDNumber,
		Email:                 info.Email,
		Phone:                 info.Phone,
		IsDisable:             info.IsDisable,
		AccessToken:           info.AccessToken,
		AccessTokenExpireTime: info.AccessTokenExpireTime,
		NickName:              info.NickName,
		RealName:              info.RealName,
		CreateTime:            info.CreateTime,
		UpdateTime:            info.UpdateTime,
	}
	anyData, err := anypb.New(detail)
	if err != nil {
		_ = level.Error(s.logger).Log("msg", "获取用户详情失败，错误[3]："+err.Error())
		return nil, status.Error(codes.FailedPrecondition, "获取详情失败")
	}
	resp.ProtoAnyData = anyData
	return resp, nil

}

// OrderInfo 获取订单详情
func (s *Server) OrderInfo(ctx context.Context, request *userPBV1.OrderInfoRequest) (*userPBV1.Response, error) {

	orderInfoResp := &userPBV1.OrderInfoResponse{}
	resp := &userPBV1.Response{}

	// 获取用户详情
	userInfo, err := s.repo.UserInfo(ctx)
	if err != nil {
		_ = level.Error(s.logger).Log("msg", "获取订单详情失败，错误[2]："+err.Error())
		return nil, status.Error(codes.FailedPrecondition, "获取订单详情失败")
	}

	// 获取订单详情
	orderInfo, err := s.repo.OrderInfo(ctx, request)
	if err != nil {
		_ = level.Error(s.logger).Log("msg", "获取订单详情失败，错误[3]："+err.Error())
		return nil, status.Error(codes.FailedPrecondition, "获取订单详情失败")
	}

	// 获取产品详情
	productInfo, err := s.repo.ProductInfo(ctx, request)
	if err != nil {
		_ = level.Error(s.logger).Log("msg", "获取订单详情失败，错误[4]："+err.Error())
		return nil, status.Error(codes.FailedPrecondition, "获取订单详情失败")
	}

	anyUserData, err := anypb.New(userInfo)
	if err != nil {
		_ = level.Error(s.logger).Log("msg", "获取订单详情失败，错误[5]："+err.Error())
		return nil, status.Error(codes.FailedPrecondition, "获取订单详情失败")
	}

	anyOrderData, err := anypb.New(orderInfo)
	if err != nil {
		_ = level.Error(s.logger).Log("msg", "获取订单详情失败，错误[6]："+err.Error())
		return nil, status.Error(codes.FailedPrecondition, "获取订单详情失败")
	}

	anyProductData, err := anypb.New(productInfo)

	if err != nil {
		_ = level.Error(s.logger).Log("msg", "获取订单详情失败，错误[7]："+err.Error())
		return nil, status.Error(codes.FailedPrecondition, "获取订单详情失败")
	}
	orderInfoResp.UserInfo = anyUserData
	orderInfoResp.ProductInfo = anyProductData
	orderInfoResp.OrderInfo = anyOrderData

	// 组合所有数据
	allInfo, err := anypb.New(orderInfoResp)
	if err != nil {
		_ = level.Error(s.logger).Log("msg", "获取订单详情失败，错误[8]："+err.Error())
		return nil, status.Error(codes.FailedPrecondition, "获取订单详情失败")
	}
	resp.ProtoAnyData = allInfo
	return resp, nil

}
