<p>
    English&nbsp ｜&nbsp <a href="README_CN.md">中文</a>
</p>

# goner/g Component

The goner/g component is one of the core components of the Gone framework, providing a series of fundamental interfaces and functionalities including log tracing, service discovery, load balancing and other core features.

## Core Features

### 1. Log Tracing (CtxLogger)

The `CtxLogger` interface is used for log tracing, which can assign a unified traceId to the same call chain for easier log tracing and issue troubleshooting.

```go
type CtxLogger interface {
    Ctx(ctx context.Context) gone.Logger
}
```

Usage example:

```go
type user struct {
    gone.Flag
    logger CtxLogger `gone:"*"` // Inject Logger
}

func (u *user) Use(ctx context.Context) (err error) {
    // Get traceId from openTelemetry context and inject into logger
    logger := u.logger.Ctx(ctx)

    logger.Infof("hello")

    return
}
```

### 2. Service Discovery (ServiceDiscovery)

The `ServiceDiscovery` interface provides service discovery and monitoring functionality, allowing clients to find available service instances and monitor instance changes.

```go
type ServiceDiscovery interface {
    // GetInstances returns all instances of the specified service
    GetInstances(serviceName string) ([]Service, error)

    // Watch creates a channel to receive updates about service instance changes
    Watch(serviceName string) (ch <-chan []Service, stop func() error, err error)
}
```

### 3. Load Balancing (LoadBalancer)

The `LoadBalancer` interface provides load balancing functionality for selecting appropriate instances from available service instances.

```go
type LoadBalancer interface {
    // GetInstance returns a service instance based on load balancing strategy
    GetInstance(ctx context.Context, serviceName string) (Service, error)
}
```

Load balancing strategy interface:

```go
type LoadBalanceStrategy interface {
    // Select chooses a service instance from the provided instance list
    Select(ctx context.Context, instances []Service) (Service, error)
}
```

### 4. Service Instance (Service)

The `Service` interface represents a service instance in the service registry, providing basic information about the service instance including identity, location, metadata and health status.

```go
type Service interface {
    // GetName returns the service name of the instance
    GetName() string

    // GetIP returns the IP address of the instance
    GetIP() string

    // GetPort returns the port number of the instance
    GetPort() int

    // GetMetadata returns metadata associated with the service instance
    GetMetadata() Metadata

    // GetWeight returns the weight of the instance
    GetWeight() float64

    // IsHealthy returns the health status of the service instance
    IsHealthy() bool
}
```

Create service instance:

```go
// NewService creates a new service instance
func NewService(name, ip string, port int, meta Metadata, healthy bool, weight float64) Service
```

## Usage Recommendations

1. When using log tracing, it's recommended to obtain `CtxLogger` instance through dependency injection and use the `Ctx()` method to inject context information when processing requests.

2. When implementing service discovery, you can choose appropriate service discovery components (such as Consul, Etcd, etc.) to implement the `ServiceDiscovery` interface.

3. When using load balancing, you can implement custom `LoadBalanceStrategy` according to actual needs, such as round-robin, random, weight-based strategies.

4. Service instance metadata can be used to store additional configuration information such as version numbers, deployment environments, etc.

## Related Components

- [goner/balancer](../balancer/README.md): Provides concrete implementation of load balancer
- [goner/consul](../consul/README.md): Service discovery implementation based on Consul
- [goner/etcd](../etcd/README.md): Service discovery implementation based on Etcd