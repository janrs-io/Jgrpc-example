package pkg

import (
	"context"

	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/natefinch/lumberjack.v2"

	"authservice/config"
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
