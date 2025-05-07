[//]: # (desc: 使用OpenTelemetry进行指标监控)

<p>
    <a href="README.md">English</a>&nbsp ｜&nbsp 中文
</p>

# 使用OpenTelemetry进行指标监控

本示例展示如何在Gone框架中集成OpenTelemetry的Meter功能，实现应用程序的指标监控。

## 项目构建步骤

### 1. 创建项目和安装依赖包

```bash
# 创建项目目录
mkdir simple-meter
cd simple-meter

# 初始化Go模块
go mod init examples/otel/meter/simple

# 安装Gone框架的OpenTelemetry Meter组件
gonectl install goner/otel/meter
```

### 2. 实现主程序

在`main.go`中实现指标监控：

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

## 运行服务

```bash
# 运行服务
go run .
```

## 查看结果
代码运行结束，会在终端输出指标监控的结果：

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