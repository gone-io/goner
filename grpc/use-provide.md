# 使用Gone V2的Provide机制改造goner/grpc组件

在[Gone V2 Provider 机制介绍](https://github.com/gone-io/gone/blob/main/docs/provider.md)中我们详细介绍了Gone框架的Provide机制及其强大的依赖注入能力。理论讲解固然重要，但真正理解一个机制的价值，还需要通过实践来检验。本文将通过一个实际案例，展示如何利用Provide机制对goner/grpc组件进行改造，让我们亲身体验Gone V2框架在简化代码、提升开发体验方面的巨大潜力。

## 目录
- [使用Gone V2的Provide机制改造goner/grpc组件](#使用gone-v2的provide机制改造gonergrpc组件)
    - [目录](#目录)
    - [现有goner/grpc组件的使用痛点](#现有gonergrpc组件的使用痛点)
        - [服务端示例代码](#服务端示例代码)
        - [客户端示例代码](#客户端示例代码)
        - [问题总结](#问题总结)
    - [服务端改造方案](#服务端改造方案)
        - [改造后的服务端业务代码](#改造后的服务端业务代码)
        - [goner/grpc/server.go的改造](#gonergrpcservergo的改造)
    - [客户端改造方案](#客户端改造方案)
        - [改造后的客户端业务代码](#改造后的客户端业务代码)
        - [goner/grpc/client.go的改造](#gonergrpcclientgo的改造)
    - [总结](#总结)

## 现有goner/grpc组件的使用痛点

在Gone中，gRPC组件的使用体验一直不够自然流畅。让我们先看看现有的实现方式，分析其中的痛点。

### 服务端示例代码

服务端实现服务时，需要通过实现`RegisterGrpcServer`接口方法来注册服务：

```go
// server/main.go
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
	// gone内置默认的配置组件只能从环境变量中读取配置，所以需要设置环境变量
	os.Setenv("GONE_SERVER_GRPC_PORT", "9091")

	gone.
		Load(&server{}).
		Loads(goneGrpc.ServerLoad).
		// 启动服务
		Serve()
}
```

### 客户端示例代码

客户端使用服务时，需要通过实现`Address`和`Stub`方法来初始化连接：

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

### 问题总结

通过分析上述代码，我们可以总结出以下几个主要痛点：

1. **手动注册机制违背DI原则**：
    - 服务端需要显式实现`RegisterGrpcServer`接口，这与Gone框架"自动装配"的核心理念相悖
    - 开发者需要手动管理gRPC服务的注册过程，而理想的依赖注入框架应该通过标签或约定自动完成这种绑定

2. **客户端实现存在大量样板代码**：
    - 每个gRPC客户端都需要实现相同模式的`Address()`和`Stub()`方法
    - 这些重复性的模板代码与Gone V2通过Provider机制消除重复代码的设计目标不符
    - 配置获取方式不够灵活，地址构建逻辑需要手写

这些问题导致开发者在使用gRPC组件时体验不佳，不符合Gone框架简洁易用的设计理念。

## 服务端改造方案

### 改造后的服务端业务代码

针对服务端的痛点，我们的改造目标是：

1. 使`*grpc.Server`能够自动注入，不再需要实现`RegisterGrpcServer`方法
2. 将服务注册流程放到`Init`方法中，使其更符合直觉和Gone的生命周期机制

改造后的服务端业务代码如下：

```go
type server struct {
	gone.Flag
	proto.UnimplementedHelloServer              // 嵌入UnimplementedHelloServer
	grpcServer                     *grpc.Server `gone:"*"` // 注入grpc.Server
}

func (s *server) Init() {
	proto.RegisterHelloServer(s.grpcServer, s) //在Init方法中完成服务注册
}

// Say 重载协议中定义的服务
func (s *server) Say(ctx context.Context, in *proto.SayRequest) (*proto.SayResponse, error) {
	log.Printf("Received: %v", in.GetName())
	return &proto.SayResponse{Message: "Hello " + in.GetName()}, nil
}
```

这种改造方式带来了以下几个明显的好处：

- 代码更加简洁，移除了不必要的接口实现
- 更符合依赖注入的思想，通过标签自动注入grpc.Server
- 服务注册逻辑放在Init方法中，符合Gone的组件生命周期管理

### goner/grpc/server.go的改造

为了支持上述服务端业务代码的改造，我们需要对`goner/grpc/server.go`进行相应的修改。改造前的完整代码可以在[v0.0.6/grpc/server.go](https://github.com/gone-io/goner/blob/v0.0.6/grpc/server.go)查看，改造后的完整代码在[grpc/server.go](https://github.com/gone-io/goner/blob/a05f194b28c9b923be61cdfa697bd51c5cd154d9/grpc/server.go)。

主要改造点包括：

1. 给server结构体增加`Provide`方法，使其成为一个**Provider**
2. 将`server.grpcServer`初始化的代码放到`Init`方法中

```go
type server struct {
	gone.Flag
	//...
	grpcServer     *grpc.Server
	listener       net.Listener
	//...
}

func (s *server) Init() error {
	err := s.initListener()
	if err != nil {
		return gone.ToError(err)
	}

	s.grpcServer = grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			s.traceInterceptor,
			s.recoveryInterceptor,
		),
	)
	return nil
}

func (s *server) Provide() (*grpc.Server, error) {
	return s.grpcServer, nil
}
```

通过这种改造，gRPC服务器组件现在能够作为一个Provider向其他组件提供`*grpc.Server`实例，极大地简化了服务注册流程。

## 客户端改造方案

### 改造后的客户端业务代码

针对客户端的痛点，我们的改造目标是：

1. 不再需要实现`Stub`和`Address`方法
2. 自动注入`*grpc.ClientConn`
3. 在`Init`方法中完成Client的初始化
4. 支持灵活的配置方式，自动从配置中读取服务地址

改造后的客户端业务代码如下：

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
	proto.HelloClient // 使用方法1：嵌入HelloClient，本组件只负责初始化，能力提供给第三方组件使用
	
	// 使用方法2：在本组件直接使用，不提供给第三方组件使用
	//hello *proto.HelloClient

	// config=${配置的key},address=${服务地址}；config优先级更高
	clientConn *grpc.ClientConn `gone:"*,config=grpc.service.hello.address"`
	
	// config和address可以一起使用，如果config没有读取到值，降级为使用address
	//clientConn1 *grpc.ClientConn `gone:"*,config=grpc.service.hello.address,address=127.0.0.1:9090""`

	// address也可以单独使用，不推荐这种方式，意味着地址硬编码
	//clientConn2 *grpc.ClientConn `gone:"*,address=127.0.0.1:9090"`
}

func (c *helloClient) Init() {
	c.HelloClient = proto.NewHelloClient(c.clientConn)
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
			// 调用Say方法，给服务端发送消息
			say, err := in.hello.Say(context.Background(), &proto.SayRequest{Name: "gone"})
			if err != nil {
				log.Printf("er:%v", err)
				return
			}
			log.Printf("say result: %s", say.Message)
		})
}
```

这种改造方式带来了以下几个明显的好处：

- 代码更加简洁，移除了不必要的接口实现
- 提供了多种灵活的配置方式，包括：
    - 仅从配置中读取地址
    - 配置与默认地址配合使用，实现降级策略
    - 直接硬编码地址（不推荐，但支持）
- 符合Gone框架的组件生命周期管理，在Init方法中完成初始化

### goner/grpc/client.go的改造

为了支持上述客户端业务代码的改造，我们需要对`goner/grpc/client.go`进行相应的修改。改造前的完整代码可以在[v0.0.6/grpc/client.go](https://github.com/gone-io/goner/blob/v0.0.6/grpc/client.go)查看，改造后的完整代码在[grpc/client.go](https://github.com/gone-io/goner/blob/a05f194b28c9b923be61cdfa697bd51c5cd154d9/grpc/client.go)。

主要改造点包括：

1. 给`clientRegister`结构体增加`Provide`方法，使其成为一个**Provider**，能够根据注入标签自动创建`*grpc.ClientConn`
2. 在`clientRegister`上注入`gone.Configure`，用于根据配置键获取服务地址
3. 实现连接缓存机制，相同地址的服务复用同一个`*grpc.ClientConn`，提高性能

主要代码如下：

```go
type clientRegister struct {
	gone.Flag
	// ...
	configure gone.Configure `gone:"configure"`
	connections map[string]*grpc.ClientConn
	// ...
}

func (s *clientRegister) getConn(address string) (conn *grpc.ClientConn, err error) {
	conn = s.connections[address]
	if conn == nil {
		if conn, err = grpc.NewClient(
			address,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithChainUnaryInterceptor(s.traceInterceptor),
		); err != nil {
			return nil, gone.ToError(err)
		}
		s.connections[address] = conn
	}
	return
}

// Provide 实现Provider接口，根据标签配置提供grpc.ClientConn
func (s *clientRegister) Provide(tagConf string) (*grpc.ClientConn, error) {
	m, _ := gone.TagStringParse(tagConf)
	address := m["address"]
	if configKey, ok := m["config"]; ok {
		err := s.configure.Get(configKey, &address, address)
		if err != nil {
			return nil, gone.ToError(err)
		}
	}
	if address == "" {
		return nil, gone.ToError("address is empty")
	}
	return s.getConn(address)
}
```

通过这种改造，gRPC客户端组件现在能够：
- 解析注入标签中的配置
- 灵活获取服务地址（支持配置优先或地址优先）
- 缓存连接以提高性能
- 自动提供grpc.ClientConn实例给需要的组件

## 总结

通过本次改造，我们利用Gone V2的Provide机制大幅提升了goner/grpc组件的使用体验：

1. **服务端改进**：
    - 移除了手动注册机制，改为自动注入
    - 使用标准的Init方法进行服务注册，符合直觉
    - 服务实现代码更加简洁明了

2. **客户端改进**：
    - 移除了重复性的模板代码（Address和Stub方法）
    - 提供灵活的配置方式，支持多种地址获取策略
    - 实现连接缓存，提高性能和资源利用率

3. **整体收益**：
    - 代码更符合依赖注入的理念
    - 减少了样板代码，提高开发效率
    - 更加符合Gone框架"约定优于配置"的设计理念

这次改造充分展示了Gone V2的Provider机制在简化组件使用、提升开发体验方面的强大潜力。开发者现在可以用更自然、更直观的方式使用gRPC功能，专注于业务逻辑而非框架细节。