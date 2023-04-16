//go:build wireinject
// +build wireinject

package server

import (
	clientV1 "authservice/service/v1/client"
	serverV1 "authservice/service/v1/server"
	"github.com/google/wire"

	"authservice/config"
)

func InitServer(cfg string) (*Server, error) {

	wire.Build(
		NewServer,
		clientV1.NewAuthClient,
		serverV1.NewServer,
		serverV1.NewRepository,
		config.NewConfig,
		NewRedis,
		NewHttpServer,
		NewGrpcServer,
		NewRunGroup,
		NewLogger,
	)

	return &Server{}, nil

}
