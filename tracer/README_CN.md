<p>
    <a href="README.md">English</a>&nbsp ｜&nbsp 中文
</p>

# Goner Tracer 组件

## 组件功能

`goner/tracer` 组件提供一个简单的功能：**让日志自动带上traceID（需要日志组件配合，已经支持的日志组件`goner/zap`），无需手动传递context**。

它解决的核心问题：
- 在Gone框架中实现日志自动关联traceID
- 无需在每个函数中显式传递`context.Context`
- 支持跨goroutine传递traceID

## 快速上手

### 安装

```bash
# 选择一种实现方式安装
gonectl install goner/tracer/gls  # 功能完整的实现
# 或
gonectl install goner/tracer/gid  # 高性能的实现
```

### 基本使用

```go
type MyService struct {
    gone.Flag
    tracer tracer.Tracer `gone:"*"`  // 注入追踪器
    logger gone.Logger   `gone:"*"`  // 注入日志器
}

func (s *MyService) DoSomething() {
    // 设置trace ID
    s.tracer.SetTraceId("", func() {  // 空字符串表示自动生成traceID
        // 日志会自动包含traceID
        s.logger.Info("处理业务逻辑...")

        // 如果需要获取当前traceID
        traceId := s.tracer.GetTraceId()
        s.logger.Infof("当前trace ID: %s", traceId)

        // 跨goroutine传递traceID - 重要！
        s.tracer.Go(func() {
            // 这里的日志也会带上相同的traceID
            s.logger.Info("子goroutine处理中...")
        })
    })
}
```

## 核心API

```go
type Tracer interface {
    // 设置traceID并执行回调函数
    SetTraceId(traceId string, fn func())

    // 获取当前goroutine的traceID
    GetTraceId() string

    // 创建新goroutine并保持traceID (必须用它替代go关键字)
    Go(fn func())
}
```

## 两种实现方式对比

| 实现方式                       | 优点           | 缺点     | 适用场景   |
| ------------------------------ | -------------- | -------- | ---------- |
| **gls实现** (goner/tracer/gls) | 无需获取gid    | 性能较低 | 一般场景   |
| **gid实现** (goner/tracer/gid) | 性能高(约10倍) | 无       | 高并发场景 |


**性能测试**

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

## 重要使用注意

1. **必须用`tracer.Go()`替代标准`go`关键字**，否则子goroutine将无法获取到traceID

2. **入口点设置traceID**：在HTTP处理器等入口点使用`SetTraceId`

3. **跨服务传递**：如果需要跨服务传递traceID，需要手动添加到请求头：
   ```go
   // 发送请求时
   req.Header.Set("X-Trace-ID", tracer.GetTraceId())
   
   // 接收请求时
   traceId := r.Header.Get("X-Trace-ID")
   tracer.SetTraceId(traceId, func() {
       // 处理请求...
   })
   ```

4. **与OpenTelemetry集成**：如果安装了`goner/otel/tracer`，traceID会由OpenTelemetry生成并自动跨进程传递，但如需完整的OpenTelemetry功能（metrics等），仍需按规范显式传递`context.Context`