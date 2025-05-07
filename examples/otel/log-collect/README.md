[//]: # (desc: Log Collection with OpenTelemetry)

<p>
    English&nbsp ｜&nbsp <a href="README_CN.md">中文</a>
</p>

# Log Collection with OpenTelemetry

This example demonstrates how to integrate OpenTelemetry with the Gone framework for log collection, enabling centralized log management and analysis.

## Project Setup Steps

### 1. Create Project and Install Dependencies

```bash
# Create project directory
mkdir log-collect
cd log-collect

# Initialize Go module
go mod init examples/otel/collect

# Install Gone framework's OpenTelemetry and log collection related components
gonectl install goner/otel/log/http    # Use oltp/http/log for log collection
gonectl install goner/otel/tracer/http # Use olte/tracer to provide traceID and oltp/http/tracer to collect trace information
gonectl install goner/zap              # Use zap for logging
gonectl install goner/viper            # Use viper for configuration
```

### 2. Configure Log Collection

First, create the configuration directory and default configuration:

```bash
mkdir config
touch config/default.yaml
```

Then, configure the service name and OpenTelemetry settings in `config/default.yaml`:

```yaml
service:
  name: &serviceName "log-collect-example"

otel:
  service:
    name: *serviceName
  log:
    http:
      endpoint: localhost:4318
      insecure: true
  tracer:
    http:
      endpoint: localhost:4318
      insecure: true

log:
  otel:
    enable: true
    log-name: *serviceName
    only: false
```

### 3. Create OpenTelemetry Collector Configuration

Create `otel-collector-config.yaml` file to configure log collection and export:

```yaml
receivers:
  otlp:
    protocols:
      grpc:
        endpoint: 0.0.0.0:4317
      http:
        endpoint: 0.0.0.0:4318
processors:
  batch:

exporters:
  otlp:
    endpoint: otelcol:4317
  file:
    path: /log/log.json

extensions:
  health_check:
  pprof:
  zpages:

service:
  extensions: [health_check, pprof, zpages]
  pipelines:
    traces:
      receivers: [otlp]
      processors: [batch]
      exporters: [file]
    logs:
      receivers: [otlp]
      processors: [batch]
      exporters: [file]
```

### 4. Create Docker Compose Configuration

Create `docker-compose.yaml` file to configure the OpenTelemetry Collector service:

```yaml
services:
  otel-collector:
    image: otel/opentelemetry-collector-contrib
    volumes:
      - ./otel-collector-config.yaml:/etc/otelcol-contrib/config.yaml
      - ./:/log/
    ports:
      - "1888:1888" # pprof extension
      - "8888:8888" # Prometheus metrics exposed by the Collector
      - "8889:8889" # Prometheus exporter metrics
      - "13133:13133" # health_check extension
      - "4317:4317" # OTLP gRPC receiver
      - "4318:4318" # OTLP http receiver
      - "55679:55679" # zpages extension
```

### 5. Create Service Entry

```bash
mkdir cmd
touch cmd/main.go
```

Then, implement logging and tracing in `cmd/main.go`:

```go
package main

import (
	"context"
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
	"go.opentelemetry.io/otel"
)

func main() {
	gone.Run(func(logger gone.Logger, ctxLogger g.CtxLogger, gTracer g.Tracer, i struct {
		name string `gone:"config,otel.service.name"`
	}) {
		//logger.Infof("service name: %s", i.name)
		//logger.Infof("hello world")
		//logger.Debugf("debug info")
		//logger.Warnf("warn info")
		//logger.Errorf("error info")

		tracer := otel.Tracer("test-tracer")
		ctx, span := tracer.Start(context.Background(), "test-run")
		defer span.End()

		log := ctxLogger.Ctx(ctx)
		log.Infof("hello world with traceId")
		log.Warnf("debug info with traceId")

		//set traceId
		gTracer.SetTraceId(span.SpanContext().TraceID().String(), func() {
			doSomething(logger, log)
		})
	})
}

func doSomething(logger gone.Logger, log gone.Logger) {
	logger.Infof("get traceId by using trace.Trace")
	log.Infof("traceId setted by ctx logger")
}
```

## Running the Service

Execute the following commands to start the OpenTelemetry Collector and the application service:

```bash
# Start OpenTelemetry Collector
docker compose up -d

# Run the service
go run ./cmd
```

## View Results

### View Collected Logs

Logs will be collected and saved to the path configured in the OpenTelemetry Collector (`/log/log.json`). You can view the log contents using the following command:

```bash
docker exec -it <collector-container-id> cat /log/log.json
```

## Log Collection Principles

This example implements log collection through the following methods:

1. Using OpenTelemetry's HTTP protocol to collect logs and trace information
2. Recording structured logs through Zap
3. Adding TraceID to logs to associate logs with traces
4. Using OpenTelemetry Collector to collect, process, and export logs

Through this approach, you can achieve centralized log management, analysis, and visualization, improving system observability.