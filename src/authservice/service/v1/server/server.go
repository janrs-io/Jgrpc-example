package serverV1

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	authPBV1 "authservice/genproto/go/v1"
)

// Server Server struct
type Server struct {
	authPBV1.UnimplementedAuthServiceServer
	repo       *Repository
	authClient authPBV1.AuthServiceClient
}

// NewServer New Server
func NewServer(repo *Repository, authClient authPBV1.AuthServiceClient) authPBV1.AuthServiceServer {
	return &Server{
		repo:       repo,
		authClient: authClient,
	}
}

// RegisterAuth 用户登录后注册新的授权
func (s *Server) RegisterAuth(ctx context.Context, req *authPBV1.RegisterAuthRequest) (*authPBV1.RegisterAuthResponse, error) {

	err := s.repo.RegisterAuthentication(ctx, req.AccessToken, req.Duration)
	resp := &authPBV1.RegisterAuthResponse{}
	resp.Success = false
	if err != nil {
		return resp, status.Error(codes.Aborted, "注册授权失败，错误："+err.Error())
	}
	resp.Success = true

	return resp, nil

}

// GetAuth 获取授权
func (s *Server) GetAuth(ctx context.Context, req *authPBV1.GetAuthRequest) (*authPBV1.GetAuthResponse, error) {

	resp := &authPBV1.GetAuthResponse{}
	resp.Success = false
	err := s.repo.GetAuthentication(ctx, req.AccessToken, req.Duration)
	if err != nil {
		return resp, status.Error(codes.Aborted, "获取授权失败，错误："+err.Error())
	}
	resp.Success = true
	return resp, nil

}

// DestroyAuth 销毁授权数据
func (s *Server) DestroyAuth(ctx context.Context, req *authPBV1.DestroyAuthRequest) (*authPBV1.DestroyAuthResponse, error) {

	resp := &authPBV1.DestroyAuthResponse{}
	resp.Success = false
	if err := s.repo.DestroyAuthentication(ctx, req.AccessToken); err != nil {
		return resp, status.Error(codes.Aborted, "销毁授权失败。错误："+err.Error())
	}
	resp.Success = true
	return resp, nil

}

// IsApiWhiteList 查询请求的接口地址是否是在白名单内
func (s *Server) IsApiWhiteList(_ context.Context, req *authPBV1.IsApiWhiteListRequest) (*authPBV1.IsApiWhiteListResponse, error) {

	resp := &authPBV1.IsApiWhiteListResponse{}
	resp.Success = s.repo.IsWhiteListApi(req.FullMethodName)
	return resp, nil

}