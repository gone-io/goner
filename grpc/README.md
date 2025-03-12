# 提供gRPC服务


首先创建一个grpc目录，在这个目录中初始化一个golang mod：
```bash
mkdir grpc
cd grpc
go mod init grpc_demo
```

## 编写proto文件，生成golang代码
- 编写协议文件
定义一个简单的Hello服务，包含一个Say方法
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
- 生成golang代码

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
protoc --go_out=. --go_opt=paths=source_relative \
--go-grpc_out=. --go-grpc_opt=paths=source_relative \
proto/hello.proto
```
>  其中，protoc的安装参考[Protocol Buffer 编译器安装](https://blog.csdn.net/waitdeng/article/details/139248507)


## 编写服务端代码
文件名：server/main.go
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

## 注册客户端
文件名：client/main.go
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

## 编写配置文件
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
- 先运行服务端：
```bash
go run server/main.go
```
程序等待请求，屏幕打印内容：
```log
2024-06-19 22:02:41.971|INFO|/Users/jim/works/gone-io/gone/goner/grpc/server.go:84||Register gRPC service *main.server
2024-06-19 22:02:41.971|INFO|/Users/jim/works/gone-io/gone/goner/grpc/server.go:88||gRPC server now listen at 127.0.0.1:9001
```

- 然后，另外开窗口启动客户端：
```bash
go run client/main.go
```
程序执行完退出，屏幕打印内容如下：
```log
2024-06-19 22:06:20.713|INFO|/Users/jim/works/gone-io/gone/goner/grpc/client.go:59||register gRPC client *main.helloClient on address 127.0.0.1:9001

2024/06/19 22:06:20 say result: Hello gone
```

- 回到服务端窗口，可以看到服务器接收到请求，新打印一行日志：
```log
2024/06/19 22:06:08 Received: gone
```

## 总结
在Gone中使用gRPC，需要完成以下几步：
- 编写服务端
  1. 编写服务端Goner，匿名嵌入proto协议生成代码的 默认实现
  2. 重载proto文件中定义的接口方法，编写提供服务的具体业务逻辑
  3. 实现`gone_grpc.Service`接口的`RegisterGrpcServer`方法，在该方法中完成服务注册
  4. 将 服务端Goner 注册到 Gone框架
  5. 启动服务

- 编写客户端
  1. 编写客户端Goner，嵌入proto协议生成代码的客户端接口
  2. 实现`gone_grpc.Client`接口的`Address`和`Stub`方法，`Address`方法返回服务端地址，`Stub`初始化客服端接口
  3. 将 客户端Goner 注册到 Gone框架
  4. 启动客户端，调用客服端接口方法

本文的代码开源在 [goner/grpc](https://github.com/gone-io/goner/tree/main/grpc/example)。