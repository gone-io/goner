module examples/otel/tracer/simple

go 1.24.1

require (
	github.com/gone-io/gone/v2 v2.2.5
	github.com/gone-io/goner/g v1.3.1
	github.com/gone-io/goner/otel/tracer v1.3.1
	go.opentelemetry.io/otel/trace v1.36.0
)

replace (
	github.com/gone-io/goner/g => ./../../../../g
	github.com/gone-io/goner/otel => ./../../../../otel
	github.com/gone-io/goner/otel/tracer => ./../../../../otel/tracer
)

require (
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/gone-io/goner/otel v1.3.1 // indirect
	github.com/google/uuid v1.6.0 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/otel v1.36.0 // indirect
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.36.0 // indirect
	go.opentelemetry.io/otel/metric v1.36.0 // indirect
	go.opentelemetry.io/otel/sdk v1.36.0 // indirect
	go.uber.org/mock v0.5.2 // indirect
	golang.org/x/sys v0.33.0 // indirect
)
