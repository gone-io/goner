<p>
    English&nbsp ｜&nbsp <a href="README_CN.md">中文</a>
</p>

# goner/gin  Component

The **goner/gin** Component is a web framework wrapper based on [gin-gonic/gin](https://github.com/gin-gonic/gin), providing HTTP service support for the Gone framework.

## Features

- **Route Management**: Complete RESTful routing support with route grouping
- **Middleware**: Built-in common middleware with support for custom middleware development
- **Parameter Injection**: Automatic parameter injection from HTTP requests into structs
- **SSE Support**: Native support for Server-Sent Events
- **Error Handling**: Unified error handling mechanism
- **Performance Optimization**: Built-in request rate limiting, connection pooling, and other optimizations
- **Observability**: Built-in request logging and distributed tracing support

## Installation

```bash
go get github.com/gone-io/goner/gin
```

## Quick Start

### Basic Routing Example

```go
package main

import (
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner"
	"github.com/gone-io/goner/gin"
)

type HelloController struct {
	gone.Flag
	gin.IRouter `gone:"*"`
}

func (h *HelloController) Mount() gin.MountError {
	h.GET("/hello", h.hello)
	return nil
}

func (h *HelloController) hello() (string, error) {
	return "Hello, Gone!", nil
}

func main() {
	gone.
		Load(&HelloController{}).
		Loads(goner.GinLoad).
		Serve()
}
```

### 2. HTTP Parameter Injection

For detailed information, please refer to [HTTP Injection Guide](./docs/http-inject.md)

```go
type UserController struct {
    gone.Flag
    gin.IRouter `gone:"*"`
    http.HttInjector `gone:"http"`
}

func (u *UserController) Mount() gin.MountError {
    u.POST("/users/:id", u.createUser)
    return nil
}

// Fields in req will be automatically injected from the HTTP request
func (u *UserController) createUser(req struct {
    ID        int64  `http:"param=id"`   // Path parameter
    Name      string `http:"query=name"` // Query parameter
    Token     string `http:"header=token"`      // Request header
    SessionID string `http:"cookie=session-id"` // Cookie
    Data      User   `http:"body"`              // Request body
}) error {
    // Handle user creation logic
    return nil
}
```

### 3. Direct Data Return without calling `context.Success`

## Middleware Usage

### 1. System Middleware

The system includes several built-in middleware components:

- Request logging
- Rate limiting
- Health check
- Distributed tracing

### 2. Custom Middleware

```go
type CustomMiddleware struct {
    gone.Flag
}

func (m *CustomMiddleware) Process(ctx *gin.Context) {
    // Pre-processing
    ctx.Next()
    // Post-processing
}
```

## SSE (Server-Sent Events)

Support for server-sent events:

```go
// SSEController is an example controller showing how to use channels to return SSE streams
type SSEController struct {
    gone.Flag
    gone.Logger `gone:"gone-logger"`
    router gin.IRouter `gone:"*"`
}

// Mount implements the Controller interface to mount routes
func (c *SSEController) Mount() gin.GinMountError {
    // Register route
    c.router.GET("/api/sse/events", c.streamEvents)

    return nil
}

// streamEvents returns a channel that will be automatically converted to an SSE stream
func (c *SSEController) streamEvents() (<-chan any, error) {
    // Create a channel for sending events
    ch := make(chan any)

    // Start a goroutine to send events
    go func() {
        defer close(ch) // Ensure channel is closed when function ends

        // Send 10 events
        for i := 1; i <= 10; i++ {
            // Create event data
            eventData := map[string]any{
                "id":      i,
                "message": fmt.Sprintf("This is event #%d", i),
                "time":    time.Now().Format(time.RFC3339),
            }

            // Send event to channel
            ch <- eventData

            // Send one event per second
            time.Sleep(1 * time.Second)
        }

        // Send an error event example
        ch <- gone.NewParameterError("This is an example error event")

        // Wait one second before ending
        time.Sleep(1 * time.Second)
    }()

    return ch, nil
}
```

## Configuration

### Server Configuration

```properties
# Basic server configuration
server.port=8080                     # Server port, default 8080
server.host=0.0.0.0                  # Server host address, default 0.0.0.0
server.mode=release                  # Server mode, options: debug, release, test, default release
server.max-wait-before-stop=5s       # Maximum wait time before server shutdown, default 5 seconds

server.address=                      # Server address in host:port format, if set, host and port are ignored
server.html-tpl-pattern=             # HTML template file pattern for loading HTML templates
```

### Logging Configuration

```properties
# Logging configuration
server.log.format=console                # Log format, default console
server.log.show-request-time=true        # Show request time, default true
server.log.show-request-log=true         # Show request log, default true
server.log.data-max-length=0             # Maximum log data length, 0 means no limit, default 0
server.log.request-id=true               # Record request ID, default true
server.log.remote-ip=true                # Record remote IP, default true
server.log.request-body=true             # Record request body, default true
server.log.user-agent=true               # Record User-Agent, default true
server.log.referer=true                  # Record Referer, default true
server.log.show-response-log=true        # Show response log, default true

# Request body log content type configuration
server.log.show-request-body-for-content-types=application/json;application/xml;application/x-www-form-urlencoded

# Response body log content type configuration
server.log.show-response-body-for-content-types=application/json;application/xml;application/x-www-form-urlencoded
```

### Rate Limiting Configuration

```properties
# Rate limiting configuration
server.req.enable-limit=false        # Enable request rate limiting, default false
server.req.limit=100                 # Requests per second limit, default 100
server.req.limit-burst=300           # Burst request limit, default 300
server.req.x-request-id-key=X-Request-Id  # Request ID header key
server.req.x-trace-id-key=X-Trace-Id      # Trace ID header key
```

### Health Check and Tracing Configuration

```properties
# Health check
server.health-check=/health          # Health check path, default /health

server.is-after-proxy=false          # Whether behind a proxy, default false; set to true if behind a reverse proxy like Nginx
```

### Proxy and Response Configuration

```properties
# Proxy statistics
server.proxy.stat=false              # Enable proxy statistics, default false

# Response wrapping
server.return.wrapped-data=true      # Wrap response data, default true
```

### Service Registration Configuration

```properties
# Service registration configuration
server.service-name=                 # Service name for registration, must be set
server.service-use-subnet=0.0.0.0/0  # Subnet for service registration, default 0.0.0.0/0, used to select IP address for registration
```

## Service Registration and Discovery

Gone Gin component supports service registration and discovery, allowing services to be registered with a service registry for easy discovery and invocation by other services.

### Service Registration Process

When a service starts, it automatically registers its information with the service registry and deregisters when the service shuts down. The registration process is as follows:

1. Get the list of local IP addresses
2. Filter IP addresses based on the configured subnet
3. Register the service using the filtered IP address and port number
4. Automatically deregister when the service shuts down

### Service Registration Example

```go
package main

import (
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner"
	"github.com/gone-io/goner/gin"
	"github.com/gone-io/goner/g"
)

// Implementing service registration requires injecting g.ServiceRegistry
type ServiceRegistry struct {
	gone.Flag
	registry g.ServiceRegistry `gone:"*"`
}

func main() {
	gone.
		NewApp(
            goner.GinLoad, // Load Gone Gin component
            nacos.RegistryLoad, // Load Nacos registry component
            viper.Load, // Load Viper configuration component
            // Load other components
            // ...
        ).
		Serve()
}
```

### Configuration Parameters

- **server.service-name**: Service name for registration, must be set
- **server.service-use-subnet**: Subnet for service registration, default 0.0.0.0/0, used to select IP address for registration

## Best Practices

1. Route Management
    - Organize controllers by business modules
    - Use route groups to manage related endpoints
    - Use HTTP methods appropriately (GET, POST, PUT, DELETE, etc.)

2. Parameter Injection
    - Use HTTP injection tags appropriately (param, query, header, cookie, body)
    - Only one body injection per request
    - Add validation rules for injected parameters

3. Error Handling
    - Use `gone.Error` for unified error handling
    - Handle exceptions uniformly in middleware
    - Define clear error codes for different types of errors

4. Performance Optimization
    - Configure rate limiting parameters appropriately
    - Set appropriate log levels
    - Use connection pools for resource management

## Performance Testing

See [Performance Test Report](./benchmark_test.md)