# Gone Elasticsearch 集成

本包为 Gone 应用程序提供 Elasticsearch 集成功能，支持低级别和类型安全的客户端。

## 功能特点

- 与 Gone 的依赖注入系统轻松集成
- 支持低级别和类型安全的 Elasticsearch 客户端
- 单例客户端实例管理
- 全面的配置选项

## 安装

```bash
go get github.com/gone-io/goner/es
```

## 配置

在项目的配置目录中创建 `default.yaml` 文件，添加以下 Elasticsearch 配置：

```yaml
es:
  addresses: http://localhost:9200   # Elasticsearch 节点地址列表
  username:   # HTTP Basic 认证的用户名
  password:   # HTTP Basic 认证的密码

  cloudID:    # Elastic Service 的端点
  aPIKey:     # Base64 编码的授权令牌
  serviceToken: # 服务令牌授权

  # 其他可选配置
  certificateFingerprint:  # SHA256 十六进制指纹
  retryOnStatus:          # 重试状态码列表（默认：502, 503, 504）
  maxRetries:             # 默认：3
  compressRequestBody:    # 默认：false
  enableMetrics:          # 启用指标收集
  enableDebugLogger:      # 启用调试日志
```

## 使用方法

### 低级别客户端

```go
package main

import (
    "bytes"
    "encoding/json"
    "github.com/elastic/go-elasticsearch/v8"
    "github.com/gone-io/gone/v2"
    "github.com/gone-io/goner/es"
    "github.com/gone-io/goner/viper"
    "io"
)

type esUser struct {
    gone.Flag
    esClient *elasticsearch.Client `gone:"*"`
    logger   gone.Logger          `gone:"*"`
}

func (s *esUser) Use() {
    // 创建索引
    create, err := s.esClient.Indices.Create("my_index")
    if err != nil {
        s.logger.Errorf("Indices.Create err:%v", err)
        return
    }

    // 创建文档
    document := struct {
        Name string `json:"name"`
    }{
        "go-elasticsearch",
    }
    data, _ := json.Marshal(document)
    index, err := s.esClient.Index("my_index", bytes.NewReader(data))
    if err != nil {
        s.logger.Errorf("Index err:%v", err)
        return
    }

    // 获取文档 ID
    var id struct {
        ID string `json:"_id"`
    }
    all, _ := io.ReadAll(index.Body)
    json.Unmarshal(all, &id)

    // 获取文档
    get, err := s.esClient.Get("my_index", id.ID)
    if err != nil {
        s.logger.Errorf("Get err:%v", err)
    }
}

func main() {
    gone.NewApp(
        viper.Load, // 加载配置
        es.Load,    // 初始化 Elasticsearch 客户端
    ).Run(func(esUser *esUser) {
        esUser.Use()
    })
}
```

### 类型安全客户端

```go
package main

import (
    "context"
    "github.com/elastic/go-elasticsearch/v8"
    "github.com/gone-io/gone/v2"
    "github.com/gone-io/goner/es"
    "github.com/gone-io/goner/viper"
)

type esUser struct {
    gone.Flag
    esClient *elasticsearch.TypedClient `gone:"*"`
    logger   gone.Logger               `gone:"*"`
}

func (s *esUser) Use() {
    ctx := context.TODO()

    // 创建索引
    create, err := s.esClient.Indices.Create("my_index").Do(ctx)
    if err != nil {
        s.logger.Errorf("Indices.Create err:%v", err)
        return
    }

    // 创建文档
    document := struct {
        Name string `json:"name"`
    }{
        "go-elasticsearch",
    }
    index, err := s.esClient.Index("my_index").Document(document).Do(ctx)
    if err != nil {
        s.logger.Errorf("Index err:%v", err)
        return
    }

    // 获取文档
    get, err := s.esClient.Get("my_index", index.Id_).Do(ctx)
    if err != nil {
        s.logger.Errorf("Get err:%v", err)
    }
}

func main() {
    gone.NewApp(
        viper.Load, // 加载配置
        es.Load,    // 初始化 Elasticsearch 客户端
    ).Run(func(esUser *esUser) {
        esUser.Use()
    })
}
```

## API 参考

### Load

```go
func Load(loader gone.Loader) error
```

向 Gone 加载器注册一个单例的低级别 Elasticsearch 客户端提供者。

### LoadTypedClient

```go
func LoadTypedClient(loader gone.Loader) error
```

向 Gone 加载器注册一个单例的类型安全 Elasticsearch 客户端提供者。

## 许可证

本项目基于 MIT 许可证 - 详见 LICENSE 文件。