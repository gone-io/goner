<p>
    <a href="README.md">English</a>&nbsp ｜&nbsp 中文
</p>

# goner/deepseek 组件

**goner/deepseek** 组件为 Gone 应用程序提供 Deepseek 集成功能，提供简单易用的客户端配置和管理。

## 功能特点

- 与 Gone 的依赖注入系统轻松集成
- 支持多个 Deepseek 客户端实例
- 单例客户端实例管理
- 全面的配置选项

## 安装

```bash
go get github.com/gone-io/goner/deepseek
```

## 配置

在项目的配置目录中创建 `default.yaml` 文件，添加以下 Deepseek 配置：

```yaml
deepseek:
  authToken: "your-api-key"      # 你的 Deepseek API 密钥
  baseURL: ""                    # 可选：自定义 API 基础 URL
  timeout: 0                     # 可选：请求超时时间（秒）
  path: ""                      # 可选：API 路径（默认为 "chat/completions"）
  proxyUrl: ""                   # 可选：HTTP 代理 URL
```

### 代理配置（proxyUrl）

`proxyUrl` 参数用于在网络受限环境下通过HTTP代理访问Deepseek API。这在以下场景特别有用：

- 当你所在的网络环境无法直接访问Deepseek API
- 需要通过企业代理或VPN访问外部API
- 为了提高API访问速度和稳定性

配置示例：

```yaml
deepseek:
  authToken: "your-api-key"
  proxyUrl: "http://proxy.example.com:8080"  # HTTP代理服务器地址和端口
```

你也可以使用带认证的代理：

```yaml
deepseek:
  authToken: "your-api-key"
  proxyUrl: "http://username:password@proxy.example.com:8080"
```

对于多客户端配置，你可以使用不同的前缀：

```yaml
deepseek:
  authToken: "default-api-key"

baidu:
  authToken: "baidu-api-key"
  baseURL: "https://custom.baidu.api.com"

aliyun:
  authToken: "aliyun-api-key"
  baseURL: "https://custom.aliyun.api.com"
```

## 使用方法

### 基本用法

```go
package main

import (
    "context"
    "github.com/gone-io/gone/v2"
    goneDeepseek "github.com/gone-io/goner/deepseek"
    "github.com/cohesion-org/deepseek-go"
)

type aiUser struct {
    gone.Flag
    client *deepseek.Client `gone:"*"` // 默认客户端
}

func (s *aiUser) Use() {
    resp, err := s.client.CreateChatCompletion(
        context.Background(),
        deepseek.ChatCompletionRequest{
            Model: "deepseek-chat",
            Messages: []deepseek.ChatCompletionMessage{
                {
                    Role:    "user",
                    Content: "你好！",
                },
            },
        },
    )
    if err != nil {
        // 处理错误
        return
    }
    // 处理响应
}

// 多客户端示例
type multiAIUser struct {
    gone.Flag
    defaultClient *deepseek.Client `gone:"*"`           // 默认客户端
    baiduClient   *deepseek.Client `gone:"*,baidu"`    // 百度客户端
    aliyunClient  *deepseek.Client `gone:"*,aliyun"`   // 阿里云客户端
}

func main() {
    gone.NewApp(goneDeepseek.Load).Run(func(u *aiUser) {
        u.Use()
    })
}
```

## 高级用法

客户端支持 `github.com/cohesion-org/deepseek-go` 包提供的所有 Deepseek API 功能，包括：

- 聊天补全
- 文本嵌入
- 模型查询

详细的 API 使用方法请参考 [deepseek-go 文档](https://github.com/cohesion-org/deepseek-go)。