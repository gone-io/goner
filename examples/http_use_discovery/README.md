[//]: # (desc: http service discovery example using gin)

<p>
    English&nbsp ｜&nbsp <a href="README_CN.md">中文</a>
</p>

# Gone Framework HTTP Service Discovery Example

## Project Overview

This example demonstrates how to use service discovery for HTTP communication in the Gone framework. The example includes a server and a client, where the server registers with the Nacos service discovery center, and the client accesses the server using a service name rather than a specific IP address.

## Features

- Demonstrates Nacos-based service registration and discovery
- Shows how HTTP clients can access services using service names
- Supports automatic load balancing (when multiple service instances are available)
- Fully integrates with Gone framework's dependency injection features

## Environment Setup

### Starting Nacos Service

This example uses Docker Compose to start the Nacos service:

```bash
# Execute in the project root directory
docker-compose up -d nacos
```

This will start a Nacos service instance listening on port 8848.

## Code Structure

```
./
├── client/         # HTTP client example
│   └── main.go     # Client main program
├── config/         # Configuration directory
│   └── default.yaml # Default configuration
├── docker-compose.yaml # Docker environment configuration
├── go.mod          # Go module definition
├── logs/           # Log directory
└── server/         # HTTP server example
    └── main.go     # Server main program
```

## Implementation

### Server Implementation

The server (`server/main.go`) uses Gone framework's Gin component to create an HTTP service and registers it with the service discovery center through the Nacos component:

```go
func main() {
	gone.
		NewApp(goner.GinLoad, nacos.RegistryLoad, viper.Load).
		Load(&HelloController{}).
		Serve()
}
```

The server defines a simple HTTP interface:

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

### Client Implementation

The client (`client/main.go`) uses Gone framework's urllib component and balancer component to access the server through the service name:

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
			// Access service through service name
			res, err := client.
				R().
				SetSuccessResult(&data).
				Get("http://user-center/hello?name=goner")
			// ...
		})
}
```

## Configuration

The configuration file (`config/default.yaml`) contains the following key configurations:

```yaml
nacos:
  client:
    # Nacos client configuration
    namespaceId: public
    # ...
  server:
    # Nacos server address configuration
    ipAddr: "127.0.0.1"
    port: 8848
    # ...
  service:
    # Service discovery related configuration
    group: DEFAULT_GROUP
    clusterName: default

# Server configuration
server:
    port: 0  # Use random port
    service-name: user-center  # Service name
```

## Running the Example

### 1. Start Nacos Service

```bash
docker-compose up -d nacos
```

### 2. Start the Server

```bash
cd server
go run main.go
```

### 3. Start the Client

```bash
cd client
go run main.go
```

The client will send 10 requests to the server and print the response results.

## Key Points

1. **Service Registration**: The server automatically registers with the Nacos service discovery center upon startup
2. **Service Discovery**: The client accesses services through service names (`user-center`) rather than IP addresses
3. **Load Balancing**: When multiple service instances are available, the client automatically performs load balancing
4. **Zero-Config Port**: The server uses `port: 0` configuration for random port assignment, avoiding port conflicts

## Further Reading

- [Gone Framework Documentation](https://github.com/gone-io/gone)
- [Nacos Service Discovery Documentation](https://nacos.io/en/docs/v2/guide/user/service-discovery.html)