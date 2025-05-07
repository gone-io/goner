otel.tracer.grpc.<p>
    English&nbsp ｜&nbsp <a href="README_CN.md">中文</a>
</p>

# goner/otel/tracer/grpc

## Overview

`goner/otel/tracer/grpc` is a gRPC exporter component in the Gone framework that supports OpenTelemetry tracing functionality. This module provides the capability to export trace data to OpenTelemetry collectors via gRPC protocol, enabling applications to conveniently collect and analyze distributed tracing data in a centralized manner.

## Key Features

- Provides OpenTelemetry trace exporter based on gRPC protocol
- Supports secure (TLS) and insecure connections
- Supports data compression (e.g., gzip)
- Supports custom gRPC headers
- Supports request timeout configuration
- Supports failure retry mechanism
- Integrates with Gone framework's lifecycle management

## Installation

```bash
# Install gRPC trace exporter
gonectl install goner/otel/tracer/grpc
```

## Configuration

| Configuration Item | Type | Description |
| --- | --- | --- |
| `otel.tracer.grpc.endpoint` | string | Address and port of OpenTelemetry collector |
| `otel.tracer.grpc.insecure` | boolean | Whether to use insecure connection (no TLS) |
| `otel.tracer.grpc.compression` | string | Data compression method, e.g., "gzip" |
| `otel.tracer.grpc.headers` | map | Custom gRPC headers |
| `otel.tracer.grpc.duration` | time | Request timeout duration |
| `otel.tracer.grpc.retry.enabled` | boolean | Whether to enable retry mechanism |
| `otel.tracer.grpc.retry.initialInterval` | time | Wait time after first failure |
| `otel.tracer.grpc.retry.maxInterval` | time | Maximum retry interval |
| `otel.tracer.grpc.retry.maxElapsedTime` | time | Maximum total time before giving up retries |

## Example
> The following example demonstrates how to export trace data using OLTP/gRPC protocol. The project includes a server and a client, both exporting trace data to Jaeger. The client calls the server via gRPC requests, passing trace information during the call process.
> Complete content: [Cross-service Tracing with gRPC](../../../examples/otel/tracer/grpc)

### Creating Application with gonectl
```bash
gonectl create -t otel/tracer/grpc grpc-demo
cd grpc-demo

# Start jaeger
# make jaeger

# Start server
# make server

# Start client
# make client
```

### View Results

After the service is running, you can view the trace data through the Jaeger UI:

1. Access Jaeger UI: http://localhost:16686
2. Select service name in the Search interface
3. Click Find Traces button to view trace data

You can see the complete call chain, including:
- Client initiating request
- Server receiving request
- Method execution
- Response returning to client

Each span contains detailed attribute information, such as request parameters, execution time, etc.

## References

- [OpenTelemetry Documentation](https://opentelemetry.io/docs/)
- [OTLP/gRPC Exporter Documentation](https://opentelemetry.io/docs/specs/otlp/#otlpgrpc)
- [Gone Framework Documentation](https://github.com/gone-io/gone)