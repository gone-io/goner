[//]: # (desc: 使用OpenTelemetry通过gRPC协议进行链路追踪)

<p>
    <a href="README.md">English</a>&nbsp ｜&nbsp 中文
</p>

# 使用OpenTelemetry通过gRPC协议进行链路追踪

本示例展示如何在Gone框架中集成OpenTelemetry与gRPC，实现分布式服务之间的链路追踪。

## 项目构建步骤

### 1. 创建项目和安装依赖包

```bash
# 创建项目目录
mkdir tracer-over-oltp-grpc
cd tracer-over-oltp-grpc

# 初始化Go模块
go mod init examples/otel/tracer-over-oltp-grpc

# 安装Gone框架的OpenTelemetry与gRPC集成组件
gonectr install goner/otel/tracer/grpc
```

### 2. 定义gRPC服务

首先，创建proto文件：

```bash
mkdir proto
touch proto/hello.proto
```

在`proto/hello.proto`中定义服务接口：

```protobuf
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

### 3. 实现服务端

在`server/services/hello.go`中实现服务：

```go
package services

import (
	"context"
	"examples/otel/tracer/oltp/grpc/proto"
	"fmt"
	"github.com/gone-io/gone/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

type server struct {
	gone.Flag
	proto.UnimplementedHelloServer              // 嵌入UnimplementedHelloServer
	grpcServer                     *grpc.Server `gone:"*"` // 注入grpc.Server
	logger                         gone.Logger  `gone:"*"`
	tracer                         trace.Tracer
}

const tracerName = "grpc-hello-server"

func (s *server) Init() {
	proto.RegisterHelloServer(s.grpcServer, s) //注册服务
	s.tracer = otel.Tracer(tracerName)
}

// Say  重载协议中定义的服务
func (s *server) Say(ctx context.Context, in *proto.SayRequest) (*proto.SayResponse, error) {
	x, span := s.tracer.Start(ctx, "Say")
	defer span.End()

	s.logger.Infof("Received: %v", in.GetName())
	span.SetAttributes(attribute.Key("Name").String(in.GetName()))

	say := s.doSay(x, in.GetName())

	return &proto.SayResponse{Message: say}, nil
}

func (s *server) doSay(ctx context.Context, name string) string {
	_, span := s.tracer.Start(ctx, "doSay")
	defer span.End()

	span.SetAttributes(attribute.Key("name").String(name))
	span.AddEvent("doSay", trace.WithAttributes(attribute.Key("name").String(name)))

	return fmt.Sprintf("Hello, %s!", name)
}
```

## 运行服务

### 1. 启动服务端

```bash
# 进入服务端目录
cd server

# 运行服务
go run ./cmd
```

### 2. 启动客户端

```bash
# 进入客户端目录
cd client

# 运行客户端
go run ./cmd
```

## 查看结果

服务运行后，可以通过Jaeger UI查看链路追踪数据：

1. 访问Jaeger UI界面：http://localhost:16686
2. 在Search界面选择服务名称：`grpc-hello-server`
3. 点击Find Traces按钮查看追踪数据
！[](./screenshot.png)

你可以看到完整的调用链路，包括：
- 客户端发起请求
- 服务端接收请求
- Say方法的执行
- doSay方法的执行
- 响应返回客户端

每个span中都包含了详细的属性信息，如请求参数、执行时间等。
