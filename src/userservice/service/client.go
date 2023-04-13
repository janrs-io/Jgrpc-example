package service

import (
	authV1 "authservice/genproto/go/v1"
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"userservice/config"
	userV1 "userservice/genproto/go/v1"
)

// NewClient New service's client
func NewClient(conf *config.Config) (userV1.UserServiceClient, error) {

	serverAddress := conf.Grpc.Host + conf.Grpc.Port
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(
		ctx, serverAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		return nil, err
	}
	client := userV1.NewUserServiceClient(conn)
	return client, nil

}

// NewAuthClient New auth service's client
func NewAuthClient(conf *config.Config) (authV1.AuthServiceClient, error) {

	serverAddress := conf.Client.AuthHost + conf.Client.AuthPort

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	conn, err := grpc.DialContext(
		ctx, serverAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		return nil, err
	}
	client := authV1.NewAuthServiceClient(conn)
	return client, nil

}
