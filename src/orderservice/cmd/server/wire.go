//go:build wireinject
// +build wireinject

package server

import (
	"github.com/google/wire"
	Jgrpc_otelspan "github.com/janrs-io/Jgrpc-otel-span"

	"orderservice/config"
	clientV1 "orderservice/service/v1/client"
	serverV1 "orderservice/service/v1/server"
)

func InitServer(cfg string) (*Server, error) {

	wire.Build(
		// 配置
		config.NewConfig,

		// 启动 grpc 以及 http 服务
		NewServer,

		// 实例化服务
		serverV1.NewServer,
		serverV1.NewRepository,

		// 实例化客户端
		clientV1.NewProductClient,
		clientV1.NewOrderClient,

		// 组件
		NewMysqlDB,
		NewHttpServer,
		NewGrpcServer,
		NewRunGroup,
		NewLogger,
		NewTrace,
		Jgrpc_otelspan.New,
	)

	return &Server{}, nil

}
