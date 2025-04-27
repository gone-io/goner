[//]: # (desc: simple example, use viper for configuration reading)

<p>
    English&nbsp ｜&nbsp <a href="README_CN.md">中文</a>
</p>

# Gone Simple Example

## Project Overview

This is a minimal example project for the Gone framework, demonstrating its basic usage and dependency injection features. Through this example, you can quickly understand:

- Basic project structure of the Gone framework
- How to use dependency injection
- How to read configuration files
- The principle of interface-implementation separation

## Project Structure

```
.
├── cmd/                # Program entry directory
│   ├── import.gone.go  # Dependency import file
│   └── main.go        # Main program entry
├── config/            # Configuration directory
│   └── default.yaml   # Default configuration file
├── implement/         # Interface implementation directory
│   ├── hello.go      # Service interface implementation
│   └── init.gone.go  # Implementation initialization file
├── service/          # Interface definition directory
│   └── interface.go  # Service interface definition
└── module.load.go    # Module loading file
```

## Features

1. **Simple Dependency Injection**: Implement interface injection through the `gone:"*"` tag
2. **Automatic Configuration Binding**: Use `gone:"config,app.name"` to automatically bind configuration items
3. **Interface-Implementation Separation**: Follow Go language best practices by separating interface definitions from implementations
4. **Modular Structure**: Clear directory structure for easy project expansion

## Usage Instructions

### 1. Install Dependencies

First, install the Gone framework's configuration component:

```bash
gonectr install github.com/gone-io/goner/viper
```

### 2. Run the Project

Execute in the project root directory:

```bash
gonectr run ./cmd
```

### 3. Expected Output

The program will output:
```
hello simple-app
```

## Core Code Explanation

### Interface Definition

```go
// service/interface.go
type Service interface {
    SayHello(name string) string
}
```

### Interface Implementation

```go
// implement/hello.go
type serviceImpl struct {
    gone.Flag
}

func (s serviceImpl) SayHello(name string) string {
    return fmt.Sprintf("hello %s", name)
}
```

### Main Program

```go
// cmd/main.go
func main() {
    gone.Run(func(in struct {
        service service.Service `gone:"*"`
        appName string          `gone:"config,app.name"`
    }) {
        println(in.service.SayHello(in.appName))
    })
}
```

## Configuration Description

Configuration file located at `config/default.yaml`:

```yaml
app:
  name: simple-app
```

## Extension Suggestions

1. Add more configuration items to demonstrate configuration binding capabilities
2. Add multiple interface implementations to showcase dependency injection flexibility
3. Add examples of real business scenarios such as database access

## Related Documentation

- [Gone Framework Documentation](https://github.com/gone-io/gone)
- [Configuration Center Example](https://github.com/gone-io/goner/tree/main/examples/config_center)
- [More Example Projects](https://github.com/gone-io/goner/tree/main/examples)