package service

import (
	"context"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"strings"

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

// Register 用户注册
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

// Login 用户登录
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

// Logout 用户退出登录
func (s *Server) Logout(ctx context.Context, empty *emptypb.Empty) (*emptypb.Empty, error) {

	resp := empty
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return resp, status.Error(codes.Aborted, "退出失败，获取 context 上下文失败")
	}
	accessTokenMD := md.Get("authorization")
	if len(accessTokenMD) == 0 || len(accessTokenMD[0]) == 0 {
		return resp, status.Error(codes.Aborted, "退出登录失败，获取 access token 失败")
	}
	// 截取 bearer 后面的 access token 字符
	accessToken := accessTokenMD[0]
	if len(accessToken) > 7 {
		bearer := accessToken[0:7]
		if strings.ToLower(bearer) == "bearer " {
			accessToken = accessToken[7:]
		}
	}
	result, err := s.repo.Logout(ctx, accessToken)
	if err != nil {
		return resp, err
	}
	if !result {
		return resp, errors.New("退出登录失败")
	}
	return resp, nil

}

// Info 用户用户详情
func (s *Server) Info(ctx context.Context, empty *emptypb.Empty) (*userV1.InfoResponse, error) {

	infoResp := &userV1.InfoResponse{}
	return infoResp, nil

}
