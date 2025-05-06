[//]: # (desc: 在Gone框架中使用Zap日志)

<p>
    <a href="README.md">English</a>&nbsp ｜&nbsp 中文
</p>

# 在Gone框架中使用Zap日志

本示例展示如何在Gone框架中使用Zap日志系统，包括基本的日志记录和自定义日志编码器的实现。

## 功能特点

- 使用Gone框架的Logger接口进行日志记录
- 支持自定义Zap Encoder实现个性化日志格式
- 与Gone框架无缝集成

## 示例说明

### 1. 使用原生的zap.Logger

在Gone框架中，我们可以直接注入`*zap.Logger`来使用Zap的原生功能。这种方式允许我们充分利用Zap的所有高级特性，包括结构化日志记录、性能优化和灵活的日志配置。

以下是一个基本示例：

```go
type UseOriginZap struct {
    gone.Flag
    zap *zap.Logger `gone:"*"`
}

func (s *UseOriginZap) PrintLog() {
    s.zap.Info("hello", zap.String("name", "gone io"))
}
```

这个示例展示了：
- 如何在结构体中注入原生的zap.Logger
- 如何使用zap.Logger记录结构化的日志信息

使用原生zap.Logger的优势：

1. **结构化日志记录**
   - 支持多种字段类型：String、Int、Bool等
   - 字段值类型安全，避免运行时错误
   - 高效的序列化性能

2. **丰富的日志方法**
   ```go
   // 不同级别的日志
   s.zap.Debug("调试信息", zap.Int("code", 100))
   s.zap.Info("普通信息", zap.String("user", "admin"))
   s.zap.Warn("警告信息", zap.Bool("critical", false))
   s.zap.Error("错误信息", zap.Error(err))

   // 带有上下文的结构化日志
   s.zap.Info("用户登录",
       zap.String("username", "admin"),
       zap.String("ip", "192.168.1.1"),
       zap.Int64("timestamp", time.Now().Unix()),
   )
   ```

3. **性能优化**
   - 避免不必要的字符串格式化
   - 内存分配最小化
   - 高效的序列化过程

### 2. 使用Gone的Logger接口

除了使用原生的zap.Logger，我们也可以直接注入`gone.Logger`接口来使用日志功能。以下是一个简单的示例：

```go
type UseGoneLogger struct {
    gone.Flag
    logger gone.Logger `gone:"*"`
}

func (u *UseGoneLogger) PrintLog() {
    u.logger.Infof("hello %s", "GONE IO")
}
```

这个示例展示了：
- 如何在结构体中注入Logger
- 如何使用Logger记录格式化的日志信息

### 2. 自定义Encoder

如果你需要自定义日志格式，可以通过实现Zap的Encoder接口来实现：

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
    //在这里自定义你的日志格式
    return e.Encoder.EncodeEntry(entry, fields)
}
```

要使用自定义的Encoder，只需要在初始化时将其加载到Gone中：

```go
package main

import (
	"github.com/gone-io/gone/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

var _ zapcore.Encoder = (*UseCustomerEncoder)(nil)

// 演示如何使用自定义的Encoder并加载到gone中
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

## 使用方法

1. 确保项目中已经安装了Gone框架和Zap日志组件：
```bash
go get github.com/gone-io/gone/v2
go get github.com/gone-io/goner/zap
```

2. 在你的代码中引入必要的包：
```go
import (
    "github.com/gone-io/gone/v2"
    "go.uber.org/zap"
)
```

3. 根据需要选择使用基本的Logger接口或实现自定义的Encoder

4. 运行你的应用程序，日志将按照配置的方式进行输出

## 注意事项

- Gone框架的Logger接口是对Zap的封装，提供了更简便的使用方式
- 自定义Encoder时需要实现`zapcore.Encoder`接口
- 建议在生产环境中使用适当的日志级别和格式化配置
