package service

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	userV1 "userservice/genproto/go/v1"
)

// Server Server struct
type Server struct {
	userV1.UnimplementedUserServiceServer
	userClient userV1.UserServiceClient
	repo       *Repository
}

// NewServer New service grpc server
func NewServer(repo *Repository, userClient userV1.UserServiceClient) userV1.UserServiceServer {
	return &Server{
		repo:       repo,
		userClient: userClient,
	}
}

// Register user register
func (s *Server) Register(ctx context.Context, req *userV1.RegisterRequest) (*emptypb.Empty, error) {

	isExists, err := s.repo.IsUsernameExists(req.Username)
	if err != nil {
		return &emptypb.Empty{}, status.Error(codes.Internal, err.Error())
	}
	if isExists {
		return &emptypb.Empty{}, status.Error(codes.FailedPrecondition, "用户名已存在")
	}

	_, err = s.repo.Register(req)
	if err != nil {
		return &emptypb.Empty{}, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil

}

// Login user login
func (s *Server) Login(ctx context.Context, req *userV1.LoginRequest) (*userV1.LoginResponse, error) {

	result, err := s.repo.Login(ctx, req)
	loginResp := &userV1.LoginResponse{}

	if err != nil {
		return loginResp, status.Error(codes.FailedPrecondition, "账号或密码错误")
	}

	// 返回数据
	loginResp.AccessToken = result.AccessToken
	loginResp.Username = req.Username
	loginResp.ExpireIn = result.AccessTokenExpireTime

	return loginResp, nil
}

// Logout user logout
func (s *Server) Logout(ctx context.Context, empty *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

// Info Get user info
func (s *Server) Info(ctx context.Context, empty *emptypb.Empty) (*userV1.InfoResponse, error) {

	_ = s.userClient
	infoResp := &userV1.InfoResponse{}
	infoResp.Username = "test============="
	return infoResp, nil

}
