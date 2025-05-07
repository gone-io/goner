[//]: # (desc: OpenTelemetry 链路追踪简单示例)

<p>
    <a href="README.md">English</a>&nbsp ｜&nbsp 中文
</p>

# OpenTelemetry 链路追踪简单示例

这是一个使用 gone-io/goner 框架集成 OpenTelemetry 进行链路追踪的简单示例。

## 功能介绍

本示例展示了如何在 gone 应用中集成 OpenTelemetry 链路追踪功能，包括：

- 在组件中注入和使用 Tracer
- 创建 Span 并添加事件
- 在函数调用间传递上下文
- 通过环境变量配置服务名称

## 代码结构

- `main.go` - 应用入口，设置服务名称并启动应用
- `your_component.go` - 示例组件，展示如何使用 OpenTelemetry Tracer
- `module.load.go` - 模块加载配置

## 代码示例

### 主程序 (main.go)

```go
package main

import (
	"context"
	"github.com/gone-io/gone/v2"
	"os"
)

func main() {
	//设置服务名称
	_ = os.Setenv("GONE_OTEL_SERVICE_NAME", "simple demo")

	gone.
		Loads(GoneModuleLoad).
		Load(&YourComponent{}).
		Run(func(c *YourComponent) {
			// 调用组件中的方法
			c.HandleRequest(context.Background())
		})
}
```

### 组件实现 (your_component.go)

```go
package main

import (
	"context"
	"github.com/gone-io/gone/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type YourComponent struct {
	gone.Flag
	tracer trace.Tracer `gone:"*,otel-tracer"` // 注入 OpenTelemetry Tracer
}

func (c *YourComponent) HandleRequest(ctx context.Context) {
	tracer := otel.Tracer("otel-tracer")

	// 创建新的 Span
	ctx, span := tracer.Start(ctx, "handle-request")
	// 确保在函数结束时结束 Span
	defer span.End()

	// 记录事件
	span.AddEvent("开始处理请求")

	// 处理业务逻辑...

	// 记录错误（如果有）
	// span.RecordError(err)
	// span.SetStatus(codes.Error, "处理请求失败")

	// 正常情况下设置状态为成功
	// span.SetStatus(codes.Ok, "")
}
```

## 运行说明

1. 确保已安装 Go 环境
2. 启动 OpenTelemetry Collector（或其他兼容的后端服务）
3. 运行示例：

```bash
go run .
```

## 依赖模块

本示例使用了以下 goner 模块：

- `github.com/gone-io/goner/otel/tracer` - OpenTelemetry 链路追踪支持

## 更多信息

有关 OpenTelemetry 的更多信息，请参考 [OpenTelemetry 官方文档](https://opentelemetry.io/docs/)。

有关 gone-io/goner 框架的更多信息，请参考 [gone-io/goner 文档](https://github.com/gone-io/gone)。

