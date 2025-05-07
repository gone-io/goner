<p>
    <a href="README.md">English</a>&nbsp ｜&nbsp 中文
</p>

# goner/consul 组件

## 组件概述

**goner/consul**组件为Gone框架提供了基于HashiCorp Consul的服务注册与发现功能。这一集成方案为分布式系统中的服务管理提供了强大的解决方案。

通过Consul组件，您可以：

- 在分布式环境中注册和发现服务
- 实现服务健康检查和自动故障转移
- 构建高可用的微服务架构

## 功能特性

### 服务注册与发现

- **服务注册**：将服务实例注册到Consul服务目录
- **服务发现**：从Consul获取可用的服务实例列表
- **服务监控**：监听服务实例变化，实时更新服务列表
- **健康检查**：基于TTL的健康检查机制，自动维护服务健康状态

## 配置参考

### 客户端配置

#### 基本配置

以下参数控制Consul客户端基本行为：

| 配置参数          | 说明                                   | 类型          | 默认值         | 示例             |
| ----------------- | -------------------------------------- | ------------- | -------------- | ---------------- |
| consul.address    | Consul服务器地址                       | string        | 127.0.0.1:8500 | "127.0.0.1:8500" |
| consul.scheme     | 连接协议                               | string        | http           | "http"           |
| consul.pathPrefix | URI前缀，用于Consul位于API网关后的场景 | string        | -              | "/consul"        |
| consul.datacenter | 数据中心名称                           | string        | -              | "dc1"            |
| consul.token      | ACL令牌                                | string        | -              | "your-token"     |
| consul.tokenFile  | 包含ACL令牌的文件路径                  | string        | -              | "/path/to/token" |
| consul.waitTime   | Watch操作的最大阻塞时间                | time.Duration | -              | "30s"            |
| consul.namespace  | 命名空间（仅Consul企业版）             | string        | -              | "my-namespace"   |
| consul.partition  | 分区（仅Consul企业版）                 | string        | -              | "my-partition"   |

#### TLS配置

以下参数用于配置Consul客户端的TLS连接：

| 配置参数                      | 说明                 | 类型   | 默认值 | 示例                 |
| ----------------------------- | -------------------- | ------ | ------ | -------------------- |
| consul.tls.address            | TLS服务器名称（SNI） | string | -      | "consul.example.com" |
| consul.tls.caFile             | CA证书文件路径       | string | -      | "/path/to/ca.pem"    |
| consul.tls.caPath             | CA证书目录路径       | string | -      | "/path/to/certs"     |
| consul.tls.certFile           | 客户端证书文件路径   | string | -      | "/path/to/cert.pem"  |
| consul.tls.keyFile            | 客户端私钥文件路径   | string | -      | "/path/to/key.pem"   |
| consul.tls.insecureSkipVerify | 是否跳过TLS主机验证  | bool   | false  | true                 |


> 给多配置，参考[consul官方文档](https://pkg.go.dev/github.com/hashicorp/consul/api#Config)

#### 支持加载配置，注入到服务

组件名为：**consul.config**

```go
    gone.
	    NewApp(
			//... 
		).
	    Loads(g.NamedThirdComponentLoadFunc("consul.config", &api.Config{
            Address: "127.0.0.1:8500",
			// 其他配置
        }))
```

### 服务注册配置

服务注册相关的配置参数：

- gin.http server

| 配置参数                   | 说明       | 类型   | 默认值    | 示例            |
| -------------------------- | ---------- | ------ | --------- | --------------- |
| service.name               | 服务名称   | string | -         | "user-service"  |
| service.host               | 服务地址   | string | -         | "192.168.1.100" |
| service.port               | 服务端口   | int    | -         | 8080            |
| service.service-use-subnet | 使用的子网 | string | 0.0.0.0/0 | 192.168.1.0/24  |

- grpc server
| 配置参数                        | 说明       | 类型   | 默认值    | 示例            |
| ------------------------------- | ---------- | ------ | --------- | --------------- |
| service.grpc.name               | 服务名称   | string | -         | "user-service"  |
| service.grpc.host               | 服务地址   | string | -         | "192.168.1.100" |
| service.grpc.port               | 服务端口   | int    | -         | 8080            |
| service.grpc.service-use-subnet | 使用的子网 | string | 0.0.0.0/0 | 192.168.1.0/24  |


## 实施指南

### 配置文件设置

在项目的config目录中创建`default.yaml`文件，定义Consul客户端连接参数：

```yaml
consul:
  # 基本配置
  address: "127.0.0.1:8500"  # Consul服务器地址
  scheme: "http"            # 连接协议
  pathPrefix: ""            # URI前缀（API网关场景）
  datacenter: "dc1"         # 数据中心
  token: ""                 # ACL令牌
  tokenFile: ""             # ACL令牌文件路径
  waitTime: "30s"           # Watch操作的最大阻塞时间
  namespace: ""             # 命名空间（仅企业版）
  partition: ""             # 分区（仅企业版）

  # TLS配置
  tls:
    address: "consul.example.com"  # TLS服务器名称
    caFile: "/path/to/ca.pem"      # CA证书文件路径
    caPath: "/path/to/certs"       # CA证书目录路径
    certFile: "/path/to/cert.pem"  # 客户端证书文件路径
    keyFile: "/path/to/key.pem"    # 客户端私钥文件路径
    insecureSkipVerify: false      # 是否跳过TLS主机验证
```

### 代码集成

#### 服务注册与发现

```go
func main() {
    // 初始化应用并加载Consul组件
    gone.NewApp(consul.RegistryLoad).Run(func(params struct {
        registry g.ServiceRegistry `gone:"*"`
        discovery g.ServiceDiscovery `gone:"*"`
    }) {
        // 创建服务实例
        service := g.NewService(
            "my-service",       // 服务名称
            "192.168.1.100",   // 服务IP
            8080,               // 服务端口
            map[string]string{  // 元数据
                "version": "1.0.0",
            },
            true,               // 是否健康
            1.0,                // 权重
        )

        // 注册后，组件会自动维护TTL健康检查，默认TTL为20秒，健康检查间隔为10秒

        // 注册服务
        err := params.registry.Register(service)
        if err != nil {
            log.Fatalf("服务注册失败: %v", err)
        }

        // 发现服务
        instances, err := params.discovery.GetInstances("target-service")
        if err != nil {
            log.Fatalf("服务发现失败: %v", err)
        }

        // 监听服务变化
        ch, stop, err := params.discovery.Watch("target-service")
        if err != nil {
            log.Fatalf("监听服务失败: %v", err)
        }

        go func() {
            for services := range ch {
                // 处理服务列表更新
                fmt.Printf("服务列表更新: %v\n", services)
            }
        }()

        // 应用退出时注销服务并停止监听
        defer func() {
            _ = params.registry.Deregister(service)
            _ = stop()
        }()
        // 应用主逻辑...
    })
}
```



## 相关链接

- [Consul官方文档](https://www.consul.io/docs)
- [Gone框架文档](https://github.com/gone-io/gone)
- [HashiCorp Consul API](https://github.com/hashicorp/consul/api)