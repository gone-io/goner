# Gone gRPC Component Usage Guide

This document explains how to use the gRPC component in the Gone framework, covering both the traditional approach and the Gone V2 Provider mechanism.

## Getting Started

First, create a grpc directory and initialize a golang mod in it:

```bash
mkdir grpc
cd grpc
go mod init grpc_demo
```

## Writing Proto File and Generating Golang Code

### Writing Protocol File

Define a simple Hello service with a Say method:

Filename: proto/hello.proto
```proto
syntax = "proto3";

option go_package="/proto";

package Business;

service Hello {
  rpc Say (SayRequest) returns (SayResponse);
}

message SayResponse {
  string Message = 1;
}

message SayRequest {
  string Name = 1;
}
```

### Generating Golang Code

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
protoc --go_out=. --go_opt=paths=source_relative \
--go-grpc_out=. --go-grpc_opt=paths=source_relative \
proto/hello.proto
```

> For protoc installation, refer to [Protocol Buffer Compiler Installation](https://blog.csdn.net/waitdeng/article/details/139248507)

## Implementation Method 1: Traditional Approach

### Server Implementation

Filename: v1_server/main.go
```go
package main

import (
  "context"
  "github.com/gone-io/gone/v2"
  "github.com/gone-io/goner"
  goneGrpc "github.com/gone-io/goner/grpc"
  "google.golang.org/grpc"
  "grpc_demo/proto"
  "log"
)

type server struct {
  gone.Flag
  proto.UnimplementedHelloServer // Embed UnimplementedHelloServer
}

// Override the service defined in the protocol
func (s *server) Say(ctx context.Context, in *proto.SayRequest) (*proto.SayResponse, error) {
  log.Printf("Received: %v", in.GetName())
  return &proto.SayResponse{Message: "Hello " + in.GetName()}, nil
}

// Implement the RegisterGrpcServer method of the goneGrpc.Service interface, which will be automatically called when the server starts
func (s *server) RegisterGrpcServer(server *grpc.Server) {
  proto.RegisterHelloServer(server, s)
}

func main() {
  gone.
    Load(&server{}).
    Loads(goner.BaseLoad, goneGrpc.ServerLoad).
    // Start the service
    Serve()
}
```

### Client Implementation

Filename: v1_client/main.go
```go
package main

import (
	"context"
	"fmt"
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner"
	gone_grpc "github.com/gone-io/goner/grpc"
	"google.golang.org/grpc"
	"grpc_demo/proto"
	"log"
)

type helloClient struct {
	gone.Flag
	proto.HelloClient // Embed HelloClient

	host string `gone:"config,server.host"`
	port string `gone:"config,server.port"`
}

// Implement the Address method of the gone_grpc.Client interface, which will be automatically called when the client starts
// This method tells the client the address of the gRPC service
func (c *helloClient) Address() string {
	return fmt.Sprintf("%s:%s", c.host, c.port)
}

// Implement the Stub method of the gone_grpc.Client interface, which will be automatically called when the client starts
// Initialize HelloClient in this method
func (c *helloClient) Stub(conn *grpc.ClientConn) {
	c.HelloClient = proto.NewHelloClient(conn)
}

func main() {
	gone.
		Load(&helloClient{}).
		Loads(goner.BaseLoad, gone_grpc.ClientRegisterLoad).
		Run(func(in struct {
			hello *helloClient `gone:"*"` // Inject helloClient in the Run method's parameters
		}) {
			// Call the Say method to send a message to the server
			say, err := in.hello.Say(context.Background(), &proto.SayRequest{Name: "gone"})
			if err != nil {
				log.Printf("er:%v", err)
				return
			}
			log.Printf("say result: %s", say.Message)
		})
}
```

## Implementation Method 2: Gone V2 Provider Mechanism

Gone V2 introduces a powerful Provider mechanism that greatly simplifies the use of gRPC components.

### Server Implementation

Filename: v2_server/main.go
```go
package main

import (
	"context"
	"github.com/gone-io/gone/v2"
	goneGrpc "github.com/gone-io/goner/grpc"
	"google.golang.org/grpc"
	"grpc_demo/proto"
	"log"
	"os"
)

type server struct {
	gone.Flag
	proto.UnimplementedHelloServer              // Embed UnimplementedHelloServer
	grpcServer                     *grpc.Server `gone:"*"` // Inject grpc.Server
}

