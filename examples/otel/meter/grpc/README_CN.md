[//]: # (desc: 使用OpenTelemetry通过OTLP/gRPC协议进行指标监控)

<p>
    <a href="README.md">English</a>&nbsp ｜&nbsp 中文
</p>

# 使用OpenTelemetry通过OTLP/gRPC协议进行指标监控

本示例展示如何在Gone框架中集成OpenTelemetry的Meter功能，并通过gRPC协议将指标数据导出到OpenTelemetry Collector。

## 项目构建步骤

### 1. 创建项目和安装依赖包

```bash
# 创建项目目录
go mod init examples/otel/meter/grpc

# 安装Gone框架的OpenTelemetry Meter gRPC组件
# 同时安装viper用于读取配置文件
gonectl install goner/otel/meter/grpc
gonectl install goner/viper
go mod tidy
```

### 2. 配置文件和主程序实现

创建配置文件`config/default.yaml`:

```yaml
otel:
  service:
    name: "meter over grpc"
  meter:
    grpc:
      endpoint: localhost:4317
      insecure: true
```

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
	_ = os.Setenv("GONE_OTEL_SERVICE_NAME", "grpc meter demo")

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

### 1. 启动OpenTelemetry Collector

创建`docker-compose.yaml`:

```yaml
services:
  otel-collector:
    image: otel/opentelemetry-collector-contrib
    volumes:
      - ./otel-collector-config.yaml:/etc/otelcol-contrib/config.yaml
      - ./:/log/
    ports:
      - "4317:4317" # OTLP gRPC receiver
      - "4318:4318" # OTLP http receiver
```

创建`otel-collector-config.yaml`:

```yaml
receivers:
  otlp:
    protocols:
      grpc:
        endpoint: 0.0.0.0:4317
      http:
        endpoint: 0.0.0.0:4318

exporters:
  file:
    path: /log/log.json

service:
  pipelines:
    metrics:
      receivers: [otlp]
      exporters: [file]
```

启动Collector:

```bash
docker-compose up -d
```

### 2. 运行应用程序

```bash
go run .
```

## 查看结果

指标数据将被收集到`log.json`文件中，内容类似如下格式：

```json
{
  "Resource": [
    {
      "Key": "service.name",
      "Value": {
        "Type": "STRING",
        "Value": "grpc meter demo"
      }
    }
  ],
  "ScopeMetrics": [
    {
      "Metrics": [
        {
          "Name": "api.counter",
          "Description": "API调用的次数",
          "Unit": "{次}",
          "Data": {
            "DataPoints": [
              {
                "Value": 5
              }
            ]
          }
        }
      ]
    }
  ]
}
```