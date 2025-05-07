[//]: # (desc: 使用OpenTelemetry日志收集简单示例)

<p>
    <a href="README.md">English</a>&nbsp ｜&nbsp 中文
</p>

# 使用OpenTelemetry进行日志收集

本示例展示如何在Gone框架中集成OpenTelemetry的Log功能，实现应用程序的日志收集。

## 项目构建步骤

### 1. 创建项目和安装依赖包

```bash
# 创建项目目录
mkdir simple-log
cd simple-log

# 初始化Go模块
go mod init examples/otel/log/simple

# 安装Gone框架的OpenTelemetry Log组件
gonectl install goner/otel/log
```

### 2. 实现主程序

在`main.go`中实现日志记录：

```go
package main

import "github.com/gone-io/gone/v2"

func main() {
	gone.
		Loads(GoneModuleLoad).
		Run(func(logger gone.Logger) {
			logger.Infof("hello world")
			logger.Errorf("error info")
		})
}
```

### 3. 配置OpenTelemetry

在`config/default.yaml`中配置OpenTelemetry的日志收集：

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

## 运行服务

```bash
# 运行服务
go run .
```

## 查看结果

运行服务后，日志将被发送到配置的OpenTelemetry收集器（endpoint: localhost:4318）。你可以在OpenTelemetry收集器的控制台或者配置的存储后端（如Elasticsearch、Loki等）中查看收集到的日志。

示例中的日志包含：
- 一条Info级别的日志："hello world"
- 一条Error级别的日志："error info"

这些日志会被OpenTelemetry收集器处理，并包含服务名称（log-collect-example）等元数据信息。