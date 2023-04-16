package clientV1

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"userservice/config"
	userPBV1 "userservice/genproto/go/v1"
)

// NewUserClient 实例化 user 客户端
func NewUserClient(conf *config.Config) (userPBV1.UserServiceClient, error) {

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
	client := userPBV1.NewUserServiceClient(conn)
	return client, nil

}
