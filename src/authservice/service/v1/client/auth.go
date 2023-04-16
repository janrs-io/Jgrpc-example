package clientV1

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"authservice/config"
	authPBV1 "authservice/genproto/go/v1"
)

// NewAuthClient 实例化 auth 客户端
func NewAuthClient(conf *config.Config) (authPBV1.AuthServiceClient, error) {

	serverAddress := conf.Grpc.Host + conf.Grpc.Port
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(
		ctx,
		serverAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		return nil, err
	}
	client := authPBV1.NewAuthServiceClient(conn)
	return client, nil

}
