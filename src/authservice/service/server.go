package service

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	authV1 "authservice/genproto/go/v1"
)

// Server Server struct
type Server struct {
	authV1.UnimplementedAuthServiceServer
	repo       *Repository
	authClient authV1.AuthServiceClient
}

// NewServer New Server
func NewServer(repo *Repository, authClient authV1.AuthServiceClient) authV1.AuthServiceServer {
	return &Server{
		repo:       repo,
		authClient: authClient,
	}
}

// RegisterAuth 用户登录后注册新的授权
func (s *Server) RegisterAuth(ctx context.Context, req *authV1.RegisterAuthRequest) (*authV1.RegisterAuthResponse, error) {

	err := s.repo.RegisterAuthentication(ctx, req.AccessToken, req.Duration)
	resp := &authV1.RegisterAuthResponse{}
	resp.Success = false
	if err != nil {
		return resp, status.Error(codes.Aborted, "注册授权失败，错误："+err.Error())
	}
	resp.Success = true

	return resp, nil

}

// GetAuth 获取授权
func (s *Server) GetAuth(ctx context.Context, req *authV1.GetAuthRequest) (*authV1.GetAuthResponse, error) {

	err := s.repo.GetAuthentication(ctx, req.AccessToken, req.Duration)
	resp := &authV1.GetAuthResponse{}
	resp.Success = false
	if err != nil {
		return resp, status.Error(codes.Aborted, "获取授权失败，错误："+err.Error())
	}
	resp.Success = true
	return resp, nil

}

// DestroyAuth 销毁授权数据
func (s *Server) DestroyAuth(ctx context.Context, req *authV1.DestroyAuthRequest) (*authV1.DestroyAuthResponse, error) {

	resp := &authV1.DestroyAuthResponse{}
	resp.Success = false
	if err := s.repo.DestroyAuthentication(ctx, req.AccessToken); err != nil {
		return resp, status.Error(codes.Aborted, "销毁授权失败。错误："+err.Error())
	}
	resp.Success = true
	return resp, nil

}
