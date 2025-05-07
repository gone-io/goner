<p>
    <a href="README.md">English</a>&nbsp ｜&nbsp 中文
</p>

# goner/otel/meter/http

## 概述

`goner/otel/meter/http` 是 Gone 框架中用于支持 OpenTelemetry 指标（Metrics）功能的 HTTP 导出器组件。该模块提供了通过 HTTP
协议将指标数据导出到 OpenTelemetry 收集器的功能，使应用程序能够方便地将性能指标数据发送到监控系统。

## 主要功能

- 提供基于 HTTP 协议的 OpenTelemetry 指标导出器
- 支持安全连接（TLS）和非安全连接
- 支持自定义 HTTP 头信息
- 支持请求超时配置
- 支持失败重试机制
- 与 Gone 框架的生命周期管理集成

## 安装

```bash
# 安装 HTTP 导出器
gonectl install goner/otel/meter/http
```

## 配置

| 配置项                     | 类型  | 说明                      |
|-------------------------|-----|-------------------------|
| `otel.meter.http.endpoint`              | 字符串 | OpenTelemetry 收集器的地址和端口 |
| `otel.meter.http.urlPath`               | 字符串 | 指标上报的 URL 路径            |
| `otel.meter.http.insecure`              | 布尔值 | 是否使用非安全连接（不使用 TLS）      |
| `otel.meter.http.headers`               | 映射  | 自定义 HTTP 头信息            |
| `otel.meter.http.duration`              | 时间  | 请求超时时间                  |
| `otel.meter.http.retry.enabled`         | 布尔值 | 是否启用重试机制                |
| `otel.meter.http.retry.initialInterval` | 时间  | 首次失败后的等待时间              |
| `otel.meter.http.retry.maxInterval`     | 时间  | 重试间隔的最大值                |
| `otel.meter.http.retry.maxElapsedTime`  | 时间  | 放弃重试前的最大总时间             |

## 例子

> 展示如何使用 OLTP/gRPC 导出器将指标数据导出到 OpenTelemetry 收集器。
> 示例所在目录：[examples/otel/meter/http](../../../examples/otel/meter/http)

- 创建示例项目：

```bash
gonectl create -t otel/meter/grpc grpc-demo
cd grpc-demo
go mod tidy
```

- 启动OpenTelemetry 收集器

```bash
docker compose up -d 
```

- 运行

```bash
go run .
```

- 结果
  在`log.json`文件中将增加一条指标信息：

```json5
{
  "resourceMetrics": [
    {
      "resource": {
        "attributes": [
          {
            "key": "service.name",
            "value": {
              "stringValue": "meter over http"
            }
            //...
          }
        ]
      },
      "scopeMetrics": [
        {
          "scope": {},
          "metrics": [
            {
              "name": "api.counter",
              "description": "API调用的次数",
              "unit": "{次}",
              "sum": {
                "dataPoints": [
                  {
                    "startTimeUnixNano": "1746606506413972000",
                    "timeUnixNano": "1746606511419301000",
                    "asInt": "5"
                  }
                ],
                "aggregationTemporality": 2,
                "isMonotonic": true
              }
            }
          ]
        }
      ],
      "schemaUrl": "https://opentelemetry.io/schemas/1.26.0"
    }
  ]
}
```