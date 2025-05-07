[//]: # (desc: Gone grpc 示例项目，展示如何使用 Gone 框架构建 gRPC 服务)

<p>
    <a href="README.md">English</a>&nbsp ｜&nbsp 中文
</p>

# Gone grpc 示例项目

## 项目简介
本示例项目展示如何基于 Gone 框架快速构建 gRPC 服务端与客户端，包含服务注册、配置管理、协议定义等典型场景，适合初学者和进阶用户参考。

## 目录结构
```
.
├── config/             # 配置文件目录
│   └── default.properties # 默认配置
├── go.mod              # Go模块定义
├── proto/              # 协议定义目录
│   ├── hello.pb.go     # 生成的协议代码
│   ├── hello.proto     # 协议定义文件
│   └── hello_grpc.pb.go# 生成的gRPC代码
├── v1_client/          # v1版本客户端
│   └── main.go         # 客户端入口
├── v1_server/          # v1版本服务端
│   └── main.go         # 服务端入口
├── v2_client/          # v2版本客户端（可选）
│   └── main.go         # 客户端入口
├── v2_server/          # v2版本服务端（可选）
│   └── main.go         # 服务端入口
└── README_CN.md        # 中文说明文档
```

## 主要依赖
- Go 1.24+
- github.com/gone-io/gone/v2 v2.1.0
- github.com/gone-io/goner/grpc v1.2.1
- github.com/gone-io/goner/viper v1.2.1
- google.golang.org/grpc v1.72.0
- google.golang.org/protobuf v1.36.6

依赖详情见 go.mod 文件。

## 如何运行
### 1. 启动服务端
进入 v1_server 目录，运行：
```shell
cd v1_server
go run main.go
```
服务端默认监听 9091 端口。

### 2. 启动客户端
另开终端，进入 v1_client 目录，运行：
```shell
cd v1_client
go run main.go
```
客户端将向服务端发送 gRPC 请求并输出响应。

### 3. 配置说明
可通过 config/default.properties 或环境变量自定义端口、主机等参数。

## 关键代码说明
### 服务端 main.go
```go
func main() {
    os.Setenv("GONE_SERVER_GRPC_PORT", "9091")
    gone.
        Load(&server{}).
        Loads(goneGrpc.ServerLoad).
        Serve()
}
```
- 通过 gone.Load 加载服务结构体，goneGrpc.ServerLoad 启动 gRPC 服务。
- 通过 RegisterGrpcServer 注册 gRPC 服务实现。

### 客户端 main.go
```go
func main() {
    gone.
        Load(&helloClient{}).
        Loads(viper.Load, gone_grpc.ClientRegisterLoad).
        Run(func(in struct {
            hello *helloClient `gone:"*"`
        }) {
            say, err := in.hello.Say(context.Background(), &proto.SayRequest{Name: "gone"})
            if err != nil {
                log.Printf("er:%v", err)
                return
            }
            log.Printf("say result: %s", say.Message)
        })
}
```
- gone.Load 加载客户端结构体，viper.Load 加载配置，gone_grpc.ClientRegisterLoad 注册 gRPC 客户端。
- 通过 Say 方法向服务端发起请求。

### 协议定义
proto/hello.proto 定义了 Say 服务及消息结构，使用 protoc 生成对应 Go 代码。

## 常见问题
1. **端口冲突**：请确保 9091 端口未被占用，或通过配置修改端口。
2. **依赖未安装**：请先执行 `go mod tidy` 安装依赖。
3. **proto 文件未生成**：请确保已用 protoc 生成 hello.pb.go 和 hello_grpc.pb.go。

## 参考文档
- [Gone 框架文档](https://github.com/gone-io/gone)
- [gRPC 官方文档](https://grpc.io/docs/)

如有疑问欢迎提 issue 或参与讨论。
