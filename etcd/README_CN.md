<p>
    <a href="README.md">English</a>&nbsp ｜&nbsp 中文
</p>

# goner/etcd 组件

## 组件概述

**goner/etcd** 组件为Gone框架提供了基于etcd的服务注册与发现功能。这一集成方案为分布式系统中的服务管理提供了可靠的解决方案。

goner/etcd 组件，您可以：

- 在分布式环境中注册和发现服务
- 实现服务健康检查和自动故障转移
- 构建高可用的微服务架构
- 利用etcd的强一致性特性进行服务协调

## 功能特性

### 服务注册与发现

- **服务注册**：将服务实例注册到etcd服务目录
- **服务发现**：从etcd获取可用的服务实例列表
- **服务监控**：监听服务实例变化，实时更新服务列表
- **健康检查**：基于TTL的健康检查机制，自动维护服务健康状态
- **强一致性**：利用etcd的Raft共识算法保证数据一致性

## 配置参考

### 客户端配置

以下参数控制etcd客户端基本行为：

| 配置参数 | 说明 | 类型 | 默认值 | 示例 |
|----------|------|------|---------|------|
| etcd.endpoints | etcd服务器地址列表 | []string | ["127.0.0.1:2379"] | ["localhost:2379"] |
| etcd.username | etcd认证用户名 | string | "" | "username" |
| etcd.password | etcd认证密码 | string | "" | "password" |
| etcd.dial-timeout | 连接超时时间 | time.Duration | "5s" | "10s" |
| etcd.keepalive-ttl | 健康检查TTL | time.Duration | "10s" | "20s" |

> 给多配置，参考[etcd官方文档](https://pkg.go.dev/go.etcd.io/etcd/client/v3#Config)

#### 支持加载配置，注入到服务
组件名为：**etcd.config**

```go
    gone.
	    NewApp(
			//... 
		).
	    Loads(g.NamedThirdComponentLoadFunc("etcd.config", &etcd.Config{
			Endpoints: []string{"localhost:2379"},
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

在项目的config目录中创建`default.yaml`文件，定义etcd客户端连接参数：

```yaml
etcd:
  # 基本配置
  endpoints: ["localhost:2379"]  # etcd服务器地址列表
  dialTimeout: "5s"              # 连接超时时间
  username: ""                   # 认证用户名
  password: ""                   # 认证密码
  autoSyncInterval: "30s"        # 端点自动同步间隔

```

### 代码集成

#### 服务注册与发现

```go
func main() {
    // 初始化应用并加载etcd组件
    gone.NewApp(etcd.RegistryLoad).Run(func(params struct {
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

        // 注册后，组件会自动维护TTL健康检查，默认TTL为10秒

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

- [etcd官方文档](https://etcd.io/docs/)
- [Gone框架文档](https://github.com/gone-io/gone)
- [etcd API](https://github.com/etcd-io/etcd/tree/main/client/v3)