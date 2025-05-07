[//]: # (desc: 使用openTelemetry收集日志)

<p>
    <a href="README.md">English</a>&nbsp ｜&nbsp 中文
</p>


# 使用OpenTelemetry收集日志

本示例展示如何在Gone框架中集成OpenTelemetry进行日志收集，实现应用日志的集中管理与分析。

## 项目构建步骤

### 1. 创建项目和安装依赖包

```bash
# 创建项目目录
mkdir log-collect
cd log-collect

# 初始化Go模块
go mod init examples/otel/collect

# 安装Gone框架的OpenTelemetry与日志收集相关组件
gonectl install goner/otel/log/http    # 使用oltp/http/log 收集日志
gonectl install goner/otel/tracer/http # 使用olte/tracer 给日志提供traceID，并且使用oltp/http/tracer 收集trace信息
gonectl install goner/zap              # 使用zap打印日志
gonectl install goner/viper            # 使用viper读取配置
```

### 2. 配置日志收集

首先，创建配置文件目录和默认配置：

```bash
mkdir config
touch config/default.yaml
```

然后，在`config/default.yaml`中配置服务名称和OpenTelemetry相关设置：

```yaml
service:
  name: &serviceName "log-collect-example"

otel:
  service:
    name: *serviceName
  log:
    http:
      endpoint: localhost:4318
      insecure: true
  tracer:
    http:
      endpoint: localhost:4318
      insecure: true

log:
  otel:
    enable: true
    log-name: *serviceName
    only: false
```

### 3. 创建OpenTelemetry Collector配置

创建`otel-collector-config.yaml`文件，配置日志收集和导出：

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
  otlp:
    endpoint: otelcol:4317
  file:
    path: /log/log.json

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

### 4. 创建Docker Compose配置

创建`docker-compose.yaml`文件，配置OpenTelemetry Collector服务：

```yaml
services:
  otel-collector:
    image: otel/opentelemetry-collector-contrib
    volumes:
      - ./otel-collector-config.yaml:/etc/otelcol-contrib/config.yaml
      - ./:/log/
    ports:
      - "1888:1888" # pprof extension
      - "8888:8888" # Prometheus metrics exposed by the Collector
      - "8889:8889" # Prometheus exporter metrics
      - "13133:13133" # health_check extension
      - "4317:4317" # OTLP gRPC receiver
      - "4318:4318" # OTLP http receiver
      - "55679:55679" # zpages extension
```

### 5. 创建服务入口

```bash
mkdir cmd
touch cmd/main.go
```

然后，在`cmd/main.go`中实现日志记录和跟踪：

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
		//logger.Infof("service name: %s", i.name)
		//logger.Infof("hello world")
		//logger.Debugf("debug info")
		//logger.Warnf("warn info")
		//logger.Errorf("error info")

		tracer := otel.Tracer("test-tracer")
		ctx, span := tracer.Start(context.Background(), "test-run")
		defer span.End()

		log := ctxLogger.Ctx(ctx)
		log.Infof("hello world with traceId")
		log.Warnf("debug info with traceId")

		//set traceId
		gTracer.SetTraceId(span.SpanContext().TraceID().String(), func() {
			doSomething(logger, log)
		})
	})
}

func doSomething(logger gone.Logger, log gone.Logger) {
	logger.Infof("get traceId by using trace.Trace")
	log.Infof("traceId setted by ctx logger")
}
```

## 运行服务

执行以下命令启动OpenTelemetry Collector和应用服务：

```bash
# 启动OpenTelemetry Collector
docker compose up -d

# 运行服务
go run ./cmd
```

## 查看结果

### 查看收集的日志

日志将被收集并保存到OpenTelemetry Collector配置的路径中（`/log/log.json`）。您可以通过以下命令查看日志内容：

```bash
docker exec -it <collector-container-id> cat /log/log.json
```

## 日志收集原理

本示例通过以下方式实现日志收集：

1. 使用OpenTelemetry的HTTP协议收集日志和跟踪信息
2. 通过Zap记录结构化日志
3. 为日志添加TraceID，实现日志与跟踪的关联
4. 使用OpenTelemetry Collector收集、处理和导出日志

通过这种方式，您可以实现应用日志的集中管理、分析和可视化，提高系统的可观测性。