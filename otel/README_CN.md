<p>
    <a href="README.md">English</a>&nbsp ｜&nbsp 中文
</p>

# goner/otel 是gone框架中用于支持OpenTelemetry的组件

## 概述
goner/otel 是一个完整的 OpenTelemetry 集成解决方案，提供了追踪（Traces）、指标（Metrics）和日志（Logs）的全方位支持。该组件集成了多种导出协议和格式，使您能够灵活地选择最适合您需求的可观测性解决方案。

## 主要功能
- 完整的 OpenTelemetry 协议支持
- 灵活的导出器选择（HTTP/gRPC）
- 与主流可观测性平台集成（如 Jaeger、Zipkin、Prometheus）
- 统一的配置管理
- 自动化的上下文传播

## 组件列表

### 核心组件
- [goner/otel](.) - 核心组件
  - 负责设置 **Propagator**
  - 提供统一的服务配置管理
  - 支持自动化的上下文传播

### 追踪组件（Traces）
- [goner/otel/tracer](./tracer) - 追踪基础组件
  - 对接 OpenTelemetry 的 **Tracer**
  - 提供基于标准输出的 Exporter
  - 支持自定义采样策略【暂未完成，支持中……】

- [goner/otel/tracer/http](./tracer/http)
  - 集成 OpenTelemetry OLTP/HTTP 协议
  - 支持将追踪数据导出到 OpenTelemetry Collector

- [goner/otel/tracer/grpc](./tracer/grpc)
  - 集成 OpenTelemetry OLTP/GRPC 协议
  - 提供高性能的 gRPC 导出支持

- [goner/otel/tracer/zipkin](./tracer/zipkin)
  - 支持 Zipkin 格式导出
  - 兼容现有 Zipkin 部署

### 指标组件（Metrics）
- [goner/otel/meter](./meter) - 指标基础组件
  - 对接 OpenTelemetry 的 **Meter**
  - 提供基于标准输出的 Exporter
  - 支持多种指标类型（计数器、仪表、直方图等）

- [goner/otel/meter/http](./meter/http)
  - 集成 OpenTelemetry OLTP/HTTP 协议
  - 支持将指标数据导出到 OpenTelemetry Collector

- [goner/otel/meter/grpc](./meter/grpc)
  - 支持 OpenTelemetry OLTP/GRPC 协议
  - 提供高性能的指标数据传输

- [goner/otel/meter/prometheus](./meter/prometheus)
  - 支持 OpenTelemetry Prometheus 的 `metric.Reader`
  - 提供 Prometheus 格式的指标导出

- [goner/otel/meter/prometheus/gin](./meter/prometheus/gin)
  - 基于 **goner/gin** Controller 的 **Gin** 中间件
  - 提供 HTTP 端点供 Prometheus 抓取指标

### 日志组件（Logs）
- [goner/otel/log](./log) - 日志基础组件
  - 对接 OpenTelemetry 的 **Log**
  - 支持结构化日志记录

- [goner/otel/log/http](./log/http)
  - 集成 OpenTelemetry OLTP/HTTP 协议
  - 支持将日志导出到 OpenTelemetry Collector

- [goner/otel/log/grpc](./log/grpc)
  - 集成 OpenTelemetry OLTP/GRPC 协议
  - 提供高性能的日志传输

## 组件依赖关系
![](./deps.png)

## 快速开始
> 下面例子展示如何使用 OLTP/HTTP 协议将 OpenTelemetry 数据发送到 Jaeger 服务。
> 完整示例代码位于：[quick-start](../examples/otel/tracer/quick-start)

### 1. 创建应用
使用 gonectl 创建示例应用：
```shell
gonectl create -t otel/tracer/quick-start quick-start
cd quick-start
go mod tidy
```

### 2. 启动 Jaeger 服务
使用 Docker 启动 Jaeger All-in-One 服务：
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

### 3. 运行应用
选择以下任一方式运行应用：

使用 gonectl：
```bash
gonectl run ./cmd
```

或使用标准 Go 命令：
```bash
go generate ./...
go run ./cmd
```

### 4. 查看追踪数据
启动完成后，访问 Jaeger UI：
- 打开浏览器访问：http://localhost:16686
- 选择服务并查看追踪数据

## 配置说明
每个组件都支持通过配置文件进行自定义设置。详细配置说明请参考各子模块的文档。

## 最佳实践
- 在生产环境中建议使用 OpenTelemetry Collector
- 根据实际需求选择合适的采样策略
- 合理配置导出间隔和批处理大小
- 注意性能开销，避免过度采集

