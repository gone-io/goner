# Gone Viper Remote

## 什么是Gone Viper Remote？

`remote`包是Gone框架的重要组件，它对Viper配置系统进行了扩展，让你的应用能够从远程配置中心（如etcd、consul等）获取配置信息。该包基于[spf13/viper/remote](https://github.com/spf13/viper/tree/master/remote)构建，专门为Gone框架优化，提供无缝集成体验。

想象一下，你有多个应用实例需要共享同一套配置，或者需要在不重启应用的情况下动态更新配置 — 这正是Gone Viper Remote的用武之地。

## 为什么选择Gone Viper Remote？

- **集中式配置管理**：所有应用实例可以从同一个配置中心获取最新配置
- **实时配置更新**：支持配置热更新，无需重启应用
- **更高的安全性**：支持加密配置，保护敏感信息
- **本地配置兜底**：当远程配置不可用时自动回退到本地配置
- **多种数据源支持**：适配多种流行的配置中心

## 开始使用

### 第一步：安装包

```bash
go get github.com/gone-io/goner/viper/remote
```

### 第二步：在应用中加载组件

```go
import (
    "github.com/gone-io/gone/v2"
    "github.com/gone-io/goner/viper/remote"
)

func main() {
    // 创建Gone应用并加载remote组件
    gone.NewApp(remote.Load).Run()
}
```

### 第三步：配置远程提供者

在你的配置文件（如`config/default.yaml`）中设置远程配置提供者：

```yaml
# config/default.yaml
viper:
  remote:
    type: yaml
    watch: true                     # 启用配置热更新
    watchDuration: 5s               # 每5秒检查一次配置变化
    useLocalConfIfKeyNotExist: true # 远程找不到时使用本地配置
    providers:
      - provider: etcd              # 提供者类型
        endpoint: localhost:2379    # 提供者地址
        path: /config/myapp         # 配置路径
        configType: json            # 配置格式
        keyring:                    # 用于加密配置的密钥(可选)
      - provider: consul
        endpoint: localhost:8500
        path: myapp/config
        configType: yaml
        keyring:
```

## 增强安全性：使用加密配置

对于敏感信息（如数据库密码、API密钥），你可以使用加密配置来提高安全性。

### 设置GPG密钥

1. 生成GPG密钥对：

```bash
# 生成GPG密钥对
gpg --gen-key

# 导出公钥(用于加密)
gpg --export > pubring.gpg

# 导出私钥(用于解密)
gpg --export-secret-keys > secring.gpg
```

2. 在配置中指定密钥文件：

```yaml
viper.remote:
  providers:
    - provider: etcd3
      endpoint: http://localhost:2379
      path: /config/secure-config
      configType: yaml
      keyring: /path/to/secring.gpg  # 指定密钥文件路径
```

### 使用加密配置的示例

```go
package main

import (
    "fmt"
    "github.com/gone-io/gone/v2"
    "github.com/gone-io/goner/viper/remote"
)

func main() {
    gone.
        NewApp(remote.Load).
        Run(func(params struct {
            apiKey string `gone:"config,secure.api.key"`
            dbPass string `gone:"config,secure.database.password"`
        }) {
            fmt.Printf("API Key: %s, DB Password: %s\n", params.apiKey, params.dbPass)
        })
}
```

## 配置详解

### Provider详细说明

每个远程提供者都由以下属性定义：

```go
type Provider struct {
    Provider   string // 提供者类型：etcd、consul等
    Endpoint   string // 提供者地址
    Path       string // 配置在提供者中的路径
    ConfigType string // 配置格式：json、yaml等
    Keyring    string // 用于加密配置的密钥(可选)
}
```

### 全局配置选项详解

| 配置项 | 说明 | 默认值 | 使用建议 |
|-------|------|-------|---------|
| viper.remote.providers | 远程配置提供者列表 | [] | 可配置多个提供者实现配置冗余 |
| viper.remote.watch | 是否启用配置热更新 | false | 生产环境建议开启 |
| viper.remote.watchDuration | 检查配置更新的间隔时间 | 5s | 根据配置变更频率调整 |
| viper.remote.useLocalConfIfKeyNotExist | 远程不存在时是否使用本地配置 | true | 建议开启，提高系统可靠性 |

## 支持的远程提供者

目前支持以下远程配置中心：

- **etcd/etcd3**：高可用的分布式键值存储，适合大规模集群
- **consul**：服务发现和配置的工具，自带健康检查
- **firestore**：Google云端的NoSQL数据库
- **nats**：高性能的分布式消息系统

## 工作原理解析

Gone Viper Remote的工作流程如下：

1. **初始化**：从本地配置文件读取远程提供者信息
2. **连接**：连接到远程配置中心
3. **加载**：从远程获取配置信息，并与本地配置合并
4. **监控**（如果启用）：周期性检查远程配置变化，及时更新

### 本地配置兜底机制

当远程配置中心不可用或某个配置键不存在时，系统会自动回退到本地配置：

1. 应用请求配置值
2. 系统先从远程获取
3. 如获取失败且`useLocalConfIfKeyNotExist`为`true`
4. 系统回退到本地配置文件
5. 如本地也不存在，则使用默认值（如有提供）

这种机制特别适合以下场景：

- **开发环境**：开发人员可在本地覆盖某些配置
- **灾备恢复**：远程配置中心不可用时，应用仍能运行
- **配置迁移**：从本地配置逐步迁移到远程配置中心

## 实用示例：完整应用

以下是一个使用etcd作为配置中心的完整示例：

```go
package main

import (
    "fmt"
    "github.com/gone-io/gone/v2"
    "github.com/gone-io/goner/viper/remote"
    "time"
)

// 定义数据库配置结构
type Database struct {
    UserName string `mapstructure:"username"`
    Pass     string `mapstructure:"password"`
}

func main() {
    gone.
        NewApp(remote.Load).
        Run(func(params struct {
            serverName string   `gone:"config,server.name"`
            serverPort int      `gone:"config,server.port"`
            dbUserName string   `gone:"config,database.username"`
            dbUserPass string   `gone:"config,database.password"`
            database   *Database `gone:"config,database"`
            key        string   `gone:"config,key.not-existed-in-etcd"`
        }) {
            // 打印配置信息
            fmt.Printf("服务名称: %s, 端口: %d\n", params.serverName, params.serverPort)
            fmt.Printf("数据库用户: %s, 密码: %s\n", params.dbUserName, params.dbUserPass)
            fmt.Printf("本地配置项: %s\n", params.key)

            // 每10秒打印一次数据库配置，演示热更新
            for i := 0; i < 10; i++ {
                fmt.Printf("数据库配置: %#+v\n", *params.database)
                time.Sleep(10 * time.Second)
            }
        })
}
```

配置文件设置：

```yaml
# config/default.yaml
viper.remote:
  type: yaml
  watch: true
  watchDuration: 5s
  useLocalConfIfKeyNotExist: true
  providers:
    - provider: etcd3
      configType: yaml
      endpoint: http://localhost:2379
      path: /config/application.yaml
    
    - provider: etcd3
      configType: yaml
      endpoint: http://localhost:2379
      path: /config/database.yaml

# 本地配置，当远程不存在时使用
key:
  not-existed-in-etcd: 1000
```

etcd中的配置内容：

```yaml
# /config/application.yaml
server.name: config-demo
server.port: 9090

# /config/database.yaml
database:
  username: config-demo
  password: config-demo-password
```

## 最佳实践建议

1. **分层配置**：将配置按功能区分，存储在不同的路径
2. **定期备份**：对远程配置中心的数据定期备份
3. **适当的监控间隔**：根据应用需求调整`watchDuration`，避免过于频繁的检查
4. **敏感信息加密**：对密码、API密钥等敏感信息使用加密存储
5. **本地配置兜底**：保持本地配置与远程配置的基本一致，作为应急措施

## 常见问题解答

1. **远程配置中心连接失败怎么办？**
    - 确保配置中心服务正常运行
    - 检查网络连接和防火墙设置
    - 系统会自动回退到本地配置

2. **如何测试配置热更新？**
    - 启动应用后，直接修改远程配置中心的值
    - 等待至少一个`watchDuration`周期
    - 观察应用日志或行为变化

3. **支持哪些配置格式？**
    - 支持JSON、YAML、TOML等主流格式
    - 由`configType`参数指定

## 许可证

本项目采用MIT许可证。详情请参阅[LICENSE](https://github.com/gone-io/goner/blob/main/LICENSE)文件。