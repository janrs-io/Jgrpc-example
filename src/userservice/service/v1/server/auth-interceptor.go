package serverV1

import (
	"context"
	"github.com/go-kit/log/level"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	authPBV1 "authservice/genproto/go/v1"
)

// AuthFuncOverride 授权验证拦截器
func (s *Server) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {

	// 查询接口是否是白名单
	isApiWhiteListReq := &authPBV1.IsApiWhiteListRequest{}
	isApiWhiteListReq.FullMethodName = fullMethodName
	isApiWhiteListResp, err := s.authClient.IsApiWhiteList(ctx, isApiWhiteListReq)
	if err != nil {
		_ = level.Info(s.logger).Log("msg", "授权失败，错误："+err.Error())
		return ctx, status.Error(codes.Unauthenticated, "请先登录")
	}

	// 如果是在白名单内，直接跳过鉴权
	if isApiWhiteListResp.Success {
		return ctx, nil
	}

	//**********************不是白名单接口，继续鉴权

	token, err := auth.AuthFromMD(ctx, "Bearer")
	if err != nil {
		_ = level.Info(s.logger).Log("msg", "授权失败，错误：未传递 access token 参数。")
		return ctx, status.Error(codes.Unauthenticated, "请先登录")
	}
	getAuthReq := &authPBV1.GetAuthRequest{}
	getAuthReq.AccessToken = token
	getAuthReq.Duration = 7 * 24 * 60 * 60
	authAuthResp, err := s.authClient.GetAuth(ctx, getAuthReq)
	if err != nil {
		_ = level.Info(s.logger).Log("msg", "授权失败，错误："+err.Error())
		return ctx, status.Error(codes.Unauthenticated, "请先登录")
	}
	if !authAuthResp.Success {
		_ = level.Info(s.logger).Log("msg", "授权失败，未知错误。请排查错误。")
		return ctx, status.Error(codes.Unauthenticated, "请先登录")
	}
	return ctx, nil

}
