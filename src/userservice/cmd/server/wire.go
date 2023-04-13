//go:build wireinject
// +build wireinject

package server

import (
	"github.com/google/wire"

	"userservice/config"
	userV1 "userservice/genproto/go/v1"
	"userservice/pkg"
	"userservice/service"
)

// InitServer Inject service's component
func InitServer(conf *config.Config) (userV1.UserServiceServer, error) {

	wire.Build(
		service.NewServer,
		service.NewRepository,
		service.NewClient,
		service.NewAuthClient,
		pkg.NewDB,
		pkg.NewRedis,
	)

	return &service.Server{}, nil

}
