# Gone OpenAI 示例

本示例展示了如何在 Gone 框架中使用 OpenAI 客户端。包含了单客户端和多客户端的使用场景，以及基本的聊天补全功能实现。

## 功能特点

- 演示单个 OpenAI 客户端的基本使用
- 演示多个客户端（OpenAI、百度、阿里云）的配置和使用
- 包含错误处理示例
- 使用 Gone 框架的依赖注入特性

## 配置说明

示例使用 `config/default.yaml` 配置文件来设置 OpenAI 客户端参数：

```yaml
openai:
  apiToken: "your-api-key"       # 你的 OpenAI API 密钥
  baseUrl: ""                    # 可选：自定义 API 基础 URL
  orgID: ""                      # 可选：组织 ID
  APIType: "openai"              # API 类型：openai, azure, anthropic
  APIVersion: ""                 # API 版本（Azure 必需）
  assistantVersion: ""           # 可选：助手 API 版本
  proxyUrl: ""                   # 可选：HTTP 代理 URL

baidu:
  apiToken: "baidu-api-key"
  baseUrl: "https://baidu-api-endpoint"

aliyun:
  apiToken: "aliyun-api-key"
  baseUrl: "https://aliyun-api-endpoint"
```

## 运行示例

1. 确保已安装 Go 1.20 或更高版本
2. 克隆仓库并进入示例目录
3. 在 `config/default.yaml` 中配置你的 API 密钥
4. 运行示例程序：

```bash
go run main.go
```

## 代码说明

- `main.go`: 包含了单客户端和多客户端使用的示例代码
- `config/default.yaml`: 配置文件，包含了客户端配置信息

### 单客户端示例

```go
type singleAIUser struct {
    gone.Flag
    client *openai.Client `gone:"*"` // 默认客户端
}
```

### 多客户端示例

```go
type multiAIUser struct {
    gone.Flag
    defaultClient *openai.Client `gone:"*"`           // 默认客户端
    baiduClient   *openai.Client `gone:"*,baidu"`    // 百度客户端
    aliyunClient  *openai.Client `gone:"*,aliyun"`   // 阿里云客户端
}
```

## 注意事项

- 运行示例前请确保已正确配置 API 密钥
- 多客户端示例中的百度和阿里云客户端仅作演示用途
- 建议参考 [Gone OpenAI 集成文档](https://github.com/gone-io/goner/tree/main/openai) 了解更多用法