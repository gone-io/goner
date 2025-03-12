# Gone Zap 组件

`gone-zap` 是 Gone 框架的日志组件，基于 [uber-go/zap](https://github.com/uber-go/zap) 实现，提供了高性能的结构化日志记录功能。通过该组件，您可以轻松地在 Gone 应用中实现统一的日志管理，支持多种输出格式和日志级别。

## 功能特性

- 与 Gone 框架无缝集成
- 高性能的结构化日志记录
- 支持多种日志级别（Debug、Info、Warn、Error、Panic、Fatal）
- 支持控制台和文件输出
- 支持日志轮转
- 支持追踪 ID 关联
- 支持自定义日志格式


## 配置说明

```properties
# 日志级别，可选值：debug、info、warn、error、panic、fatal，默认为 info
log.level=info

# 是否启用追踪 ID，默认为 true
log.enable-trace-id=true

# 是否禁用堆栈跟踪，默认为 false
log.disable-stacktrace=false

# 堆栈跟踪级别，默认为 error
log.stacktrace-level=error

# 是否报告调用者信息，默认为 true
log.report-caller=true

# 日志编码器，可选值：console、json，默认为 console
log.encoder=console

# 日志输出路径，默认为 stdout
log.output=stdout

# 日志文件配置（当 log.output 设置为文件路径时有效）
log.filename=app.log
log.max-size=100  # 单个日志文件最大大小，单位为 MB
log.max-age=30    # 日志文件保留天数
log.max-backups=5 # 最大备份数量
log.compress=true # 是否压缩备份文件
```

## 快速开始

### 1. 加载日志组件

```go
package main

import (
    "github.com/gone-io/gone/v2"
    "github.com/gone-io/goner/zap"
)

func main() {
    gone.Loads(
        zap.Load,  // 加载日志组件
        // 其他组件...
    )
}
```

### 2. 使用日志记录

```go
type MyService struct {
    gone.Flag
    logger gone.Logger `gone:"*"`  // 注入日志器
}

func (s *MyService) DoSomething() error {
    // 记录不同级别的日志
    s.logger.Debug("调试信息")
    s.logger.Info("普通信息")
    s.logger.Warn("警告信息")
    s.logger.Error("错误信息")
    
    // 使用格式化日志
    s.logger.Infof("用户 %s 登录成功", "admin")
    
    // 记录带有上下文的日志
    s.logger.With("user_id", 123).Info("用户操作")
    
    // 记录带有错误的日志
    err := errors.New("操作失败")
    s.logger.WithError(err).Error("处理请求时出错")
    
    return nil
}
```

### 3. 创建命名日志器

```go
type UserService struct {
    gone.Flag
    logger gone.Logger `gone:"*"`
}

func (s *UserService) Init() {
    // 创建带有模块名称的日志器
    s.logger = s.logger.Named("user-service")
}

func (s *UserService) CreateUser() {
    // 日志输出会包含模块名称前缀
    s.logger.Info("创建用户")
    // 输出: [user-service] 创建用户
}
```

## API 参考

### Logger 接口

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

## 最佳实践

1. 为不同的模块创建命名日志器，便于日志分类和过滤
2. 在生产环境中使用 JSON 格式的日志，便于日志收集和分析
3. 合理设置日志级别，避免过多的调试日志影响性能
4. 使用结构化日志记录关键信息，如用户 ID、请求 ID 等
5. 与 Tracer 组件结合使用，在日志中包含追踪 ID
6. 对于敏感信息，如密码、令牌等，避免直接记录到日志中

## 注意事项

1. `Panic` 和 `Fatal` 级别的日志会导致程序终止，请谨慎使用
2. 在高并发场景下，过多的日志记录可能会影响性能，建议适当调整日志级别
3. 日志文件轮转功能依赖于 [lumberjack](https://github.com/natefinch/lumberjack) 库，确保正确配置相关参数