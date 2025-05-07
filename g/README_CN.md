<p>
    <a href="README.md">English</a>&nbsp ｜&nbsp 中文
</p>

# goner/g 组件

goner/g 组件是 gone 框架的核心组件之一，提供了一系列基础接口和功能，包括日志追踪、服务发现、负载均衡等核心功能。

## 核心功能

### 1. 日志追踪 (CtxLogger)

`CtxLogger` 接口用于日志追踪，可以为同一调用链路分配统一的 traceId，方便日志追踪和问题排查。

```go
type CtxLogger interface {
    Ctx(ctx context.Context) gone.Logger
}
```

使用示例：

```go
type user struct {
    gone.Flag
    logger CtxLogger `gone:"*"` //注入 Logger
}

func (u *user) Use(ctx context.Context) (err error) {
    // 从 openTelemetry context 中获取 traceId 并注入到 logger 中
    logger := u.logger.Ctx(ctx)

    logger.Infof("hello")

    return
}
```

### 2. 服务发现 (ServiceDiscovery)

`ServiceDiscovery` 接口提供服务发现和监控功能，允许客户端查找可用的服务实例并监听实例变化。

```go
type ServiceDiscovery interface {
    // GetInstances 返回指定服务的所有实例
    GetInstances(serviceName string) ([]Service, error)

    // Watch 创建一个通道，用于接收服务实例变化的更新
    Watch(serviceName string) (ch <-chan []Service, stop func() error, err error)
}
```

### 3. 负载均衡 (LoadBalancer)

`LoadBalancer` 接口提供负载均衡功能，用于从可用的服务实例中选择合适的实例。

```go
type LoadBalancer interface {
    // GetInstance 根据负载均衡策略返回一个服务实例
    GetInstance(ctx context.Context, serviceName string) (Service, error)
}
```

负载均衡策略接口：

```go
type LoadBalanceStrategy interface {
    // Select 从提供的实例列表中选择一个服务实例
    Select(ctx context.Context, instances []Service) (Service, error)
}
```

### 4. 服务实例 (Service)

`Service` 接口表示服务注册表中的一个服务实例，提供了服务实例的基本信息，包括身份标识、位置、元数据和健康状态等。

```go
type Service interface {
    // GetName 返回实例的服务名称
    GetName() string

    // GetIP 返回实例的 IP 地址
    GetIP() string

    // GetPort 返回实例的端口号
    GetPort() int

    // GetMetadata 返回与服务实例关联的元数据
    GetMetadata() Metadata

    // GetWeight 返回实例的权重
    GetWeight() float64

    // IsHealthy 返回服务实例的健康状态
    IsHealthy() bool
}
```

创建服务实例：

```go
// NewService 创建一个新的服务实例
func NewService(name, ip string, port int, meta Metadata, healthy bool, weight float64) Service
```

## 使用建议

1. 在使用日志追踪时，建议通过依赖注入获取 `CtxLogger` 实例，并在处理请求时使用 `Ctx()` 方法注入上下文信息。

2. 在实现服务发现时，可以选择合适的服务发现组件（如 Consul、Etcd 等）来实现 `ServiceDiscovery` 接口。

3. 在使用负载均衡时，可以根据实际需求实现自定义的 `LoadBalanceStrategy`，如轮询、随机、权重等策略。

4. 服务实例的元数据可以用于存储额外的配置信息，如版本号、部署环境等。

## 相关组件

- [goner/balancer](../balancer/README_CN.md): 提供了负载均衡器的具体实现
- [goner/consul](../consul/README_CN.md): 基于 Consul 的服务发现实现
- [goner/etcd](../etcd/README_CN.md): 基于 Etcd 的服务发现实现
