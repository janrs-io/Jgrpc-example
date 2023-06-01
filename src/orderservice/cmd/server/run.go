package server

import (
	"context"
	Jgrpc_pgv_interceptor "github.com/janrs-io/Jgrpc-pgv-interceptor"
	"net"
	"net/http"
	"os"
	"syscall"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	Jgrpc_response "github.com/janrs-io/Jgrpc-response"
	"github.com/oklog/run"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/gorm"

	"orderservice/config"
	orderPBV1 "orderservice/genproto/go/v1"
	"orderservice/service/model"
)

// Server 启动微服务所需要的所有依赖
type Server struct {
	conf       *config.Config
	runGroup   *run.Group
	logger     log.Logger
	httpServer *http.Server
	grpcServer *grpc.Server
	mysqlDB    *gorm.DB
	trace      *sdktrace.TracerProvider
}

// NewServer 实例化 Server
func NewServer(
	conf *config.Config,
	runGroup *run.Group,
	logger log.Logger,
	httpServer *http.Server,
	grpcServer *grpc.Server,
	mysqlDB *gorm.DB,
	trace *sdktrace.TracerProvider,
) *Server {
	return &Server{
		conf:       conf,
		runGroup:   runGroup,
		logger:     logger,
		httpServer: httpServer,
		grpcServer: grpcServer,
		mysqlDB:    mysqlDB,
		trace:      trace,
	}
}

// NewGrpcServer 实例化 Grpc 服务
func NewGrpcServer(orderServerV1 orderPBV1.OrderServiceServer) *grpc.Server {

	grpcServer := grpc.NewServer(
		grpc.ChainStreamInterceptor(
			// otel 链路追踪
			otelgrpc.StreamServerInterceptor(),
		),
		grpc.ChainUnaryInterceptor(
			// otel 链路追踪
			otelgrpc.UnaryServerInterceptor(),
			// PGV 中间件
			Jgrpc_pgv_interceptor.ValidationUnaryInterceptor,
		),
	)
	orderPBV1.RegisterOrderServiceServer(grpcServer, orderServerV1)
	return grpcServer

}

// NewHttpServer 实例化 Http 服务
func NewHttpServer(conf *config.Config) *http.Server {

	mux := runtime.NewServeMux(
		runtime.WithErrorHandler(Jgrpc_response.HttpErrorHandler),
		runtime.WithForwardResponseOption(Jgrpc_response.HttpSuccessResponseModifier),
		runtime.WithMarshalerOption("*", &Jgrpc_response.CustomMarshaller{}),
	)
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	if err := orderPBV1.RegisterOrderServiceHandlerFromEndpoint(
		context.Background(),
		mux,
		conf.Grpc.Port,
		opts,
	); err != nil {
		panic("register service handler failed.[ERROR]=>" + err.Error())
	}

	httpServer := &http.Server{
		Addr:    conf.Http.Port,
		Handler: mux,
	}
	return httpServer

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

	// 启动 http 服务
	s.runGroup.Add(func() error {

		_ = level.Info(s.logger).Log("msg", "starting HTTP server", "addr", s.httpServer.Addr)
		return s.httpServer.ListenAndServe()

	}, func(err error) {
		if err = s.httpServer.Close(); err != nil {
			_ = level.Error(s.logger).Log("msg", "failed to stop web server", "err", err)
		}
	})
	s.runGroup.Add(run.SignalHandler(context.Background(), syscall.SIGINT, syscall.SIGTERM))

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

	// 执行 migrate
	model.Migrate(server.mysqlDB)

	// 上报 trace 数据
	defer func() {
		if err = server.trace.Shutdown(context.Background()); err != nil {
			_ = level.Info(server.logger).Log("msg", "shutdown trace provider failed", "err", err)
		}
	}()

	// 启动 grpc 以及 http 服务
	server.RunServer()

}
