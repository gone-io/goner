[//]: # (desc: Metrics monitoring with OpenTelemetry via OTLP/gRPC protocol)

<p>
    English&nbsp ｜&nbsp <a href="README_CN.md">中文</a>
</p>

# Metrics Monitoring with OpenTelemetry via OTLP/gRPC Protocol

This example demonstrates how to integrate OpenTelemetry Meter functionality in the Gone framework and export metrics data to OpenTelemetry Collector via gRPC protocol.

## Project Setup Steps

### 1. Create Project and Install Dependencies

```bash
# Create project directory
go mod init examples/otel/meter/grpc

# Install Gone framework's OpenTelemetry Meter gRPC component
# Also install viper for reading configuration files
gonectl install goner/otel/meter/grpc
gonectl install goner/viper
go mod tidy
```

### 2. Configuration File and Main Program Implementation

Create configuration file `config/default.yaml`:

```yaml
otel:
  service:
    name: "meter over grpc"
  meter:
    grpc:
      endpoint: localhost:4317
      insecure: true
```

Implement metrics monitoring in `main.go`:

```go
package main

import (
	"context"
	"github.com/gone-io/gone/v2"
	"go.opentelemetry.io/otel/metric"
	"os"
	"time"
)

func main() {
	_ = os.Setenv("GONE_OTEL_SERVICE_NAME", "grpc meter demo")

	gone.
		NewApp(GoneModuleLoad).
		Run(func(meter metric.Meter, logger gone.Logger) {
			apiCounter, err := meter.Int64Counter(
				"api.counter",
				metric.WithDescription("Number of API calls"),
				metric.WithUnit("{count}"),
			)
			if err != nil {
				logger.Errorf("create meter err:%v", err)
				return
			}

			for i := 0; i < 5; i++ {
				time.Sleep(1 * time.Second)
				apiCounter.Add(context.Background(), 1)
			}
		})
}
```

## Running the Service

### 1. Start OpenTelemetry Collector

Create `docker-compose.yaml`:

```yaml
services:
  otel-collector:
    image: otel/opentelemetry-collector-contrib
    volumes:
      - ./otel-collector-config.yaml:/etc/otelcol-contrib/config.yaml
      - ./:/log/
    ports:
      - "4317:4317" # OTLP gRPC receiver
      - "4318:4318" # OTLP http receiver
```

Create `otel-collector-config.yaml`:

```yaml
receivers:
  otlp:
    protocols:
      grpc:
        endpoint: 0.0.0.0:4317
      http:
        endpoint: 0.0.0.0:4318

exporters:
  file:
    path: /log/log.json

service:
  pipelines:
    metrics:
      receivers: [otlp]
      exporters: [file]
```

Start Collector:

```bash
docker-compose up -d
```

### 2. Run the Application

```bash
go run .
```

## View Results

Metrics data will be collected in the `log.json` file, with content similar to:

```json
{
  "Resource": [
    {
      "Key": "service.name",
      "Value": {
        "Type": "STRING",
        "Value": "grpc meter demo"
      }
    }
  ],
  "ScopeMetrics": [
    {
      "Metrics": [
        {
          "Name": "api.counter",
          "Description": "Number of API calls",
          "Unit": "{count}",
          "Data": {
            "DataPoints": [
              {
                "Value": 5
              }
            ]
          }
        }
      ]
    }
  ]
}
```