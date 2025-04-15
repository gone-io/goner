# Gone框架 HTTP服务发现示例

## 项目概述

本示例展示了如何在Gone框架中使用服务发现功能进行HTTP通信。示例包含一个服务端和一个客户端，服务端注册到Nacos服务发现中心，客户端通过服务名称而非具体IP地址来访问服务端。

## 功能特点

- 演示基于Nacos的服务注册与发现
- 展示HTTP客户端如何通过服务名称访问服务
- 支持自动负载均衡（当有多个服务实例时）
- 完全集成Gone框架的依赖注入特性

## 环境准备

### 启动Nacos服务

本示例使用Docker Compose来启动Nacos服务：

```bash
# 在项目根目录下执行
docker-compose up -d nacos
```

这将启动一个Nacos服务实例，监听在8848端口。

## 代码结构

```
./
├── client/         # HTTP客户端示例
│   └── main.go     # 客户端主程序
├── config/         # 配置文件目录
│   └── default.yaml # 默认配置
├── docker-compose.yaml # Docker环境配置
├── go.mod          # Go模块定义
├── logs/           # 日志目录
└── server/         # HTTP服务端示例
    └── main.go     # 服务端主程序
```

## 代码实现

### 服务端实现

服务端（`server/main.go`）使用Gone框架的Gin组件创建HTTP服务，并通过Nacos组件注册到服务发现中心：

```go
func main() {
	gone.
		NewApp(goner.GinLoad, nacos.RegistryLoad, viper.Load).
		Load(&HelloController{}).
		Serve()
}
```

服务端定义了一个简单的HTTP接口：

```go
func (c *HelloController) Mount() gin.MountError {
	c.GET("/hello", func(in struct {
		name string `gone:"http,query"`
	}) string {
		return fmt.Sprintf("hello, %s", in.name)
	})
	return nil
}
```

### 客户端实现

客户端（`client/main.go`）使用Gone框架的urllib组件和balancer组件，通过服务名称访问服务端：

```go
func main() {
	gone.
		NewApp(
			nacos.RegistryLoad,
			balancer.Load,
			viper.Load,
			urllib.Load,
		).
		Run(func(client urllib.Client, logger gone.Logger) {
			// 通过服务名称访问服务
			res, err := client.
				R().
				SetSuccessResult(&data).
				Get("http://user-center/hello?name=goner")
			// ...
		})
}
```

## 配置说明

配置文件（`config/default.yaml`）包含以下关键配置：

```yaml
nacos:
  client:
    # Nacos客户端配置
    namespaceId: public
    # ...
  server:
    # Nacos服务端地址配置
    ipAddr: "127.0.0.1"
    port: 8848
    # ...
  service:
    # 服务发现相关配置
    group: DEFAULT_GROUP
    clusterName: default

# 服务端配置
server:
    port: 0  # 使用随机端口
    service-name: user-center  # 服务名称
```

## 运行示例

### 1. 启动Nacos服务

```bash
docker-compose up -d nacos
```

### 2. 启动服务端

```bash
cd server
go run main.go
```

### 3. 启动客户端

```bash
cd client
go run main.go
```

客户端将发送10次请求到服务端，并打印响应结果。

## 关键点说明

1. **服务注册**：服务端启动时会自动注册到Nacos服务发现中心
2. **服务发现**：客户端通过服务名称（`user-center`）而非IP地址访问服务
3. **负载均衡**：当有多个服务实例时，客户端会自动进行负载均衡
4. **零配置端口**：服务端使用`port: 0`配置随机端口，避免端口冲突

## 扩展阅读

- [Gone框架文档](https://github.com/gone-io/gone)
- [Nacos服务发现文档](https://nacos.io/zh-cn/docs/v2/guide/user/service-discovery.html)