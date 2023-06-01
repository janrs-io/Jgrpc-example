module userservice

go 1.19

replace authservice => ../authservice

replace orderservice => ../orderservice

replace productservice => ../productservice

require (
	github.com/envoyproxy/protoc-gen-validate v1.0.1
	github.com/go-kit/log v0.2.1
	github.com/google/uuid v1.3.0
	github.com/google/wire v0.5.0
	github.com/grpc-ecosystem/go-grpc-middleware/v2 v2.0.0-rc.5
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.15.2
	github.com/janrs-io/Jgrpc-otel-span v0.0.3
	github.com/janrs-io/Jgrpc-pgv-interceptor v0.0.1
	github.com/janrs-io/Jgrpc-response v0.0.2
	github.com/oklog/run v1.1.0
	github.com/redis/go-redis/v9 v9.0.4
	github.com/spf13/viper v1.15.0
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.36.0
	go.opentelemetry.io/otel v1.16.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.16.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.16.0
	go.opentelemetry.io/otel/sdk v1.16.0
	golang.org/x/crypto v0.0.0-20220525230936-793ad666bf5e
	google.golang.org/grpc v1.55.0
	google.golang.org/protobuf v1.30.0
	gorm.io/driver/mysql v1.5.1
	gorm.io/gorm v1.25.1
	orderservice v0.0.0
	productservice v0.0.0
)

require (
	github.com/cenkalti/backoff/v4 v4.2.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	github.com/go-logfmt/logfmt v0.5.1 // indirect
	github.com/go-logr/logr v1.2.4 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-sql-driver/mysql v1.7.0 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/pelletier/go-toml/v2 v2.0.6 // indirect
	github.com/spf13/afero v1.9.3 // indirect
	github.com/spf13/cast v1.5.0 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/subosito/gotenv v1.4.2 // indirect
	go.opentelemetry.io/otel/exporters/otlp/internal/retry v1.16.0 // indirect
	go.opentelemetry.io/otel/metric v1.16.0 // indirect
	go.opentelemetry.io/otel/trace v1.16.0 // indirect
	go.opentelemetry.io/proto/otlp v0.19.0 // indirect
	golang.org/x/net v0.10.0 // indirect
	golang.org/x/sys v0.8.0 // indirect
	golang.org/x/text v0.9.0 // indirect
	google.golang.org/genproto v0.0.0-20230410155749-daa745c078e1 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
