package server

import (
	"fmt"
	"log"
	"net"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"google.golang.org/grpc"

	"userservice/config"
	userV1 "userservice/genproto/go/v1"
	"userservice/pkg"
)

// RunGrpcServer Run grpc server
func RunGrpcServer(server userV1.UserServiceServer, conf *config.Config) {

	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_zap.StreamServerInterceptor(pkg.ZapInterceptor(conf)),
			//grpc_auth.StreamServerInterceptor(pkg.AuthInterceptor),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_zap.UnaryServerInterceptor(pkg.ZapInterceptor(conf)),
			//grpc_auth.UnaryServerInterceptor(pkg.AuthInterceptor),
		)),
		grpc.ChainUnaryInterceptor(pkg.ServerValidationUnaryInterceptor),
	)

	userV1.RegisterUserServiceServer(grpcServer, server)

	fmt.Println("Listening grpc server on port" + conf.Grpc.Port)
	listen, err := net.Listen("tcp", conf.Grpc.Port)
	if err != nil {
		panic("listen grpc tcp failed.[ERROR]=>" + err.Error())
	}

	go func() {
		if err = grpcServer.Serve(listen); err != nil {
			log.Fatal("grpc serve failed", err)
		}
	}()

	conf.Grpc.Server = grpcServer

}
