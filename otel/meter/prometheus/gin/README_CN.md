<p>
    <a href="README.md">English</a>&nbsp ｜&nbsp 中文
</p>

# goner/otel/meter/prometheus/gin

## 概述

`goner/otel/meter/prometheus/gin` 是 Gone 框架中用于支持 OpenTelemetry 指标（Metrics）功能的 Prometheus 导出器的 Gin 集成组件。该模块提供了将应用程序的指标数据以 Prometheus 格式通过 Gin 路由暴露的功能，使基于 Gin 的应用程序能够方便地与 Prometheus 监控系统集成，实现指标数据的收集、存储和可视化。

## 主要功能

- 提供基于 Gin 的 Prometheus 指标暴露端点
- 自动注册 `/metrics` 路由（可配置）
- 与 Gone 框架的生命周期管理集成
- 无需额外配置即可使用的简便接口

## 安装

```bash
# 安装 Prometheus 指标导出器的 Gin 集成组件
gonectl install goner/otel/meter/prometheus/gin
```

## 基本使用

### 在应用中加载模块

```go
func main() {
    gone.
		NewApp(gin.Load).
		Serve()
}
```

### 配置指标端点

在配置文件中添加以下配置来自定义指标暴露路径：

```yaml
otel:
  meter:
    prometheus:
      path: "/metrics"  # Prometheus 抓取端点，默认为 /metrics
```

## 示例

> 下面例子展示如何在 Gone 框架中集成 OpenTelemetry 与 Prometheus，实现应用指标的监控与可视化。项目包括一个基于 Gin 的 HTTP 服务，通过 Prometheus 抓取指标数据，并使用 Grafana 进行可视化展示。
> 完整内容：[Prometheus指标监控](../../../../examples/otel/meter/prometheus)

### 创建项目和安装依赖

```bash
# 创建项目目录
mkdir prometheus-demo
cd prometheus-demo

# 初始化Go模块
go mod init examples/prometheus-demo

# 安装 Prometheus 的 Gin 集成组件
gonectl install goner/otel/meter/prometheus/gin
```

### 实现自定义指标

创建控制器文件，实现 API 访问计数器：

```go
package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

type ctr struct {
	gone.Flag
	r g.IRoutes `gone:"*"`
}

func (c *ctr) Mount() (err g.MountError) {
	var meter = otel.Meter("my-service-meter")
	apiCounter, err := meter.Int64Counter(
		"api.counter",
		metric.WithDescription("API调用的次数"),
		metric.WithUnit("{次}"),
	)
	if err != nil {
		return gone.ToErrorWithMsg(err, "创建api.counter失败")
	}

	c.r.GET("/hello", func(ctx *gin.Context) string {
		apiCounter.Add(ctx, 1)
		return "hello, world"
	})
	return
}
```

### 查看结果

服务运行后，可以通过以下方式查看指标数据：

1. 访问指标端点：http://localhost:2112/metrics
2. 访问 Prometheus UI：http://localhost:9090
   - 在 Graph 界面输入指标名称（如 `api_counter`）
   - 点击 Execute 按钮查看指标数据
3. 访问 Grafana 界面：http://localhost:3000
   - 导入预设的 Dashboard
   - 查看指标可视化面板

## 工作原理

`goner/otel/meter/prometheus/gin` 模块通过以下方式工作：

1. 在应用启动时，自动注册一个 Gin 路由处理器，用于暴露 Prometheus 格式的指标数据
2. 当 Prometheus 服务器访问配置的端点（默认为 `/metrics`）时，模块会收集当前应用的所有指标数据并以 Prometheus 格式返回
3. 应用程序可以使用 OpenTelemetry API 创建和更新各种类型的指标（计数器、仪表、直方图等）

## 参考资料

- [OpenTelemetry 官方文档](https://opentelemetry.io/docs/)
- [Prometheus 官方文档](https://prometheus.io/docs/)
- [Gone 框架文档](https://github.com/gone-io/gone)
- [Gin Web 框架](https://github.com/gin-gonic/gin)