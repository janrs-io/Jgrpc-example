//go:build wireinject
// +build wireinject

package server

import (
	"github.com/google/wire"

	"userservice/config"
	clientV1 "userservice/service/v1/client"
	serverV1 "userservice/service/v1/server"
)

func InitServer(cfg string) (*Server, error) {

	wire.Build(
		// 启动服务
		NewServer,

		// 实例化服务
		serverV1.NewServer,
		serverV1.NewRepository,

		// 客户端
		clientV1.NewAuthClient,
		clientV1.NewUserClient,
		clientV1.NewOrderClient,
		clientV1.NewProductClient,

		// 配置
		config.NewConfig,

		// 组件
		NewRedis,
		NewDB,
		NewHttpServer,
		NewGrpcServer,
		NewRunGroup,
		NewLogger,
	)

	return &Server{}, nil

}
