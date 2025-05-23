<p>
    <a href="README.md">English</a>&nbsp ｜&nbsp 中文
</p>

# goner/gin 组件

**goner/gin** 组件 是一个基于 [gin-gonic/gin](https://github.com/gin-gonic/gin) 的 Web 框架封装，为 Gone 框架提供 HTTP 服务支持。

## 功能特性

- **路由管理**：完整的 RESTful 路由支持，支持路由分组
- **中间件**：内置常用中间件，支持自定义中间件开发
- **参数注入**：自动从 HTTP 请求中注入参数到结构体
- **SSE 支持**：原生支持 Server-Sent Events 服务器推送
- **错误处理**：统一的错误处理机制
- **性能优化**：内置请求限流、连接池等优化
- **可观测性**：内置请求日志、分布式追踪支持

## 安装

```bash
go get github.com/gone-io/goner/gin
```

## 快速开始

### 基础路由示例

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

### 2. HTTP 参数注入

详细说明请参考 [HTTP 注入说明](./docs/http-inject_CN.md)

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

// req 中的字段会自动从 HTTP 请求中注入
func (u *UserController) createUser(req struct {
    ID        int64  `http:"param=id"`   // 路径参数
    Name      string `http:"query=name"` // 查询参数
    Token     string `http:"header=token"`      // 请求头
    SessionID string `http:"cookie=session-id"` // Cookie
    Data      User   `http:"body"`              // 请求体
}) error {
    // 处理创建用户的逻辑
    return nil
}
```

### 3.直接返回数据，不需要调用`context.Success`

## 中间件使用

### 1. 系统中间件

系统已内置了一些常用中间件，包括：

- 请求日志记录
- 限流
- 健康检查
- 分布式追踪

### 2. 自定义中间件

```go
type CustomMiddleware struct {
    gone.Flag
}

func (m *CustomMiddleware) Process(ctx *gin.Context) {
    // 前置处理
    ctx.Next()
    // 后置处理
}
```

## SSE（Server-Sent Events）

支持服务器发送事件：

```go
// SSEController 是一个示例控制器，展示如何使用channel返回SSE流
type SSEController struct {
    gone.Flag
    gone.Logger `gone:"gone-logger"`
    router gin.IRouter `gone:"*"`
}

// Mount 实现Controller接口，挂载路由
func (c *SSEController) Mount() gin.GinMountError {
    // 注册路由
    c.router.GET("/api/sse/events", c.streamEvents)

    return nil
}

// streamEvents 返回一个channel，系统会自动将其转换为SSE流
func (c *SSEController) streamEvents() (<-chan any, error) {
    // 创建一个channel用于发送事件
    ch := make(chan any)

    // 启动一个goroutine来发送事件
    go func() {
        defer close(ch) // 确保在函数结束时关闭channel

        // 发送10个事件
        for i := 1; i <= 10; i++ {
            // 创建事件数据
            eventData := map[string]any{
                "id":      i,
                "message": fmt.Sprintf("这是第%d个事件", i),
                "time":    time.Now().Format(time.RFC3339),
            }

            // 发送事件到channel
            ch <- eventData

            // 每秒发送一个事件
            time.Sleep(1 * time.Second)
        }

        // 发送一个错误事件示例
        ch <- gone.NewParameterError("这是一个错误事件示例")

        // 等待一秒后结束
        time.Sleep(1 * time.Second)
    }()

    return ch, nil
}
```

## 配置说明

### 服务器配置

```properties
# 服务器基本配置
server.port=8080                     # 服务器端口，默认8080
server.host=0.0.0.0                  # 服务器主机地址，默认0.0.0.0
server.mode=release                  # 服务器模式，可选：debug, release, test，默认release
server.max-wait-before-stop=5s       # 服务器关闭前最大等待时间，默认5秒

server.address=                      # 服务器地址，格式为host:port，如果设置了此项，则忽略host和port
server.html-tpl-pattern=             # HTML模板文件匹配模式，用于加载HTML模板
```

### 日志配置

```properties
# 日志配置
server.log.format=console                # 日志格式，默认console
server.log.show-request-time=true        # 是否显示请求时间，默认true
server.log.show-request-log=true         # 是否显示请求日志，默认true
server.log.data-max-length=0             # 日志数据最大长度，0表示不限制，默认0
server.log.request-id=true               # 是否记录请求ID，默认true
server.log.remote-ip=true                # 是否记录远程IP，默认true
server.log.request-body=true             # 是否记录请求体，默认true
server.log.user-agent=true               # 是否记录User-Agent，默认true
server.log.referer=true                  # 是否记录Referer，默认true
server.log.show-response-log=true        # 是否显示响应日志，默认true

# 请求体日志内容类型配置
server.log.show-request-body-for-content-types=application/json;application/xml;application/x-www-form-urlencoded

# 响应体日志内容类型配置
server.log.show-response-body-for-content-types=application/json;application/xml;application/x-www-form-urlencoded
```

### 限流配置

```properties
# 限流配置
server.req.enable-limit=false        # 是否启用请求限流，默认false
server.req.limit=100                 # 每秒请求限制数，默认100
server.req.limit-burst=300           # 突发请求限制数，默认300
server.req.x-request-id-key=X-Request-Id  # 请求ID的Header键名
server.req.x-trace-id-key=X-Trace-Id      # 追踪ID的Header键名
```

### 健康检查与追踪配置

```properties
# 健康检查
server.health-check=/health          # 健康检查路径，默认为/health

server.is-after-proxy=false          # 是否在代理后面，默认false；如果存在反向代理，比如Nginx，则需要设置为true
```

### 代理与响应配置

```properties
# 代理统计
server.proxy.stat=false              # 是否启用代理统计，默认false

# 响应包装
server.return.wrapped-data=true      # 是否包装响应数据，默认true
```

### 服务注册与发现配置

```properties
# 服务注册配置
server.service-name=                 # 服务名称，用于服务注册，必须设置
server.service-use-subnet=0.0.0.0/0  # 服务使用的子网，默认0.0.0.0/0，用于选择注册的IP地址
```

## 服务注册与发现

Gone Gin 组件支持服务注册与发现功能，可以将服务注册到服务注册中心，方便其他服务发现和调用。

### 服务注册原理

服务启动时，会自动将服务信息注册到服务注册中心，服务关闭时会自动注销。注册过程如下：

1. 获取本地IP地址列表
2. 根据配置的子网过滤IP地址
3. 使用过滤后的IP地址和端口号注册服务
4. 服务关闭时自动注销

### 服务注册示例

```go
package main

import (
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner"
	"github.com/gone-io/goner/gin"
	"github.com/gone-io/goner/g"
)

// 实现服务注册需要注入 g.ServiceRegistry
type ServiceRegistry struct {
	gone.Flag
	registry g.ServiceRegistry `gone:"*"`
}

func main() {
	gone.
		NewApp(
            goner.GinLoad, // 加载Gone Gin组件
            nacos.RegistryLoad, // 加载Nacos注册中心组件
            viper.Load, // 加载Viper配置组件
            // 加载其他组件
            // ...
        ).
		Serve()
}
```

### 配置说明

- **server.service-name**：服务名称，用于服务注册，必须设置
- **server.service-use-subnet**：服务使用的子网，默认0.0.0.0/0，用于选择注册的IP地址

## 最佳实践

1. 路由管理
    - 按业务模块划分控制器
    - 使用路由分组管理相关接口
    - 合理使用 HTTP 方法（GET、POST、PUT、DELETE 等）

2. 参数注入
    - 合理使用 HTTP 注入标签（param、query、header、cookie、body）
    - 一个请求只能注入一次 body
    - 为注入的参数添加验证规则

3. 错误处理
    - 使用 `gone.Error` 统一处理错误
    - 在中间件中统一处理异常
    - 为不同类型的错误定义清晰的错误码

4. 性能优化
    - 合理配置限流参数
    - 适当配置日志级别
    - 使用连接池管理资源

4. **错误处理**：
   - 统一错误格式
   - 为错误定义清晰的错误码

## 性能测试

详见 [性能测试报告](./benchmark_test_CN.md)
