[//]: # (desc: OpenTelemetry 链路追踪快速启动示例)

<p>
    <a href="README.md">English</a>&nbsp ｜&nbsp 中文
</p>

# OpenTelemetry 链路追踪快速启动示例

这是一个使用 gone-io/goner 框架集成 OpenTelemetry 进行链路追踪的快速启动示例。

## 功能介绍

本示例展示了如何在 gone 应用中快速集成 OpenTelemetry 链路追踪功能，包括：

- 创建和使用 Tracer
- 创建 Span 并添加事件
- 在函数调用间传递上下文
- 配置 OTLP HTTP 导出器

## 配置说明

在 `config/default.yaml` 中配置 OpenTelemetry：

```yaml
otel:
  service:
    name: "http-hello-client"  # 服务名称
  tracer:
    http:
      endpoint: localhost:4318  # OTLP HTTP 导出端点
      insecure: true            # 是否使用非安全连接
```

## 代码示例

主要代码位于 `cmd/main.go`：

```go
package main

import (
	"context"
	"github.com/gone-io/gone/v2"
	"go.opentelemetry.io/otel"
)

const tracerName = "demo"

//go:generate gonectl generate -m . -s ..
func main() {
	gone.Run(func() {
		tracer := otel.Tracer(tracerName)
		ctx, span := tracer.Start(context.Background(), "run demo")
		defer span.End()
		span.AddEvent("x event")
		doSomething(ctx)
	})
}

func doSomething(ctx context.Context) {
	tracer := otel.Tracer(tracerName)
	_, span := tracer.Start(ctx, "doSomething")
	defer span.End()
}
```

## 运行说明

- 启动Jaeger服务
```bash
docker run -d --name jaeger \
  -e COLLECTOR_ZIPKIN_HOST_PORT=:9411 \
  -p 6831:6831/udp \
  -p 6832:6832/udp \
  -p 5778:5778 \
  -p 16686:16686 \
  -p 4317:4317 \
  -p 4318:4318 \
  -p 14250:14250 \
  -p 14268:14268 \
  -p 14269:14269 \
  -p 9411:9411 \
  jaegertracing/all-in-one:1.55
```

- 运行代码
```bash
gonectl run ./cmd
```
或者
```bash
go generate ./...
go run ./cmd
```

- 查看追踪数据
启动完成后，访问 Jaeger UI：
- 打开浏览器访问：http://localhost:16686
- 选择服务并查看追踪数据

## 依赖模块

本示例使用了以下 goner 模块：

- `github.com/gone-io/goner/otel/tracer/http` - OpenTelemetry HTTP 导出器
- `github.com/gone-io/goner/viper` - 配置管理

## 更多信息

有关 OpenTelemetry 的更多信息，请参考 [OpenTelemetry 官方文档](https://opentelemetry.io/docs/)。

