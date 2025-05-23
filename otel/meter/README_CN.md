<p>
    <a href="README.md">English</a>&nbsp ｜&nbsp 中文
</p>

# goner/otel/meter

## 概述

`goner/otel/meter` 是 Gone 框架中用于支持 OpenTelemetry 指标收集功能的组件。该模块提供了与 OpenTelemetry
指标系统的集成，使应用程序能够创建、记录和导出各种指标数据，帮助开发者监控应用性能、资源使用情况和业务指标。

## 主要功能

- 提供 OpenTelemetry 指标系统的初始化和配置
- 支持创建和管理各类指标（Counter、Gauge、Histogram等）
- 默认提供标准输出（stdout）的指标数据导出器
- 与 Gone 框架的生命周期管理集成，确保应用关闭时正确刷新和关闭指标系统
- 支持多种指标导出格式和协议

## 子模块

- `goner/otel/meter/http`: 提供基于 HTTP 协议的 OpenTelemetry 指标数据导出器
- `goner/otel/meter/grpc`: 提供基于 gRPC 协议的 OpenTelemetry 指标数据导出器
- `goner/otel/meter/prometheus`: 提供与 Prometheus 兼容的指标数据导出器
  - `goner/otel/meter/prometheus/gin`: 提供 Gin 框架集成的 Prometheus 指标暴露接口

## 安装方法

```bash
# 安装基础指标模块
gonectl install goner/otel/meter

# 安装 HTTP 导出器
gonectl install goner/otel/meter/http

# 安装 gRPC 导出器
gonectl install goner/otel/meter/grpc

# 安装 Prometheus 导出器
gonectl install goner/otel/meter/prometheus

# 安装 Prometheus Gin 集成
gonectl install goner/otel/meter/prometheus/gin
```

> 如果只安装`goner/otel/meter`, 则默认使用标准输出（stdout）的指标数据导出器；`goner/otel/meter/http`、
`goner/otel/meter/grpc`、`goner/otel/meter/prometheus`只需要根据实际情况安装其中一个即可，并且已经依赖了`goner/otel/meter`
> ，可以不用再手动安装`goner/otel/meter`。

## 简单例子

> 展示通过`goner/otel/meter`组件将指标数据导出到标准输出（stdout）

### 执行下面命令

```bash
gonectl create -t otel/meter/simple simple-demo
cd simple-demo
go run .
```

### 项目目录结构

```log
.
├── go.mod
├── go.sum
├── main.go 
└── module.load.go
```

### 代码

- [module.load.go](../../examples/otel/meter/simple/module.load.go)，通过运行`gonectl install goner/otel/meter` 安装生成。

```go
// Code generated by gonectl. DO NOT EDIT.
package main

import(
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
	"github.com/gone-io/goner/otel/meter"
)

// load installed gone module LoadFunc
var loaders = []gone.LoadFunc{
	meter.Register,
}

func GoneModuleLoad(loader gone.Loader) error {
	var ops []*g.LoadOp
	for _, f := range loaders {
		ops = append(ops, g.F(f))
	}
	return g.BuildOnceLoadFunc(ops...)(loader)
}
```

- [main.go](../../examples/otel/meter/simple/main.go)，函数入口

```go
package main

import (
	"context"
	"github.com/gone-io/gone/v2"
	"go.opentelemetry.io/otel/metric"
	"os"
	"time"
)

func main() {
	_ = os.Setenv("GONE_OTEL_SERVICE_NAME", "simple meter demo")

	gone.
		NewApp(GoneModuleLoad).
		Run(func(meter metric.Meter, logger gone.Logger) {
			apiCounter, err := meter.Int64Counter(
				"api.counter",
				metric.WithDescription("API调用的次数"),
				metric.WithUnit("{次}"),
			)
			if err != nil {
				logger.Errorf("create meter err:%v", err)
				return
			}

			for i := 0; i < 5; i++ {
				time.Sleep(1 * time.Second)
				apiCounter.Add(context.Background(), 1)
			}
		})
}
```

### 运行结果

```json
{
	"Resource": [
		{
			"Key": "service.name",
			"Value": {
				"Type": "STRING",
				"Value": "simple meter demo"
			}
		},
		{
			"Key": "telemetry.sdk.language",
			"Value": {
				"Type": "STRING",
				"Value": "go"
			}
		},
		{
			"Key": "telemetry.sdk.name",
			"Value": {
				"Type": "STRING",
				"Value": "opentelemetry"
			}
		},
		{
			"Key": "telemetry.sdk.version",
			"Value": {
				"Type": "STRING",
				"Value": "1.35.0"
			}
		}
	],
	"ScopeMetrics": [
		{
			"Scope": {
				"Name": "",
				"Version": "",
				"SchemaURL": "",
				"Attributes": null
			},
			"Metrics": [
				{
					"Name": "api.counter",
					"Description": "API调用的次数",
					"Unit": "{次}",
					"Data": {
						"DataPoints": [
							{
								"Attributes": [],
								"StartTime": "0001-01-01T00:00:00Z",
								"Time": "0001-01-01T00:00:00Z",
								"Value": 5
							}
						],
						"Temporality": "CumulativeTemporality",
						"IsMonotonic": true
					}
				}
			]
		}
	]
}
```

## 使用其他Exporter
- [使用OLTP/HTTP](../../examples/otel/meter/http)
- [使用OLTP/gRPC](../../examples/otel/meter/grpc)
- [使用Prometheus](../../examples/otel/meter/prometheus)

## 高级配置

### 自定义导出器

如果需要自定义指标数据导出器，可以实现 `metric.Exporter` 接口并注入到框架中：

```go
var _ metric.Exporter = (*CustomExporter)(nil)

type CustomExporter struct {
// 实现 metric.Exporter 接口
}

// 定义加载函数
func RegisterCustomExporter(loader gone.Loader) {
    return loader.Load(&CustomExporter{})
}

//....
// 在启动过程中加载
// gone.Loads(RegisterCustomExporter)
//...
```

### 在自定义组件中检查程序是否使用OpenTelemetry/Meter，在没有使用的情况下降级处理
```go
type YourComponent struct {
    isOtelMeterLoaded g.IsOtelMeterLoaded `gone:"*"`
}

func(s*YourComponent)BusinessLogic(){
	if s.isOtelMeterLoaded{
		// 处理指标收集逻辑
    } else{
	    // 降级处理，不使用指标收集	
    }
}
```

## 最佳实践

1. **合理命名指标**：使用有意义的名称，反映指标的实际含义。
2. **设置适当的单位和描述**：为指标添加清晰的单位和描述信息，便于理解。
3. **选择合适的指标类型**：根据业务需求选择合适的指标类型（Counter、Gauge、Histogram等）。
4. **添加有意义的标签**：使用标签对指标进行分类和过滤，但避免过多的标签组合。
5. **定期刷新指标数据**：确保指标数据能够及时导出和查看。

## 常见问题

### 指标数据未正确导出

- 检查收集器地址和端口是否正确
- 确认网络连接是否正常
- 查看应用日志中是否有导出错误信息

### 指标数据不完整

- 检查指标创建和记录是否正确
- 确保应用运行足够长的时间以收集数据
- 验证指标收集的频率是否合适

### 性能问题

- 减少高频率记录的指标数量
- 优化批处理配置
- 减少不必要的标签和属性

## 参考资料

- [OpenTelemetry 官方文档](https://opentelemetry.io/docs/)
- [指标系统原理](https://opentelemetry.io/docs/concepts/signals/metrics/)
- [Gone 框架文档](https://github.com/gone-io/gone)