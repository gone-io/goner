module examples/otel/tracer/oltp/grpc/client

go 1.24.1

require (
	examples/otel/tracer/oltp/grpc v1.1.1
	github.com/gone-io/gone/v2 v2.2.6
	github.com/gone-io/goner/g v1.3.6
	github.com/gone-io/goner/grpc v1.3.6
	github.com/gone-io/goner/otel v1.3.6 // indirect
	github.com/gone-io/goner/otel/tracer v1.3.6 // indirect
	github.com/gone-io/goner/otel/tracer/grpc v1.3.6
	github.com/gone-io/goner/viper v1.3.6
)

replace (
	examples/otel/tracer/oltp/grpc => ../
	github.com/gone-io/goner/g => ./../../../../../g
	github.com/gone-io/goner/grpc => ./../../../../../grpc
	github.com/gone-io/goner/otel => ./../../../../../otel
	github.com/gone-io/goner/otel/tracer => ./../../../../../otel/tracer
	github.com/gone-io/goner/otel/tracer/grpc => ./../../../../../otel/tracer/grpc
	github.com/gone-io/goner/viper => ./../../../../../viper
)

require (
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-viper/encoding/javaproperties v0.1.0 // indirect
	github.com/go-viper/mapstructure/v2 v2.4.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.27.2 // indirect
	github.com/magiconair/properties v1.8.10 // indirect
	github.com/pelletier/go-toml/v2 v2.2.4 // indirect
	github.com/sagikazarmark/locafero v0.11.0 // indirect
	github.com/sourcegraph/conc v0.3.1-0.20240121214520-5f936abd7ae8 // indirect
	github.com/spf13/afero v1.15.0 // indirect
	github.com/spf13/cast v1.10.0 // indirect
	github.com/spf13/pflag v1.0.10 // indirect
	github.com/spf13/viper v1.21.0 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/otel v1.38.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.38.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.38.0 // indirect
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.38.0 // indirect
	go.opentelemetry.io/otel/metric v1.38.0 // indirect
	go.opentelemetry.io/otel/sdk v1.38.0 // indirect
	go.opentelemetry.io/otel/trace v1.38.0
	go.opentelemetry.io/proto/otlp v1.8.0 // indirect
	go.uber.org/mock v0.6.0 // indirect
	golang.org/x/net v0.44.0 // indirect
	golang.org/x/sys v0.36.0 // indirect
	golang.org/x/text v0.29.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20250908214217-97024824d090 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250908214217-97024824d090 // indirect
	google.golang.org/grpc v1.75.0
	google.golang.org/protobuf v1.36.9 // indirect
)

require (
	github.com/cenkalti/backoff/v5 v5.0.3 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.63.0 // indirect
	go.yaml.in/yaml/v3 v3.0.4 // indirect
)
