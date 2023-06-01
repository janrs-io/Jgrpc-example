package server

import (
	"context"
	"os"

	"github.com/go-kit/log"
	"github.com/oklog/run"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.9.0"
	"google.golang.org/grpc"

	"authservice/config"
)

// NewRedis 实例化 redis 组件
func NewRedis(conf *config.Config) *redis.Client {

	rdb := redis.NewClient(&redis.Options{
		Addr:         conf.Redis.Host + conf.Redis.Port,
		Username:     conf.Redis.Username,
		Password:     conf.Redis.Password,
		DB:           conf.Redis.Database,
		DialTimeout:  conf.Redis.DialTimeout,
		ReadTimeout:  conf.Redis.ReadTimeout,
		WriteTimeout: conf.Redis.WriteTimeout,
		PoolSize:     conf.Redis.PoolSize,
		PoolTimeout:  conf.Redis.PoolTimeout,
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		panic("redis connect failed [ERROR]=> " + err.Error())
	}
	return rdb

}

// NewLogger 实例化 logger 组件
func NewLogger() log.Logger {
	return log.NewLogfmtLogger(os.Stderr)
}

// NewRunGroup 实例化 run.Group 组件
func NewRunGroup() *run.Group {
	return &run.Group{}
}

// NewTrace 实例化 Trace
func NewTrace(conf *config.Config) (*sdktrace.TracerProvider, error) {

	ctx := context.Background()
	traceClient := otlptracegrpc.NewClient(
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(conf.Trace.EndPoint),
		otlptracegrpc.WithDialOption(grpc.WithBlock()),
	)
	traceExp, err := otlptrace.New(ctx, traceClient)
	if err != nil {
		return nil, err
	}

	res, err := resource.New(ctx,
		resource.WithFromEnv(),
		resource.WithProcess(),
		resource.WithTelemetrySDK(),
		resource.WithHost(),
		resource.WithAttributes(semconv.ServiceNameKey.String(conf.Trace.ServiceName)))
	if err != nil {
		return nil, err
	}

	bsp := sdktrace.NewBatchSpanProcessor(traceExp)

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)

	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(propagation.TraceContext{},
			propagation.Baggage{}),
	)
	otel.SetTracerProvider(tracerProvider)

	return tracerProvider, nil

}
