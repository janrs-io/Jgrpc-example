package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"userservice/config"
	userV1 "userservice/genproto/go/v1"
	"userservice/pkg"
)

// RunHttpServer Run http server
func RunHttpServer(conf *config.Config) {

	mux := runtime.NewServeMux(
		runtime.WithErrorHandler(pkg.HttpErrorHandler),
		runtime.WithForwardResponseOption(pkg.HttpSuccessResponseModifier),
		runtime.WithMarshalerOption("*", &pkg.CustomMarshaler{}),
	)

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	if err := userV1.RegisterUserServiceHandlerFromEndpoint(
		context.Background(),
		mux,
		conf.Grpc.Port,
		opts,
	); err != nil {
		panic("register service handler failed.[ERROR]=>" + err.Error())
	}

	httpServer := &http.Server{
		Addr:    conf.Http.Port,
		Handler: mux,
	}
	fmt.Println("Listening http server on port" + conf.Http.Port)

	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			fmt.Println("listen http server failed.[ERROR]=>" + err.Error())
		}
	}()

	conf.Http.Server = httpServer

}
