# Gone Viper 组件

`gone-viper` 是 Gone 框架的配置管理组件，基于 [spf13/viper](https://github.com/spf13/viper) 实现，提供了灵活的配置管理功能。通过该组件，您可以轻松地在 Gone 应用中管理各种配置，支持多种配置源和格式。

## 功能特性

- 与 Gone 框架无缝集成
- 支持多种配置源：文件、环境变量、命令行参数等
- 支持多种配置格式：JSON、YAML、TOML、Properties 等
- 支持配置热重载
- 支持配置层级和默认值
- 支持环境变量覆盖配置


## 配置文件位置

组件默认会按照以下顺序查找配置文件：

1. `./config/default.properties`：默认配置文件
2. `./config/${env}.properties`：环境特定配置文件，其中 `${env}` 是由环境变量 `GONE_ENV` 指定的环境名称

您可以通过环境变量 `GONE_CONFIG_PATH` 自定义配置文件路径。

## 快速开始

### 1. 加载配置组件

```go
package main

import (
	"github.com/gone-io/v2"
	"github.com/gone-io/goner/viper"
)

func main() {
	gone.
		Loads(
			viper.Load, //加载配置组件
			// 其他组件...
		).
		// 或者 Serve()
		Run()
}

```

### 2. 使用配置注入

```go
type MyService struct {
    gone.Flag
    
    // 通过 gone:"config" 标签注入配置
    ServerHost string `gone:"config,server.host,default=localhost"`
    ServerPort int    `gone:"config,server.port,default=8080"`
    DbURL      string `gone:"config,db.url"`
}

func (s *MyService) Start() error {
    // 使用注入的配置值
    fmt.Printf("Server running at %s:%d\n", s.ServerHost, s.ServerPort)
    return nil
}
```

### 3. 手动获取配置

```go
type MyComponent struct {
    gone.Flag
    conf gone.Configure `gone:"*"`  // 注入配置管理器
}

func (c *MyComponent) DoSomething() error {
    // 获取字符串配置
    var host string
    err := c.conf.Get("server.host", &host, "localhost")
    if err != nil {
        return err
    }
    
    // 获取整数配置
    var port int
    err = c.conf.Get("server.port", &port, "8080")
    if err != nil {
        return err
    }
    
    // 获取复杂结构体配置
    var dbConfig struct {
        URL      string
        Username string
        Password string
    }
    err = c.conf.Get("db", &dbConfig, "")
    if err != nil {
        return err
    }
    
    return nil
}
```

## 配置格式

### Properties 格式示例

```properties
# 服务器配置
server.host=localhost
server.port=8080

# 数据库配置
db.url=mysql://localhost:3306/mydb
db.username=root
db.password=secret

# 日志配置
log.level=info
log.path=/var/log/myapp
```

### YAML 格式示例

```yaml
server:
  host: localhost
  port: 8080

db:
  url: mysql://localhost:3306/mydb
  username: root
  password: secret

log:
  level: info
  path: /var/log/myapp
```

## 环境变量覆盖

您可以使用环境变量覆盖配置文件中的值。环境变量名称格式为 `GONE_配置键名`，其中配置键名中的点号 `.` 替换为下划线 `_`。

例如，要覆盖 `server.port` 配置，可以设置环境变量 `GONE_SERVER_PORT=9090`。

## API 参考

### Configure 接口

```go
type Configure interface {
    Get(key string, v any, defaultVal string) error
}
```

用于获取配置值的接口，参数说明：

- `key`：配置键名，支持点号分隔的层级结构
- `v`：用于存储配置值的变量指针
- `defaultVal`：默认值，当配置不存在时使用

## 最佳实践

1. 使用层级结构组织配置，提高可读性和可维护性
2. 为所有配置提供合理的默认值，确保应用在缺少配置时仍能正常运行
3. 敏感信息（如密码、API 密钥等）应通过环境变量注入，避免硬编码在配置文件中
4. 使用不同的环境配置文件（如 `dev.properties`、`prod.properties`）管理不同环境的配置
5. 配置键名使用小写字母和点号，遵循一致的命名规范