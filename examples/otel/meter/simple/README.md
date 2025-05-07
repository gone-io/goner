[//]: # (desc: Metrics Monitoring with OpenTelemetry)

<p>
    English&nbsp ｜&nbsp <a href="README_CN.md">中文</a>
</p>

# Metrics Monitoring with OpenTelemetry

This example demonstrates how to integrate OpenTelemetry's Meter functionality in the Gone framework to implement application metrics monitoring.

## Project Setup Steps

### 1. Create Project and Install Dependencies

```bash
# Create project directory
mkdir simple-meter
cd simple-meter

# Initialize Go module
go mod init examples/otel/meter/simple

# Install Gone framework's OpenTelemetry Meter component
gonectl install goner/otel/meter
```

### 2. Implement Main Program

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
	_ = os.Setenv("GONE_OTEL_SERVICE_NAME", "simple meter demo")

	gone.
		NewApp(GoneModuleLoad).
		Run(func(meter metric.Meter, logger gone.Logger) {
			apiCounter, err := meter.Int64Counter(
				"api.counter",
				metric.WithDescription("API call count"),
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

## Run Service

```bash
# Run service
go run .
```

## View Results
After execution, the terminal will output metrics monitoring results:

```json
{
	"Resource": [
		{
			"Key": "service.name",
			"Value": {
				"Type": "STRING",
				"Value": "simple meter demo"
			}
		},
		{
			"Key": "telemetry.sdk.language",
			"Value": {
				"Type": "STRING",
				"Value": "go"
			}
		},
		{
			"Key": "telemetry.sdk.name",
			"Value": {
				"Type": "STRING",
				"Value": "opentelemetry"
			}
		},
		{
			"Key": "telemetry.sdk.version",
			"Value": {
				"Type": "STRING",
				"Value": "1.35.0"
			}
		}
	],
	"ScopeMetrics": [
		{
			"Scope": {
				"Name": "",
				"Version": "",
				"SchemaURL": "",
				"Attributes": null
			},
			"Metrics": [
				{
					"Name": "api.counter",
					"Description": "API call count",
					"Unit": "{count}",
					"Data": {
						"DataPoints": [
							{
								"Attributes": [],
								"StartTime": "0001-01-01T00:00:00Z",
								"Time": "0001-01-01T00:00:00Z",
								"Value": 5
							}
						],
						"Temporality": "CumulativeTemporality",
						"IsMonotonic": true
					}
				}
			]
		}
	]
}
```