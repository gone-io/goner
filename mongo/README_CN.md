<p align="center">
    <a href="README.md">English</a>&nbsp ｜&nbsp 中文
</p>

# goner/mongo 组件和 Gone MongoDB 集成

此包为 Gone 应用程序提供 MongoDB 集成，提供易于使用的客户端配置和管理。

## 特性

- 与 Gone 依赖注入系统轻松集成
- 支持多个 MongoDB 客户端实例
- 单例客户端实例管理
- 全面的配置选项
- 连接池和超时管理
- 身份验证支持

## 安装

```bash
gonectl install goner/mongo
```

## 配置

在项目的配置目录中创建 `default.yaml` 文件，包含以下 MongoDB 配置：

```yaml
mongo:
  uri: "mongodb://localhost:27017"     # MongoDB 连接 URI
  database: "myapp"                    # 默认数据库名称
  username: ""                         # 可选：身份验证用户名
  password: ""                         # 可选：身份验证密码
  authSource: "admin"                  # 可选：身份验证数据库
  maxPoolSize: 100                     # 可选：最大连接池大小
  minPoolSize: 0                       # 可选：最小连接池大小
  maxConnIdleTime: "30m"               # 可选：最大连接空闲时间
  connectTimeout: "10s"                # 可选：连接超时
  socketTimeout: "30s"                 # 可选：套接字超时
  serverSelectionTimeout: "30s"        # 可选：服务器选择超时
```

### 配置选项

- **uri**: MongoDB 连接字符串。可以包含主机、端口、数据库和其他选项
- **database**: 要使用的默认数据库名称
- **username/password**: 身份验证凭据
- **authSource**: 要进行身份验证的数据库（默认："admin"）
- **maxPoolSize**: 连接池中的最大连接数
- **minPoolSize**: 连接池中的最小连接数
- **maxConnIdleTime**: 连接可以保持空闲的最长时间
- **connectTimeout**: 建立连接的超时时间
- **socketTimeout**: 套接字操作的超时时间
- **serverSelectionTimeout**: 服务器选择的超时时间

## 使用方法

### 基本用法

```go
package main

import (
    "context"
    "github.com/gone-io/gone/v2"
    "github.com/gone-io/goner/mongo"
    "go.mongodb.org/mongo-driver/bson"
    mongoDriver "go.mongodb.org/mongo-driver/mongo"
)

type UserService struct {
    gone.Flag
    mongoClient *mongoDriver.Client `gone:"*"`
    logger      gone.Logger          `gone:"*"`
}

func (s *UserService) CreateUser(name, email string) error {
    collection := s.mongoClient.Database("myapp").Collection("users")
    
    user := bson.M{
        "name":  name,
        "email": email,
    }
    
    _, err := collection.InsertOne(context.Background(), user)
    if err != nil {
        s.logger.Errorf("创建用户失败: %v", err)
        return err
    }
    
    s.logger.Infof("用户创建成功: %s", name)
    return nil
}

func (s *UserService) GetUser(email string) (bson.M, error) {
    collection := s.mongoClient.Database("myapp").Collection("users")
    
    var user bson.M
    err := collection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
    if err != nil {
        s.logger.Errorf("获取用户失败: %v", err)
        return nil, err
    }
    
    return user, nil
}

func main() {
    gone.NewApp(
        gone.Load(mongo.Load),
        gone.Load(&UserService{}),
    ).Run()
}
```

### 多数据库连接

您可以配置多个 MongoDB 连接：

```yaml
mongo:
  uri: "mongodb://localhost:27017"
  database: "main"

mongo-analytics:
  uri: "mongodb://analytics-server:27017"
  database: "analytics"
  username: "analytics_user"
  password: "analytics_pass"
```

```go
type AnalyticsService struct {
    gone.Flag
    mainClient      *mongoDriver.Client `gone:"*"`
    analyticsClient *mongoDriver.Client `gone:"*,mongo-analytics"`
}
```

## 错误处理

组件包含全面的错误处理：

- 连接失败会被正确报告
- 配置错误会详细说明
- 连接池错误会被优雅处理

## 最佳实践

1. **连接池**: 根据应用程序的需求配置适当的池大小
2. **超时**: 设置合理的超时以避免挂起操作
3. **身份验证**: 在生产环境中使用身份验证
4. **数据库选择**: 在配置中指定数据库名称以保持清晰
5. **错误处理**: 始终处理 MongoDB 操作返回的错误

## 依赖项

- [go.mongodb.org/mongo-driver](https://github.com/mongodb/mongo-go-driver) - 官方 MongoDB Go 驱动程序
- [github.com/gone-io/gone/v2](https://github.com/gone-io/gone) - Gone 框架