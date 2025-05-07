<p>
    English&nbsp ｜&nbsp <a href="README_CN.md">中文</a>
</p>

# goner/otel/tracer/http

## Overview

`goner/otel/tracer/http` is an HTTP exporter component for OpenTelemetry tracing functionality in the Gone framework. This module provides the capability to export trace data to OpenTelemetry collectors via HTTP protocol, enabling applications to conveniently collect and analyze distributed tracing data.

## Key Features

- Provides OpenTelemetry trace exporter based on HTTP protocol
- Supports both secure (TLS) and non-secure connections
- Supports custom HTTP headers
- Supports request timeout configuration
- Supports failure retry mechanism
- Integrates with Gone framework's lifecycle management

## Installation

```bash
# Install HTTP trace exporter
gonectl install goner/otel/tracer/http
```

## Configuration

| Configuration Item | Type | Description |
| --- | --- | --- |
| `otel.tracer.http.endpoint` | string | Address and port of the OpenTelemetry collector |
| `otel.tracer.http.urlPath` | string | URL path for trace reporting |
| `otel.tracer.http.insecure` | boolean | Whether to use non-secure connection (without TLS) |
| `otel.tracer.http.headers` | map | Custom HTTP headers |
| `otel.tracer.http.duration` | time | Request timeout duration |
| `otel.tracer.http.retry.enabled` | boolean | Whether to enable retry mechanism |
| `otel.tracer.http.retry.initialInterval` | time | Wait time after first failure |
| `otel.tracer.http.retry.maxInterval` | time | Maximum retry interval |
| `otel.tracer.http.retry.maxElapsedTime` | time | Maximum total time before giving up retries |

## Example
> The following example demonstrates how to export trace data using OLTP/HTTP protocol. The project includes a server and a client, with both trace data exported to Jaeger; the client calls the server via HTTP requests, passing trace information during the process.
> Complete content: [HTTP Cross-Service Tracing](../../../examples/otel/tracer/http)

### Create Application Using gonectl
```bash
gonectl create -t otel/tracer/http http-demo
cd http-demo

# Start Jaeger
# make jaeger

# Start server
# make server

# Start client
# make client
```

### View Results

After the service is running, you can view the trace data through the Jaeger UI:

1. Access the Jaeger UI interface: http://localhost:16686
2. Select the service name in the Search interface
3. Click the Find Traces button to view trace data

You can see the complete call chain, including:
- Client initiating request
- Server receiving request
- Method execution
- Response returning to client

Each span contains detailed attribute information, such as request parameters, execution time, etc.

## References

- [OpenTelemetry Official Documentation](https://opentelemetry.io/docs/)
- [OTLP/HTTP Exporter Documentation](https://opentelemetry.io/docs/specs/otlp/#otlphttp)
- [Gone Framework Documentation](https://github.com/gone-io/gone)