func (s *server) Init() {
	proto.RegisterHelloServer(s.grpcServer, s) // Register service
}

// Say  Override the service defined in the protocol
func (s *server) Say(ctx context.Context, in *proto.SayRequest) (*proto.SayResponse, error) {
	log.Printf("Received: %v", in.GetName())
	return &proto.SayResponse{Message: "Hello " + in.GetName()}, nil
}

func main() {
	// Gone's built-in default configuration component can only read configuration from environment variables
	os.Setenv("GONE_SERVER_GRPC_PORT", "9091")

	gone.
		Load(&server{}).
		Loads(goneGrpc.ServerLoad).
		// Start the service
		Serve()
}
```

### Client Implementation

Filename: v2_client/main.go
```go
package main

import (
	"context"
	"github.com/gone-io/gone/v2"
	gone_grpc "github.com/gone-io/goner/grpc"
	"google.golang.org/grpc"
	"grpc_demo/proto"
	"log"
	"os"
)

type helloClient struct {
	gone.Flag
	proto.HelloClient // Method 1: Embed HelloClient, this component only handles initialization and provides capabilities to third-party components

	// Method 2: Use directly in this component, not providing to third-party components
	//hello *proto.HelloClient

	// config=${config key},address=${service address}; //config has higher priority
	clientConn *grpc.ClientConn `gone:"*,config=grpc.service.hello.address"`

	// config and address can be used together, if config value is not found, fallback to using address
	//clientConn1 *grpc.ClientConn `gone:"*,config=grpc.service.hello.address,address=127.0.0.1:9091"`

	// address can also be used alone, not recommended as it means hardcoding
	//clientConn2 *grpc.ClientConn `gone:"*,address=127.0.0.1:9091"`
}

func (c *helloClient) Init() {
	c.HelloClient = proto.NewHelloClient(c.clientConn)
	//c.hello = &c.HelloClient
}

func main() {
	// Gone's built-in default configuration component can only read configuration from environment variables
	os.Setenv("GONE_GRPC_SERVICE_HELLO_ADDRESS", "127.0.0.1:9091")

	gone.
		Load(&helloClient{}).
		Loads(gone_grpc.ClientRegisterLoad).
		Run(func(in struct {
			hello *helloClient `gone:"*"` // Inject helloClient in the Run method's parameters
		}) {
			// Call the Say method to send a message to the server
			say, err := in.hello.Say(context.Background(), &proto.SayRequest{Name: "gone"})
			if err != nil {
				log.Printf("er:%v", err)
				return
			}
			log.Printf("say result: %s", say.Message)
		})
}
```

## Configuration File

Filename: config/default.properties
```properties
# Set grpc service port and host
server.port=9001
server.host=127.0.0.1

# Set grpc service port and host used by the client
server.grpc.port=${server.port}
server.grpc.host=${server.host}
```

## Testing

### Running the Server

```bash
go run v2_server/main.go  # or v1_server/main.go
```

The program waits for requests, and the screen displays:
```log
2024-06-19 22:02:41.971|INFO|/Users/jim/works/gone-io/gone/goner/grpc/server.go:84||Register gRPC service *main.server
2024-06-19 22:02:41.971|INFO|/Users/jim/works/gone-io/gone/goner/grpc/server.go:88||gRPC server now listen at 127.0.0.1:9091
```

### Running the Client

```bash
go run v2_client/main.go  # or v1_client/main.go
```

The program exits after execution, and the screen displays:
```log
2024-06-19 22:06:20.713|INFO|/Users/jim/works/gone-io/gone/goner/grpc/client.go:59||register gRPC client *main.helloClient on address 127.0.0.1:9091

