package clientV1

import (
	"context"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"productservice/config"
	productPBV1 "productservice/genproto/go/v1"
)

// NewProductClient 实例化 Product 客户端
func NewProductClient(conf *config.Config) (productPBV1.ProductServiceClient, error) {

	serverAddress := conf.Grpc.Host + conf.Grpc.Port
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(
		ctx, serverAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),
	)
	if err != nil {
		return nil, err
	}
	client := productPBV1.NewProductServiceClient(conn)
	return client, nil

}
