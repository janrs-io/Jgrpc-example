// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package server

import (
	"authservice/config"
	"authservice/genproto/go/v1"
	"authservice/pkg"
	"authservice/service"
)

// Injectors from wire.go:

// InitServer Inject service's component
func InitServer(conf *config.Config) (authV1.AuthServiceServer, error) {
	client := pkg.NewRedis(conf)
	repository := service.NewRepository(client)
	authServiceClient, err := service.NewClient(conf)
	if err != nil {
		return nil, err
	}
	authServiceServer := service.NewServer(repository, authServiceClient)
	return authServiceServer, nil
}