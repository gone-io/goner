<p>
    English&nbsp ｜&nbsp <a href="README_CN.md">中文</a>
</p>

# OpenTelemetry HTTP Log Collector

This module provides OpenTelemetry log collection functionality based on HTTP protocol, allowing applications to send log data via OTLP HTTP protocol to OpenTelemetry Collector for centralized processing and analysis.

## Features

- Supports exporting log data to OpenTelemetry Collector via HTTP protocol
- Provides flexible configuration options including endpoint, TLS, retry policies, etc.
- Seamless integration with Gone framework, easy to use
- Supports correlation with tracing (Trace) to improve log observability

## Installation

Install using Gone framework's package management tool:

```bash
gonectl install goner/otel/log/http
```

## Configuration

| Configuration Item | Type | Description | Default Value |
| --- | --- | --- | --- |
| `otel.log.http.endpoint` | string | OpenTelemetry Collector address and port | `http://localhost:4318` |
| `otel.service.name` | string | Service name to identify log source | - |
| `log.otel.enable` | boolean | Whether to enable OpenTelemetry log collection | `false` |
| `log.otel.log-name` | string | Log name to distinguish different log streams | Service name |
| `log.otel.only` | boolean | Whether to use only OpenTelemetry log collector | `false` |

## Usage Examples

### Configuration File

Add the following configuration in project's config file (e.g. `config/default.yaml`):

```yaml
service:
  name: &serviceName "my-service"

otel:
  service:
    name: *serviceName

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
    logger.Infof("Using trace.Trace to get trace ID")
    log.Infof("Using ctx logger to set trace ID")
}
```

### Expected Output

After starting the application, logs will be collected and exported via OpenTelemetry Collector. Each log will contain the following information:

- Service Name
- Log Level
- Timestamp
- Trace ID (if exists)
- Message

You can use various OpenTelemetry Collector exporters to send logs to different backend systems such as Elasticsearch, Loki or file systems.