module github.com/gone-io/goner/otel/log/grpc

go 1.24.1

require (
	github.com/gone-io/goner/g v1.3.3
	github.com/gone-io/goner/otel v1.3.3 // indirect
	github.com/gone-io/goner/otel/log v1.3.3
	go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc v0.11.0
	go.opentelemetry.io/otel/sdk/log v0.11.0
	google.golang.org/grpc v1.72.0
)

replace (
	github.com/gone-io/goner/g => ../../../g
	github.com/gone-io/goner/otel => ../../
	github.com/gone-io/goner/otel/log => ../
)

require (
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/gone-io/gone/v2 v2.2.6
	github.com/google/uuid v1.6.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.26.1 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/otel v1.36.0 // indirect
	go.opentelemetry.io/otel/exporters/stdout/stdoutlog v0.11.0 // indirect
	go.opentelemetry.io/otel/log v0.11.0 // indirect
	go.opentelemetry.io/otel/metric v1.36.0 // indirect
	go.opentelemetry.io/otel/sdk v1.36.0 // indirect
	go.opentelemetry.io/otel/trace v1.36.0 // indirect
	go.opentelemetry.io/proto/otlp v1.5.0 // indirect
	golang.org/x/net v0.40.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.25.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20250218202821-56aae31c358a // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250218202821-56aae31c358a // indirect
	google.golang.org/protobuf v1.36.6 // indirect
)

require github.com/stretchr/testify v1.10.0

require (
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	go.uber.org/mock v0.5.2 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
