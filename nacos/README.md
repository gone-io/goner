# Gone Framework Nacos Component

## Overview

The Nacos component empowers the Gone framework with dynamic configuration management and service discovery capabilities by leveraging Alibaba's Nacos. This integration provides a robust solution for managing application configurations and service discovery in distributed systems.

With the Nacos component, you can:

- Centralize configuration management across your application ecosystem
- Implement real-time configuration updates without service restarts
- Support multiple configuration formats (JSON, YAML, Properties, TOML)
- Organize configurations with logical grouping and namespaces
- Maintain configuration version control and change history
- Register and discover services in a distributed environment
- Implement load balancing and service routing
- Monitor service health and availability

## Configuration Reference

### Client Configuration

The following parameters control the Nacos client behavior:

| Configuration Parameter | Description | Type | Default Value | Example |
|------------------------|-------------|------|--------------|---------|
| nacos.client.namespaceId | Namespace identifier for isolating configuration environments | string | public | "public" |
| nacos.client.timeoutMs | Request timeout in milliseconds | uint64 | 10000 | 10000 |
| nacos.client.logLevel | Client logging verbosity | string | info | "info" |
| nacos.client.logDir | Directory for client logs | string | /tmp/nacos/log | "/tmp/nacos/log" |
| nacos.client.cacheDir | Directory for client cache | string | /tmp/nacos/cache | "/tmp/nacos/cache" |

### Server Configuration

These settings define the connection to your Nacos server:

| Configuration Parameter | Description | Type | Default Value | Example |
|------------------------|-------------|------|--------------|---------|
| nacos.server.ipAddr | Nacos server address | string | - | "127.0.0.1" |
| nacos.server.contextPath | Server context path | string | /nacos | "/nacos" |
| nacos.server.port | Server port number | uint64 | 8848 | 8848 |
| nacos.server.scheme | Connection protocol | string | http | "http" |

### Configuration Properties

General configuration behavior settings:

| Configuration Parameter | Description | Type | Default Value | Example |
|------------------------|-------------|------|--------------|---------|
| nacos.dataId | Configuration data identifier | string | - | "user-center" |
| nacos.watch | Enable configuration change monitoring | bool | false | true |
| nacos.useLocalConfIfKeyNotExist | Fallback to local configuration when key not found in Nacos | bool | true | true |

### Group Configuration

Settings for organizing configurations into logical groups:

| Configuration Parameter | Description | Type | Default Value | Example |
|------------------------|-------------|------|--------------|---------|
| nacos.groups[].group | Configuration group name | string | - | "DEFAULT_GROUP" |
| nacos.groups[].format | Configuration file format | string | - | "properties" |

Supported configuration formats:
- json
- yaml/yml
- properties
- toml

## Implementation Guide

### Configuration File Setup

Create a `default.yaml` file in your project's config directory to define the Nacos client connection parameters:

```yaml
nacos:
  client:
    namespaceId: public        # Namespace identifier
  server:
    ipAddr: "127.0.0.1"        # Nacos server address
    contextPath: /nacos        # Context path
    port: 8848                 # Server port
    scheme: http               # Connection protocol
  dataId: user-center          # Configuration data identifier
  watch: true                  # Enable configuration change monitoring
  useLocalConfIfKeyNotExist: true  # Fallback to local configuration when key not found
  groups:                      # Configuration group definitions
    - group: DEFAULT_GROUP     # Default group
      format: properties       # Configuration format
    - group: database          # Database configuration group
      format: yaml            # Configuration format
```

### Code Integration

```go
func main() {
    // Initialize application with Nacos configuration loader
    gone.NewApp(nacos.Load).
        Run(func(params struct {
            // Bind individual configuration values
            serverName string `gone:"config,server.name"`    // Server name configuration
            serverPort int    `gone:"config,server.port"`    // Server port configuration
            
            // Database credentials
            dbUserName string `gone:"config,database.username"` // Database username
            dbUserPass string `gone:"config,database.password"` // Database password
            
            // Bind entire configuration section to a struct
            database *Database `gone:"config,database"`  // Complete database configuration
        }) {
            // Use the configuration values in your application
            fmt.Printf("serverName=%s, serverPort=%d\n", params.serverName, params.serverPort)
            fmt.Printf("database: %#+v\n", *params.database)
        })
}
```

### Configuration Binding

The Nacos component provides flexible configuration binding capabilities:

- Use the `gone:"config,key"` tag to mark configuration fields
- Support for binding primitive types and complex structures
- Automatic configuration hot-reloading - changes in Nacos propagate automatically to your application
- Hierarchical configuration structure with dot notation for nested properties
- Type conversion handled automatically between configuration formats and Go types