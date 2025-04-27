[//]: # (desc: gRPC使用服务发现的例子)

<p>
    <a href="README.md">English</a>&nbsp ｜&nbsp 中文
</p>

# Gone框架 gRPC服务发现示例

## 项目概述

本示例展示了如何在Gone框架中使用服务发现功能进行gRPC通信。示例包含一个服务端和一个客户端，服务端注册到Nacos服务发现中心，客户端通过服务名称而非具体IP地址来访问服务端。

## 功能特点

- 演示基于Nacos的gRPC服务注册与发现
- 展示gRPC客户端如何通过服务名称访问服务
- 支持自动负载均衡（当有多个服务实例时）
- 完全集成Gone框架的依赖注入特性

## 项目结构

```
.
├── client/             # 客户端代码
│   └── main.go         # 客户端入口文件
├── config/             # 配置文件目录
│   └── default.yaml    # 默认配置文件
├── docker-compose.yaml # Docker环境配置文件
├── go.mod              # Go模块定义
├── logs/               # 日志目录
├── proto/              # 协议定义目录
│   ├── hello.pb.go     # 生成的协议代码
│   ├── hello.proto     # 协议定义文件
│   └── hello_grpc.pb.go# 生成的gRPC代码
└── server/             # 服务端代码
    └── main.go         # 服务端入口文件
```

## 工作原理

### 服务发现流程

1. 服务端启动时，通过Nacos注册中心注册自己的服务信息（服务名、IP地址、端口等）
2. 客户端启动时，通过服务名称向Nacos查询服务地址
3. Nacos返回服务的地址信息给客户端
4. 客户端使用获取到的地址信息建立gRPC连接
5. 当服务端实例发生变化时，客户端能够自动感知并更新连接

### 关键组件

- **服务端**：实现gRPC服务并注册到Nacos
- **客户端**：通过服务名称发现服务并建立连接
- **Nacos**：提供服务注册与发现功能
- **gRPC**：提供高性能的RPC通信

## 配置说明

### 服务端配置

服务端在`config/default.yaml`中配置：

```yaml
server:
  grpc:
    port: 0  # 使用0表示随机端口
    service-name: user-center  # 服务名称
```

### 客户端配置

客户端通过依赖注入配置服务连接：

```go
// 使用方法1：通过配置文件中的服务名称连接
clientConn *grpc.ClientConn `gone:"*,config=grpc.service.hello.address"`

// 使用方法2：配置和地址一起使用，配置优先
//clientConn1 *grpc.ClientConn `gone:"*,config=grpc.service.hello.address,address=127.0.0.1:9091"`

// 使用方法3：直接指定地址（不推荐，硬编码）
//clientConn2 *grpc.ClientConn `gone:"*,address=127.0.0.1:9091"`
```

### Nacos配置

```yaml
nacos:
  client:
    namespaceId: public
    asyncUpdateService: false
    logLevel: debug
    logDir: ./logs/
  server:
    ipAddr: "127.0.0.1"
    contextPath: /nacos
    port: 8848
    scheme: http

  service:
    group: DEFAULT_GROUP
    clusterName: default
```

## 运行示例

### 前置条件

- 安装Docker和Docker Compose
- 安装Go 1.16+

### 启动Nacos

```bash
docker-compose up -d nacos
```

### 启动服务端

```bash
cd server
go run main.go
```

### 启动客户端

```bash
cd client
go run main.go
```

### 预期输出

客户端输出示例：
```
2023/xx/xx xx:xx:xx say result: Hello gone
2023/xx/xx xx:xx:xx say result: Hello gone
...
```

服务端输出示例：
```
2023/xx/xx xx:xx:xx Received: gone
2023/xx/xx xx:xx:xx Received: gone
...
```

## 代码解析

### 服务端实现

```go
type server struct {
	gone.Flag
	proto.UnimplementedHelloServer              // 嵌入UnimplementedHelloServer
	grpcServer                     *grpc.Server `gone:"*"` // 注入grpc.Server
}

func (s *server) Init() {
	proto.RegisterHelloServer(s.grpcServer, s) //注册服务
}

// Say 实现协议中定义的服务
func (s *server) Say(ctx context.Context, in *proto.SayRequest) (*proto.SayResponse, error) {
	log.Printf("Received: %v", in.GetName())
	return &proto.SayResponse{Message: "Hello " + in.GetName()}, nil
}
```

### 客户端实现

```go
type helloClient struct {
	gone.Flag
	proto.HelloClient // 使用方法1：嵌入HelloClient

	// 通过服务名称注入连接
	clientConn *grpc.ClientConn `gone:"*,config=grpc.service.hello.address"`
}

func (c *helloClient) Init() {
	c.HelloClient = proto.NewHelloClient(c.clientConn)
}
```

## 扩展阅读

- [Gone框架文档](https://github.com/gone-io/gone)
- [Nacos服务发现文档](https://nacos.io/zh-cn/docs/v2/guide/user/service-discovery.html)
- [gRPC官方文档](https://grpc.io/docs/)