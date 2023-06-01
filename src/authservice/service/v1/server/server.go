package serverV1

import (
	"encoding/json"
	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	typev3 "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/redis/go-redis/v9"
	"google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"

	"authservice/config"
	"golang.org/x/net/context"
)

type Server struct {
	authv3.UnimplementedAuthorizationServer
	conf   *config.Config
	redis  *redis.Client
	repo   *Repository
	logger log.Logger
}

func NewServer(
	conf *config.Config,
	redis *redis.Client,
	repo *Repository,
	logger log.Logger,
) authv3.AuthorizationServer {
	return &Server{
		conf:   conf,
		redis:  redis,
		repo:   repo,
		logger: logger,
	}
}

var (
	UnauthorizedMsg = "没有权限"
	ForbiddenMsg    = "没有权限"
)

// Response 返回 HTTP Body 数据
type Response struct {
	Code int64    `json:"code"`
	Msg  string   `json:"msg"`
	Data struct{} `json:"data"`
}

// Check istio-grpc 外部鉴权方法
func (s *Server) Check(ctx context.Context, req *authv3.CheckRequest) (*authv3.CheckResponse, error) {
	attrs := req.GetAttributes()
	httpHeaders := attrs.GetRequest().GetHttp().GetHeaders()
	// 获取请求路径
	path, exists := httpHeaders[":path"]
	if !exists {
		_ = level.Info(s.logger).Log("msg", "获取不到 :path 字段")
		return s.Unauthorized(), nil
	}
	// 判断是否是白名单
	if s.repo.IsWhiteListApi(path) {
		return s.Allow(), nil
	}
	// 获取头部 token
	token, exists := httpHeaders["authorization"]
	duration := 7 * 24 * 60 * 60
	if !exists {
		_ = level.Info(s.logger).Log("msg", "未传递头部 authorization 字段")
		return s.Unauthorized(), nil
	}
	// 去除头部 "Bearer "字符串
	if len(token) <= 7 {
		_ = level.Info(s.logger).Log("msg", "authorization 数据格式错误。没有设置 Bearer 前缀")
		return s.Unauthorized(), nil
	}
	// 截取后面的 token 字符串
	token = token[7:]

	// 验证 token
	if err := s.repo.GetAuthentication(ctx, token, int64(duration)); err != nil {
		_ = level.Info(s.logger).Log("msg", "access token 不存在")
		return s.Unauthorized(), nil
	}
	return s.Allow(), nil
}

// Allow 通过鉴权。返回 200
func (s *Server) Allow() *authv3.CheckResponse {
	return &authv3.CheckResponse{
		Status: &status.Status{Code: int32(codes.OK)},
		HttpResponse: &authv3.CheckResponse_OkResponse{
			OkResponse: &authv3.OkHttpResponse{},
		},
	}
}

// Unauthorized Unauthorized 未授权 401
func (s *Server) Unauthorized() *authv3.CheckResponse {
	resp := &Response{
		Code: int64(typev3.StatusCode_Unauthorized),
		Msg:  UnauthorizedMsg,
		Data: struct{}{},
	}
	respJson, err := json.Marshal(resp)
	httpBody := ""
	if err == nil {
		httpBody = string(respJson)
	}
	return &authv3.CheckResponse{
		Status: &status.Status{Code: int32(codes.Unauthenticated)},
		HttpResponse: &authv3.CheckResponse_DeniedResponse{
			DeniedResponse: &authv3.DeniedHttpResponse{
				Status: &typev3.HttpStatus{Code: typev3.StatusCode_Unauthorized},
				Body:   httpBody,
			},
		},
	}
}

// Forbidden Forbidden 没有权限 403
func (s *Server) Forbidden() *authv3.CheckResponse {
	resp := &Response{
		Code: int64(typev3.StatusCode_Forbidden),
		Msg:  ForbiddenMsg,
		Data: struct{}{},
	}
	respJson, err := json.Marshal(resp)
	httpBody := ""
	if err == nil {
		httpBody = string(respJson)
	}

	return &authv3.CheckResponse{
		Status: &status.Status{Code: int32(codes.PermissionDenied)},
		HttpResponse: &authv3.CheckResponse_DeniedResponse{
			DeniedResponse: &authv3.DeniedHttpResponse{
				Status: &typev3.HttpStatus{Code: typev3.StatusCode_Forbidden},
				Body:   httpBody,
			},
		},
	}
}
