module examples/otel/log/collect_traceid

go 1.24.1

require (
	github.com/gone-io/gone/v2 v2.2.1
	github.com/gone-io/goner/g v1.2.1
	github.com/gone-io/goner/otel v1.2.1 // indirect
	github.com/gone-io/goner/otel/log v1.2.1 // indirect
	github.com/gone-io/goner/otel/log/http v1.2.1
	github.com/gone-io/goner/otel/tracer v1.2.1 // indirect
	github.com/gone-io/goner/otel/tracer/http v1.2.1
	github.com/gone-io/goner/tracer/gid v1.2.1
	github.com/gone-io/goner/viper v1.2.1
	github.com/gone-io/goner/zap v1.2.1
)

replace (
	github.com/gone-io/goner/g => ./../../../../g
	github.com/gone-io/goner/otel => ./../../../../otel
	github.com/gone-io/goner/otel/log => ./../../../../otel/log
	github.com/gone-io/goner/otel/log/http => ./../../../../otel/log/http
	github.com/gone-io/goner/otel/tracer => ./../../../../otel/tracer
	github.com/gone-io/goner/otel/tracer/http => ./../../../../otel/tracer/http
	github.com/gone-io/goner/tracer/gid => ./../../../../tracer/gid
	github.com/gone-io/goner/viper => ./../../../../viper
	github.com/gone-io/goner/zap => ./../../../../zap
)

require (
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/otel v1.35.0
	go.opentelemetry.io/otel/metric v1.35.0 // indirect
	go.opentelemetry.io/otel/trace v1.35.0 // indirect

)

require (
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/fsnotify/fsnotify v1.8.0 // indirect
	github.com/go-viper/encoding/javaproperties v0.1.0 // indirect
	github.com/go-viper/mapstructure/v2 v2.2.1 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.26.1 // indirect
	github.com/magiconair/properties v1.8.9 // indirect
	github.com/pelletier/go-toml/v2 v2.2.3 // indirect
	github.com/petermattis/goid v0.0.0-20250319124200-ccd6737f222a // indirect
	github.com/sagikazarmark/locafero v0.9.0 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/spf13/afero v1.14.0 // indirect
	github.com/spf13/cast v1.7.1 // indirect
	github.com/spf13/pflag v1.0.6 // indirect
	github.com/spf13/viper v1.20.1 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	go.opentelemetry.io/contrib/bridges/otelzap v0.10.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp v0.11.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.35.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.35.0 // indirect
	go.opentelemetry.io/otel/exporters/stdout/stdoutlog v0.11.0 // indirect
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.35.0 // indirect
	go.opentelemetry.io/otel/log v0.11.0 // indirect
	go.opentelemetry.io/otel/sdk v1.35.0 // indirect
	go.opentelemetry.io/otel/sdk/log v0.11.0 // indirect
	go.opentelemetry.io/proto/otlp v1.5.0 // indirect
	go.uber.org/mock v0.5.2 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	golang.org/x/net v0.38.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.23.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20250218202821-56aae31c358a // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250218202821-56aae31c358a // indirect
	google.golang.org/grpc v1.71.0 // indirect
	google.golang.org/protobuf v1.36.5 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
