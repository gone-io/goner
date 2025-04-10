# 使用 gone.WrapFunctionProvider 快速接入第三方服务—— LLM接入支持 openAI 和 deepseek

- [使用 gone.WrapFunctionProvider 快速接入第三方服务—— LLM接入支持 openAI 和 deepseek](#使用-gonewrapfunctionprovider-快速接入第三方服务-llm接入支持-openai-和-deepseek)
  - [1. gone.WrapFunctionProvider 简介](#1-gonewrapfunctionprovider-简介)
  - [2. 配置注入实现](#2-配置注入实现)
  - [3. 实战示例：OpenAI 和 Deepseek 集成](#3-实战示例openai-和-deepseek-集成)
    - [3.1 Deepseek 集成](#31-deepseek-集成)
    - [3.2 OpenAI 集成](#32-openai-集成)
  - [4. 配置示例](#4-配置示例)
    - [4.1 Deepseek 配置](#41-deepseek-配置)
    - [4.2 OpenAI 配置](#42-openai-配置)
    - [4.3 多客户端配置](#43-多客户端配置)
  - [5. 使用方式](#5-使用方式)
    - [5.1 基本使用](#51-基本使用)
    - [5.2 多客户端使用](#52-多客户端使用)
  - [6. 最佳实践](#6-最佳实践)
    - [6.1 配置分离](#61-配置分离)
    - [6.2 单例管理](#62-单例管理)
    - [6.3 错误处理](#63-错误处理)
    - [6.4 代理配置](#64-代理配置)
  - [7. 总结](#7-总结)


本文将介绍如何使用 gone.WrapFunctionProvider 和配置注入来快速接入大型语言模型(LLM)服务，包括 OpenAI 和 Deepseek。我们将详细说明这种方式的实现原理和最佳实践。

## 1. gone.WrapFunctionProvider 简介

Gone 框架提供了 `gone.WrapFunctionProvider` 这个强大的工具函数，它可以将一个普通的函数包装成 Provider。这种方式特别适合于：

- 需要注入配置的场景
- 需要创建单例的场景
- 需要延迟初始化的场景
- 需要错误处理的场景

对于 LLM 服务集成，这些特性尤其重要，因为我们通常需要管理 API 密钥、自定义基础 URL、设置超时和代理等配置。

## 2. 配置注入实现

在 Gone 框架中，配置注入是通过结构体标签（struct tag）实现的。以 Deepseek 为例：

```go
param struct {
    config deepseek.Config `gone:"config,deepseek"`
}
```

这里的 `gone:"config,deepseek"` 标签表示：
- `config` 表示这是一个配置项
- `deepseek` 是配置的命名空间

## 3. 实战示例：OpenAI 和 Deepseek 集成

### 3.1 Deepseek 集成

让我们看一个完整的示例，展示如何使用 gone.WrapFunctionProvider 来集成 Deepseek：

```go
func Load(loader gone.Loader) error {
    var load = gone.OnceLoad(func(loader gone.Loader) error {
        var single *deepseek.Client

        getSingleDeepseek := func(
            tagConf string,
            param struct {
                config Config `gone:"config,deepseek"`
                doer   g.HTTPDoer `gone:"*,optional"`
            },
        ) (*deepseek.Client, error) {
            var err error
            if single == nil {
                options := param.config.ToDeepseekOptions(param.doer)
                single, err = deepseek.NewClientWithOptions(param.config.AuthToken, options...)
                if err != nil {
                    return nil, gone.ToError(err)
                }
            }
            return single, nil
        }
        provider := gone.WrapFunctionProvider(getSingleDeepseek)
        return loader.Load(provider)
    })
    return load(loader)
}
```

### 3.2 OpenAI 集成

类似地，OpenAI 的集成也可以使用相同的模式：

```go
func Load(loader gone.Loader) error {
    var load = gone.OnceLoad(func(loader gone.Loader) error {
        var single *openai.Client

        getSingleOpenAI := func(
            tagConf string,
            param struct {
                config Config `gone:"config,openai"`
                doer   g.HTTPDoer `gone:"*,optional"`
            },
        ) (*openai.Client, error) {
            var err error
            if single == nil {
                options := param.config.ToOpenAIOptions(param.doer)
                single, err = openai.NewClientWithOptions(param.config.AuthToken, options...)
                if err != nil {
                    return nil, gone.ToError(err)
                }
            }
            return single, nil
        }
        provider := gone.WrapFunctionProvider(getSingleOpenAI)
        return loader.Load(provider)
    })
    return load(loader)
}
```

这些代码实现了以下功能：

1. **单例模式**：通过闭包变量 `single` 确保只创建一个客户端实例
2. **配置注入**：通过结构体标签自动注入 LLM 服务配置
3. **错误处理**：使用 `gone.ToError` 统一错误处理
4. **HTTP 客户端注入**：支持可选的 HTTP 客户端注入，便于自定义网络行为

## 4. 配置示例

### 4.1 Deepseek 配置

在项目的配置目录中创建 `default.yaml` 文件，添加以下 Deepseek 配置：

```yaml
deepseek:
  authToken: "your-api-key"      # 你的 Deepseek API 密钥
  baseURL: ""                    # 可选：自定义 API 基础 URL
  timeout: 0                     # 可选：请求超时时间（秒）
  path: ""                       # 可选：API 路径（默认为 "chat/completions"）
  proxyUrl: ""                   # 可选：HTTP 代理 URL
```

### 4.2 OpenAI 配置

类似地，OpenAI 的配置可以如下：

```yaml
openai:
  authToken: "your-api-key"      # 你的 OpenAI API 密钥
  baseURL: ""                    # 可选：自定义 API 基础 URL
  timeout: 0                     # 可选：请求超时时间（秒）
  proxyUrl: ""                   # 可选：HTTP 代理 URL
```

### 4.3 多客户端配置

对于需要连接多个不同服务提供商的场景，可以使用不同的命名空间：

```yaml
deepseek:
  authToken: "default-deepseek-key"

openai:
  authToken: "default-openai-key"

baidu:
  authToken: "baidu-api-key"
  baseURL: "https://custom.baidu.api.com"

aliyun:
  authToken: "aliyun-api-key"
  baseURL: "https://custom.aliyun.api.com"
```

## 5. 使用方式

### 5.1 基本使用

在应用中使用这些 Provider 非常简单：

```go
type llmUser struct {
    gone.Flag
    deepseekClient *deepseek.Client `gone:"*"`
    openaiClient   *openai.Client   `gone:"*"`
}

func (s *llmUser) Use() {
    // 使用 Deepseek 客户端
    deepseekResp, err := s.deepseekClient.CreateChatCompletion(
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
    
    // 使用 OpenAI 客户端
    openaiResp, err := s.openaiClient.CreateChatCompletion(
        context.Background(),
        openai.ChatCompletionRequest{
            Model: "gpt-3.5-turbo",
            Messages: []openai.ChatCompletionMessage{
                {
                    Role:    "user",
                    Content: "Hello!",
                },
            },
        },
    )
    // ...
}
```

### 5.2 多客户端使用

对于需要使用多个不同配置的客户端的场景：

```go
type multiLLMUser struct {
    gone.Flag
    // Deepseek 客户端
    defaultDeepseek *deepseek.Client `gone:"*"`           
    baiduDeepseek   *deepseek.Client `gone:"*,baidu"`    
    
    // OpenAI 客户端
    defaultOpenAI   *openai.Client   `gone:"*"`           
    aliyunOpenAI    *openai.Client   `gone:"*,aliyun"`   
}

func main() {
    gone.NewApp(
        goneDeepseek.Load,
        goneOpenAi.Load,
    ).Run(func(u *multiLLMUser) {
        // 使用不同的客户端...
    })
}
```

## 6. 最佳实践

### 6.1 配置分离

- 将配置放在独立的配置文件中
- 使用命名空间避免配置冲突
- 敏感信息如 API 密钥应通过环境变量注入

### 6.2 单例管理

- 对于 LLM 客户端，始终使用单例模式以避免资源浪费
- 使用 `gone.OnceLoad` 确保安全的单例初始化

### 6.3 错误处理

- 使用 `gone.ToError` 包装错误
- 在初始化时进行充分的错误检查
- 对 API 调用进行适当的重试和超时处理

### 6.4 代理配置

对于网络受限环境，合理配置代理是必要的：

```yaml
deepseek:
  authToken: "your-api-key"
  proxyUrl: "http://username:password@proxy.example.com:8080"
```

## 7. 总结

使用 `gone.WrapFunctionProvider` 和配置注入是一种优雅且高效的 LLM 服务接入方式。它具有以下优势：

- 代码简洁，易于维护
- 配置灵活，支持动态注入
- 资源管理合理，支持单例模式
- 错误处理统一，便于排查问题
- 支持多种 LLM 服务提供商

通过这种方式，我们可以快速且规范地集成 OpenAI、Deepseek 等各种 LLM 服务，提高开发效率，同时保持代码的可维护性和可扩展性。无论是构建聊天机器人、内容生成工具还是其他 AI 驱动的应用，这种集成方式都能提供稳定可靠的基础设施支持。