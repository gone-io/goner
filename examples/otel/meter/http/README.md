[//]: # (desc: OpenTelemetry Metrics Monitoring via OTLP/HTTP)

<p>
    English&nbsp ｜&nbsp <a href="README_CN.md">中文</a>
</p>

# OpenTelemetry Metrics Monitoring via OTLP/HTTP

This example demonstrates how to integrate OpenTelemetry Meter functionality in the Gone framework and export metrics data to OpenTelemetry Collector via HTTP protocol.

## Project Setup Steps

### 1. Create Project and Install Dependencies

```bash
# Create project directory
mkdir http-meter
go mod init examples/otel/meter/http

# Install Gone framework's OpenTelemetry Meter HTTP component
# Also install viper for reading configuration files
gonectl install goner/otel/meter/http
gonectl install goner/viper
go mod tidy
```

### 2. Configuration File and Main Program Implementation

Create configuration file `config/default.yaml`:

```yaml
otel:
  service:
    name: "meter over http"
  meter:
    http:
      endpoint: localhost:4318
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
	_ = os.Setenv("GONE_OTEL_SERVICE_NAME", "http meter demo")

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
      - "4318:4318" # OTLP http receiver
```

Create `otel-collector-config.yaml`:

```yaml
receivers:
  otlp:
    protocols:
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

## Viewing Results

Metrics data will be collected in `log.json` file with content similar to:

```json
{
  "Resource": [
    {
      "Key": "service.name",
      "Value": {
        "Type": "STRING",
        "Value": "http meter demo"
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