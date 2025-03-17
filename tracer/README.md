# Gone Tracer 组件

`gone-tracer` 是 Gone 框架的分布式追踪组件，用于提供了统一的trace ID。通过该组件，您可以轻松地在 Gone 应用中实现分布式追踪，跟踪请求在多个服务和多个 goroutine 之间的传递，便于问题排查和性能分析。

## 功能特性

- 与 Gone 框架无缝集成
- 自动生成和传递trace ID
- 支持跨 goroutine 的trace ID 传递
- 提供两种实现方式：基于`github.com/jtolds/gls`和基于`github.com/petermattis/goid`映射
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
    gone.
    Loads(
        tracer.Load,  // 加载追踪组件
        // 其他组件...
    ).
    Run(func() {
        // 启动应用
    })
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
    // 设置trace ID 并执行函数
    s.tracer.SetTraceId("", func() {
        // 获取当前 goroutine 的trace ID
        traceId := s.tracer.GetTraceId()
        s.logger.Infof("当前trace ID: %s", traceId)

        // 在新的 goroutine 中保持trace ID
        s.tracer.Go(func() {
            // 这里的 GetTraceId() 会返回与父 goroutine 相同的trace ID
            s.logger.Infof("子 goroutine 的trace ID: %s", s.tracer.GetTraceId())
        })
    })
}
```

## 实现方式

Gone Tracer 组件提供了两种实现方式：

1. **基于goroutine本地存储 (tracer)**：使用 `github.com/jtolds/gls` 库实现，通过goroutine本地存储机制保存和传递追踪ID。

2. **基于goroutine ID映射 (tracerOverGid)**：使用 `github.com/petermattis/goid` 库获取goroutine ID，并通过sync.Map维护goroutine ID与追踪ID的映射关系。

默认情况下，组件使用第一种实现方式。如果您希望使用基于goroutine ID的实现，可以在加载组件时指定：

```go
tracer.LoadOverGid  // 加载基于goroutine ID的追踪组件
```

## 性能测试

```log
➜  goner go test -bench=. -benchmem ./tracer
goos: darwin
goarch: arm64
pkg: github.com/gone-io/goner/tracer
cpu: Apple M1 Pro
BenchmarkTracer_SetTraceId-8             1479470               833.5 ns/op           976 B/op         11 allocs/op
BenchmarkTracerOverGid_SetTraceId-8     17734480                67.41 ns/op           64 B/op          2 allocs/op
BenchmarkTracer_GetTraceId-8             1533403               783.8 ns/op           128 B/op          1 allocs/op
BenchmarkTracerOverGid_GetTraceId-8     120586562                9.443 ns/op           0 B/op          0 allocs/op
BenchmarkTracer_Go-8                      365029              4421 ns/op             987 B/op         12 allocs/op
BenchmarkTracerOverGid_Go-8              1693129               709.2 ns/op           157 B/op          5 allocs/op
BenchmarkTracer_Concurrent-8               30535             39531 ns/op           12665 B/op        148 allocs/op
BenchmarkTracerOverGid_Concurrent-8       252675              4841 ns/op            1222 B/op         41 allocs/op
BenchmarkTracer_Nested-8                  120984              9931 ns/op            2592 B/op         28 allocs/op
BenchmarkTracerOverGid_Nested-8          5440866               221.3 ns/op           168 B/op          7 allocs/op
PASS
ok      github.com/gone-io/goner/tracer 17.094s
```

## API 参考

### Tracer 接口

```go
type Tracer interface {
    // SetTraceId 设置trace ID，如果 traceId 为空字符串，则自动生成一个
    // 通过回调函数 fn 执行业务逻辑，在 fn 内部可以通过 GetTraceId 获取设置的trace ID
    SetTraceId(traceId string, fn func())

    // GetTraceId 获取当前 goroutine 的trace ID
    GetTraceId() string

    // Go 启动一个新的 goroutine，并传递当前的trace ID
    // 这个方法可以替代标准的 go 关键字，确保子 goroutine 能够继承父 goroutine 的trace ID
    Go(fn func())
}
```

## 高级用法

### 在HTTP服务中使用

结合Gone的gin组件，可以轻松实现HTTP请求的追踪：

```go
func setupRouter(router gin.Router, tracer tracer.Tracer) {
    // 添加中间件，为每个请求设置追踪ID
    router.Use(func(c *gin.Context) {
        // 从请求头获取追踪ID，如果没有则生成新的
        traceId := c.GetHeader("X-Trace-ID")
        tracer.SetTraceId(traceId, func() {
            // 将追踪ID设置到响应头
            c.Header("X-Trace-ID", tracer.GetTraceId())
            c.Next()
        })
    })

    // 路由处理
    router.GET("/api/example", func(c *gin.Context) {
        // 在处理函数中可以直接获取追踪ID
        traceId := tracer.GetTraceId()
        // 处理业务逻辑...
    })
}
```

### 在微服务间传递追踪ID

```go
// 客户端发送请求
func (c *Client) CallService() {
    c.tracer.SetTraceId("", func() {
        // 创建HTTP请求
        req, _ := http.NewRequest("GET", "http://service-b/api", nil)
        // 将追踪ID添加到请求头
        req.Header.Set("X-Trace-ID", c.tracer.GetTraceId())
        // 发送请求
        c.httpClient.Do(req)
    })
}

// 服务端接收请求
func (s *Server) HandleRequest(w http.ResponseWriter, r *http.Request) {
    // 从请求头获取追踪ID
    traceId := r.Header.Get("X-Trace-ID")
    s.tracer.SetTraceId(traceId, func() {
        // 处理请求
        // ...
    })
}
```

## 最佳实践

1. **在服务入口点设置追踪ID**：在HTTP处理器、gRPC服务方法等入口点使用`SetTraceId`设置追踪ID

2. **使用`tracer.Go`代替标准的`go`关键字**：确保子goroutine能够继承父goroutine的追踪ID

3. **在日志中包含追踪ID**：便于关联同一请求的不同日志条目

4. **在微服务调用中传递追踪ID**：通过HTTP头或gRPC元数据传递追踪ID，实现跨服务的请求追踪

5. **选择合适的实现方式**：根据应用场景选择性能更优的实现方式

6. **结合日志组件使用**：将追踪ID自动添加到日志字段中，提高日志的可追踪性

## 注意事项

1. 追踪ID不会自动跨越进程边界，需要手动在服务间请求中传递

2. 在微服务架构中，建议使用统一的头部字段（如`X-Trace-ID`）传递追踪ID

3. 同一个goroutine中不要重复设置traceId，如果尝试重复设置，将会保留第一次设置的值

4. 当traceId参数为空字符串时，会自动生成一个UUID作为traceId

5. 在高并发场景下，tracerOverGid实现可能会有更好的性能表现