module github.com/gone-io/goner/otel/tracer/zpkkin

go 1.24.1

require (
	github.com/gone-io/gone/v2 v2.2.5
	github.com/gone-io/goner/g v1.3.2
	github.com/gone-io/goner/otel v1.3.2 // indirect
	github.com/gone-io/goner/otel/tracer v1.3.2
	go.opentelemetry.io/otel/exporters/zipkin v1.36.0
	go.opentelemetry.io/otel/sdk v1.36.0
)

replace (
	github.com/gone-io/goner/g => ../../../g
	github.com/gone-io/goner/otel => ../../
	github.com/gone-io/goner/otel/tracer => ../
)

require github.com/stretchr/testify v1.10.0

require (
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/openzipkin/zipkin-go v0.4.3 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/otel v1.36.0 // indirect
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.36.0 // indirect
	go.opentelemetry.io/otel/metric v1.36.0 // indirect
	go.opentelemetry.io/otel/trace v1.36.0 // indirect
	go.uber.org/mock v0.5.2 // indirect
	golang.org/x/sys v0.33.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
