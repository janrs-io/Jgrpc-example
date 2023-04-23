//go:build wireinject
// +build wireinject

package server

import (
	"github.com/google/wire"

	"orderservice/config"
	clientV1 "orderservice/service/v1/client"
	serverV1 "orderservice/service/v1/server"
)

func InitServer(cfg string) (*Server, error) {

	wire.Build(
		// run server
		NewServer,
		// server
		serverV1.NewServer,
		// client
		serverV1.NewRepository,
		clientV1.NewAuthClient,
		clientV1.NewProductClient,
		clientV1.NewOrderClient,
		// config
		config.NewConfig,
		// component
		NewDB,
		NewHttpServer,
		NewGrpcServer,
		NewRunGroup,
		NewLogger,
	)

	return &Server{}, nil

}
