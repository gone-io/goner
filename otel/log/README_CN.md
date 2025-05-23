<p>
    <a href="README.md">English</a>&nbsp ｜&nbsp 中文
</p>

# goner/otel/log

## 概述

`goner/otel/log` 是 Gone 框架中用于支持 OpenTelemetry 日志功能的组件。该模块提供了与 OpenTelemetry
日志系统的集成，使应用程序能够创建、记录和导出日志数据，帮助开发者更好地监控和诊断分布式系统中的问题。

## 主要功能

- 提供 OpenTelemetry 日志系统的初始化和配置
- 支持创建和管理日志记录
- 与 Gone 框架的生命周期管理集成，确保应用关闭时正确刷新和关闭日志系统
- 支持多种日志导出方式
- 与 OpenTelemetry 的 Trace 和 Metrics 无缝集成

## 子模块

- `goner/otel/log/http`: 提供基于 HTTP 协议的 OpenTelemetry 日志数据导出器
- `goner/otel/log/grpc`: 提供基于 gRPC 协议的 OpenTelemetry 日志数据导出器

## 安装方法

```bash
# 安装基础日志模块
gonectl install goner/otel/log

# 安装 HTTP 导出器
gonectl install goner/otel/log/http

# 安装 gRPC 导出器
gonectl install goner/otel/log/grpc
```

## 简单例子

> 展示通过`goner/otel/log`组件将日志数据导出到 OpenTelemetry Collector

### 执行下面命令

```bash
gonectl create -t otel/log/simple simple-demo
cd simple-demo
go run .
```

### 项目目录结构

```log
.
├── config/
│   └── default.yaml
├── go.mod
├── go.sum
├── main.go
└── module.load.go
```

### 代码

- [module.load.go](../../examples/otel/log/simple/module.load.go)，通过运行`gonectl install goner/otel/log` 安装生成。

```go
// Code generated by gonectl. DO NOT EDIT.
package main

import(
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
	"github.com/gone-io/goner/otel/log"
	"github.com/gone-io/goner/viper"
	zap "github.com/gone-io/goner/zap"
)

// load installed gone module LoadFunc
var loaders = []gone.LoadFunc{
	log.Register,
	viper.Load,
	zap.Load,
}

func GoneModuleLoad(loader gone.Loader) error {
	var ops []*g.LoadOp
	for _, f := range loaders {
		ops = append(ops, g.F(f))
	}
	return g.BuildOnceLoadFunc(ops...)(loader)
}
```

- [main.go](../../examples/otel/log/simple/main.go)，函数入口

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

- [config/default.yaml](../../examples/otel/log/simple/config/default.yaml)，配置文件

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

## 高级配置

### 配置说明

在配置文件中，可以通过以下配置项来控制日志的行为：

```yaml
log:
  otel:
    enable: true       # 是否启用 OpenTelemetry 日志收集
    log-name: "app"    # 日志名称
    only: false        # 是否仅使用 OpenTelemetry 日志（不输出到控制台）
```

### 在自定义组件中检查程序是否使用OpenTelemetry/Log，在没有使用的情况下降级处理
```go
type YourComponent struct {
    isOtelLogLoaded g.IsOtelLogLoaded `gone:"*"`
}

func(s*YourComponent)BusinessLogic(){
	if s.isOtelLogLoaded{
		// 处理日志逻辑
    } else{
	    // 降级处理，不使用 OpenTelemetry 日志
    }
}
```

## 最佳实践

1. **合理配置日志级别**：根据环境（开发、测试、生产）设置合适的日志级别
2. **结构化日志**：使用结构化的方式记录日志，便于后续分析
3. **关联追踪信息**：在分布式系统中，将日志与追踪信息关联
4. **合理使用日志字段**：添加有助于问题诊断的字段，但避免过多无用信息
5. **性能考虑**：在高并发场景下合理使用日志，避免性能瓶颈

## 常见问题

### 日志数据未正确导出

- 检查 Collector 地址和端口是否正确
- 确认网络连接是否正常
- 查看应用日志中是否有导出错误信息

### 日志与追踪信息未关联

- 确保正确配置了服务名称
- 检查 Trace 上下文是否正确传播
- 验证日志记录时是否包含了追踪信息

### 性能问题

- 调整日志级别
- 优化日志内容和结构
- 考虑使用异步日志处理

## 参考资料

- [OpenTelemetry 官方文档](https://opentelemetry.io/docs/)
- [OpenTelemetry 日志规范](https://opentelemetry.io/docs/specs/otel/logs/)
- [Gone 框架文档](https://github.com/gone-io/gone)