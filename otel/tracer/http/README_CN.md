<p>
    <a href="README.md">English</a>&nbsp ｜&nbsp 中文
</p>

# goner/otel/tracer/http

## 概述

`goner/otel/tracer/http` 是 Gone 框架中用于支持 OpenTelemetry 追踪（Traces）功能的 HTTP 导出器组件。该模块提供了通过 HTTP 协议将追踪数据导出到 OpenTelemetry 收集器的功能，使应用程序能够方便地将分布式追踪数据集中收集和分析。

## 主要功能

- 提供基于 HTTP 协议的 OpenTelemetry 追踪导出器
- 支持安全连接（TLS）和非安全连接
- 支持自定义 HTTP 头信息
- 支持请求超时配置
- 支持失败重试机制
- 与 Gone 框架的生命周期管理集成

## 安装

```bash
# 安装 HTTP 追踪导出器
gonectl install goner/otel/tracer/http
```

## 配置

| 配置项 | 类型 | 说明 |
| --- | --- | --- |
| `otel.tracer.http.endpoint` | 字符串 | OpenTelemetry 收集器的地址和端口 |
| `otel.tracer.http.urlPath` | 字符串 | 追踪上报的 URL 路径 |
| `otel.tracer.http.insecure` | 布尔值 | 是否使用非安全连接（不使用 TLS） |
| `otel.tracer.http.headers` | 映射 | 自定义 HTTP 头信息 |
| `otel.tracer.http.duration` | 时间 | 请求超时时间 |
| `otel.tracer.http.retry.enabled` | 布尔值 | 是否启用重试机制 |
| `otel.tracer.http.retry.initialInterval` | 时间 | 首次失败后的等待时间 |
| `otel.tracer.http.retry.maxInterval` | 时间 | 重试间隔的最大值 |
| `otel.tracer.http.retry.maxElapsedTime` | 时间 | 放弃重试前的最大总时间 |

## 例子
> 下面例子，展示如何使用OLTP/HTTP协议导出追踪数据。项目包括一个服务端和一个客户端，服务的和客户端的追踪数据导出到Jaeger；客户端通过http请求调用服务端，调用过程中传递追踪信息。
> 完整内容：[http跨服务追踪](../../../examples/otel/tracer/http)

### 使用gonectl创建应用
```bash
gonectl create -t otel/tracer/http http-demo
cd http-demo

# 启动 jaeger
# make jaeger

# 启动服务的
# make server

# 启动客户端
# make client
```

### 查看结果

服务运行后，可以通过Jaeger UI查看链路追踪数据：

1. 访问Jaeger UI界面：http://localhost:16686
2. 在Search界面选择服务名称
3. 点击Find Traces按钮查看追踪数据

你可以看到完整的调用链路，包括：
- 客户端发起请求
- 服务端接收请求
- 方法的执行
- 响应返回客户端

每个span中都包含了详细的属性信息，如请求参数、执行时间等。

## 参考资料

- [OpenTelemetry 官方文档](https://opentelemetry.io/docs/)
- [OTLP/HTTP 导出器文档](https://opentelemetry.io/docs/specs/otlp/#otlphttp)
- [Gone 框架文档](https://github.com/gone-io/gone)
