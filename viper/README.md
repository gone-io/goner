# Gone Viper 组件

`gone-viper` 是 Gone 框架的配置管理组件，基于 [spf13/viper](https://github.com/spf13/viper) 实现。它为您的 Gone 应用提供灵活而强大的配置管理能力，支持多种配置源和格式，让您的应用配置变得简单高效。

## 功能特性

- 与 Gone 框架无缝集成，开箱即用
- 多种配置源支持：文件、环境变量、命令行参数
- 多格式配置文件：JSON、YAML、TOML、Properties
- 层级化配置结构和默认值机制
- 环境变量覆盖功能，提高灵活性

## 配置文件查找机制

组件会按照以下优先顺序自动查找配置文件：

1. 可执行文件所在目录
2. 可执行文件所在目录下的 `config` 子目录
3. 当前工作目录
4. 当前工作目录下的 `config` 子目录
5. 使用 `Test` 函数启动 Gone 时的附加路径：
    - go.mod 文件所在目录下的 `config` 目录
    - 当前工作目录下的 `testdata` 目录
    - 当前工作目录下的 `testdata/config` 目录
6. 如果设置了环境变量 `CONF` 或启动时使用了 `-conf` 选项，还会查找指定的配置文件路径

## 配置文件加载顺序

在同一目录中，存在多个配置文件时，组件会按以下顺序加载并合并配置：

### 默认配置文件（按顺序）
1. "default.json"
2. "default.toml"
3. "default.yaml"
4. "default.yml"
5. "default.properties"

### 环境相关配置文件
默认环境是 `local`，可通过环境变量 `ENV` 或启动参数 `-env` 修改。环境配置文件按以下顺序加载：

1. "${env}.json"
2. "${env}.toml"
3. "${env}.yaml"
4. "${env}.yml"
5. "${env}.properties"

### 测试专用配置文件
当使用 `Test` 函数启动 Gone 时，还会额外加载测试专用配置文件 `${default|env}_test.${ext}`。

例如，在以下测试代码中：

```go
func TestCase(t *testing.T){
    gone.
        Loads(
            viper.Load, // 加载配置组件
            // 其他组件...
        ).
        Test(func(){
            // 测试代码
        })
}
```

系统会按顺序加载存在的配置文件（不存在的文件会被忽略）：
1. "default.json"
2. "default.toml"
3. "default.yaml"
4. "default.yml"
5. "default.properties"
6. "default_test.json"
7. "default_test.toml"
8. "default_test.yaml"
9. "default_test.yml"
10. "default_test.properties"
11. "local.json"
12. "local.toml"
13. "local.yaml"
14. "local.yml"
15. "local.properties"
16. "local_test.json"
17. "local_test.toml"
18. "local_test.yaml"
19. "local_test.yml"
20. "local_test.properties"

**重要说明：**
1. 多个配置文件的内容会自动合并，相同键名时后加载的配置会覆盖先加载的配置
2. properties 文件的变量替换仅在同一配置文件内有效，不支持跨文件替换

## 快速开始

### 1. 安装组件

```bash
go install github.com/gone-io/goner/viper
```

### 2. 在应用中加载配置组件

```go
package main

import (
    "github.com/gone-io/v2"
    "github.com/gone-io/goner/viper"
)

func main() {
    gone.
        Loads(
            viper.Load, // 加载配置组件
            // 其他组件...
        ).
        Run() // 或使用 Serve()
}
```

### 3. 通过标签注入配置

```go
type MyService struct {
    gone.Flag
    
    // 通过 gone:"config" 标签注入配置，支持默认值
    ServerHost string `gone:"config,server.host,default=localhost"`
    ServerPort int    `gone:"config,server.port,default=8080"`
    DbURL      string `gone:"config,db.url"`
}

func (s *MyService) Start() error {
    // 使用注入的配置值
    fmt.Printf("服务运行于 %s:%d\n", s.ServerHost, s.ServerPort)
    return nil
}
```

### 4. 手动获取配置值

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

## 配置文件格式示例

### Properties 格式

```properties
# 服务器配置
server.host=localhost
server.port=8080

# 数据库配置
db.username=root
db.password=secret
db.database=mydb

# 同文件内支持变量替换
db.url=mysql://localhost:3306/${db.database}

# 日志配置
log.level=info
log.path=/var/log/myapp
```

### YAML 格式

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

## 使用环境变量覆盖配置

您可以通过环境变量覆盖配置文件中的值，提高部署灵活性。环境变量命名规则为 `GONE_配置键名`，其中配置键名中的点号 `.` 需替换为下划线 `_`。

例如，要覆盖 `server.port` 配置，可以设置环境变量：
```
GONE_SERVER_PORT=9090
```

## API 参考

### Configure 接口

```go
type Configure interface {
    Get(key string, v any, defaultVal string) error
}
```

获取配置值的核心接口，参数说明：

- `key`：配置键名，支持使用点号分隔的层级结构
- `v`：用于接收配置值的变量指针
- `defaultVal`：默认值，当配置不存在时使用

## 最佳实践建议

1. 使用层级结构组织配置，提高可读性和可维护性
2. 为关键配置提供合理的默认值，确保应用在缺少配置时仍能正常运行
3. 将敏感信息（如密码、API 密钥等）通过环境变量注入，避免硬编码到配置文件中
4. 利用不同环境配置文件（如 `dev.yaml`、`prod.yaml`）管理不同环境的配置
5. 配置键名使用小写字母和点号分隔，保持一致的命名规范