# Gone框架 Etcd配置中心示例

本示例展示了如何使用Gone框架与Etcd配置中心集成，实现配置的集中管理和动态更新。

## 项目概述

本示例演示了以下功能：

- 使用Gone框架从Etcd配置中心读取配置
- 配置的自动监听和动态更新
- 结构化配置注入
- 本地配置与远程配置的混合使用

## 环境准备

本示例使用Docker Compose启动Etcd服务和Etcd管理界面(etcdkeeper)：

```yaml
services:
  Etcd:
    image: 'bitnami/etcd:latest'
    environment:
      - ALLOW_NONE_AUTHENTICATION=yes
      - ETCD_ADVERTISE_CLIENT_URLS=http://0.0.0.0:2379
    ports:
      - "2379:2379"
      - "2380:2380"
  etcdKeeper:
    image: evildecay/etcdkeeper
    environment:
      HOST: "0.0.0.0"
    ports:
      - "12000:8080"
    depends_on:
      - Etcd
```

### 启动环境

```bash
docker-compose up -d
```

启动后可以通过 http://localhost:12000 访问etcdkeeper管理界面。

## 配置文件结构

### 本地配置文件

`config/default.yaml` 包含Viper远程配置的设置和本地配置：

```yaml
viper.remote:
  type: yaml
  watch: true
  watchDuration: 5s
  useLocalConfIfKeyNotExist: true
  providers:
    - provider: etcd3
      configType: yaml
      endpoint: http://localhost:2379
      path:  /config/application.yaml
      keyring:

    - provider: etcd3
      configType: yaml
      endpoint: http://localhost:2379
      path:  /config/database.yaml
      keyring:


key:
  not-existed-in-etcd: 1000
```

配置说明：
- `watch: true` - 启用配置变更监听
- `watchDuration: 5s` - 监听间隔为5秒
- `useLocalConfIfKeyNotExist: true` - 当远程配置中不存在某个键时，使用本地配置
- 配置了两个Etcd配置源：`/config/application.yaml`和`/config/database.yaml`

### Etcd中的配置文件

需要将`etcd-config-files`目录下的配置文件导入到Etcd中：

**application.yaml**：
```yaml
# /config/application.yaml

server.name: config-demo
server.port: 9090
```

**database.yaml**：
```yaml
# /config/database.yaml

database:
  username: config-demo
  password: config-demo-password
```

## 代码实现

主程序`main.go`展示了如何使用Gone框架注入配置：

```go
package main

import (
	"fmt"
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/viper/remote"
	"time"
)

type Database struct {
	UserName string `mapstructure:"username"`
	Pass     string `mapstructure:"password"`
}

func main() {
	gone.
		NewApp(remote.Load).
		Run(func(params struct {
			serverName string `gone:"config,server.name"`
			serverPort int    `gone:"config,server.port"`

			dbUserName string `gone:"config,database.username"`
			dbUserPass string `gone:"config,database.password"`

			database *Database `gone:"config,database"`

			key string `gone:"config,key.not-existed-in-etcd"`
		}) {
			fmt.Printf("serverName=%s, serverPort=%d, dbUserName=%s, dbUserPass=%s, key=%s\n", params.serverName, params.serverPort, params.dbUserName, params.dbUserPass, params.key)

			for i := 0; i < 10; i++ {
				fmt.Printf("database: %#+v\n", *params.database)
				time.Sleep(10 * time.Second)
			}
		})
}
```

代码说明：

1. 通过`remote.Load`加载远程配置组件
2. 使用`gone:"config,xxx"`标签注入配置项：
   - 基本类型配置：`serverName`、`serverPort`等
   - 结构体配置：`database`
   - 本地配置：`key.not-existed-in-etcd`（仅存在于本地配置中）
3. 程序每10秒打印一次数据库配置，用于演示配置动态更新

## 运行示例

### 1. 导入配置到Etcd

使用etcdkeeper（http://localhost:12000）将`etcd-config-files`目录下的配置文件导入到Etcd中：

- 创建`/config/application.yaml`键，值为application.yaml的内容
- 创建`/config/database.yaml`键，值为database.yaml的内容

### 2. 运行示例程序

```bash
go run main.go
```

### 3. 测试配置动态更新

1. 运行程序后，通过etcdkeeper修改`/config/database.yaml`中的配置
2. 观察程序输出，约5秒后配置将自动更新

## 配置优先级

1. 远程配置（Etcd）优先级高于本地配置
2. 当远程配置中不存在某个键时，会使用本地配置（由`useLocalConfIfKeyNotExist: true`控制）

## 总结

本示例展示了Gone框架如何与Etcd配置中心集成，实现配置的集中管理和动态更新。通过简单的配置和少量代码，即可实现强大的配置管理功能，为微服务架构提供灵活的配置解决方案。