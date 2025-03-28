# Gone Zap Component

`gone-zap` is a logging component for the Gone framework, implemented based on [uber-go/zap](https://github.com/uber-go/zap), providing high-performance structured logging functionality. With this component, you can easily implement unified log management in Gone applications, supporting various output formats and log levels.

## Features

- Seamless integration with Gone framework
- High-performance structured logging
- Support for multiple log levels (Debug, Info, Warn, Error, Panic, Fatal)
- Console and file output support
- Log rotation support
- Trace ID correlation support
- Custom log format support

## Configuration

```properties
# Log level, available values: debug, info, warn, error, panic, fatal, default is info
log.level=info

# Whether to disable stack trace, default is false
log.disable-stacktrace=false

# Stack trace level, default is error
log.stacktrace-level=error

# Whether to report caller information, default is true
log.report-caller=true

# Log encoder, available values: console, json, default is console
log.encoder=console

# Log output path, default is stdout
log.output=stdout

# Log file configuration (effective when log.output is set to a file path)
log.filename=app.log
log.max-size=100  # Maximum size of single log file in MB
log.max-age=30    # Number of days to retain log files
log.max-backups=5 # Maximum number of backup files
log.compress=true # Whether to compress backup files
```

## Quick Start

### 1. Load the Logger Component

```go
package main

import (
    "github.com/gone-io/gone/v2"
    "github.com/gone-io/goner/zap"
)

func main() {
    gone.Loads(
        zap.Load,  // Load logger component
        // Other components...
    )
}
```

### 2. Using the Logger

```go
type MyService struct {
    gone.Flag
    logger gone.Logger `gone:"*"` // Inject logger
}

func (s *MyService) DoSomething() error {
    // Log different levels
    s.logger.Debug("Debug message")
    s.logger.Info("Normal message")
    s.logger.Warn("Warning message")
    s.logger.Error("Error message")
    
    // Use formatted logging
    s.logger.Infof("User %s logged in successfully", "admin")
    
    // Log with context
    s.logger.With("user_id", 123).Info("User operation")
    
    // Log with error
    err := errors.New("operation failed")
    s.logger.WithError(err).Error("Error processing request")
    
    return nil
}
```

### 3. Creating Named Logger

```go
type UserService struct {
    gone.Flag
    logger gone.Logger `gone:"*"`
}

func (s *UserService) Init() {
    // Create logger with module name
    s.logger = s.logger.Named("user-service")
}

func (s *UserService) CreateUser() {
    // Log output will include module name prefix
    s.logger.Info("Creating user")
    // Output: [user-service] Creating user
}
```

## API Reference

### Logger Interface

```go
type Logger interface {
    Debug(args ...interface{})
    Info(args ...interface{})
    Warn(args ...interface{})
    Error(args ...interface{})
    Panic(args ...interface{})
    Fatal(args ...interface{})
    
    Debugf(template string, args ...interface{})
    Infof(template string, args ...interface{})
    Warnf(template string, args ...interface{})
    Errorf(template string, args ...interface{})
    Panicf(template string, args ...interface{})
    Fatalf(template string, args ...interface{})
    
    With(args ...interface{}) Logger
    WithError(err error) Logger
    Named(name string) Logger
}
```

## Best Practices

1. Create named loggers for different modules to facilitate log categorization and filtering
2. Use JSON format logs in production environment for easier log collection and analysis
3. Set appropriate log levels to avoid performance impact from excessive debug logs
4. Use structured logging to record key information such as user ID, request ID, etc.
5. Combine with Tracer component to include trace IDs in logs
6. Avoid logging sensitive information such as passwords and tokens directly

## Notes

1. `Panic` and `Fatal` level logs will terminate the program, use with caution
2. In high-concurrency scenarios, excessive logging may impact performance, consider adjusting log levels accordingly
3. Log rotation functionality depends on the [lumberjack](https://github.com/natefinch/lumberjack) library, ensure proper configuration of related parameters