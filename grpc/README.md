# Gone gRPC 组件使用指南

本文档介绍如何在 Gone 框架中使用 gRPC 组件，包括传统方式和 Gone V2 的 Provider 机制两种实现方式。

## 准备工作

首先创建一个 grpc 目录，在这个目录中初始化一个 golang mod：

```bash
mkdir grpc
cd grpc
go mod init grpc_demo
```

## 编写 proto 文件，生成 golang 代码

### 编写协议文件

定义一个简单的 Hello 服务，包含一个 Say 方法：

文件名：proto/hello.proto
```proto
syntax = "proto3";

option go_package="/proto";

package Business;

service Hello {
  rpc Say (SayRequest) returns (SayResponse);
}

message SayResponse {
  string Message = 1;
}

message SayRequest {
  string Name = 1;
}
```

### 生成 golang 代码

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
protoc --go_out=. --go_opt=paths=source_relative \
--go-grpc_out=. --go-grpc_opt=paths=source_relative \
proto/hello.proto
```

> 其中，protoc 的安装参考 [Protocol Buffer 编译器安装](https://blog.csdn.net/waitdeng/article/details/139248507)

## 实现方式一：传统方式

### 服务端实现

文件名：v1_server/main.go
```go
package main

import (
  "context"
  "github.com/gone-io/gone/v2"
  "github.com/gone-io/goner"
  goneGrpc "github.com/gone-io/goner/grpc"
  "google.golang.org/grpc"
  "grpc_demo/proto"
  "log"
)

type server struct {
  gone.Flag
  proto.UnimplementedHelloServer // 嵌入UnimplementedHelloServer
}

// 重载协议中定义的服务
func (s *server) Say(ctx context.Context, in *proto.SayRequest) (*proto.SayResponse, error) {
  log.Printf("Received: %v", in.GetName())
  return &proto.SayResponse{Message: "Hello " + in.GetName()}, nil
}

// 实现 goneGrpc.Service接口的RegisterGrpcServer方法，该方法在服务器启动时会被自动调用
func (s *server) RegisterGrpcServer(server *grpc.Server) {
  proto.RegisterHelloServer(server, s)
}

func main() {
  gone.
    Load(&server{}).
    Loads(goner.BaseLoad, goneGrpc.ServerLoad).
    // 启动服务
    Serve()
}
```

### 客户端实现

文件名：v1_client/main.go
```go
package main

import (
	"context"
	"fmt"
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner"
	gone_grpc "github.com/gone-io/goner/grpc"
	"google.golang.org/grpc"
	"grpc_demo/proto"
	"log"
)

type helloClient struct {
	gone.Flag
	proto.HelloClient // 嵌入HelloClient

	host string `gone:"config,server.host"`
	port string `gone:"config,server.port"`
}

// 实现 gone_grpc.Client接口的Address方法，该方法在客户端启动时会被自动调用
// 该方法的作用是告诉客户端gRPC服务的地址
func (c *helloClient) Address() string {
	return fmt.Sprintf("%s:%s", c.host, c.port)
}

// 实现 gone_grpc.Client接口的Stub方法，该方法在客户端启动时会被自动调用
// 在该方法中，完成 HelloClient的初始化
func (c *helloClient) Stub(conn *grpc.ClientConn) {
	c.HelloClient = proto.NewHelloClient(conn)
}

func main() {
	gone.
		Load(&helloClient{}).
		Loads(goner.BaseLoad, gone_grpc.ClientRegisterLoad).
		Run(func(in struct {
			hello *helloClient `gone:"*"` // 在Run方法的参数中，注入 helloClient
		}) {
			// 调用Say方法，给服务段发送消息
			say, err := in.hello.Say(context.Background(), &proto.SayRequest{Name: "gone"})
			if err != nil {
				log.Printf("er:%v", err)
				return
			}
			log.Printf("say result: %s", say.Message)
		})
}
```

## 实现方式二：Gone V2 Provider 机制

Gone V2 引入了强大的 Provider 机制，可以大幅简化 gRPC 组件的使用。

### 服务端实现

文件名：v2_server/main.go
```go
package main

import (
	"context"
	"github.com/gone-io/gone/v2"
	goneGrpc "github.com/gone-io/goner/grpc"
	"google.golang.org/grpc"
	"grpc_demo/proto"
	"log"
	"os"
)

type server struct {
	gone.Flag
	proto.UnimplementedHelloServer              // 嵌入UnimplementedHelloServer
	grpcServer                     *grpc.Server `gone:"*"` // 注入grpc.Server
}

func (s *server) Init() {
	proto.RegisterHelloServer(s.grpcServer, s) //注册服务
}

// Say  重载协议中定义的服务
func (s *server) Say(ctx context.Context, in *proto.SayRequest) (*proto.SayResponse, error) {
	log.Printf("Received: %v", in.GetName())
	return &proto.SayResponse{Message: "Hello " + in.GetName()}, nil
}

