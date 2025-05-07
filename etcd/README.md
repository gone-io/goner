<p>
    English&nbsp ｜&nbsp <a href="README_CN.md">中文</a>
</p>

# goner/etcd Component

## Component Overview

The **goner/etcd** component provides service registration and discovery capabilities for the Gone framework based on etcd. This integration offers a reliable solution for service management in distributed systems.

With the **goner/etcd** component, you can:

- Register and discover services in distributed environments
- Implement service health checks and automatic failover
- Build highly available microservice architectures
- Leverage etcd's strong consistency for service coordination

## Features

### Service Registration & Discovery

- **Service Registration**: Register service instances to etcd service directory
- **Service Discovery**: Get available service instances from etcd
- **Service Monitoring**: Watch service instance changes and update service lists in real-time
- **Health Check**: TTL-based health check mechanism to automatically maintain service health status
- **Strong Consistency**: Utilize etcd's Raft consensus algorithm for data consistency

## Configuration Reference

### Client Configuration

Basic parameters for etcd client behavior:

| Parameter | Description | Type | Default | Example |
|-----------|-------------|------|---------|---------|
| etcd.endpoints | etcd server addresses | []string | ["127.0.0.1:2379"] | ["localhost:2379"] |
| etcd.username | etcd authentication username | string | "" | "username" |
| etcd.password | etcd authentication password | string | "" | "password" |
| etcd.dial-timeout | Connection timeout | time.Duration | "5s" | "10s" |
| etcd.keepalive-ttl | Health check TTL | time.Duration | "10s" | "20s" |

> For more configurations, refer to [etcd official documentation](https://pkg.go.dev/go.etcd.io/etcd/client/v3#Config)

#### Supported configuration loading and injection
Component name: **etcd.config**

```go
    gone.
        NewApp(
            //... 
        ).
        Loads(g.NamedThirdComponentLoadFunc("etcd.config", &etcd.Config{
            Endpoints: []string{"localhost:2379"},
            // Other configurations
        }))
```

### Service Registration Configuration

Service registration parameters:

- gin.http server

| Parameter | Description | Type | Default | Example |
|-----------|-------------|------|---------|---------|
| service.name | Service name | string | - | "user-service" |
| service.host | Service address | string | - | "192.168.1.100" |
| service.port | Service port | int | - | 8080 |
| service.service-use-subnet | Subnet used | string | 0.0.0.0/0 | 192.168.1.0/24 |

- grpc server
| Parameter | Description | Type | Default | Example |
|-----------|-------------|------|---------|---------|
| service.grpc.name | Service name | string | - | "user-service" |
| service.grpc.host | Service address | string | - | "192.168.1.100" |
| service.grpc.port | Service port | int | - | 8080 |
| service.grpc.service-use-subnet | Subnet used | string | 0.0.0.0/0 | 192.168.1.0/24 |

## Implementation Guide

### Configuration File Setup

Create `default.yaml` in project's config directory to define etcd client connection parameters:

```yaml
etcd:
  # Basic configurations
  endpoints: ["localhost:2379"]  # etcd server addresses
  dialTimeout: "5s"              # Connection timeout
  username: ""                   # Authentication username
  password: ""                   # Authentication password
  autoSyncInterval: "30s"        # Endpoint auto-sync interval
```

### Code Integration

#### Service Registration & Discovery

```go
func main() {
    // Initialize app and load etcd component
    gone.NewApp(etcd.RegistryLoad).Run(func(params struct {
        registry g.ServiceRegistry `gone:"*"`
        discovery g.ServiceDiscovery `gone:"*"`
    }) {
        // Create service instance
        service := g.NewService(
            "my-service",       // Service name
            "192.168.1.100",   // Service IP
            8080,               // Service port
            map[string]string{  // Metadata
                "version": "1.0.0",
            },
            true,               // Health status
            1.0,                // Weight
        )

        // After registration, component will automatically maintain TTL health check (default TTL: 10s)

        // Register service
        err := params.registry.Register(service)
        if err != nil {
            log.Fatalf("Service registration failed: %v", err)
        }

        // Discover services
        instances, err := params.discovery.GetInstances("target-service")
        if err != nil {
            log.Fatalf("Service discovery failed: %v", err)
        }

        // Watch service changes
        ch, stop, err := params.discovery.Watch("target-service")
        if err != nil {
            log.Fatalf("Failed to watch service: %v", err)
        }

        go func() {
            for services := range ch {
                // Handle service list updates
                fmt.Printf("Service list updated: %v\n", services)
            }
        }()

        // Deregister service and stop watching when app exits
        defer func() {
            _ = params.registry.Deregister(service)
            _ = stop()
        }()
        // Main application logic...
    })
}
```

## Related Links

- [etcd official documentation](https://etcd.io/docs/)
- [Gone framework documentation](https://github.com/gone-io/gone)
- [etcd API](https://github.com/etcd-io/etcd/tree/main/client/v3)