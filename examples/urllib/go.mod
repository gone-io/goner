module examples/urllib

go 1.24.1

require (
	github.com/gone-io/gone/v2 v2.2.6
	github.com/gone-io/goner/urllib v1.3.5
	github.com/imroc/req/v3 v3.55.0
)

replace github.com/gone-io/goner/urllib => ../../urllib

replace github.com/gone-io/goner/g => ../../g

require (
	github.com/andybalholm/brotli v1.2.0 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/gone-io/goner/g v1.3.5 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/icholy/digest v1.1.0 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/quic-go/qpack v0.5.1 // indirect
	github.com/quic-go/quic-go v0.54.1 // indirect
	github.com/refraction-networking/utls v1.8.0 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace v0.63.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.63.0 // indirect
	go.opentelemetry.io/otel v1.38.0 // indirect
	go.opentelemetry.io/otel/metric v1.38.0 // indirect
	go.opentelemetry.io/otel/trace v1.38.0 // indirect
	go.uber.org/mock v0.6.0 // indirect
	golang.org/x/crypto v0.42.0 // indirect
	golang.org/x/mod v0.28.0 // indirect
	golang.org/x/net v0.44.0 // indirect
	golang.org/x/sync v0.17.0 // indirect
	golang.org/x/sys v0.36.0 // indirect
	golang.org/x/text v0.29.0 // indirect
	golang.org/x/tools v0.36.0 // indirect
)
