<p>
    English&nbsp ｜&nbsp <a href="README_CN.md">中文</a>
</p>

# OpenTelemetry gRPC Log Collector

This module provides OpenTelemetry log collection functionality based on the gRPC protocol, allowing application log data to be sent to the OpenTelemetry Collector via OTLP gRPC protocol for centralized processing and analysis. Compared to HTTP protocol, gRPC offers higher performance and stronger reliability.

## Features

- Supports exporting log data to OpenTelemetry Collector via gRPC protocol
- Provides high-performance bidirectional streaming
- Supports message compression and flow control
- Offers flexible configuration options including endpoint, TLS, retry policies, etc.
- Seamless integration with Gone framework, easy to use
- Supports correlation with tracing (Trace) to improve log observability

## Installation

Install using Gone framework's package management tool:

```bash
gonectl install goner/otel/log/grpc
```

## Configuration

| Configuration Item | Type | Description | Default Value |
| --- | --- | --- | --- |
| `otel.log.grpc.endpoint` | string | OpenTelemetry Collector address and port | `localhost:4317` |
| `otel.log.grpc.compression` | string | gRPC compressor type (gzip, snappy, zstd) | - |
| `otel.log.grpc.retry.enabled` | boolean | Whether to enable retry mechanism | `true` |
| `otel.log.grpc.retry.max-attempts` | integer | Maximum retry attempts | `5` |
| `otel.log.grpc.retry.initial-interval` | string | Initial retry interval | `5s` |
| `otel.log.grpc.retry.max-interval` | string | Maximum retry interval | `30s` |
| `otel.service.name` | string | Service name for identifying log source | - |
| `log.otel.enable` | boolean | Whether to enable OpenTelemetry log collection | `false` |
| `log.otel.log-name` | string | Log name for distinguishing different log streams | Service name |
| `log.otel.only` | boolean | Whether to use only OpenTelemetry log collector | `false` |

## Usage Examples

### Configuration File

Add the following configuration to your project's configuration file (e.g. `config/default.yaml`):

```yaml
service:
  name: &serviceName "my-service"

otel:
  service:
    name: *serviceName
  log:
    grpc:
      endpoint: "localhost:4317"
      compression: "gzip"
      retry:
        enabled: true
        max-attempts: 5
        initial-interval: "5s"
        max-interval: "30s"

log:
  otel:
    enable: true
    log-name: *serviceName
    only: true
```

### Code Example

```go
func main() {
    gone.Run(func(logger gone.Logger, ctxLogger g.CtxLogger, gTracer g.Tracer, i struct {
        name string `gone:"config,otel.service.name"`
    }) {
        // Create a new trace
        tracer := otel.Tracer("my-tracer")
        ctx, span := tracer.Start(context.Background(), "my-operation")
        defer span.End()

        // Use context logger to record logs with trace ID
        log := ctxLogger.Ctx(ctx)
        log.Infof("This is a log with trace ID")

        // Manually set trace ID
        gTracer.SetTraceId(span.SpanContext().TraceID().String(), func() {
            doSomething(logger, log)
        })
    })
}

func doSomething(logger gone.Logger, log gone.Logger) {
    logger.Infof("Get trace ID using trace.Trace")
    log.Infof("Set trace ID using ctx logger")
}
```

### Expected Results

After starting the application, logs will be collected and exported via OpenTelemetry Collector. Each log will contain the following information:

- Service Name
- Log Level
- Timestamp
- Trace ID (if exists)
- Message

You can use various OpenTelemetry Collector exporters to send logs to different backend systems such as Elasticsearch, Loki, or file systems. Compared to HTTP protocol, gRPC protocol provides better performance and reliability, making it particularly suitable for high-throughput log collection scenarios.