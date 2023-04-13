package pkg

import (
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"permissionservice/config"

	"go.uber.org/zap"
)

// ZapInterceptor Grpc logger middleware
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

	/*
		logger, err := zap.NewDevelopment()
		if err != nil {
			log.Fatalf("failed to initialize zap logger: %v", err)
		}
		grpc_zap.ReplaceGrpcLoggerV2(logger)
		return logger
	*/

}
