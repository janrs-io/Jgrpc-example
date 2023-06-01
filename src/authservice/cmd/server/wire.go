//go:build wireinject
// +build wireinject

package server

import (
	serverV1 "authservice/service/v1/server"
	"github.com/google/wire"

	"authservice/config"
)

func InitServer(cfg string) (*Server, error) {

	wire.Build(
		// 配置
		config.NewConfig,

		// 实例化 grpc 以及 http 服务
		NewServer,

		// 实例化服务
		serverV1.NewServer,
		serverV1.NewRepository,

		// 组件
		NewRedis,
		NewGrpcServer,
		NewRunGroup,
		NewLogger,
		NewTrace,
	)

	return &Server{}, nil

}
