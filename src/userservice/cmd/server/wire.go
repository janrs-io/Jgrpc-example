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
		NewServer,
		clientV1.NewAuthClient,
		clientV1.NewUserClient,
		serverV1.NewServer,
		serverV1.NewRepository,
		config.NewConfig,
		NewRedis,
		NewDB,
		NewHttpServer,
		NewGrpcServer,
		NewRunGroup,
		NewLogger,
	)

	return &Server{}, nil

}
