[//]: # (desc: Example of using Provide function to integrate third-party components)

<p>
    English&nbsp ｜&nbsp <a href="README_CN.md">中文</a>
</p>

# Gone Framework Provide Function Example

This example demonstrates how to elegantly integrate third-party components using the **Provide function** of the Gone framework. Through simple tag configuration and dependency injection mechanisms, you can easily incorporate third-party components into the Gone framework, achieving flexible configuration and efficient management of components.

## Core Features

- Standardized **Provide function** implementation pattern
- Powerful tag configuration parsing mechanism
- Automatic dependency injection and logging management
- Optimized for third-party component integration

## Technical Implementation

### Component Definition

```go
// ThirdComponent simulates a third-party component
type ThirdComponent struct {
}
```

### Provide Function Implementation

```go
func provide(tagConf string, in struct {
	logger gone.Logger `gone:"*"`
}) (*ThirdComponent, error) {
	confMap, confKeys := gone.TagStringParse(tagConf)
	in.logger.Infof("confMap => %#v\nconfKeys=>%#v", confMap, confKeys)

	// Create third-party component based on different configurations
	return &ThirdComponent{}, nil
}
```

**Provide Function Parameters:**
- `tagConf string`: Receives component tag configuration, supports key=value format and positional parameters
- `in struct`: Dependency injection struct, used to inject required system components
  - `logger gone.Logger`: System logging component, injected via `gone:"*"` tag

**Return Values:**
- `*ThirdComponent`: Returns the created component instance
- `error`: Returns possible error information

### Component Registration

```go
func Load(loader gone.Loader) error {
	return loader.Load(gone.WrapFunctionProvider(provide))
}
```

Use `gone.WrapFunctionProvider` to wrap the Provide function as a standard Gone framework component.

## Usage Guide

### Basic Usage

```go
type YourComponent struct {
    third *ThirdComponent `gone:"*"`
}
```

### Configuration Driven

```go
type YourComponent struct {
    third *ThirdComponent `gone:"*,key=value,another=123"`
}
```

**Tag Configuration Examples:**
- Basic injection: `gone:"*"`
- Configuration with name: `gone:"*,name=myComponent"`
- Multiple configurations: `gone:"*,host=localhost,port=8080,timeout=5s"`

## Best Practices

### Use Cases

- Dynamic configuration for component initialization
- Standardized integration of third-party components
- Encapsulation of complex component initialization logic
- Precise management of component lifecycle

### Configuration Management

- Use `gone.TagStringParse` to parse tag configurations
  - Returns `confMap`: Contains all key=value format configurations
  - Returns `confKeys`: Maintains the original order of configuration items
- Supports default values and required field validation
- Configuration items support multiple formats:
  - Positional parameters: `gone:"*,redis"`
  - Key-value pairs: `gone:"*,driver=redis"`
  - Mixed usage: `gone:"*,redis,port=6379"`

### Error Handling

- Configuration Validation
  - Check required configuration items
  - Validate configuration value format and range
  - Handle configuration conflicts

- Component Initialization
  - Capture and wrap errors from third-party components
  - Provide clear error messages
  - Implement error recovery mechanisms

## Additional Resources

- [Gone Framework Official Documentation](https://github.com/gone-io/gone)
- [Dependency Injection Best Practices Guide](https://github.com/gone-io/gone/blob/main/docs/wrap-function-provider.md)