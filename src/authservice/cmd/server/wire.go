//go:build wireinject
// +build wireinject

package server

import (
	"github.com/google/wire"

	"authservice/config"
	authV1 "authservice/genproto/go/v1"
	"authservice/pkg"
	"authservice/service"
)

// InitServer Inject service's component
func InitServer(conf *config.Config) (authV1.AuthServiceServer, error) {

	wire.Build(
		service.NewClient,
		service.NewServer,
		service.NewRepository,
		pkg.NewRedis,
	)

	return &service.Server{}, nil

}
