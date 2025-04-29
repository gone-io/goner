module github.com/gone-io/goner/otel/meter/prometheus

go 1.24.1

require (
	github.com/gone-io/gone/v2 v2.1.0
	github.com/gone-io/goner/g v1.1.1 // indirect
	github.com/gone-io/goner/otel v1.1.1 // indirect
	github.com/gone-io/goner/otel/meter v1.1.1
	go.opentelemetry.io/otel/exporters/prometheus v0.57.0
	go.opentelemetry.io/otel/sdk/metric v1.35.0
)

replace (
	github.com/gone-io/goner/g => ../../../g
	github.com/gone-io/goner/otel => ../../
	github.com/gone-io/goner/otel/meter => ../
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/prometheus/client_golang v1.20.5 // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/common v0.62.0 // indirect
	github.com/prometheus/procfs v0.15.1 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/otel v1.35.0 // indirect
	go.opentelemetry.io/otel/exporters/stdout/stdoutmetric v1.35.0 // indirect
	go.opentelemetry.io/otel/metric v1.35.0 // indirect
	go.opentelemetry.io/otel/sdk v1.35.0 // indirect
	go.opentelemetry.io/otel/trace v1.35.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
	google.golang.org/protobuf v1.36.5 // indirect
)
