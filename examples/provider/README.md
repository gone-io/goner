[//]: # (desc: Define Provider components for third-party integration)

<p>
    English&nbsp ｜&nbsp <a href="README_CN.md">中文</a>
</p>

# Gone Framework Standard Provider Component Example

This example demonstrates how to elegantly integrate third-party components using the Gone framework's **Provider Interface**. By implementing the standard Provider interface, you can easily incorporate third-party components into the Gone framework, achieving flexible configuration and efficient management.

## Core Features

- Standardized **Provider Interface** implementation pattern
- Support for Provider implementations with and without parameters
- Automatic dependency injection and component management
- Optimized for third-party component integration

## Technical Implementation

### Component Definition

```go
// ThirdComponent1 simulates a third-party component
type ThirdComponent1 struct {
}

// ThirdComponent2 simulates a third-party component
type ThirdComponent2 struct {
}
```

### Provider Interface Implementation

#### Provider Implementation with Parameters

```go
type provider struct {
    gone.Flag
}

func (p *provider) Provide(tagConf string) (*ThirdComponent1, error) {
    // Create third-party component based on configuration
    return &ThirdComponent1{}, nil
}
```

**Provider Interface Description:**
- `Provide(tagConf string)`: Accepts component tag configuration, supports key=value format and positional parameters
- Return value: Returns the created component instance and possible error information

#### Provider Implementation without Parameters

```go
type noneParamProvider struct {
    gone.Flag
}

func (p noneParamProvider) Provide() (*ThirdComponent2, error) {
    // Create third-party component
    return &ThirdComponent2{}, nil
}
```

### Component Registration

```go
func Load(loader gone.Loader) error {
    loader.
        MustLoad(&provider{}).
        MustLoad(&noneParamProvider{})
    return nil
}
```

## Usage Guide

### Basic Usage

```go
type YourComponent struct {
    comp1 *ThirdComponent1 `gone:"*"`
    comp2 *ThirdComponent2 `gone:"*"`
}
```

### Configuration-Driven

```go
type YourComponent struct {
    comp1 *ThirdComponent1 `gone:"*,key=value,another=123"`
    comp2 *ThirdComponent2 `gone:"*"`
}
```

**Tag Configuration Examples:**
- Basic injection: `gone:"*"`
- Configuration with injection: `gone:"*,name=myComponent"`
- Multiple configuration items: `gone:"*,host=localhost,port=8080,timeout=5s"`

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
- [Provider Mechanism Introduction](https://github.com/gone-io/gone/blob/main/docs/provider.md)