<p>
    English&nbsp ｜&nbsp <a href="README_CN.md">中文</a>
</p>

# goner/cmux Component

**goner/cmux** is a component for multiplexing multiple protocols on the same port, allowing you to handle different protocol requests like HTTP and gRPC on a single port. This component is implemented based on [soheilhy/cmux](https://github.com/soheilhy/cmux).

## Features

- Supports handling multiple protocols on the same port
- Supports multiplexing of HTTP and gRPC protocols
- Automatic protocol detection and distribution
- Seamless integration with the gone framework

## Installation

```bash
go get github.com/gone-io/goner/cmux
```

## Configuration Parameters

The following parameters can be set in the configuration file:

```properties
# Server network type, default is tcp
server.network=tcp

# Server address, if not set, will use host and port combination
server.address=

# Server hostname, default is empty
server.host=

# Server port number, default is 8080
server.port=8080
```

## Basic Usage

1. First, load the cmux component in your application:

```go
package main

import (
    "github.com/gone-io/gone/v2"
    "github.com/gone-io/goner/cmux"
)

func main() {
    gone.Run(
        cmux.Load,
        // ... other components
    )
}
```

2. Using cmux in HTTP service:

```go
type server struct {
    gone.Flag
    keeper gone.GonerKeeper `gone:"*"`
    // ...
}

func (s *server) initListener() error {
    goner := s.keeper.GetGonerByName(cmux.Name)
    if goner != nil {
        if muxServer, ok := goner.(cmux.CMuxServer); ok {
            s.listener = muxServer.Match(cmux.HTTP1Fast())
            s.address = muxServer.GetAddress()
            return nil
        }
    }
    // Fallback to normal TCP listening
    return s.createListener()
}
```

3. Using cmux in gRPC service:

```go
func (s *server) initListener() error {
    goner := s.keeper.GetGonerByName(cmux.Name)
    if goner != nil {
        if muxServer, ok := goner.(cmux.CMuxServer); ok {
            s.listener = muxServer.MatchWithWriters(
                cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/grpc"),
            )
            s.address = muxServer.GetAddress()
            return nil
        }
    }
    // Fallback to normal TCP listening
    return s.createListener()
}
```

## API Interface

### CMuxServer

```go
type CMuxServer interface {
    // Match gets the corresponding listener based on the matcher
    Match(matcher ...cmux.Matcher) net.Listener
    
    // MatchWithWriters gets the corresponding listener based on the writer matcher
    MatchWithWriters(matcher ...cmux.MatchWriter) net.Listener
    
    // GetAddress gets the server address
    GetAddress() string
}
```

## Best Practices

1. Priority setting: The cmux component uses `gone.HighStartPriority()` to ensure it starts before other services.

2. Error handling: It is recommended to implement appropriate fallback strategies when cmux is unavailable.

3. Protocol matching order: When setting multiple protocol matchers, it is recommended to follow this order:
   - gRPC (HTTP/2)
   - HTTP/1.x
   - Other protocols

4. Monitoring and logging: The cmux component integrates with gone's logging and tracing system for easy monitoring and debugging.

## Notes

1. Ensure all necessary parameters are correctly configured when using cmux.

2. Remember to call the Stop method when stopping the service to ensure proper resource release.

3. Pay special attention to protocol detection configuration when using TLS.

4. It is recommended to conduct thorough testing in the development environment to ensure all protocols work correctly.