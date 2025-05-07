<p>
    English&nbsp ｜&nbsp <a href="README_CN.md">中文</a>
</p>

# goner/zap Component

**goner/zap** is a Gone framework component that integrates [uber-go/zap](https://github.com/uber-go/zap), providing high-performance structured logging functionality.

Main features include:

- Configuration support
- Provides `*zap.Logger` injection using Gone's Provider mechanism
- Provides `gone.Logger` implementation based on zap, enhancing Gone's logging capabilities
- Integrates with [openTelemetry](https://github.com/open-telemetry/opentelemetry-go) and `goner/tracer` to provide log tracing functionality

## Configuration

| Configuration Item        | Description                                                                                                                                                                | Default Value    |
| ------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | --------------- |
| log.output                | Log output path                                                                                                                                                             | `stdout`        |
| log.error-output          | Error log output path, if not configured, error output will reuse `log.output`                                                                                                | empty           |
| log.level                 | Log level, supports `debug`,`info`,`warn`,`error`,`panic`,`fatal`, default is `info`, supports dynamic configuration via configuration center                                 | `info`          |
| log.encoder               | Log encoding format, supports `console` and `json`, default is `console`; if `zapcore.Encoder` injection is provided, this configuration will be invalid                        | `console`       |
| log.disable-stacktrace    | Whether to disable stack trace                                                                                                                                                 | `false`         |
| log.stacktrace-level      | Log level that triggers stack trace                                                                                                                                            | `error`         |
| log.report-caller         | Whether to report caller information in logs                                                                                                                                   | `true`          |
| log.rotation.output       | Log rotation output file path                                                                                                                                                  | empty           |
| log.rotation.error-output | Error log rotation output file path                                                                                                                                          | empty           |
| log.rotation.max-size     | Maximum size of log rotation files (MB)                                                                                                                                      | `100`           |
| log.rotation.max-files    | Maximum number of log rotation files to retain                                                                                                                                | `10`            |
| log.rotation.max-age      | Maximum retention days for log rotation files                                                                                                                                 | `30`            |
| log.rotation.local-time   | Whether to use local time for log rotation                                                                                                                                    | `true`          |
| log.rotation.compress     | Whether to compress old files in log rotation                                                                                                                                  | `false`         |
| log.otel.enable          | Whether to enable OpenTelemetry log integration                                                                                                                                | `false`         |
| log.otel.only            | Whether to only use OpenTelemetry for logging, without file output                                                                                                            | `true`          |
| log.otel.log-name        | OpenTelemetry log name                                                                                                                                                        | `zap`           |


## Installation
```bash
gonectl install goner/zap
```

## Using `*zap.Logger` for Logging
```go
package main

import (
	"github.com/gone-io/gone/v2"
	"go.uber.org/zap"
)

type UseOriginZap struct {
	gone.Flag
	zap *zap.Logger `gone:"*"`
}

func (s *UseOriginZap) PrintLog() {
	s.zap.Info("hello", zap.String("name", "gone io"))
}
```

## Using `goner.Logger` for Logging
```go
package main

import "github.com/gone-io/gone/v2"

type UseGoneLogger struct {
	gone.Flag
	logger gone.Logger `gone:"*"`
}

func (u *UseGoneLogger) PrintLog() {
	u.logger.Infof("hello %s", "GONE IO")
}
```

## Using `g.tracer` to Provide traceId for Logs

- Install `g.tracer` implementation
```bash
gonectl install goner/tracer/gls

# or
# gonectl install goner/tracer/gid
```

- Print logs
```go
package main

import (
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
	"go.uber.org/zap"
)

type UseTracer struct {
	gone.Flag
	logger gone.Logger `gone:"*"`
	zap    *zap.Logger `gone:"*"`
	tracer g.Tracer    `gone:"*"`
}

func (s *UseTracer) PrintLog() {
	s.tracer.SetTraceId("", func() {
		s.logger.Infof("hello with traceId")
		s.zap.Info("hello with traceId")
	})
}
```

## Custom `zapcore.Encoder`

```go
package main

import (
	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

var _ zapcore.Encoder = (*UseCustomerEncoder)(nil)

func init() {
	gone.Load(NewUseCustomerEncoder())
}

func NewUseCustomerEncoder() *UseCustomerEncoder {
	return &UseCustomerEncoder{
		Encoder: zapcore.NewConsoleEncoder(zap.NewProductionEncoderConfig()),
	}
}

type UseCustomerEncoder struct {
	zapcore.Encoder
	gone.Flag
}

func (e *UseCustomerEncoder) EncodeEntry(entry zapcore.Entry, fields []zap.Field) (*buffer.Buffer, error) {
	//do something
	return e.Encoder.EncodeEntry(entry, fields)
}
```

## Integration with OpenTelemetry
### Features
- Use OpenTelemetry tracer to provide tracerId for logs
- Use OpenTelemetry log/oltp protocol to collect logs

### Steps

#### 1. Component Installation

```bash
# Install OpenTelemetry related components
gonectl install goner/otel/log/http    # Use oltp/http/log to collect logs
gonectl install goner/otel/tracer/http # Use olte/tracer to provide traceID for logs and use oltp/http/tracer to collect trace information
gonectl install goner/zap              # Use zap for logging
gonectl install goner/viper            # Use viper for configuration
```

#### 2. Configuration Settings

Add OpenTelemetry and logging related configurations to the configuration file:

```yaml
service:
  name: &serviceName "your-service-name"

otel:
  service:
    name: *serviceName
  log:
    http:
      endpoint: localhost:4318  # OpenTelemetry Collector HTTP endpoint
      insecure: true           # Whether to use insecure connection
  tracer:
    http:
      endpoint: localhost:4318  # OpenTelemetry Collector HTTP endpoint
      insecure: true           # Whether to use insecure connection

log:
  otel:
    enable: true               # Enable OpenTelemetry log integration
    log-name: *serviceName     # Log name, usually same as service name
    only: false                # Whether to only use OpenTelemetry for logging, without file output
```

#### 3. Usage Example

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
		// Create OpenTelemetry tracer
		tracer := otel.Tracer("your-tracer-name")
		ctx, span := tracer.Start(context.Background(), "operation-name")
		defer span.End()

		// Use logger with context, automatically includes traceId
		log := ctxLogger.Ctx(ctx)
		log.Infof("Log message with traceId")

		// Use tracer to set traceId
		gTracer.SetTraceId(span.SpanContext().TraceID().String(), func() {
			// All logs within this function will include traceId
			logger.Infof("Log with traceId set via g.Tracer")
		})
	})
}
```

#### 4. Configure OpenTelemetry Collector

To collect and process logs, you need to set up OpenTelemetry Collector. Here's a basic configuration example:

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
  file:
    path: /log/log.json  # Log output path

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

You can use Docker Compose to start OpenTelemetry Collector:

```yaml
services:
  otel-collector:
    image: otel/opentelemetry-collector-contrib
    volumes:
      - ./otel-collector-config.yaml:/etc/otelcol-contrib/config.yaml
      - ./:/log/
    ports:
      - "4317:4317"  # OTLP gRPC receiver
      - "4318:4318"  # OTLP HTTP receiver
```

#### 5. View Collected Logs

Logs will be collected and saved to the path configured in OpenTelemetry Collector. You can view logs through:

- Direct log file viewing
- Export logs to log management systems like Elasticsearch, Loki
- Use tools like Grafana for visualization