module examples/otel/meter/prometheus

go 1.24.1

replace (
	github.com/gone-io/goner/g => ./../../../../g
	github.com/gone-io/goner/gin => ./../../../../gin
	github.com/gone-io/goner/otel => ./../../../../otel
	github.com/gone-io/goner/otel/meter => ./../../../../otel/meter
	github.com/gone-io/goner/otel/meter/prometheus => ./../../../../otel/meter/prometheus
	github.com/gone-io/goner/otel/meter/prometheus/gin => ./../../../../otel/meter/prometheus/gin
)

require (
	github.com/gin-gonic/gin v1.10.1
	github.com/gone-io/gone/v2 v2.2.6
	github.com/gone-io/goner/g v1.3.3
	github.com/gone-io/goner/otel/meter/prometheus/gin v1.3.3
	go.opentelemetry.io/otel v1.36.0
	go.opentelemetry.io/otel/metric v1.36.0
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/bytedance/sonic v1.13.3 // indirect
	github.com/bytedance/sonic/loader v0.2.4 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/cloudwego/base64x v0.1.5 // indirect
	github.com/gabriel-vasile/mimetype v1.4.9 // indirect
	github.com/gin-contrib/sse v1.1.0 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.26.0 // indirect
	github.com/goccy/go-json v0.10.5 // indirect
	github.com/gone-io/goner/gin v1.3.3 // indirect
	github.com/gone-io/goner/otel v1.3.3 // indirect
	github.com/gone-io/goner/otel/meter v1.3.3 // indirect
	github.com/gone-io/goner/otel/meter/prometheus v1.3.3 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/cpuid/v2 v2.2.10 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/pelletier/go-toml/v2 v2.2.4 // indirect
	github.com/prometheus/client_golang v1.22.0 // indirect
	github.com/prometheus/client_model v0.6.2 // indirect
	github.com/prometheus/common v0.63.0 // indirect
	github.com/prometheus/procfs v0.16.1 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.2.14 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin v0.60.0 // indirect
	go.opentelemetry.io/otel/exporters/prometheus v0.57.0 // indirect
	go.opentelemetry.io/otel/exporters/stdout/stdoutmetric v1.36.0 // indirect
	go.opentelemetry.io/otel/sdk v1.36.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v1.36.0 // indirect
	go.opentelemetry.io/otel/trace v1.36.0 // indirect
	go.uber.org/mock v0.5.2 // indirect
	golang.org/x/arch v0.17.0 // indirect
	golang.org/x/crypto v0.38.0 // indirect
	golang.org/x/net v0.40.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.25.0 // indirect
	golang.org/x/time v0.11.0 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
