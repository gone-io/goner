module examples/zap

go 1.24.1

replace (
	github.com/gone-io/goner/g => ../../g
	github.com/gone-io/goner/tracer/gid => ../../tracer/gid
	github.com/gone-io/goner/zap => ../../zap

)

require (
	github.com/gone-io/gone/v2 v2.2.6
	github.com/gone-io/goner/g v1.3.6
	github.com/gone-io/goner/tracer/gid v1.3.6-00010101000000-000000000000
	github.com/gone-io/goner/zap v1.3.6-00010101000000-000000000000
	go.uber.org/zap v1.27.0
)

require (
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/petermattis/goid v0.0.0-20250904145737-900bdf8bb490 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/contrib/bridges/otelzap v0.13.0 // indirect
	go.opentelemetry.io/otel v1.38.0 // indirect
	go.opentelemetry.io/otel/log v0.14.0 // indirect
	go.opentelemetry.io/otel/metric v1.38.0 // indirect
	go.opentelemetry.io/otel/trace v1.38.0 // indirect
	go.uber.org/mock v0.6.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
)
