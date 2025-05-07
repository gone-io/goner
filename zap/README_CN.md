<p>
    <a href="README.md">English</a>&nbsp ｜&nbsp 中文
</p>

# goner/zap 组件

**goner/zap** 组件 是整合[uber-go/zap](https://github.com/uber-go/zap)的 Gone 框架组件，提供了高性能的结构化日志记录功能。

主要包括功能：

- 配置支持
- 利用Gone的Provider机制，提供`*zap.Logger`的注入
- 提供给予zap实现的`gone.Logger`，增强Gone的日志记录功能。
- 整合[openTelemetry](https://github.com/open-telemetry/opentelemetry-go)和`goner/tracer`，提供日志追踪功能。

## 配置

| 配置项                    | 说明                                                                                                | 默认值            |
| ------------------------- | --------------------------------------------------------------------------------------------------- | ----------------- |
| log.output                | 日志输出路径                                                                                        | `stdout`,标准输出 |
| log.error-output          | 错误日志输出路径，如果不配置，错误输出将复用`log.output`                                            | 空                |
| log.level                 | 日志级别，支持`debug`,`info`,`warn`,`error`,`panic`,`fatal`，默认为`info`，支持配置中心动态配置     | `info`            |
| log.encoder               | 日志编码格式，支持`console`和`json`，默认为`console`；如果提供了`zapcore.Encoder`注入，该配置将无效 | `console`         |
| log.disable-stacktrace    | 是否禁用堆栈跟踪                                                                                    | `false`           |
| log.stacktrace-level      | 触发堆栈跟踪的日志级别                                                                              | `error`           |
| log.report-caller         | 是否在日志中报告调用者信息                                                                          | `true`            |
| log.rotation.output       | 日志轮转输出文件路径                                                                                | 空                |
| log.rotation.error-output | 错误日志轮转输出文件路径                                                                            | 空                |
| log.rotation.max-size     | 日志轮转文件的最大大小（MB）                                                                        | `100`             |
| log.rotation.max-files    | 日志轮转保留的最大文件数                                                                            | `10`              |
| log.rotation.max-age      | 日志轮转文件的最大保留天数                                                                          | `30`              |
| log.rotation.local-time   | 日志轮转是否使用本地时间                                                                            | `true`            |
| log.rotation.compress     | 日志轮转是否压缩旧文件                                                                              | `false`           |
| log.otel.enable           | 是否启用OpenTelemetry日志集成                                                                       | `false`           |
| log.otel.only             | 是否仅使用OpenTelemetry记录日志，不输出到文件                                                       | `true`            |
| log.otel.log-name         | OpenTelemetry日志名称                                                                               | `zap`             |


## 安装
```bash
gonectl install goner/zap
```

## 使用`*zap.Logger`来输出日志
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

## 使用`goner.Logger`来输出日志
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

## 使用 `g.tracer`给日志提供traceId

- 安装`g.tracer`实现
```bash
gonectl install goner/tracer/gls

# 或者
# gonectl install goner/tracer/gid
```

- 打印日志
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

## 自定义`zapcore.Encoder`

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

## 和 OpenTelemetry 集成
### 功能
- 使用 OpenTelemetry tracer 给日志提供 tracerId
- 使用 OpenTelemetry log/oltp 协议收集日志

### 操作步骤

#### 1. 组件安装

```bash
# 安装 OpenTelemetry 相关组件
gonectl install goner/otel/log/http    # 使用 oltp/http/log 收集日志
gonectl install goner/otel/tracer/http # 使用 olte/tracer 给日志提供 traceID，并使用 oltp/http/tracer 收集 trace 信息
gonectl install goner/zap              # 使用 zap 打印日志
gonectl install goner/viper            # 使用 viper 配置文件
```

#### 2. 配置设置

在配置文件中添加 OpenTelemetry 和日志相关配置：

```yaml
service:
  name: &serviceName "your-service-name"

otel:
  service:
    name: *serviceName
  log:
    http:
      endpoint: localhost:4318  # OpenTelemetry Collector 的 HTTP 端点
      insecure: true            # 是否使用非安全连接
  tracer:
    http:
      endpoint: localhost:4318  # OpenTelemetry Collector 的 HTTP 端点
      insecure: true            # 是否使用非安全连接

log:
  otel:
    enable: true                # 启用 OpenTelemetry 日志集成
    log-name: *serviceName      # 日志名称，通常与服务名称相同
    only: false                 # 是否仅使用 OpenTelemetry 记录日志，不输出到文件
```

#### 3. 使用示例

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
		// 创建 OpenTelemetry tracer
		tracer := otel.Tracer("your-tracer-name")
		ctx, span := tracer.Start(context.Background(), "operation-name")
		defer span.End()

		// 使用带有 context 的日志记录器，自动包含 traceId
		log := ctxLogger.Ctx(ctx)
		log.Infof("日志消息带有 traceId")

		// 使用 tracer 设置 traceId
		gTracer.SetTraceId(span.SpanContext().TraceID().String(), func() {
			// 在此函数内的所有日志都将包含 traceId
			logger.Infof("通过 g.Tracer 设置 traceId 的日志")
		})
	})
}
```

#### 4. 配置 OpenTelemetry Collector

要收集和处理日志，您需要设置 OpenTelemetry Collector。以下是一个基本的配置示例：

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
    path: /log/log.json  # 日志输出路径

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

您可以使用 Docker Compose 启动 OpenTelemetry Collector：

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

#### 5. 查看收集的日志

日志将被收集并保存到 OpenTelemetry Collector 配置的路径中。您可以通过以下方式查看日志：

- 直接查看日志文件
- 将日志导出到 Elasticsearch、Loki 等日志管理系统
- 使用 Grafana 等工具进行可视化展示
