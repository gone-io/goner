# Gone Apollo 组件

## 简介

Gone Apollo 组件是基于 [Apollo](https://www.apolloconfig.com/) 配置中心的 Gone 框架集成组件，提供了配置的动态获取和实时更新功能。Apollo 是携程开源的分布式配置中心，能够集中管理应用不同环境、不同集群的配置，配置修改后能够实时推送到应用端，并且具备规范的权限、流程治理等特性。

## 快速开始

### 1. 加载 Apollo 配置组件

```go
package main

import (
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/apollo"
)

func main() {
	gone.
		Loads(
			apollo.Load, // 加载 Apollo 配置组件
			// 其他组件...
		).
		// 或者 Serve()
		Run()
}
```

### 2. 配置 Apollo 连接信息

在项目的配置文件中（如 `config/default.yaml`）添加以下配置：

```yaml
apollo.appId: YourAppId           # Apollo 应用 ID
apollo.cluster: default           # 集群名称，默认为 default
apollo.ip: http://apollo-server:8080  # Apollo 配置中心地址
apollo.namespace: application     # 命名空间，默认为 application
apollo.secret: YourSecretKey      # 访问密钥（如果启用了访问密钥验证）
apollo.isBackupConfig: true       # 是否开启备份配置
apollo.watch: true                # 是否监听配置变更
apollo.useLocalConfIfKeyNotExist: true  # 如果 Apollo 配置中不存在某个 key，是否使用本地配置文件中的值
```

### 3. 使用配置

在 Gone 组件中注入配置：

```go
type YourComponent struct {
	gone.Flag
	
	// 方式一：直接注入配置值
	DbUrl string `gone:"config,database.url"`
	
	// 方式二：通过 Configure 接口获取配置
	configure gone.Configure `gone:"*"`
}

func (c *YourComponent) AfterProp() {
	// 方式二：动态获取配置
	var port int
	err := c.configure.Get("server.port", &port, "8080")
	if err != nil {
		// 处理错误
	}
}
```

## 配置动态更新

当 `apollo.watch` 设置为 `true` 时，Apollo 组件会监听配置变更，并自动更新已注册的配置项。
**注意**：需要动态更新的字段，**必须使用指针类型**才有效。

要使配置项支持动态更新，需要在获取配置时将配置项注册到变更监听器中：

```go
type YourComponent struct {
	gone.Flag
	
	// 这些配置项将支持动态更新
	ServerPort *int    `gone:"config,server.port"`
	DbUrl      *string `gone:"config,database.url"`
}

// 配置变更后，ServerPort 和 DbUrl 的值会自动更新
```

## 配置项说明

| 配置项 | 说明 | 默认值 |
| --- | --- | --- |
| apollo.appId | Apollo 应用 ID，必须与 Apollo 配置中心中的应用 ID 一致 | - |
| apollo.cluster | 集群名称 | default |
| apollo.ip | Apollo 配置中心地址 | - |
| apollo.namespace | 命名空间 | application |
| apollo.secret | 访问密钥，用于验证客户端身份 | - |
| apollo.isBackupConfig | 是否开启备份配置，开启后会将配置保存到本地 | true |
| apollo.watch | 是否监听配置变更，开启后配置变更时会自动更新 | false |
|apollo.useLocalConfIfKeyNotExist|如果 Apollo 配置中不存在某个 key，是否使用本地配置文件中的值|true|



## 高级用法

### 多命名空间支持

Apollo 支持多个命名空间，默认使用 `application` 命名空间。如果需要使用多个命名空间，可以在配置中指定：

```yaml
apollo.namespace: application,common,custom
```

### 本地缓存配置

当 `apollo.isBackupConfig` 设置为 `true` 时，Apollo 客户端会将配置缓存到本地，当 Apollo 服务不可用时，会使用本地缓存的配置。

## 注意事项

1. 确保 Apollo 配置中心已正确部署并可访问
2. 配置项的类型转换由 Gone 框架处理，支持基本类型（如 string、int、bool 等）
3. 对于复杂类型（如结构体、数组等），Apollo 客户端会尝试将配置值解析为 JSON
4. 配置变更监听功能需要设置 `apollo.watch: true`

## 参考资料

- [Apollo 官方文档](https://www.apolloconfig.com/)
- [Gone 框架文档](https://github.com/gone-io/gone)