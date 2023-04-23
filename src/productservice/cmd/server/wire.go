//go:build wireinject
// +build wireinject

package server

import (
	"github.com/google/wire"

	"productservice/config"
	clientV1 "productservice/service/v1/client"
	"productservice/service/v1/server"
)

func InitServer(cfg string) (*Server, error) {

	wire.Build(
		NewServer,
		serverV1.NewServer,
		serverV1.NewRepository,
		clientV1.NewAuthClient,
		config.NewConfig,
		NewDB,
		NewHttpServer,
		NewGrpcServer,
		NewRunGroup,
		NewLogger,
	)

	return &Server{}, nil

}
