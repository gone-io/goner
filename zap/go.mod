module github.com/gone-io/goner/zap

go 1.24.1

require (
	github.com/gone-io/gone/v2 v2.2.6
	github.com/stretchr/testify v1.10.0
	go.uber.org/zap v1.27.0
	gopkg.in/natefinch/lumberjack.v2 v2.2.1
)

require (
	github.com/gone-io/goner/g v1.3.3
	go.opentelemetry.io/contrib/bridges/otelzap v0.10.0
	go.opentelemetry.io/otel/log v0.11.0
	go.opentelemetry.io/otel/trace v1.35.0
	go.uber.org/mock v0.5.2
)

replace github.com/gone-io/goner/g => ../g

require (
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/otel v1.35.0 // indirect
	go.opentelemetry.io/otel/metric v1.35.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
