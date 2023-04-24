package clientV1

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderPBV1 "orderservice/genproto/go/v1"
	"userservice/config"
)

// NewOrderClient 实例化 order 客户端
func NewOrderClient(conf *config.Config) (orderPBV1.OrderServiceClient, error) {

	serverAddress := conf.Client.OrderHost + conf.Client.OrderPort

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	conn, err := grpc.DialContext(
		ctx, serverAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		return nil, err
	}
	client := orderPBV1.NewOrderServiceClient(conn)
	return client, nil

}
