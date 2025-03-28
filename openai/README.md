# Gone OpenAI Integration

This package provides OpenAI integration for Gone applications, offering easy-to-use client configuration and management.

## Features

- Easy integration with Gone's dependency injection system
- Support for multiple OpenAI client instances
- Singleton client instance management
- Comprehensive configuration options

## Installation

```bash
go get github.com/gone-io/goner/openai
```

## Configuration

Create a `default.yaml` file in your project's config directory with the following OpenAI configuration:

```yaml
openai:
  apiToken: "your-api-key"       # Your OpenAI API key
  baseUrl: ""                    # Optional: Custom API base URL
  orgID: ""                      # Optional: Organization ID
  APIType: "openai"              # API type: openai, azure, anthropic
  APIVersion: ""                 # API version (required for Azure)
  assistantVersion: ""           # Optional: Assistant API version
  proxyUrl: ""                   # Optional: HTTP proxy URL
```

### Proxy Configuration (proxyUrl)

The `proxyUrl` parameter is used to access the OpenAI API through an HTTP proxy in network-restricted environments. This is particularly useful in the following scenarios:

- When your network environment cannot directly access the OpenAI API
- When you need to access external APIs through corporate proxies or VPNs
- To improve API access speed and stability

Configuration example:

```yaml
openai:
  apiToken: "your-api-key"
  proxyUrl: "http://proxy.example.com:8080"  # HTTP proxy server address and port
```

You can also use a proxy with authentication:

```yaml
openai:
  apiToken: "your-api-key"
  proxyUrl: "http://username:password@proxy.example.com:8080"
```

For multiple client configurations, you can use different prefixes:

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

## Usage

### Basic Usage

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
    client *openai.Client `gone:"*"` // Default client
}

func (s *aiUser) Use() {
    resp, err := s.client.CreateChatCompletion(
        context.Background(),
        openai.ChatCompletionRequest{
            Model: openai.GPT3Dot5Turbo,
            Messages: []openai.ChatCompletionMessage{
                {
                    Role:    openai.ChatMessageRoleUser,
                    Content: "Hello!",
                },
            },
        },
    )
    if err != nil {
        // Handle error
        return
    }
    // Process response
}

// Multiple clients example
type multiAIUser struct {
    gone.Flag
    defaultClient *openai.Client `gone:"*"`           // Default client
    baiduClient   *openai.Client `gone:"*,baidu"`    // Baidu client
    aliyunClient  *openai.Client `gone:"*,aliyun"`   // Aliyun client
}

func main() {
    gone.NewApp(goneOpenAi.Load).Run(func(u *aiUser) {
        u.Use()
    })
}
```

## Advanced Usage

The client supports all OpenAI API features provided by the `github.com/sashabaranov/go-openai` package, including:

- Chat Completions
- Completions
- Edits
- Images
- Embeddings
- Audio
- Files
- Fine-tunes
- Moderations

Refer to the [go-openai documentation](https://github.com/sashabaranov/go-openai) for detailed API usage.