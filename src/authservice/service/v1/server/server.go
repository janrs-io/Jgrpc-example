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
func (s *Server) RegisterAuth(ctx context.Context, req *authPBV1.RegisterAuthRequest) (*authPBV1.Response, error) {

	err := s.repo.RegisterAuthentication(ctx, req.AccessToken, req.Duration)
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, "注册授权失败，错误："+err.Error())
	}
	return &authPBV1.Response{}, nil

}

// GetAuth 获取授权
func (s *Server) GetAuth(ctx context.Context, req *authPBV1.GetAuthRequest) (*authPBV1.Response, error) {

	err := s.repo.GetAuthentication(ctx, req.AccessToken, req.Duration)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "获取授权失败，错误："+err.Error())
	}
	return &authPBV1.Response{}, nil

}

// DestroyAuth 销毁授权数据
func (s *Server) DestroyAuth(ctx context.Context, req *authPBV1.DestroyAuthRequest) (*authPBV1.Response, error) {

	if err := s.repo.DestroyAuthentication(ctx, req.AccessToken); err != nil {
		return nil, status.Error(codes.FailedPrecondition, "销毁授权失败。错误："+err.Error())
	}
	return &authPBV1.Response{}, nil

}

// IsApiWhiteList 查询请求的接口地址是否是在白名单内
func (s *Server) IsApiWhiteList(_ context.Context, req *authPBV1.IsApiWhiteListRequest) (*authPBV1.Response, error) {

	if s.repo.IsWhiteListApi(req.FullMethodName) {
		return &authPBV1.Response{}, nil
	}
	return nil, status.Error(codes.FailedPrecondition, "非白名单接口")

}
