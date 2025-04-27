[//]: # (desc: OpenAI model client integration)

<p>
    English&nbsp ｜&nbsp <a href="README_CN.md">中文</a>
</p>

# Gone OpenAI Example

This example demonstrates how to use the OpenAI client in the Gone framework. It includes scenarios for both single-client and multi-client usage, as well as basic chat completion functionality implementation.

## Features

- Demonstrates basic usage of a single OpenAI client
- Shows configuration and usage of multiple clients (OpenAI, Baidu, Aliyun)
- Includes error handling examples
- Utilizes Gone framework's dependency injection features

## Configuration

The example uses `config/default.yaml` configuration file to set OpenAI client parameters:

```yaml
openai:
  apiToken: "your-api-key"       # Your OpenAI API key
  baseUrl: ""                    # Optional: Custom API base URL
  orgID: ""                      # Optional: Organization ID
  APIType: "openai"              # API type: openai, azure, anthropic
  APIVersion: ""                 # API version (Required for Azure)
  assistantVersion: ""           # Optional: Assistant API version
  proxyUrl: ""                   # Optional: HTTP proxy URL

baidu:
  apiToken: "baidu-api-key"
  baseUrl: "https://baidu-api-endpoint"

aliyun:
  apiToken: "aliyun-api-key"
  baseUrl: "https://aliyun-api-endpoint"
```

## Running the Example

1. Ensure Go 1.20 or higher is installed
2. Clone the repository and navigate to the example directory
3. Configure your API key in `config/default.yaml`
4. Run the example program:

```bash
go run main.go
```

## Code Overview

- `main.go`: Contains example code for single-client and multi-client usage
- `config/default.yaml`: Configuration file containing client settings

### Single Client Example

```go
type singleAIUser struct {
    gone.Flag
    client *openai.Client `gone:"*"` // Default client
}
```

### Multi-Client Example

```go
type multiAIUser struct {
    gone.Flag
    defaultClient *openai.Client `gone:"*"`           // Default client
    baiduClient   *openai.Client `gone:"*,baidu"`    // Baidu client
    aliyunClient  *openai.Client `gone:"*,aliyun"`   // Aliyun client
}
```

## Notes

- Ensure API keys are properly configured before running the example
- Baidu and Aliyun clients in the multi-client example are for demonstration purposes only
- For more usage details, refer to the [Gone OpenAI Integration Documentation](https://github.com/gone-io/goner/tree/main/openai)