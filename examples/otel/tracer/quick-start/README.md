[//]: # (desc: OpenTelemetry Quick Start Example for Distributed Tracing)

<p>
    English&nbsp ｜&nbsp <a href="README_CN.md">中文</a>
</p>

# OpenTelemetry Quick Start Example for Distributed Tracing

This is a quick start example demonstrating how to integrate OpenTelemetry distributed tracing with the gone-io/goner framework.

## Features

This example shows how to quickly integrate OpenTelemetry distributed tracing in a gone application, including:

- Creating and using Tracer
- Creating Spans and adding events
- Passing context between function calls
- Configuring OTLP HTTP exporter

## Configuration

Configure OpenTelemetry in `config/default.yaml`:

```yaml
otel:
  service:
    name: "http-hello-client"  # Service name
  tracer:
    http:
      endpoint: localhost:4318  # OTLP HTTP export endpoint
      insecure: true           # Whether to use insecure connection
```

## Code Example

Main code is located in `cmd/main.go`:

```go
package main

import (
	"context"
	"github.com/gone-io/gone/v2"
	"go.opentelemetry.io/otel"
)

const tracerName = "demo"

//go:generate gonectl generate -m . -s ..
func main() {
	gone.Run(func() {
		tracer := otel.Tracer(tracerName)
		ctx, span := tracer.Start(context.Background(), "run demo")
		defer span.End()
		span.AddEvent("x event")
		doSomething(ctx)
	})
}

func doSomething(ctx context.Context) {
	tracer := otel.Tracer(tracerName)
	_, span := tracer.Start(ctx, "doSomething")
	defer span.End()
}
```

## Running Instructions

- Start Jaeger service
```bash
docker run -d --name jaeger \
  -e COLLECTOR_ZIPKIN_HOST_PORT=:9411 \
  -p 6831:6831/udp \
  -p 6832:6832/udp \
  -p 5778:5778 \
  -p 16686:16686 \
  -p 4317:4317 \
  -p 4318:4318 \
  -p 14250:14250 \
  -p 14268:14268 \
  -p 14269:14269 \
  -p 9411:9411 \
  jaegertracing/all-in-one:1.55
```

- Run the code
```bash
gonectl run ./cmd
```
or
```bash
go generate ./...
go run ./cmd
```

- View trace data
After startup, access Jaeger UI:
- Open browser and visit: http://localhost:16686
- Select service and view trace data

## Dependencies

This example uses the following goner modules:

- `github.com/gone-io/goner/otel/tracer/http` - OpenTelemetry HTTP exporter
- `github.com/gone-io/goner/viper` - Configuration management

## More Information

For more information about OpenTelemetry, please refer to the [OpenTelemetry Documentation](https://opentelemetry.io/docs/).