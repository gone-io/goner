[//]: # (desc: Using Zap Logger in Gone Framework)

<p>
    English&nbsp ｜&nbsp <a href="README_CN.md">中文</a>
</p>

# Using Zap Logger in Gone Framework

This example demonstrates how to use the Zap logging system in the Gone framework, including basic logging and implementation of custom log encoders.

## Features

- Use Gone framework's Logger interface for logging
- Support custom Zap Encoder for personalized log formats
- Seamless integration with Gone framework

## Examples

### 1. Using Native zap.Logger

In the Gone framework, we can directly inject `*zap.Logger` to use Zap's native functionality. This approach allows us to fully utilize all of Zap's advanced features, including structured logging, performance optimization, and flexible log configuration.

Here's a basic example:

```go
type UseOriginZap struct {
    gone.Flag
    zap *zap.Logger `gone:"*"`
}

func (s *UseOriginZap) PrintLog() {
    s.zap.Info("hello", zap.String("name", "gone io"))
}
```

This example shows:
- How to inject the native zap.Logger into a struct
- How to use zap.Logger for structured logging

Advantages of using native zap.Logger:

1. **Structured Logging**
   - Supports various field types: String, Int, Bool, etc.
   - Type-safe field values, avoiding runtime errors
   - Efficient serialization performance

2. **Rich Logging Methods**
   ```go
   // Different log levels
   s.zap.Debug("Debug message", zap.Int("code", 100))
   s.zap.Info("Info message", zap.String("user", "admin"))
   s.zap.Warn("Warning message", zap.Bool("critical", false))
   s.zap.Error("Error message", zap.Error(err))

   // Structured logging with context
   s.zap.Info("User login",
       zap.String("username", "admin"),
       zap.String("ip", "192.168.1.1"),
       zap.Int64("timestamp", time.Now().Unix()),
   )
   ```

3. **Performance Optimization**
   - Avoids unnecessary string formatting
   - Minimizes memory allocation
   - Efficient serialization process

### 2. Using Gone's Logger Interface

Besides using the native zap.Logger, we can also directly inject the `gone.Logger` interface to use logging functionality. Here's a simple example:

```go
type UseGoneLogger struct {
    gone.Flag
    logger gone.Logger `gone:"*"`
}

func (u *UseGoneLogger) PrintLog() {
    u.logger.Infof("hello %s", "GONE IO")
}
```

This example shows:
- How to inject Logger into a struct
- How to use Logger for formatted logging

### 3. Custom Encoder

If you need to customize the log format, you can implement Zap's Encoder interface:

```go
type UseCustomerEncoder struct {
    zapcore.Encoder
    gone.Flag
}

func NewUseCustomerEncoder() *UseCustomerEncoder {
    return &UseCustomerEncoder{
        Encoder: zapcore.NewConsoleEncoder(zap.NewProductionEncoderConfig()),
    }
}

func (e *UseCustomerEncoder) EncodeEntry(entry zapcore.Entry, fields []zap.Field) (*buffer.Buffer, error) {
    //Customize your log format here
    return e.Encoder.EncodeEntry(entry, fields)
}
```

To use the custom Encoder, simply load it into Gone during initialization:

```go
package main

import (
	"github.com/gone-io/gone/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

var _ zapcore.Encoder = (*UseCustomerEncoder)(nil)

// Demonstrate how to use custom Encoder and load it into gone
func init() {
	gone.Load(NewUseCustomerEncoder())
}

func NewUseCustomerEncoder() *UseCustomerEncoder {
	return &UseCustomerEncoder{
		Encoder: zapcore.NewConsoleEncoder(zap.NewProductionEncoderConfig()),
	}
}

type UseCustomerEncoder struct {
	zapcore.Encoder
	gone.Flag
}

func (e *UseCustomerEncoder) EncodeEntry(entry zapcore.Entry, fields []zap.Field) (*buffer.Buffer, error) {
	//do something
	return e.Encoder.EncodeEntry(entry, fields)
}
```

## Usage

1. Ensure Gone framework and Zap logging components are installed in your project:
```bash
go get github.com/gone-io/gone/v2
go get github.com/gone-io/goner/zap
```

2. Import the necessary packages in your code:
```go
import (
    "github.com/gone-io/gone/v2"
    "go.uber.org/zap"
)
```

3. Choose to use either the basic Logger interface or implement a custom Encoder as needed

4. Run your application, and logs will be output according to the configured settings

## Notes

- Gone framework's Logger interface is a wrapper for Zap, providing a more convenient way to use it
- When customizing Encoder, you need to implement the `zapcore.Encoder` interface
- It's recommended to use appropriate log levels and formatting configurations in production environments