func main() {
	// gone内置默认的配置组件只能从环境变量中读取配置，所以需要设置环境变量
	os.Setenv("GONE_SERVER_GRPC_PORT", "9091")

	gone.
		Load(&server{}).
		Loads(goneGrpc.ServerLoad).
		// 启动服务
		Serve()
}
```

### 客户端实现

文件名：v2_client/main.go
```go
package main

import (
	"context"
	"github.com/gone-io/gone/v2"
	gone_grpc "github.com/gone-io/goner/grpc"
	"google.golang.org/grpc"
	"grpc_demo/proto"
	"log"
	"os"
)

type helloClient struct {
	gone.Flag
	proto.HelloClient // 使用方法1：嵌入HelloClient，本组件只负载初始化，能力提供给第三方组件使用

	// 使用方法2：在本组件直接使用，不提供第三方组件使用
	//hello *proto.HelloClient

	// config=${配置的key},address=${服务地址}； //config优先级更高
	clientConn *grpc.ClientConn `gone:"*,config=grpc.service.hello.address"`

	// config 和 address 可以一起使用，如果config没有读取到值，降级为使用 address
	//clientConn1 *grpc.ClientConn `gone:"*,config=grpc.service.hello.address,address=127.0.0.1:9091"`

	// address 也可以单独使用，不推荐这种方式，意味着写死了
	//clientConn2 *grpc.ClientConn `gone:"*,address=127.0.0.1:9091"`
}

func (c *helloClient) Init() {
	c.HelloClient = proto.NewHelloClient(c.clientConn)
	//c.hello = &c.HelloClient
}

func main() {
	// gone内置默认的配置组件只能从环境变量中读取配置，所以需要设置环境变量
	os.Setenv("GONE_GRPC_SERVICE_HELLO_ADDRESS", "127.0.0.1:9091")

	gone.
		Load(&helloClient{}).
		Loads(gone_grpc.ClientRegisterLoad).
		Run(func(in struct {
			hello *helloClient `gone:"*"` // 在Run方法的参数中，注入 helloClient
		}) {
			// 调用Say方法，给服务段发送消息
			say, err := in.hello.Say(context.Background(), &proto.SayRequest{Name: "gone"})
			if err != nil {
				log.Printf("er:%v", err)
				return
			}
			log.Printf("say result: %s", say.Message)
		})
}
```

## 配置文件

文件名：config/default.properties
```properties
# 设置grpc服务的端口和host
server.port=9001
server.host=127.0.0.1

# 设置客户端使用的grpc服务端口和host
server.grpc.port=${server.port}
server.grpc.host=${server.host}
```

## 测试

### 运行服务端

```bash
go run v2_server/main.go  # 或 v1_server/main.go
```

程序等待请求，屏幕打印内容：
```log
2024-06-19 22:02:41.971|INFO|/Users/jim/works/gone-io/gone/goner/grpc/server.go:84||Register gRPC service *main.server
2024-06-19 22:02:41.971|INFO|/Users/jim/works/gone-io/gone/goner/grpc/server.go:88||gRPC server now listen at 127.0.0.1:9091
```

### 运行客户端

```bash
go run v2_client/main.go  # 或 v1_client/main.go
```

程序执行完退出，屏幕打印内容如下：
```log
2024-06-19 22:06:20.713|INFO|/Users/jim/works/gone-io/gone/goner/grpc/client.go:59||register gRPC client *main.helloClient on address 127.0.0.1:9091

2024/06/19 22:06:20 say result: Hello gone
```

回到服务端窗口，可以看到服务器接收到请求，新打印一行日志：
```log
2024/06/19 22:06:08 Received: gone
```

## 两种实现方式对比

### 传统方式

**服务端**：
1. 需要实现 `RegisterGrpcServer` 接口方法来注册服务
2. 手动管理 gRPC 服务的注册过程

**客户端**：
1. 需要实现 `Address` 和 `Stub` 方法来初始化连接
2. 配置获取方式不够灵活，地址构建逻辑需要手写

### Provider 机制

**服务端**：
1. 通过标签自动注入 `*grpc.Server`
2. 在 `Init` 方法中完成服务注册，符合 Gone 的组件生命周期管理

**客户端**：
1. 不再需要实现 `Address` 和 `Stub` 方法
2. 支持灵活的配置方式，包括：
   - 仅从配置中读取地址
   - 配置与默认地址配合使用，实现降级策略
   - 直接硬编码地址（不推荐，但支持）

## 总结

Gone V2 的 Provider 机制大幅提升了 gRPC 组件的使用体验：

1. **代码更简洁**：移除了不必要的接口实现和重复性的模板代码
2. **更符合依赖注入思想**：通过标签自动注入所需组件
3. **配置更灵活**：支持多种地址获取策略，提高了代码的可维