<p>
    English&nbsp ｜&nbsp <a href="README_CN.md">中文</a>
</p>

# goner/otel/meter/grpc

## Overview

`goner/otel/meter/grpc` is a gRPC exporter component in the Gone framework that supports OpenTelemetry Metrics functionality. This module provides the capability to export metrics data to OpenTelemetry Collector via gRPC protocol, which offers more efficient binary transmission and streaming capabilities compared to HTTP.

## Key Features

- Provides OpenTelemetry metrics exporter based on gRPC protocol
- Supports secure (TLS) and insecure connections
- Supports data compression (e.g., gzip)
- Supports custom gRPC headers
- Supports request timeout configuration
- Supports failure retry mechanism
- Integrates with Gone framework's lifecycle management

## Installation

```bash
# Install gRPC exporter
gonectl install goner/otel/meter/grpc
```

## Configuration

| Configuration Item          | Type   | Description                     |
|-------------------------|--------|---------------------------------|
| `otel.meter.grpc.endpoint`              | string | OpenTelemetry Collector address and port |
| `otel.meter.grpc.endpointUrl`           | string | Complete endpoint URL (alternative to endpoint) |
| `otel.meter.grpc.compressor`            | string | Data compression method, e.g. "gzip" |
| `otel.meter.grpc.headers`               | map    | Custom gRPC headers             |
| `otel.meter.grpc.duration`              | time   | Request timeout duration        |
| `otel.meter.grpc.retry.enabled`         | bool   | Whether to enable retry mechanism |
| `otel.meter.grpc.retry.initialInterval` | time   | Initial wait time after first failure |
| `otel.meter.grpc.retry.maxInterval`     | time   | Maximum retry interval          |
| `otel.meter.grpc.retry.maxElapsedTime`  | time   | Maximum total time before giving up retries |

## Examples

> Demonstrates how to use OLTP/gRPC exporter to export metrics data to OpenTelemetry Collector.
> Example directory: [examples/otel/meter/grpc](../../../examples/otel/meter/grpc)

- Create example project:

```bash
gonectl create -t otel/meter/grpc grpc-demo
cd grpc-demo
go mod tidy
```

- Start OpenTelemetry Collector

```bash
docker compose up -d 
```

- Run

```bash
go run .
```

- Result
  The `log.json` file will contain a new metrics entry:

```json5
{
  "resourceMetrics": [
    {
      "resource": {
        "attributes": [
          {
            "key": "service.name",
            "value": {
              "stringValue": "meter over grpc"
            }
            //...
          }
        ]
      },
      "scopeMetrics": [
        {
          "scope": {},
          "metrics": [
            {
              "name": "api.counter",
              "description": "Number of API calls",
              "unit": "{count}",
              "sum": {
                "dataPoints": [
                  {
                    "startTimeUnixNano": "1746606506413972000",
                    "timeUnixNano": "1746606511419301000",
                    "asInt": "5"
                  }
                ],
                "aggregationTemporality": 2,
                "isMonotonic": true
              }
            }
          ]
        }
      ],
      "schemaUrl": "https://opentelemetry.io/schemas/1.26.0"
    }
  ]
}
```