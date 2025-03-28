# Gone Tracer Component

`gone-tracer` is a distributed tracing component for the Gone framework, providing unified trace IDs. With this component, you can easily implement distributed tracing in Gone applications, tracking requests across multiple services and goroutines, facilitating problem diagnosis and performance analysis.

## Features

- Seamless integration with Gone framework
- Automatic generation and propagation of trace IDs
- Support for trace ID propagation across goroutines
- Two implementation approaches: based on `github.com/jtolds/gls` and `github.com/petermattis/goid` mapping
- Simplified log correlation and request tracing
- Lightweight design with low performance overhead

## Installation

```bash
go get github.com/gone-io/goner
```

## Quick Start

### 1. Load the Tracing Component

```go
package main

import (
    "github.com/gone-io/gone/v2"
    "github.com/gone-io/goner/tracer"
)

func main() {
    gone.
    Loads(
        tracer.Load,  // Load tracing component
        // Other components...
    ).
    Run(func() {
        // Start application
    })
}
```

### 2. Use Tracing Features

```go
type MyService struct {
    gone.Flag
    tracer tracer.Tracer `gone:"*"`  // Inject tracer
    logger gone.Logger   `gone:"*"`  // Inject logger
}

func (s *MyService) DoSomething() {
    // Set trace ID and execute function
    s.tracer.SetTraceId("", func() {
        // Get trace ID of current goroutine
        traceId := s.tracer.GetTraceId()
        s.logger.Infof("Current trace ID: %s", traceId)

        // Maintain trace ID in new goroutine
        s.tracer.Go(func() {
            // GetTraceId() here will return the same trace ID as the parent goroutine
            s.logger.Infof("Child goroutine trace ID: %s", s.tracer.GetTraceId())
        })
    })
}
```

## Implementation Approaches

Gone Tracer component provides two implementation approaches:

1. **Goroutine Local Storage Based (tracer)**: Uses the `github.com/jtolds/gls` library, implementing trace ID storage and propagation through goroutine local storage mechanism.

2. **Goroutine ID Mapping Based (tracerOverGid)**: Uses the `github.com/petermattis/goid` library to get goroutine ID and maintains trace ID mapping through sync.Map.

By default, the component uses the first approach. If you prefer to use the goroutine ID-based implementation, specify it when loading the component:

```go
tracer.LoadOverGid  // Load goroutine ID-based tracing component
```

## Performance Tests

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

## API Reference

### Tracer Interface

```go
type Tracer interface {
    // SetTraceId sets the trace ID. If traceId is empty string, it automatically generates one
    // Business logic is executed through callback function fn, within which GetTraceId can be used to get the set trace ID
    SetTraceId(traceId string, fn func())

    // GetTraceId gets the trace ID of the current goroutine
    GetTraceId() string

    // Go starts a new goroutine and propagates the current trace ID
    // This method can replace the standard go keyword to ensure child goroutines inherit parent goroutine's trace ID
    Go(fn func())
}
```

## Advanced Usage

### Using in HTTP Services

Combined with Gone's gin component, HTTP request tracing can be easily implemented:

```go
func setupRouter(router gin.Router, tracer tracer.Tracer) {
    // Add middleware to set trace ID for each request
    router.Use(func(c *gin.Context) {
        // Get trace ID from request header, generate new if not exists
        traceId := c.GetHeader("X-Trace-ID")
        tracer.SetTraceId(traceId, func() {
            // Set trace ID to response header
            c.Header("X-Trace-ID", tracer.GetTraceId())
            c.Next()
        })
    })

    // Route handling
    router.GET("/api/example", func(c *gin.Context) {
        // Get trace ID directly in handler
        traceId := tracer.GetTraceId()
        // Handle business logic...
    })
}
```

### Propagating Trace ID Between Microservices

```go
// Client sending request
func (c *Client) CallService() {
    c.tracer.SetTraceId("", func() {
        // Create HTTP request
        req, _ := http.NewRequest("GET", "http://service-b/api", nil)
        // Add trace ID to request header
        req.Header.Set("X-Trace-ID", c.tracer.GetTraceId())
        // Send request
        c.httpClient.Do(req)
    })
}

// Server receiving request
func (s *Server) HandleRequest(w http.ResponseWriter, r *http.Request) {
    // Get trace ID from request header
    traceId := r.Header.Get("X-Trace-ID")
    s.tracer.SetTraceId(traceId, func() {
        // Handle request
        // ...
    })
}
```

## Best Practices

1. **Set Trace ID at Service Entry Points**: Use `SetTraceId` at entry points like HTTP handlers, gRPC service methods

2. **Use `tracer.Go` Instead of Standard `go` Keyword**: Ensure child goroutines inherit parent goroutine's trace ID

3. **Include Trace ID in Logs**: Facilitate correlation of different log entries for the same request

4. **Propagate Trace ID in Microservice Calls**: Pass trace ID via HTTP headers or gRPC metadata for cross-service request tracing

5. **Choose Appropriate Implementation**: Select the most performant implementation based on your use case

6. **Integrate with Logging Component**: Automatically add trace ID to log fields for improved traceability

## Important Notes

1. Trace IDs don't automatically cross process boundaries; manual propagation between services is required

2. In microservice architecture, use consistent header fields (like `X-Trace-ID`) for trace ID propagation

3. Avoid setting traceId multiple times in the same goroutine; the first set value will be retained

4. When traceId parameter is empty string, a UUID will be automatically generated as traceId

5. In high-concurrency scenarios, tracerOverGid implementation may offer better performance