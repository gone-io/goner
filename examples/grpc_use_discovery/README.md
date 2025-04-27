[//]: # (desc: Example of gRPC with Service Discovery)

<p>
    English&nbsp ｜&nbsp <a href="README_CN.md">中文</a>
</p>

# Gone Framework gRPC Service Discovery Example

## Project Overview

This example demonstrates how to use service discovery functionality with gRPC communication in the Gone framework. The example includes a server and a client, where the server registers with the Nacos service discovery center, and the client accesses the server using a service name rather than a specific IP address.

## Features

- Demonstrates gRPC service registration and discovery based on Nacos
- Shows how gRPC clients can access services using service names
- Supports automatic load balancing (when multiple service instances are available)
- Fully integrates with Gone framework's dependency injection features

## Project Structure

```
.
├── client/             # Client code
│   └── main.go         # Client entry file
├── config/             # Configuration directory
│   └── default.yaml    # Default configuration file
├── docker-compose.yaml # Docker environment configuration
├── go.mod              # Go module definition
├── logs/               # Log directory
├── proto/              # Protocol definition directory
│   ├── hello.pb.go     # Generated protocol code
│   ├── hello.proto     # Protocol definition file
│   └── hello_grpc.pb.go# Generated gRPC code
└── server/             # Server code
    └── main.go         # Server entry file
```

## How It Works

### Service Discovery Flow

1. When the server starts, it registers its service information (service name, IP address, port, etc.) with the Nacos registry center
2. When the client starts, it queries Nacos for the service address using the service name
3. Nacos returns the service address information to the client
4. The client establishes a gRPC connection using the obtained address information
5. When server instances change, the client automatically detects and updates connections

### Key Components

- **Server**: Implements gRPC service and registers with Nacos
- **Client**: Discovers and connects to services using service names
- **Nacos**: Provides service registration and discovery functionality
- **gRPC**: Provides high-performance RPC communication

## Configuration Guide

### Server Configuration

Server configuration in `config/default.yaml`:

```yaml
server:
  grpc:
    port: 0  # Use 0 for random port
    service-name: user-center  # Service name
```

### Client Configuration

Client configures service connection through dependency injection:

```go
// Method 1: Connect using service name from configuration file
clientConn *grpc.ClientConn `gone:"*,config=grpc.service.hello.address"`

// Method 2: Use both configuration and address, configuration takes precedence
//clientConn1 *grpc.ClientConn `gone:"*,config=grpc.service.hello.address,address=127.0.0.1:9091"`

// Method 3: Directly specify address (not recommended, hardcoded)
//clientConn2 *grpc.ClientConn `gone:"*,address=127.0.0.1:9091"`
```

### Nacos Configuration

```yaml
nacos:
  client:
    namespaceId: public
    asyncUpdateService: false
    logLevel: debug
    logDir: ./logs/
  server:
    ipAddr: "127.0.0.1"
    contextPath: /nacos
    port: 8848
    scheme: http

  service:
    group: DEFAULT_GROUP
    clusterName: default
```

## Running the Example

### Prerequisites

- Docker and Docker Compose installed
- Go 1.16+

### Start Nacos

```bash
docker-compose up -d nacos
```

### Start the Server

```bash
cd server
go run main.go
```

### Start the Client

```bash
cd client
go run main.go
```

### Expected Output

Client output example:
```
2023/xx/xx xx:xx:xx say result: Hello gone
2023/xx/xx xx:xx:xx say result: Hello gone
...
```

Server output example:
```
2023/xx/xx xx:xx:xx Received: gone
2023/xx/xx xx:xx:xx Received: gone
...
```

## Code Analysis

### Server Implementation

```go
type server struct {
	gone.Flag
	proto.UnimplementedHelloServer              // Embed UnimplementedHelloServer
	grpcServer                     *grpc.Server `gone:"*"` // Inject grpc.Server
}

func (s *server) Init() {
	proto.RegisterHelloServer(s.grpcServer, s) // Register service
}

// Say implements the service defined in the protocol
func (s *server) Say(ctx context.Context, in *proto.SayRequest) (*proto.SayResponse, error) {
	log.Printf("Received: %v", in.GetName())
	return &proto.SayResponse{Message: "Hello " + in.GetName()}, nil
}
```

### Client Implementation

```go
type helloClient struct {
	gone.Flag
	proto.HelloClient // Method 1: Embed HelloClient

	// Inject connection using service name
	clientConn *grpc.ClientConn `gone:"*,config=grpc.service.hello.address"`
}

func (c *helloClient) Init() {
	c.HelloClient = proto.NewHelloClient(c.clientConn)
}
```

## Further Reading

- [Gone Framework Documentation](https://github.com/gone-io/gone)
- [Nacos Service Discovery Documentation](https://nacos.io/en/docs/v2/guide/user/service-discovery.html)
- [gRPC Official Documentation](https://grpc.io/docs/)