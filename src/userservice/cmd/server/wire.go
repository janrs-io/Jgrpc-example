//go:build wireinject
// +build wireinject

package server

import (
	"github.com/google/wire"
	Jgrpc_otelspan "github.com/janrs-io/Jgrpc-otel-span"

	"userservice/config"
	clientV1 "userservice/service/v1/client"
	serverV1 "userservice/service/v1/server"
)

func InitServer(cfg string) (*Server, error) {

	wire.Build(
		// 获取配置
		config.NewConfig,
		// 启动服务
		NewServer,

		// 实例化服务
		serverV1.NewServer,
		serverV1.NewRepository,

		// 客户端
		clientV1.NewUserClient,
		clientV1.NewOrderClient,
		clientV1.NewProductClient,

		// 组件
		NewRedis,
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
