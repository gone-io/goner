# Gone Framework Consul Component

## Component Overview

The Consul component provides service registration and discovery capabilities for the Gone framework based on HashiCorp Consul. This integration offers a robust solution for service management in distributed systems.

With the Consul component, you can:

- Register and discover services in distributed environments
- Implement service health checks and automatic failover
- Build highly available microservice architectures

## Features

### Service Registration and Discovery

- **Service Registration**: Register service instances to Consul service catalog
- **Service Discovery**: Retrieve available service instances from Consul
- **Service Monitoring**: Monitor service instance changes and update service lists in real-time
- **Health Check**: TTL-based health check mechanism to automatically maintain service health status

## Configuration Reference

### Client Configuration

#### Basic Configuration

The following parameters control basic Consul client behavior:

| Parameter          | Description                          | Type          | Default       | Example           |
| ------------------ | ------------------------------------ | ------------- | ------------- | ----------------- |
| consul.address    | Consul server address                | string        | 127.0.0.1:8500 | "127.0.0.1:8500" |
| consul.scheme     | Connection protocol                  | string        | http          | "http"           |
| consul.pathPrefix | URI prefix for Consul behind API gateways | string        | -             | "/consul"        |
| consul.datacenter | Datacenter name                      | string        | -             | "dc1"            |
| consul.token      | ACL token                           | string        | -             | "your-token"     |
| consul.tokenFile  | File path containing ACL token       | string        | -             | "/path/to/token" |
| consul.waitTime   | Maximum blocking time for Watch operations | time.Duration | -             | "30s"            |
| consul.namespace  | Namespace (Consul Enterprise only)   | string        | -             | "my-namespace"   |
| consul.partition  | Partition (Consul Enterprise only)   | string        | -             | "my-partition"   |

#### TLS Configuration

The following parameters configure TLS connections for Consul client:

| Parameter                      | Description                 | Type   | Default | Example                 |
| ------------------------------ | --------------------------- | ------ | ------- | ----------------------- |
| consul.tls.address            | TLS server name (SNI)       | string | -       | "consul.example.com"   |
| consul.tls.caFile             | CA certificate file path    | string | -       | "/path/to/ca.pem"      |
| consul.tls.caPath             | CA certificate directory path | string | -       | "/path/to/certs"       |
| consul.tls.certFile           | Client certificate file path | string | -       | "/path/to/cert.pem"    |
| consul.tls.keyFile            | Client private key file path | string | -       | "/path/to/key.pem"     |
| consul.tls.insecureSkipVerify | Skip TLS host verification  | bool   | false   | true                    |

> For more configurations, refer to [Consul official documentation](https://pkg.go.dev/github.com/hashicorp/consul/api#Config)

#### Support loading configuration and injecting into services
Component name: **consul.config**

```go
    gone.
        NewApp(
            //... 
        ).
        Loads(g.NamedThirdComponentLoadFunc("consul.config", &api.Config{
            Address: "127.0.0.1:8500",
            // Other configurations
        }))
```

### Service Registration Configuration

Service registration related parameters:

- gin.http server

| Parameter                   | Description       | Type   | Default     | Example          |
| --------------------------- | ----------------- | ------ | ----------- | ---------------- |
| service.name               | Service name      | string | -           | "user-service"  |
| service.host               | Service address   | string | -           | "192.168.1.100" |
| service.port               | Service port      | int    | -           | 8080             |
| service.service-use-subnet | Subnet to use     | string | 0.0.0.0/0   | 192.168.1.0/24   |

- grpc server
| Parameter                        | Description       | Type   | Default     | Example          |
| -------------------------------- | ----------------- | ------ | ----------- | ---------------- |
| service.grpc.name               | Service name      | string | -           | "user-service"  |
| service.grpc.host               | Service address   | string | -           | "192.168.1.100" |
| service.grpc.port               | Service port      | int    | -           | 8080             |
| service.grpc.service-use-subnet | Subnet to use     | string | 0.0.0.0/0   | 192.168.1.0/24   |


## Implementation Guide

### Configuration File Setup

Create a `default.yaml` file in the project's config directory to define Consul client connection parameters:

```yaml
consul:
  # Basic configuration
  address: "127.0.0.1:8500"  # Consul server address
  scheme: "http"            # Connection protocol
  pathPrefix: ""            # URI prefix (for API gateway scenarios)
  datacenter: "dc1"         # Datacenter
  token: ""                 # ACL token
  tokenFile: ""             # ACL token file path
  waitTime: "30s"           # Maximum blocking time for Watch operations
  namespace: ""             # Namespace (Enterprise only)
  partition: ""             # Partition (Enterprise only)

  # TLS configuration
  tls:
    address: "consul.example.com"  # TLS server name
    caFile: "/path/to/ca.pem"      # CA certificate file path
    caPath: "/path/to/certs"       # CA certificate directory path
    certFile: "/path/to/cert.pem"  # Client certificate file path
    keyFile: "/path/to/key.pem"    # Client private key file path
    insecureSkipVerify: false      # Whether to skip TLS host verification
```

### Code Integration

#### Service Registration and Discovery

```go
func main() {
    // Initialize application and load Consul component
    gone.NewApp(consul.RegistryLoad).Run(func(params struct {
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
            true,               // Healthy status
            1.0,                // Weight
        )

        // After registration, the component will automatically maintain TTL health check
        // Default TTL is 20 seconds, health check interval is 10 seconds

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

        // Watch for service changes
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

        // Deregister service and stop watching when application exits
        defer func() {
            _ = params.registry.Deregister(service)
            _ = stop()
        }()
        // Main application logic...
    })
}
```


## Related Links

- [Consul Official Documentation](https://www.consul.io/docs)
- [Gone Framework Documentation](https://github.com/gone-io/gone)
- [HashiCorp Consul API](https://github.com/hashicorp/consul/api)