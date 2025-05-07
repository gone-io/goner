<p>
    <a href="README.md">English</a>&nbsp ｜&nbsp 中文
</p>

# goner/otel/tracer/zipkin

## 概述

`goner/otel/tracer/zipkin` 是 Gone 框架中用于支持 Zipkin 追踪系统的 OpenTelemetry 导出器组件。该模块提供了将 OpenTelemetry 追踪数据导出到 Zipkin 追踪系统的功能，使应用程序能够与现有的 Zipkin 生态系统集成。

## 主要功能

- 提供 Zipkin 格式的 OpenTelemetry 追踪导出器
- 支持自定义 HTTP 头信息
- 与 Gone 框架的生命周期管理集成
- 兼容 Zipkin 追踪系统的数据格式和 API

## 安装

```bash
# 安装 Zipkin 追踪导出器
gonectl install goner/otel/tracer/zipkin
```

## 配置

在应用程序中使用 Zipkin 追踪导出器，需要在 Gone 框架的配置文件中添加相关配置：

```yaml
otel:
  service:
    name: "your-service-name"  # 设置服务名称
  tracer:
    zipkin:
      url: "http://your-zipkin-endpoint/api/v2/spans"  # Zipkin 接收端点
      headers:  # 可选，自定义 HTTP 头信息
        Authorization: "Bearer your-token"
```

### 配置选项说明

| 配置项 | 类型 | 说明 |
| --- | --- | --- |
| `url` | 字符串 | Zipkin 接收端点的完整 URL |
| `headers` | 映射 | 自定义 HTTP 头信息 |



## 参考资料

- [OpenTelemetry 官方文档](https://opentelemetry.io/docs/)
- [Zipkin 官方文档](https://zipkin.io/)
- [Gone 框架文档](https://github.com/gone-io/gone)