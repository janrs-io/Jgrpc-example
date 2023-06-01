package server

import (
	"context"
	"net"
	"os"
	"syscall"

	serverV1 "authservice/service/v1/server"
	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/oklog/run"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc"

	"authservice/config"
)

// Server 启动微服务所需要的所有依赖
type Server struct {
	repo       *serverV1.Repository
	conf       *config.Config
	runGroup   *run.Group
	logger     log.Logger
	grpcServer *grpc.Server
	trace      *sdktrace.TracerProvider
}

// NewServer 实例化 Server
func NewServer(
	repo *serverV1.Repository,
	conf *config.Config,
	runGroup *run.Group,
	logger log.Logger,
	grpcServer *grpc.Server,
	trace *sdktrace.TracerProvider,
) *Server {
	return &Server{
		conf:       conf,
		repo:       repo,
		runGroup:   runGroup,
		logger:     logger,
		grpcServer: grpcServer,
		trace:      trace,
	}
}

// NewGrpcServer 实例化 Grpc 服务
func NewGrpcServer(authServerV1 authv3.AuthorizationServer) *grpc.Server {

	grpcServer := grpc.NewServer()
	authv3.RegisterAuthorizationServer(grpcServer, authServerV1)
	return grpcServer

}

// RunServer 启动 http 以及 grpc 服务
func (s *Server) RunServer() {

	// 启动 grpc 服务
	s.runGroup.Add(func() error {
		l, err := net.Listen("tcp", s.conf.Grpc.Port)
		if err != nil {
			return err
		}
		_ = level.Info(s.logger).Log("msg", "starting gRPC server", "addr", l.Addr().String())
		return s.grpcServer.Serve(l)
	}, func(err error) {
		s.grpcServer.GracefulStop()
		s.grpcServer.Stop()
	})

	// 监听退出信号
	s.runGroup.Add(run.SignalHandler(context.Background(), syscall.SIGINT, syscall.SIGTERM))

	// 顺序启动服务
	if err := s.runGroup.Run(); err != nil {
		_ = level.Error(s.logger).Log("err", err)
		os.Exit(1)
	}

}

// Run 启动服务
func Run(cfg string) {

	server, err := InitServer(cfg)
	if err != nil {
		panic("run server failed.[ERROR]=>" + err.Error())
	}

	// 上报链路 trace 数据
	defer func() {
		if err = server.trace.Shutdown(context.Background()); err != nil {
			_ = level.Info(server.logger).Log("msg", "shutdown trace provider failed", "err", err)
		}
	}()

	server.RunServer()

}
