[//]: # (desc: simple example, use viper for configuration reading)

<p>
    English&nbsp ｜&nbsp <a href="README_CN.md">中文</a>
</p>

# Gone Simple Example

## Project Overview

This is a simple example project for the Gone framework, demonstrating the basic usage and dependency injection features of the Gone framework. Through this example, you can quickly understand:

- Standard project structure of the Gone framework
- Usage of dependency injection
- Configuration file reading (using Viper)
- Separation of interfaces and implementations
- CLI and Server application modes

## Project Structure

```
.
├── cmd/                    # Program entry directory
│   ├── cli/               # Command line application entry
│   │   ├── import.gone.go # Dependency import file
│   │   └── main.go        # CLI main program
│   └── server/            # Server application entry
│       ├── import.gone.go # Dependency import file
│       ├── init.gone.go   # Initialization file
│       └── main.go        # Server main program
├── config/                # Configuration file directory
│   ├── default.yaml       # Default configuration file
│   ├── dev.yaml          # Development environment configuration
│   ├── local.yaml        # Local environment configuration
│   ├── prod.yaml         # Production environment configuration
│   └── test.yaml         # Test environment configuration
├── internal/              # Internal code directory
│   ├── controller/        # Controller directory
│   │   └── hello.go      # Hello controller
│   ├── interface/         # Interface definition directory
│   │   ├── entity/       # Entity definitions
│   │   ├── mock/         # Mock implementations
│   │   └── service/      # Service interfaces
│   │       └── i_server.go # IService interface definition
│   ├── module/           # Module implementation directory
│   │   ├── hello/        # Hello module
│   │   │   └── hello.go  # Service interface implementation
│   │   └── user/         # User module
│   ├── pkg/              # Internal utility packages
│   │   ├── e/           # Error definitions
│   │   └── utils/       # Utility functions
│   └── router/           # Route definitions
│       ├── auth_router.go # Authentication routes
│       └── pub_router.go  # Public routes
├── asserts/              # Static assets directory
├── docker-compose.yaml   # Docker Compose configuration
├── Dockerfile           # Docker image build file
├── Makefile            # Make build script
├── module.load.go      # Module loading file
└── pacakge.go          # Package definition file
```

## Features

1. **Simple Dependency Injection**: Interface injection through `gone:"*"` tags
2. **Automatic Configuration Binding**: Automatic binding of configuration items using `gone:"config,app.name"`
3. **Interface and Implementation Separation**: Following Go language best practices, separating interface definitions and implementations
4. **Modular Structure**: Using internal directory structure with clear module division
5. **Multi-Environment Configuration**: Support for different environment configuration files
6. **Dual Mode Support**: Supporting both CLI and Server running modes

## Usage Instructions

### 1. Install Dependencies

The project has been configured with necessary dependencies, including:
- `github.com/gone-io/gone/v2` - Gone core framework
- `github.com/gone-io/goner/viper` - Viper configuration component
- `github.com/gone-io/goner/gin` - Gin web framework component

### 2. Run CLI Application

Execute in the project root directory:

```bash
go run ./cmd/cli
```

### 3. Run Server Application

Execute in the project root directory:

```bash
go run ./cmd/server
```

### 4. Expected Output

**CLI Application Output:**
```
hello root-app
```

**Server Application Output:**
```
after server start
press `ctr + c` to stop!
```

## Core Code Explanation

### Interface Definition

```go
// internal/interface/service/i_server.go
type IService interface {
    SayHello(name string) string
}
```

### Interface Implementation

```go
// internal/module/hello/hello.go
type serviceImpl struct {
    gone.Flag
}

func (s serviceImpl) SayHello(name string) string {
    return fmt.Sprintf("hello %s", name)
}
```

### CLI Main Program

```go
// cmd/cli/main.go
func main() {
    gone.Run(func(in struct {
        service service.IService `gone:"*"`
        appName string           `gone:"config,app.name"`
    }) {
        println(in.service.SayHello(in.appName))
    })
}
```

### Server Main Program

```go
// cmd/server/main.go
func main() {
    gone.Serve()
}
```

### Module Loading Configuration

```go
// module.load.go
var loaders = []gone.LoadFunc{
    viper.Load,
}

func GoneModuleLoad(loader gone.Loader) error {
    var ops []*g.LoadOp
    for _, f := range loaders {
        ops = append(ops, g.F(f))
    }
    return g.BuildOnceLoadFunc(ops...)(loader)
}
```

## Configuration

Configuration file located at `config/default.yaml`:

```yaml
app:
  name: root-app
  env: dev
```

The project supports multi-environment configuration, allowing different configuration files to be specified through environment variables or startup parameters.

## Docker Support

The project includes Docker-related configurations:

- `Dockerfile` - For building application images
- `docker-compose.yaml` - For local development and testing
- `.dockerignore` - Files to ignore during Docker build

## Extension Suggestions

1. **Add Web Interfaces**: HTTP interfaces can be added in the controller directory
2. **Database Integration**: GORM or XORM can be integrated for database operations
3. **Middleware Support**: Authentication, logging, rate limiting and other middleware can be added
4. **Configuration Extension**: More configuration items can be added to demonstrate configuration binding capabilities
5. **Test Coverage**: Unit tests and integration tests can be added

## Related Documentation

- [Gone Framework Documentation](https://github.com/gone-io/gone)
- [Goner Component Library](https://github.com/gone-io/goner)
- [More Example Projects](https://github.com/gone-io/goner/tree/main/examples)
