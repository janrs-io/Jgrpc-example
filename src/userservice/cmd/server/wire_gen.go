// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package server

import (
	"github.com/janrs-io/Jgrpc-otel-span"
	"userservice/config"
	"userservice/service/v1/client"
	"userservice/service/v1/server"
)

// Injectors from wire.go:

func InitServer(cfg string) (*Server, error) {
	configConfig := config.NewConfig(cfg)
	db := NewMysqlDB(configConfig)
	client := NewRedis(configConfig)
	orderServiceClient, err := clientV1.NewOrderClient(configConfig)
	if err != nil {
		return nil, err
	}
	productServiceClient, err := clientV1.NewProductClient(configConfig)
	if err != nil {
		return nil, err
	}
	userServiceClient, err := clientV1.NewUserClient(configConfig)
	if err != nil {
		return nil, err
	}
	tracerProvider, err := NewTrace(configConfig)
	if err != nil {
		return nil, err
	}
	otelSpan := Jgrpc_otelspan.New(tracerProvider)
	repository := serverV1.NewRepository(db, client, orderServiceClient, productServiceClient, userServiceClient, otelSpan, configConfig)
	group := NewRunGroup()
	logger := NewLogger()
	server := NewHttpServer(configConfig)
	userServiceServer := serverV1.NewServer(repository, logger, userServiceClient, orderServiceClient, productServiceClient)
	grpcServer := NewGrpcServer(userServiceServer)
	serverServer := NewServer(repository, configConfig, group, logger, server, grpcServer, db, tracerProvider)
	return serverServer, nil
}
