package server

import (
	"context"
	"net"
	"net/http"
	"os"
	"syscall"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	Jgrpc_response "github.com/janrs-io/Jgrpc-response"
	"github.com/oklog/run"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	stdout "go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"authservice/config"
	authPBV1 "authservice/genproto/go/v1"
	serverV1 "authservice/service/v1/server"
)

// Server 启动微服务所需要的所有依赖
type Server struct {
	repo       *serverV1.Repository
	conf       *config.Config
	g          *run.Group
	logger     log.Logger
	httpServer *http.Server
	grpcServer *grpc.Server
}

// NewServer 实例化 Server
func NewServer(
	repo *serverV1.Repository,
	conf *config.Config,
	g *run.Group,
	logger log.Logger,
	httpServer *http.Server,
	grpcServer *grpc.Server,
) *Server {
	return &Server{
		conf:       conf,
		repo:       repo,
		g:          g,
		logger:     logger,
		httpServer: httpServer,
		grpcServer: grpcServer,
	}
}

// NewGrpcServer 实例化 Grpc 服务
func NewGrpcServer(logger log.Logger, authServer authPBV1.AuthServiceServer) *grpc.Server {

	logTraceID := func(ctx context.Context) logging.Fields {
		if span := trace.SpanContextFromContext(ctx); span.IsSampled() {
			return logging.Fields{"traceID", span.TraceID().String()}
		}
		return nil
	}

	// Set up OTLP tracing (stdout for debug).
	exporter, err := stdout.New(stdout.WithPrettyPrint())
	if err != nil {
		_ = level.Error(logger).Log("err", err)
		os.Exit(1)
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	defer func() { _ = exporter.Shutdown(context.Background()) }()

	grpcServer := grpc.NewServer(
		grpc.ChainStreamInterceptor(
			otelgrpc.StreamServerInterceptor(),
			logging.StreamServerInterceptor(LoggerInterceptor(logger), logging.WithFieldsFromContext(logTraceID)),
		),
		grpc.ChainUnaryInterceptor(
			otelgrpc.UnaryServerInterceptor(),
			logging.UnaryServerInterceptor(LoggerInterceptor(logger), logging.WithFieldsFromContext(logTraceID)),
			// PGV 中间件
			ValidationUnaryInterceptor,
		),
	)
	authPBV1.RegisterAuthServiceServer(grpcServer, authServer)
	return grpcServer

}

// NewHttpServer 实例化 Http 服务
func NewHttpServer(conf *config.Config) *http.Server {

	mux := runtime.NewServeMux(
		runtime.WithErrorHandler(Jgrpc_response.HttpErrorHandler),
		runtime.WithForwardResponseOption(Jgrpc_response.HttpSuccessResponseModifier),
		runtime.WithMarshalerOption("*", &Jgrpc_response.CustomMarshaler{}),
	)
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	if err := authPBV1.RegisterAuthServiceHandlerFromEndpoint(
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

	// 执行 migrate

	// 启动 grpc 服务
	s.g.Add(func() error {
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
	s.g.Add(func() error {

		_ = level.Info(s.logger).Log("msg", "starting HTTP server", "addr", s.httpServer.Addr)
		return s.httpServer.ListenAndServe()

	}, func(err error) {
		if err = s.httpServer.Close(); err != nil {
			_ = level.Error(s.logger).Log("msg", "failed to stop web server", "err", err)
		}
	})
	s.g.Add(run.SignalHandler(context.Background(), syscall.SIGINT, syscall.SIGTERM))

	if err := s.g.Run(); err != nil {
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
	server.RunServer()

}
