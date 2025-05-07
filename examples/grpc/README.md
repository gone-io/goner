[//]: # (desc: Gone gRPC example project demonstrating how to build gRPC services with Gone framework)

<p>
    English&nbsp ｜&nbsp <a href="README_CN.md">中文</a>
</p>

# Gone gRPC Example Project

## Project Overview
This example project demonstrates how to quickly build gRPC server and client based on Gone framework, including service registration, configuration management, protocol definition and other typical scenarios. Suitable for both beginners and advanced users.

## Directory Structure
```
.
├── config/             # Configuration directory
│   └── default.properties # Default configuration
├── go.mod              # Go module definition
├── proto/              # Protocol definition directory
│   ├── hello.pb.go     # Generated protocol code
│   ├── hello.proto     # Protocol definition file
│   └── hello_grpc.pb.go# Generated gRPC code
├── v1_client/          # v1 client
│   └── main.go         # Client entry
├── v1_server/          # v1 server
│   └── main.go         # Server entry
├── v2_client/          # v2 client (optional)
│   └── main.go         # Client entry
├── v2_server/          # v2 server (optional)
│   └── main.go         # Server entry
└── README_CN.md        # Chinese documentation
```

## Main Dependencies
- Go 1.24+
- github.com/gone-io/gone/v2 v2.1.0
- github.com/gone-io/goner/grpc v1.2.1
- github.com/gone-io/goner/viper v1.2.1
- google.golang.org/grpc v1.72.0
- google.golang.org/protobuf v1.36.6

See go.mod for dependency details.

## How to Run
### 1. Start Server
Go to v1_server directory and run:
```shell
cd v1_server
go run main.go
```
Server listens on port 9091 by default.

### 2. Start Client
Open another terminal, go to v1_client directory and run:
```shell
cd v1_client
go run main.go
```
Client will send gRPC request to server and output response.

### 3. Configuration
Customize port, host and other parameters through config/default.properties or environment variables.

## Key Code Explanation
### Server main.go
```go
func main() {
    os.Setenv("GONE_SERVER_GRPC_PORT", "9091")
    gone.
        Load(&server{}).
        Loads(goneGrpc.ServerLoad).
        Serve()
}
```
- Load service struct with gone.Load, start gRPC service with goneGrpc.ServerLoad.
- Register gRPC service implementation with RegisterGrpcServer.

### Client main.go
```go
func main() {
    gone.
        Load(&helloClient{}).
        Loads(viper.Load, gone_grpc.ClientRegisterLoad).
        Run(func(in struct {
            hello *helloClient `gone:"*"`
        }) {
            say, err := in.hello.Say(context.Background(), &proto.SayRequest{Name: "gone"})
            if err != nil {
                log.Printf("er:%v", err)
                return
            }
            log.Printf("say result: %s", say.Message)
        })
}
```
- Load client struct with gone.Load, load configuration with viper.Load, register gRPC client with gone_grpc.ClientRegisterLoad.
- Send request to server with Say method.

### Protocol Definition
proto/hello.proto defines Say service and message structure, use protoc to generate corresponding Go code.

## FAQ
1. **Port Conflict**: Make sure port 9091 is available, or modify port through configuration.
2. **Dependencies Not Installed**: Run `go mod tidy` to install dependencies first.
3. **Proto Files Not Generated**: Make sure hello.pb.go and hello_grpc.pb.go are generated with protoc.

## References
- [Gone Framework Documentation](https://github.com/gone-io/gone)
- [gRPC Official Documentation](https://grpc.io/docs/)

Feel free to open issues or join discussions if you have any questions.