package clientV1

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	productPBV1 "productservice/genproto/go/v1"
	"userservice/config"
)

// NewProductClient 实例化 product 客户端
func NewProductClient(conf *config.Config) (productPBV1.ProductServiceClient, error) {

	serverAddress := conf.Client.ProductHost + conf.Client.ProductPort

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	conn, err := grpc.DialContext(
		ctx, serverAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		return nil, err
	}
	client := productPBV1.NewProductServiceClient(conn)
	return client, nil

}
