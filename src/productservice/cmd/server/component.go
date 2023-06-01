package server

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.9.0"
	"google.golang.org/grpc"
	"os"
	"strconv"

	"github.com/go-kit/log"
	"github.com/oklog/run"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"productservice/config"
)

// NewLogger 实例化 go-kit/log 组件
func NewLogger() log.Logger {
	return log.NewLogfmtLogger(os.Stderr)
}

// NewRunGroup 实例化 run.Group 组件
func NewRunGroup() *run.Group {
	return &run.Group{}
}

// NewMysqlDB 初始化 mysql 连接
func NewMysqlDB(conf *config.Config) *gorm.DB {

	// Database configuration
	dbConf := conf.Database
	if dbConf.Mysql.Database == "" {
		panic("database config is empty.")
	}

	// Database connection dsn
	dsn := dbConf.Mysql.UserName + ":" +
		dbConf.Mysql.Password +
		"@tcp(" + dbConf.Mysql.Host + ":" + strconv.Itoa(dbConf.Mysql.Port) + ")/" + dbConf.Mysql.Database +
		"?charset=" + dbConf.Mysql.Charset + "&parseTime=True&loc=Local"

	// Set Config
	mysqlConfig := mysql.Config{
		DSN:                       dsn,   // DSN data source name
		DefaultStringSize:         255,   // Default length of the string type field
		DisableDatetimePrecision:  true,  // Disable datetime precision, not supported on databases prior to MySQL 5.6
		DontSupportRenameIndex:    true,  // Renaming indexes is done by deleting and creating new ones.
		DontSupportRenameColumn:   true,  // Rename columns with `change`, renaming columns is not supported in databases prior to MySQL 8 and MariaDB
		SkipInitializeWithVersion: false, // Automatic configuration based on version
	}

	// New mysql with config
	newMysql := mysql.New(mysqlConfig)

	// Connect mysql
	conn, err := gorm.Open(
		newMysql,
		&gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true, // Disable automatic foreign key creation constraints
			SkipDefaultTransaction:                   true, // Close global open transactions
		})

	if err != nil {
		panic("mysql connect failed [ERROR]=> " + err.Error())
	}

	sqlDB, _ := conn.DB()
	sqlDB.SetMaxIdleConns(dbConf.Mysql.MaxIdleCons)
	sqlDB.SetMaxOpenConns(dbConf.Mysql.MaxOpenCons)

	return conn

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
