package server

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"userservice/config"
	userV1 "userservice/genproto/go/v1"
	"userservice/migrate"
	"userservice/pkg"
)

// Run Run service server
func Run(cfg string) {

	conf := config.NewConfig(cfg)
	// migrate
	Migrate(conf)
	// run grpc server
	RunGrpcServer(initServer(conf), conf)
	// run http server
	RunHttpServer(conf)
	// listen exit server event
	HandleExitServer(conf)

}

// SetServer Wire inject service's component
func initServer(conf *config.Config) userV1.UserServiceServer {
	server, err := InitServer(conf)
	if err != nil {
		panic("run server failed.[ERROR]=>" + err.Error())
	}
	return server
}

// HandleExitServer Handle service exit event
func HandleExitServer(conf *config.Config) {

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	conf.Grpc.Server.GracefulStop()
	if err := conf.Http.Server.Shutdown(ctx); err != nil {
		panic("shutdown service failed.[ERROR]=>" + err.Error())
	}
	<-ctx.Done()
	close(ch)
	fmt.Println("Graceful shutdown http & grpc server.")

}

// Migrate Migrate tables
func Migrate(conf *config.Config) {
	migrate.Migrate(pkg.NewDB(conf), conf)
}