# Gone Tracer 组件

`gone-tracer` 是 Gone 框架的分布式追踪组件，提供了统一的追踪 ID 管理和日志关联功能。通过该组件，您可以轻松地在 Gone 应用中实现分布式追踪，跟踪请求在多个服务和多个 goroutine 之间的传递，便于问题排查和性能分析。

## 功能特性

- 与 Gone 框架无缝集成
- 自动生成和传递追踪 ID
- 支持跨 goroutine 的追踪 ID 传递
- 提供 panic 恢复机制
- 简化日志关联和请求追踪
- 轻量级设计，低性能开销

## 安装

```bash
go get github.com/gone-io/goner
```

## 快速开始

### 1. 加载追踪组件

```go
package main

import (
    "github.com/gone-io/gone/v2"
    "github.com/gone-io/goner/tracer"
)

func main() {
    gone.Run(
        tracer.Load,  // 加载追踪组件
        // 其他组件...
    )
}
```

### 2. 使用追踪功能

```go
type MyService struct {
    gone.Flag
    tracer tracer.Tracer `gone:"*"`  // 注入追踪器
    logger gone.Logger   `gone:"*"`  // 注入日志器
}

func (s *MyService) DoSomething() {
    // 设置追踪 ID 并执行函数
    s.tracer.SetTraceId("", func() {
        // 获取当前 goroutine 的追踪 ID
        traceId := s.tracer.GetTraceId()
        s.logger.Infof("当前追踪 ID: %s", traceId)

        // 在新的 goroutine 中保持追踪 ID
        s.tracer.Go(func() {
            // 这里的 GetTraceId() 会返回与父 goroutine 相同的追踪 ID
            s.logger.Infof("子 goroutine 的追踪 ID: %s", s.tracer.GetTraceId())
        })
    })
}
```

## API 参考

### Tracer 接口

```go
type Tracer interface {
    // SetTraceId 设置追踪 ID，如果 traceId 为空字符串，则自动生成一个
    // 通过回调函数 fn 执行业务逻辑，在 fn 内部可以通过 GetTraceId 获取设置的追踪 ID
    SetTraceId(traceId string, fn func())

    // GetTraceId 获取当前 goroutine 的追踪 ID
    GetTraceId() string

    // Go 启动一个新的 goroutine，并传递当前的追踪 ID
    // 这个方法可以替代标准的 go 关键字，确保子 goroutine 能够继承父 goroutine 的追踪 ID
    Go(fn func())
}
```

## 最佳实践

1. 在服务入口点（如 HTTP 处理器、gRPC 服务方法等）使用 `SetTraceId` 设置追踪 ID
2. 使用 `tracer.Go` 代替标准的 `go` 关键字启动 goroutine，以保持追踪 ID 的传递
3. 在日志记录中包含追踪 ID，便于关联同一请求的不同日志条目
4. 在微服务调用中传递追踪 ID，实现跨服务的请求追踪

## 注意事项

1. 追踪 ID 是基于 goroutine 本地存储（使用 github.com/jtolds/gls 库）实现的，不会自动跨越进程边界
2. 在微服务架构中，需要手动在服务间请求中传递追踪 ID（如通过 HTTP 头或 gRPC 元数据）
3. 同一个 goroutine 中不能重复设置 traceId，如果尝试重复设置，将会保留第一次设置的值并记录警告日志
4. 当 traceId 参数为空字符串时，会自动生成一个 UUID 作为 traceId