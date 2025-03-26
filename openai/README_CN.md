# Gone OpenAI 集成

本包为 Gone 应用程序提供 OpenAI 集成功能，提供简单易用的客户端配置和管理。

## 功能特点

- 与 Gone 的依赖注入系统轻松集成
- 支持多个 OpenAI 客户端实例
- 单例客户端实例管理
- 全面的配置选项

## 安装

```bash
go get github.com/gone-io/goner/openai
```

## 配置

在项目的配置目录中创建 `default.yaml` 文件，添加以下 OpenAI 配置：

```yaml
openai:
  apiToken: "your-api-key"       # 你的 OpenAI API 密钥
  baseUrl: ""                    # 可选：自定义 API 基础 URL
  orgID: ""                      # 可选：组织 ID
  APIType: "openai"              # API 类型：openai, azure, anthropic
  APIVersion: ""                 # API 版本（Azure 必需）
  assistantVersion: ""           # 可选：助手 API 版本
  proxyUrl: ""                   # 可选：HTTP 代理 URL
```

对于多客户端配置，你可以使用不同的前缀：

```yaml
openai:
  apiToken: "default-api-key"
  APIType: "openai"

azure:
  apiToken: "azure-api-key"
  baseUrl: "https://your-resource.openai.azure.com"
  APIType: "azure"
  APIVersion: "2023-12-01-preview"

anthropic:
  apiToken: "anthropic-api-key"
  baseUrl: "https://api.anthropic.com"
  APIType: "anthropic"
```

## 使用方法

### 基本用法

```go
package main

import (
    "context"
    "github.com/gone-io/gone/v2"
    goneOpenAi "github.com/gone-io/goner/openai"
    "github.com/sashabaranov/go-openai"
)

type aiUser struct {
    gone.Flag
    client *openai.Client `gone:"*"` // 默认客户端
}

func (s *aiUser) Use() {
    resp, err := s.client.CreateChatCompletion(
        context.Background(),
        openai.ChatCompletionRequest{
            Model: openai.GPT3Dot5Turbo,
            Messages: []openai.ChatCompletionMessage{
                {
                    Role:    openai.ChatMessageRoleUser,
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
    defaultClient *openai.Client `gone:"*"`           // 默认客户端
    baiduClient   *openai.Client `gone:"*,baidu"`    // 百度客户端
    aliyunClient  *openai.Client `gone:"*,aliyun"`   // 阿里云客户端
}

func main() {
    gone.NewApp(goneOpenAi.Load).Run(func(u *aiUser) {
        u.Use()
    })
}
```

## 高级用法

客户端支持 `github.com/sashabaranov/go-openai` 包提供的所有 OpenAI API 功能，包括：

- 聊天补全
- 文本补全
- 文本编辑
- 图像生成
- 文本嵌入
- 音频处理
- 文件操作
- 微调模型
- 内容审核

详细的 API 使用方法请参考 [go-openai 文档](https://github.com/sashabaranov/go-openai)。