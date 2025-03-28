# CMux

`cmux`是一个用于在同一端口上复用多种协议的组件，它允许你在一个端口上同时处理HTTP和gRPC等不同协议的请求。这个组件是基于[soheilhy/cmux](https://github.com/soheilhy/cmux)实现的。

## 功能特点

- 支持在同一端口上处理多种协议
- 支持HTTP和gRPC协议的复用
- 自动协议检测和分发
- 与gone框架无缝集成

## 安装

```bash
go get github.com/gone-io/goner/cmux
```

## 配置参数

在配置文件中可以设置以下参数：

```properties
# 服务器网络类型，默认为tcp
server.network=tcp

# 服务器地址，如果不设置，将使用host和port组合
server.address=

# 服务器主机名，默认为空
server.host=

# 服务器端口号，默认为8080
server.port=8080
```

## 基础使用

1. 首先，在你的应用中加载cmux组件：

```go
package main

import (
    "github.com/gone-io/gone/v2"
    "github.com/gone-io/goner/cmux"
)

func main() {
    gone.Run(
        cmux.Load,
        // ... 其他组件
    )
}
```

2. 在HTTP服务中使用cmux：

```go
type server struct {
    gone.Flag
    keeper gone.GonerKeeper `gone:"*"`
    // ...
}

func (s *server) initListener() error {
    goner := s.keeper.GetGonerByName(cmux.Name)
    if goner != nil {
        if muxServer, ok := goner.(cmux.CMuxServer); ok {
            s.listener = muxServer.Match(cmux.HTTP1Fast())
            s.address = muxServer.GetAddress()
            return nil
        }
    }
    // 降级为普通TCP监听
    return s.createListener()
}
```

3. 在gRPC服务中使用cmux：

```go
func (s *server) initListener() error {
    goner := s.keeper.GetGonerByName(cmux.Name)
    if goner != nil {
        if muxServer, ok := goner.(cmux.CMuxServer); ok {
            s.listener = muxServer.MatchWithWriters(
                cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/grpc"),
            )
            s.address = muxServer.GetAddress()
            return nil
        }
    }
    // 降级为普通TCP监听
    return s.createListener()
}
```

## API接口

### CMuxServer

```go
type CMuxServer interface {
    // Match 根据匹配器获取对应的监听器
    Match(matcher ...cmux.Matcher) net.Listener
    
    // MatchWithWriters 根据写入匹配器获取对应的监听器
    MatchWithWriters(matcher ...cmux.MatchWriter) net.Listener
    
    // GetAddress 获取服务器地址
    GetAddress() string
}
```

## 最佳实践

1. 优先级设置：cmux组件使用`gone.HighStartPriority()`确保在其他服务之前启动。

2. 错误处理：建议实现合适的降级策略，当cmux不可用时可以回退到普通的TCP监听。

3. 协议匹配顺序：在设置多个协议匹配器时，建议按照以下顺序：
   - gRPC (HTTP/2)
   - HTTP/1.x
   - 其他协议

4. 监控和日志：cmux组件集成了gone的日志和追踪系统，可以方便地进行监控和调试。

## 注意事项

1. 确保在使用cmux时正确配置所有必要的参数。

2. 在停止服务时，记得调用Stop方法以确保资源正确释放。

3. 当使用TLS时，需要特别注意协议检测的配置。

4. 建议在开发环境中进行充分的测试，确保所有协议都能正确工作。