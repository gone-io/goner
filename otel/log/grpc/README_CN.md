<p>
    <a href="README.md">English</a>&nbsp ｜&nbsp 中文
</p>

# OpenTelemetry gRPC日志收集器

本模块提供了基于gRPC协议的OpenTelemetry日志收集功能，可以将应用程序的日志数据通过OTLP gRPC协议发送到OpenTelemetry Collector进行集中处理和分析。相比HTTP协议，gRPC提供了更高的性能和更强的可靠性。

## 功能特点

- 支持通过gRPC协议导出日志数据到OpenTelemetry Collector
- 提供高性能的双向流式传输
- 支持消息压缩和流控制
- 提供灵活的配置选项，包括端点、TLS、重试策略等
- 与Gone框架无缝集成，易于使用
- 支持与链路追踪（Trace）关联，提高日志的可观测性

## 安装

使用Gone框架的包管理工具安装：

```bash
gonectl install goner/otel/log/grpc
```

## 配置

| 配置项 | 类型 | 说明 | 默认值 |
| --- | --- | --- | --- |
| `otel.log.grpc.endpoint` | 字符串 | OpenTelemetry 收集器的地址和端口 | `localhost:4317` |
| `otel.log.grpc.compression` | 字符串 | gRPC压缩器类型（gzip, snappy, zstd） | - |
| `otel.log.grpc.retry.enabled` | 布尔值 | 是否启用重试机制 | `true` |
| `otel.log.grpc.retry.max-attempts` | 整数 | 最大重试次数 | `5` |
| `otel.log.grpc.retry.initial-interval` | 字符串 | 初始重试间隔 | `5s` |
| `otel.log.grpc.retry.max-interval` | 字符串 | 最大重试间隔 | `30s` |
| `otel.service.name` | 字符串 | 服务名称，用于标识日志来源 | - |
| `log.otel.enable` | 布尔值 | 是否启用OpenTelemetry日志收集 | `false` |
| `log.otel.log-name` | 字符串 | 日志名称，用于区分不同的日志流 | 服务名称 |
| `log.otel.only` | 布尔值 | 是否仅使用OpenTelemetry日志收集器 | `false` |

## 使用示例

### 配置文件

在项目的配置文件（如`config/default.yaml`）中添加以下配置：

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

### 代码示例

```go
func main() {
    gone.Run(func(logger gone.Logger, ctxLogger g.CtxLogger, gTracer g.Tracer, i struct {
        name string `gone:"config,otel.service.name"`
    }) {
        // 创建一个新的追踪
        tracer := otel.Tracer("my-tracer")
        ctx, span := tracer.Start(context.Background(), "my-operation")
        defer span.End()

        // 使用上下文记录器记录带有追踪ID的日志
        log := ctxLogger.Ctx(ctx)
        log.Infof("这是一条带有追踪ID的日志")

        // 手动设置追踪ID
        gTracer.SetTraceId(span.SpanContext().TraceID().String(), func() {
            doSomething(logger, log)
        })
    })
}

func doSomething(logger gone.Logger, log gone.Logger) {
    logger.Infof("使用trace.Trace获取追踪ID")
    log.Infof("使用ctx logger设置的追踪ID")
}
```

### 运行效果

启动应用程序后，日志将通过OpenTelemetry Collector收集并导出。每条日志都会包含以下信息：

- 服务名称（Service Name）
- 日志级别（Log Level）
- 时间戳（Timestamp）
- 追踪ID（Trace ID，如果存在）
- 日志消息（Message）

您可以使用OpenTelemetry Collector的各种导出器将日志发送到不同的后端系统，如Elasticsearch、Loki或文件系统等。相比HTTP协议，gRPC协议提供了更好的性能和可靠性，特别适合高吞吐量的日志收集场景。