# Gone Deepseek Integration

This package provides Deepseek integration for Gone applications, offering easy-to-use client configuration and management.

## Features

- Easy integration with Gone's dependency injection system
- Support for multiple Deepseek client instances
- Singleton client instance management
- Comprehensive configuration options

## Installation

```bash
go get github.com/gone-io/goner/deepseek
```

## Configuration

Create a `default.yaml` file in your project's config directory with the following Deepseek configuration:

```yaml
deepseek:
  authToken: "your-api-key"      # Your Deepseek API key
  baseURL: ""                    # Optional: Custom API base URL
  timeout: 0                     # Optional: Request timeout in seconds
  path: ""                      # Optional: API path (defaults to "chat/completions")
  proxyUrl: ""                   # Optional: HTTP proxy URL
```

### Proxy Configuration (proxyUrl)

The `proxyUrl` parameter is used to access the Deepseek API through an HTTP proxy in network-restricted environments. This is particularly useful in the following scenarios:

- When your network environment cannot directly access the Deepseek API
- When you need to access external APIs through corporate proxies or VPNs
- To improve API access speed and stability

Configuration example:

```yaml
deepseek:
  authToken: "your-api-key"
  proxyUrl: "http://proxy.example.com:8080"  # HTTP proxy server address and port
```

You can also use a proxy with authentication:

```yaml
deepseek:
  authToken: "your-api-key"
  proxyUrl: "http://username:password@proxy.example.com:8080"
```

For multiple client configurations, you can use different prefixes:

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

## Usage

### Basic Usage

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
    client *deepseek.Client `gone:"*"` // Default client
}

func (s *aiUser) Use() {
    resp, err := s.client.CreateChatCompletion(
        context.Background(),
        deepseek.ChatCompletionRequest{
            Model: "deepseek-chat",
            Messages: []deepseek.ChatCompletionMessage{
                {
                    Role:    "user",
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
    defaultClient *deepseek.Client `gone:"*"`           // Default client
    baiduClient   *deepseek.Client `gone:"*,baidu"`    // Baidu client
    aliyunClient  *deepseek.Client `gone:"*,aliyun"`   // Aliyun client
}

func main() {
    gone.NewApp(goneDeepseek.Load).Run(func(u *aiUser) {
        u.Use()
    })
}
```

## Advanced Usage

The client supports all Deepseek API features provided by the `github.com/cohesion-org/deepseek-go` package, including:

- Chat Completions
- Embeddings
- Models

Refer to the [deepseek-go documentation](https://github.com/cohesion-org/deepseek-go) for detailed API usage.