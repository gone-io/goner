<p>
    English&nbsp ｜&nbsp <a href="README_CN.md">中文</a>
</p>

# goner/otel/meter/http

## Overview

`goner/otel/meter/http` is an HTTP exporter component in the Gone framework that supports OpenTelemetry Metrics functionality. This module provides the capability to export metric data to OpenTelemetry collectors via HTTP protocol, enabling applications to easily send performance metrics to monitoring systems.

## Key Features

- Provides OpenTelemetry metrics exporter based on HTTP protocol
- Supports both secure (TLS) and insecure connections
- Supports custom HTTP headers
- Configurable request timeout
- Supports failure retry mechanism
- Integrated with Gone framework's lifecycle management

## Installation

```bash
# Install HTTP exporter
gonectl install goner/otel/meter/http
```

## Configuration

| Configuration Item         | Type   | Description                      |
|-------------------------|--------|---------------------------------|
| `otel.meter.http.endpoint`              | string | Address and port of OpenTelemetry collector |
| `otel.meter.http.urlPath`               | string | URL path for metrics reporting  |
| `otel.meter.http.insecure`              | bool   | Whether to use insecure connection (without TLS) |
| `otel.meter.http.headers`               | map    | Custom HTTP headers             |
| `otel.meter.http.duration`              | time   | Request timeout duration        |
| `otel.meter.http.retry.enabled`         | bool   | Whether to enable retry mechanism |
| `otel.meter.http.retry.initialInterval` | time   | Wait time after first failure   |
| `otel.meter.http.retry.maxInterval`     | time   | Maximum retry interval          |
| `otel.meter.http.retry.maxElapsedTime`  | time   | Maximum total time before giving up retries |

## Example

> Demonstrates how to use OLTP/gRPC exporter to export metrics data to OpenTelemetry collector.
> Example directory: [examples/otel/meter/http](../../../examples/otel/meter/http)

- Create example project:

```bash
gonectl create -t otel/meter/grpc grpc-demo
cd grpc-demo
go mod tidy
```

- Start OpenTelemetry collector

```bash
docker compose up -d 
```

- Run

```bash
go run .
```

- Result
  The `log.json` file will contain an additional metric entry:

```json5
{
  "resourceMetrics": [
    {
      "resource": {
        "attributes": [
          {
            "key": "service.name",
            "value": {
              "stringValue": "meter over http"
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
              "description": "API call count",
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