2024/06/19 22:06:20 say result: Hello gone
```

Back in the server window, you can see that the server received the request, with a new log line:
```log
2024/06/19 22:06:08 Received: gone
```

## Comparison of Two Implementation Methods

### Traditional Approach

**Server Side**:
1. Need to implement the `RegisterGrpcServer` interface method to register services
2. Manual management of gRPC service registration process

**Client Side**:
1. Need to implement `Address` and `Stub` methods to initialize connection
2. Configuration retrieval is not flexible, address construction logic needs to be written manually

### Provider Mechanism

**Server Side**:
1. Automatically inject `*grpc.Server` through tags
2. Complete service registration in the `Init` method, conforming to Gone's component lifecycle management

**Client Side**:
1. No longer need to implement `Address` and `Stub` methods
2. Support flexible configuration methods, including:
   - Read address from configuration only
   - Use configuration with default address for fallback strategy
   - Direct hardcoded address (not recommended, but supported)

## Summary

Gone V2's Provider mechanism greatly improves the gRPC component usage experience:

1. **More Concise Code**: Removes unnecessary interface implementations and repetitive template code
2. **Better Alignment with Dependency Injection**: Automatically injects required components through tags
3. **More Flexible Configuration**: Supports multiple address acquisition strategies, improving code maintainability

## Implementation Method 3: Service Registration and Discovery

Gone framework provides service registration and discovery functionality, allowing gRPC services to be deployed and called more flexibly.

### Service Registration

The server can register itself with a service discovery center (such as Nacos), allowing clients to access services by service name rather than specific IP address.

#### Server Implementation

```go
package main

import (
	"context"
	"github.com/gone-io/gone/v2"
	goneGrpc "github.com/gone-io/goner/grpc"
	"github.com/gone-io/goner/nacos" // Import nacos component
	"github.com/gone-io/goner/viper" // Import configuration component
	"google.golang.org/grpc"
	"grpc_demo/proto"
	"log"
)

type server struct {
	gone.Flag
	proto.UnimplementedHelloServer              // Embed UnimplementedHelloServer
	grpcServer                     *grpc.Server `gone:"*"` // Inject grpc.Server
}

func (s *server) Init() {
	proto.RegisterHelloServer(s.grpcServer, s) // Register service
}

// Say Override the service defined in the protocol
func (s *server) Say(ctx context.Context, in *proto.SayRequest) (*proto.SayResponse, error) {
	log.Printf("Received: %v", in.GetName())
	return &proto.SayResponse{Message: "Hello " + in.GetName()}, nil
}

func main() {
	gone.
		NewApp(
			goneGrpc.ServerLoad,
			nacos.RegistryLoad, // Load nacos registry center
			viper.Load,       // Load configuration component
		).
		Load(&server{}).
		// Start the service
		Serve()
}
```

#### Server Configuration

The server needs to set the service name and other Nacos-related configurations in the configuration file:

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

server:
  grpc:
    port: 0  # Use 0 to indicate a random port
    service-name: user-center  # Service name
```

### Service Discovery

The client can obtain the service address from the service discovery center by service name, without hardcoding the service address.

#### Client Implementation

```go
package main

import (
	"context"
	"github.com/gone-io/gone/v2"
	gone_grpc "github.com/gone-io/goner/grpc"
	"github.com/gone-io/goner/nacos" // Import nacos component
	"github.com/gone-io/goner/viper" // Import configuration component
	"google.golang.org/grpc"
	"grpc_demo/proto"
	"log"
)

type helloClient struct {
	gone.Flag
	proto.HelloClient // Embed HelloClient

	// Connect to the service by service name
	clientConn *grpc.ClientConn `gone:"*,config=grpc.service.hello.address"`
}

func (c *helloClient) Init() {
	c.HelloClient = proto.NewHelloClient(c.clientConn)
}

func main() {
	gone.
		NewApp(
			gone_grpc.ClientRegisterLoad,
			viper.Load,       // Load configuration component
			nacos.RegistryLoad, // Load nacos registry center
		).
		Load(&helloClient{}).
		Run(func(in struct {
			hello *helloClient `gone:"*"` // Inject helloClient in the Run method's parameters
		}) {
			// Call the Say method to send a message to the server
			say, err := in.hello.Say(context.Background(), &proto.SayRequest{Name: "gone"})
			if err != nil {
				log.Printf("err: %v", err)
				return
			}
			log.Printf("say result: %s", say.Message)
		})
}
```

#### Client Configuration

The client needs to set the service name and other Nacos-related configurations in the configuration file:

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

grpc:
  service:
    hello:
      address: user-center  # Service name, not a specific IP address
```

### Advantages of Service Registration and Discovery

1. **Service Decoupling**: The client does not need to know the specific address of the server, only the service name
2. **Dynamic Scaling**: Service instances can be dynamically added or reduced, and the client automatically perceives this
3. **Load Balancing**: When there are multiple service instances, the client can automatically perform load balancing
4. **High Availability**: When a service instance fails, the client can automatically switch to other available instances
5. **Unified Management**: All services can be uniformly managed and monitored in the service discovery center