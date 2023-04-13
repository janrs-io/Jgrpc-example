package pkg

import (
	"context"
	"fmt"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/natefinch/lumberjack.v2"

	"userservice/config"
)

// ZapInterceptor Grpc zap logger interceptor
func ZapInterceptor(conf *config.Config) *zap.Logger {

	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:  conf.Logger.Path,
		MaxSize:   conf.Logger.MaxSize,
		LocalTime: conf.Logger.LocalTime,
	})

	zapConfig := zap.NewProductionEncoderConfig()
	zapConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zapConfig),
		w,
		zap.NewAtomicLevel(),
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	grpc_zap.ReplaceGrpcLoggerV2(logger)
	return logger

}

// Validator PGV validator interceptor
type Validator interface {
	ValidateAll() error
}

// ServerValidationUnaryInterceptor PGV validator interceptor
func ServerValidationUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {

	if r, ok := req.(Validator); ok {
		if err = r.ValidateAll(); err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
	}
	return handler(ctx, req)

}

// AuthInterceptor Grpc gateway server auth interceptor
func AuthInterceptor(ctx context.Context) (context.Context, error) {
	authToken, err := grpc_auth.AuthFromMD(ctx, "Bearer")
	if err != nil {
		fmt.Println(err)
		return nil, status.Error(codes.Unauthenticated, "请先登录")
	}

	auth := NewAuth()

	username, err := auth.ValidateAuth(authToken)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "请先登录")
	}

	newCtx := context.WithValue(ctx, "username", username)
	return newCtx, nil
}
