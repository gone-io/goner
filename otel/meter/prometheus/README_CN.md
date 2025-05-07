<p>
    <a href="README.md">English</a>&nbsp ｜&nbsp 中文
</p>

# goner/otel/meter/prometheus

## 概述

`goner/otel/meter/prometheus` 是 Gone 框架中用于支持 OpenTelemetry 指标（Metrics）功能的 Prometheus 导出器组件。该模块提供了将应用程序的指标数据以 Prometheus 格式暴露的功能，使应用程序能够方便地与 Prometheus 监控系统集成，实现指标数据的收集、存储和可视化。

## 主要功能

- 提供 Prometheus 格式的指标读取器
- 支持自定义指标名称和标签
- 支持多种指标类型（计数器、仪表、直方图等）
- 支持指标单位配置
- 支持指标描述信息
- 与 Gone 框架的生命周期管理集成

## 安装

```bash
# 安装 Prometheus 指标导出器
gonectl install goner/otel/meter/prometheus
```


## 例子
> 下面例子展示如何在Gone框架中集成OpenTelemetry与Prometheus，实现应用指标的监控与可视化。项目包括一个HTTP服务，通过Prometheus抓取指标数据，并使用Grafana进行可视化展示。
> 完整内容：[Prometheus指标监控](../../../examples/otel/meter/prometheus)

### 使用gonectl创建应用
```bash
gonectl create -t otel/meter/prometheus prometheus-demo
cd prometheus-demo

# 启动 Prometheus
# make prometheus

# 启动应用
# make run
```

### 查看结果

服务运行后，可以通过以下方式查看指标数据：

1. 访问指标端点：http://localhost:2112/metrics
2. 访问Prometheus UI：http://localhost:9090
   - 在Graph界面输入指标名称
   - 点击Execute按钮查看指标数据
3. 访问Grafana界面：http://localhost:3000
   - 导入预设的Dashboard
   - 查看指标可视化面板

你可以看到完整的指标数据，包括：
- 应用程序的基本指标
- 自定义业务指标
- 系统资源使用情况
- 请求处理统计等

每个指标都包含详细的标签信息和帮助说明，方便进行数据分析和告警配置。

## 参考资料

- [OpenTelemetry 官方文档](https://opentelemetry.io/docs/)
- [Prometheus 官方文档](https://prometheus.io/docs/)
- [Gone 框架文档](https://github.com/gone-io/gone)