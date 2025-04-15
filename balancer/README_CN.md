# Balancer 负载均衡器

[![Go Report Card](https://goreportcard.com/badge/github.com/gone-io/goner/balancer)](https://goreportcard.com/report/github.com/gone-io/goner/balancer)
[![GoDoc](https://godoc.org/github.com/gone-io/goner/balancer?status.svg)](https://godoc.org/github.com/gone-io/goner/balancer)

`balancer`是Gone框架的客户端负载均衡器组件，提供服务发现和负载均衡功能，支持多种负载均衡策略。该组件与`goner/urllib`无缝集成，为`urllib`提供负载均衡能力。**注意**：服务端和客户端需要使用相同的服务注册/发现组件才能正常工作。

## 功能特性

- 与Gone框架无缝集成
- 支持多种负载均衡策略（轮询、随机、权重）
- 自动服务发现和实例监控
- 实例缓存和自动更新机制

## 安装

```bash
go get github.com/gone-io/goner/balancer
```

## 快速开始

### 1. 引入模块

在应用中引入balancer模块：

```go
import (
    "github.com/gone-io/gone/v2"
    "github.com/gone-io/goner/balancer"
    // 其他导入
)

func main() {
    // 创建Gone应用
    app := gone.NewApp(
        balancer.Load,  // 加载balancer模块
        // 加载其他模块...
    )
    
    // 运行应用
    if err := app.Run(); err != nil {
        panic(err)
    }
}
```

### 2. 使用负载均衡器

在服务中注入并使用负载均衡器：

```go
type MyService struct {
    gone.Flag
    balancer g.LoadBalancer `gone:"*"`
}

func (s *MyService) CallRemoteService(ctx context.Context) (interface{}, error) {
    // 获取服务实例
    instance, err := s.balancer.GetInstance(ctx, "remote-service-name")
    if err != nil {
        return nil, err
    }
    
    // 使用获取到的服务实例
    serviceAddr := instance.GetHost() + ":" + instance.GetPort()
    // 调用远程服务...
    
    return result, nil
}
```

### 3. 与urllib集成

`balancer`可以与`urllib`无缝集成，为HTTP请求提供负载均衡能力：

```go
func main() {
    gone.
        NewApp(
            nacos.RegistryLoad,  // 加载服务发现组件
            balancer.Load,       // 加载balancer模块
            viper.Load,          // 加载配置模块
            urllib.Load,         // 加载urllib模块
        ).
        Run(func(client urllib.Client, logger gone.Logger) {
            // 直接使用服务名作为主机名，balancer会自动处理负载均衡
            res, err := client.
                R().
                SetSuccessResult(&data).
                Get("http://user-center/hello?name=goner")
                
            if err != nil {
                logger.Errorf("client request err: %v", err)
                return
            }

            if res.IsSuccessState() {
                logger.Infof("res=> %#v", data)
            }
        })
}
```

## 负载均衡策略

balancer模块提供了以下几种负载均衡策略：

### 1. 轮询策略 (RoundRobinStrategy)

按顺序轮流选择服务实例，是默认的负载均衡策略。

```go
// 默认已加载，无需额外配置
// 在load.go中已经默认加载了RoundRobinStrategy
```

### 2. 随机策略 (RandomStrategy)

随机选择一个服务实例。

```go
// 方法1：使用提供的加载函数
import (
    "github.com/gone-io/gone/v2"
    "github.com/gone-io/goner/balancer"
)

func main() {
    app := gone.NewApp(
        balancer.Load,               // 先加载基础balancer模块
        balancer.LoadRandomStrategy, // 再加载随机策略，会替换默认策略
        // 加载其他模块...
    )
    
    // 运行应用
    if err := app.Run(); err != nil {
        panic(err)
    }
}
```

### 3. 权重策略 (WeightStrategy)

根据服务实例的权重进行选择，权重越高被选中的概率越大。

```go
// 方法1：使用提供的加载函数
import (
    "github.com/gone-io/gone/v2"
    "github.com/gone-io/goner/balancer"
)

func main() {
    app := gone.NewApp(
        balancer.Load,               // 先加载基础balancer模块
        balancer.LoadWeightStrategy, // 再加载权重策略，会替换默认策略
        // 加载其他模块...
    )
    
    // 运行应用
    if err := app.Run(); err != nil {
        panic(err)
    }
}
```

## 高级用法

### 自定义负载均衡策略

您可以通过实现`g.LoadBalanceStrategy`接口来创建自定义的负载均衡策略：

```go
type MyCustomStrategy struct {
    gone.Flag
    // 自定义字段
}

func (s *MyCustomStrategy) Select(ctx context.Context, instances []g.Service) (g.Service, error) {
    // 实现自定义的选择逻辑
    // ...
    return selectedInstance, nil
}

func main() {
    myStrategy := &MyCustomStrategy{}

    app := gone.NewApp(
        balancer.Load, // 先加载基础balancer模块
            // 加载其他模块...
        ).
		Load(balancer.LoadCustomerStrategy(myStrategy)) // 使用LoadCustomerStrategy加载自定义策略

    // 运行应用
    if err := app.Run(); err != nil {
        panic(err)
    }
}

```

### 与服务发现集成

balancer模块需要与服务发现组件配合使用，例如与nacos集成：

```go
import (
    "github.com/gone-io/gone/v2"
    "github.com/gone-io/goner/balancer"
    "github.com/gone-io/goner/nacos"
    "github.com/gone-io/goner/urllib"
)

func main() {
    app := gone.NewApp(
        nacos.RegistryLoad,  // 加载nacos服务发现模块
        balancer.Load,       // 加载balancer模块
        urllib.Load,         // 加载urllib模块
        // 加载其他模块...
    )

    // 运行应用
    if err := app.Run(); err != nil {
        panic(err)
    }
}
```

## 实现原理

balancer模块的核心功能包括：

1. **服务发现**：通过注入的`g.ServiceDiscovery`接口获取服务实例列表
2. **实例缓存**：缓存已获取的服务实例，提高性能
3. **实例监控**：监听服务实例变化，自动更新缓存
4. **负载均衡**：根据选定的策略从可用实例中选择一个实例


## 许可证

[Apache License 2.0](https://github.com/gone-io/gone/blob/main/LICENSE)