<p>
    English&nbsp ｜&nbsp <a href="README_CN.md">中文</a>
</p>

# goner/balancer Component, for Load Balancer 

[![Go Report Card](https://goreportcard.com/badge/github.com/gone-io/goner/balancer)](https://goreportcard.com/report/github.com/gone-io/goner/balancer)
[![GoDoc](https://godoc.org/github.com/gone-io/goner/balancer?status.svg)](https://godoc.org/github.com/gone-io/goner/balancer)

**goner/balancer** is a client-side load balancer component for the Gone framework, providing service discovery and load balancing functionality with support for multiple load balancing strategies. This component seamlessly integrates with `goner/urllib` to provide load balancing capabilities. **Note**: The server and client need to use the same service registration/discovery component to work properly.

## Features

- Seamless integration with Gone framework
- Support for multiple load balancing strategies (round-robin, random, weighted)
- Automatic service discovery and instance monitoring
- Instance caching and automatic update mechanism

## Installation

```bash
go get github.com/gone-io/goner/balancer
```

## Quick Start

### 1. Import the Module

Import the balancer module in your application:

```go
import (
    "github.com/gone-io/gone/v2"
    "github.com/gone-io/goner/balancer"
    // other imports
)

func main() {
    // Create Gone application
    app := gone.NewApp(
        balancer.Load,  // Load balancer module
        // Load other modules...
    )
    
    // Run the application
    if err := app.Run(); err != nil {
        panic(err)
    }
}
```

### 2. Use the Load Balancer

Inject and use the load balancer in your service:

```go
type MyService struct {
    gone.Flag
    balancer g.LoadBalancer `gone:"*"`
}

func (s *MyService) CallRemoteService(ctx context.Context) (interface{}, error) {
    // Get service instance
    instance, err := s.balancer.GetInstance(ctx, "remote-service-name")
    if err != nil {
        return nil, err
    }
    
    // Use the obtained service instance
    serviceAddr := instance.GetHost() + ":" + instance.GetPort()
    // Call remote service...
    
    return result, nil
}
```

### 3. Integration with urllib

`balancer` can be seamlessly integrated with `urllib` to provide load balancing capabilities for HTTP requests:

```go
func main() {
    gone.
        NewApp(
            nacos.RegistryLoad,  // Load service discovery component
            balancer.Load,       // Load balancer module
            viper.Load,          // Load configuration module
            urllib.Load,         // Load urllib module
        ).
        Run(func(client urllib.Client, logger gone.Logger) {
            // Use service name directly as hostname, balancer will handle load balancing automatically
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

## Load Balancing Strategies

The balancer module provides the following load balancing strategies:

### 1. Round Robin Strategy (RoundRobinStrategy)

Selects service instances in sequence, which is the default load balancing strategy.

```go
// Already loaded by default, no additional configuration needed
// RoundRobinStrategy is loaded by default in load.go
```

### 2. Random Strategy (RandomStrategy)

Randomly selects a service instance.

```go
// Method 1: Use the provided loading function
import (
    "github.com/gone-io/gone/v2"
    "github.com/gone-io/goner/balancer"
)

func main() {
    app := gone.NewApp(
        balancer.Load,               // First load the basic balancer module
        balancer.LoadRandomStrategy, // Then load the random strategy, which will replace the default strategy
        // Load other modules...
    )
    
    // Run the application
    if err := app.Run(); err != nil {
        panic(err)
    }
}
```

### 3. Weight Strategy (WeightStrategy)

Selects service instances based on their weight, with higher weight instances having a higher probability of being selected.

```go
// Method 1: Use the provided loading function
import (
    "github.com/gone-io/gone/v2"
    "github.com/gone-io/goner/balancer"
)

func main() {
    app := gone.NewApp(
        balancer.Load,               // First load the basic balancer module
        balancer.LoadWeightStrategy, // Then load the weight strategy, which will replace the default strategy
        // Load other modules...
    )

    // Run the application
    if err := app.Run(); err != nil {
        panic(err)
    }
}
```

## Advanced Usage

### Custom Load Balancing Strategy

You can create a custom load balancing strategy by implementing the `g.LoadBalanceStrategy` interface:

```go
type MyCustomStrategy struct {
    gone.Flag
    // Custom fields
}

func (s *MyCustomStrategy) Select(ctx context.Context, instances []g.Service) (g.Service, error) {
    // Implement custom selection logic
    // ...
    return selectedInstance, nil
}

func main() {
    myStrategy := &MyCustomStrategy{}

    app := gone.NewApp(
        balancer.Load, // First load the basic balancer module
            // Load other modules...
        ).
		Load(balancer.LoadCustomerStrategy(myStrategy)) // Use LoadCustomerStrategy to load custom strategy

    // Run the application
    if err := app.Run(); err != nil {
        panic(err)
    }
}

```

### Integration with Service Discovery

The balancer module needs to be used with a service discovery component, for example, integration with nacos:

```go
import (
    "github.com/gone-io/gone/v2"
    "github.com/gone-io/goner/balancer"
    "github.com/gone-io/goner/nacos"
    "github.com/gone-io/goner/urllib"
)

func main() {
    app := gone.NewApp(
        nacos.RegistryLoad,  // Load nacos service discovery module
        balancer.Load,       // Load balancer module
        urllib.Load,         // Load urllib module
        // Load other modules...
    )
    
    // Run the application
    if err := app.Run(); err != nil {
        panic(err)
    }
}
```

## Implementation Principles

The core functions of the balancer module include:

1. **Service Discovery**: Obtain service instance lists through the injected `g.ServiceDiscovery` interface
2. **Instance Caching**: Cache obtained service instances to improve performance
3. **Instance Monitoring**: Monitor service instance changes and automatically update the cache
4. **Load Balancing**: Select an instance from available instances based on the selected strategy

## Contributing

Contributions of issues and code are welcome, please refer to the [Gone project contribution guidelines](https://github.com/gone-io/gone/blob/main/CONTRIBUTING.md).

## License

[Apache License 2.0](https://github.com/gone-io/gone/blob/main/LICENSE)