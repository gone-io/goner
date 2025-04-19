# Gone MCP 组件使用指南

> Gone MCP是基于MCP协议封装的Go组件，支持快速构建AI模型与业务系统的集成。组件提供服务端和客户端实现，支持两种定义方式(Goner Define和Functional Define)来定义Tool、Prompt和Resource。服务端支持Stdio和SSE两种通信通道，可配置Hooks和ContextFunc，并支持多实例。客户端支持多实例配置、Stdio/SSE传输及自定义Transport。通过简洁的API设计，开发者可以轻松实现AI模型与外部系统的交互。
>

- [Gone MCP 组件使用指南](#gone-mcp-组件使用指南)
	- [简介](#简介)
	- [快速开始](#快速开始)
	- [Gone MCP 组件 特性](#gone-mcp-组件-特性)
	- [服务端使用](#服务端使用)
		- [基本概念](#基本概念)
		- [Gone MCP 组件 支持 两种方式定义 **Tool** 、 **Prompt** 和 **Resource**](#gone-mcp-组件-支持-两种方式定义-tool--prompt-和-resource)
			- [Goner Define，基于Goner实现接口 `goneMcp.ITool`、`goneMcp.IPrompt`、`goneMcp.IResource`](#goner-define基于goner实现接口-gonemcpitoolgonemcpipromptgonemcpiresource)
			- [Functional Define，基于`github.com/mark3labs/mcp-go`中的函数式的定义](#functional-define基于githubcommark3labsmcp-go中的函数式的定义)
		- [服务端配置](#服务端配置)
		- [支持Stdio和SSE两种通道](#支持stdio和sse两种通道)
			- [Stdio](#stdio)
			- [SSE](#sse)
		- [支持加载`*server.Hooks`到gone框架，MCP 组件会使用相关的Hook实现功能拦截](#支持加载serverhooks到gone框架mcp-组件会使用相关的hook实现功能拦截)
		- [支持加载 `server.StdioContextFunc`(使用stdio通道时) 或 `server.SseContextFunc`（使用sse通道时），MCP 组件会使用相关的ContextFunc设置Context](#支持加载-serverstdiocontextfunc使用stdio通道时-或-serverssecontextfunc使用sse通道时mcp-组件会使用相关的contextfunc设置context)
		- [多实例支持](#多实例支持)
	- [客户端使用](#客户端使用)
		- [注入客户端](#注入客户端)
		- [配置](#配置)
			- [Stdio客户端配置](#stdio客户端配置)
			- [SSE客户端配置](#sse客户端配置)
		- [使用客户端](#使用客户端)
		- [示例](#示例)
		- [高级用法](#高级用法)
			- [自定义Transport](#自定义transport)


## 简介

MCP（Model Context Protocol，模型上下文协议）是一个由 Anthropic 公司（Claude 模型的开发者）主导的开放协议。它的主要目标是为 AI 模型提供一个标准化的框架，使其能够方便地与外部数据源、工具和服务进行交互。

Gone MCP 组件是基于 `github.com/mark3labs/mcp-go` 进行封装的工具包，它能帮助开发者快速构建 MCP 的服务端和客户端应用。通过使用 Gone MCP 组件，您可以轻松地将 AI 模型与您的业务系统进行集成。

## 快速开始

完整代码参考：[mcp/quick_start](../examples/mcp/quick_start)

<details>
<summary>1.创建项目，安装依赖</summary>

```bash
go mod init quickstart
go get github.com/gone-io/goner/mcp
```

</details>



<details>
<summary>2.基于stdio创建服务端，创建文件 `server/main.go` </summary>

```go
// 这是一个基于MCP(Model-Controller-Provider)的服务器端示例程序
// 展示了如何创建一个简单的MCP工具，该工具可以接收一个名字参数并返回问候语

package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/gone-io/gone/v2"
	goneMcp "github.com/gone-io/goner/mcp"
	"github.com/mark3labs/mcp-go/mcp"
)

// main 函数初始化并启动MCP服务器
// 加载必要的组件并注册helloTool工具
func main() {
	_ = os.Setenv("GONE_MCP", `{"name":"quick-start", "version": "0.0.1"}`)
	gone.
		Loads(goneMcp.ServerLoad). // 加载MCP服务器组件
		Load(&helloTool{}). // 加载helloTool工具
		Serve() // 启动服务
}

// helloTool 实现了一个简单的问候工具
// 继承gone.Flag以支持MCP工具的基本功能
type helloTool struct {
	gone.Flag
}

// Define 定义了工具的名称、描述和参数
// 返回值: goneMcp.Tool - 工具的定义信息
func (h helloTool) Define() goneMcp.Tool {
	return mcp.NewTool("hello_world",
		mcp.WithDescription("Say hello to someone"), // 工具描述
		mcp.WithString("name", // 定义字符串类型的name参数
			mcp.Required(),                                 // 将name参数设置为必填
			mcp.Description("Name of the person to greet"), // name参数的描述
		),
	)
}

// Handler 处理工具的调用请求
// 参数:
//   - ctx: 上下文信息
//   - request: 调用请求，包含参数信息
//
// 返回值:
//   - *mcp.CallToolResult: 调用结果
//   - error: 错误信息
func (h helloTool) Handler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// 获取并验证name参数
	name, ok := request.Params.Arguments["name"].(string)
	if !ok {
		return nil, errors.New("name must be a string") // 参数类型错误时返回错误
	}
	// 返回格式化的问候语
	return mcp.NewToolResultText(fmt.Sprintf("Hello, %s!", name)), nil
}

// 编译时类型检查，确保helloTool实现了goneMcp.ITool接口
var _ goneMcp.ITool = (*helloTool)(nil)

```

</details>


<details>
<summary>3.基于stdio创建客户端，创建文件 `client/main.go` </summary>

```go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gone-io/gone/v2"
	goneMcp "github.com/gone-io/goner/mcp"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
)

// main 函数是程序的入口点
// 这是一个基于MCP（Message Communication Protocol）的客户端示例程序
// 演示了如何初始化MCP客户端、连接服务器、列出可用工具并调用工具
func main() {
	gone.
		NewApp(goneMcp.ClientLoad). // 创建一个新的Gone应用，使用MCP客户端加载器
		Run(func(in struct {
			// 注入MCP客户端，使用stdio类型通信，并通过本地服务器运行，运行命令为：go run ./server
			client *client.Client `gone:"*,type=stdio,param=go run ./server"`
		}) {
			c := in.client
			// 创建一个带30秒超时的上下文
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			// 启动客户端
			err := c.Start(ctx)
			if err != nil {
				log.Fatalf("Failed to start client: %v", err)
			}

			// 准备初始化请求
			initRequest := mcp.InitializeRequest{}
			// 设置协议版本和客户端信息
			initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
			initRequest.Params.ClientInfo = mcp.Implementation{
				Name:    "example-client",
				Version: "1.0.0",
			}

			// 发送初始化请求到服务器
			initResult, err := c.Initialize(ctx, initRequest)
			if err != nil {
				log.Fatalf("Failed to initialize: %v", err)
			}
			// 打印服务器信息
			fmt.Printf(
				"Initialized with server: %s %s\n\n",
				initResult.ServerInfo.Name,
				initResult.ServerInfo.Version,
			)

			// 列出服务器提供的所有可用工具
			fmt.Println("Listing available tools...")
			toolsRequest := mcp.ListToolsRequest{}
			tools, err := c.ListTools(ctx, toolsRequest)
			if err != nil {
				log.Fatalf("Failed to list tools: %v", err)
			}
			// 遍历并打印每个工具的名称和描述
			for _, tool := range tools.Tools {
				fmt.Printf("- %s: %s\n", tool.Name, tool.Description)
			}
			fmt.Println()

			// 调用hello_world工具示例
			fmt.Println("Calling `hello_world`")
			request := mcp.CallToolRequest{}
			// 设置要调用的工具名称和参数
			request.Params.Name = "hello_world"
			request.Params.Arguments = map[string]any{
				"name": "John",
			}
			// 执行工具调用
			tool, err := c.CallTool(ctx, request)
			if err != nil {
				log.Fatalf("Failed to call tool: %v", err)
			}
			// 打印工具执行结果
			printToolResult(tool)
		})
}

// printToolResult 辅助函数，用于打印工具执行的结果
// 支持文本内容和其他类型内容的格式化输出
func printToolResult(result *mcp.CallToolResult) {
	for _, content := range result.Content {
		// 如果是文本内容，直接打印
		if textContent, ok := content.(mcp.TextContent); ok {
			fmt.Println(textContent.Text)
		} else {
			// 其他类型内容，转换为JSON格式打印
			jsonBytes, _ := json.MarshalIndent(content, "", "  ")
			fmt.Println(string(jsonBytes))
		}
	}
}

```

</details>

<details>
<summary>4.运行：`go run client/main.go`</summary>

```bash
go mod tidy
go run client/main.go
```

运行结果：

```
Initialized with server: quick-start 0.0.1

Listing available tools...
- hello_world: Say hello to someone

Calling `hello_world`
Hello, John!
```

</details>

## Gone MCP 组件 特性

- 服务端
    - 支持 `Goner Define` 和 `Functional Define` 定义 MCP 的 Tool、Prompt、Resource
    - 支持配置文件
    - 支持Stdio和SSE两种通道
    - 支持定义Hooks注入
    - 支持定义ContextFunc注入
    - 支持多实例
- 客户端
    - 支持多实例，根据不同的`gone`标签配置获取不同的实例
    - 支持从配置中读取参数
    - 支持Stdio和SSE
    - SSE支持设置header
    - 支持自定义`transport.Interface`

## 服务端使用

### 基本概念

MCP (Model Control Protocol) 是一个用于控制AI模型的协议，Gone MCP组件提供了对该协议的支持。服务端主要包含以下几个核心概念：

- **Tool**: 定义可供AI模型调用的工具
- **Prompt**: 定义提示词模板
- **Resource**: 定义可供AI模型访问的资源

### Gone MCP 组件 支持 两种方式定义 **Tool** 、 **Prompt** 和 **Resource**

#### Goner Define，基于Goner实现接口 `goneMcp.ITool`、`goneMcp.IPrompt`、`goneMcp.IResource`

**例子，代码所在位置: [examples/mcp/stdio/server/goner_define](../examples/mcp/stdio/server/goner_define)：**

<details>
<summary>1. 定义Tool</summary>

```go
package goner_define

import (
	"context"
	"errors"
	"fmt"
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/mcp"
	"github.com/mark3labs/mcp-go/mcp"
)

type helloTool struct {
	gone.Flag
}

func (h helloTool) Define() goneMcp.Tool {
	return mcp.NewTool("hello_world",
		mcp.WithDescription("Say hello to someone"),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("Name of the person to greet"),
		),
	)
}

func (h helloTool) Handler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, ok := request.Params.Arguments["name"].(string)
	if !ok {
		return nil, errors.New("name must be a string")
	}
	return mcp.NewToolResultText(fmt.Sprintf("Hello, %s!", name)), nil
}

var _ goneMcp.ITool = (*helloTool)(nil)

```

</details>

<details>
<summary>2. 定义Prompt</summary>

```go
package goner_define

import (
	"context"
	"fmt"
	"github.com/gone-io/gone/v2"
	goneMcp "github.com/gone-io/goner/mcp"
	"github.com/mark3labs/mcp-go/mcp"
)

type greeting struct {
	gone.Flag
}

func (g greeting) Define() goneMcp.Prompt {
	return mcp.NewPrompt("greeting",
		mcp.WithPromptDescription("A friendly greeting prompt"),
		mcp.WithArgument("name",
			mcp.ArgumentDescription("Name of the person to greet"),
		),
	)
}

func (g greeting) Handler(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	name := request.Params.Arguments["name"]
	if name == "" {
		name = "friend"
	}

	return mcp.NewGetPromptResult(
		"A friendly greeting",
		[]mcp.PromptMessage{
			mcp.NewPromptMessage(
				mcp.RoleAssistant,
				mcp.NewTextContent(fmt.Sprintf("Hello, %s! How can I help you today?", name)),
			),
		},
	), nil
}

var _ goneMcp.IPrompt = (*greeting)(nil)

```

</details>

<details>
<summary>3. 定义Resource</summary>

```go
package goner_define

import (
	"context"
	"github.com/gone-io/gone/v2"
	goneMcp "github.com/gone-io/goner/mcp"
	"github.com/mark3labs/mcp-go/mcp"
)

type readmeResource struct {
	gone.Flag
}

func (r readmeResource) Define() goneMcp.Resource {
	return mcp.NewResource(
		"docs://readme",
		"Project README",
		mcp.WithResourceDescription("The project's README file"),
		mcp.WithMIMEType("text/markdown"),
	)
}

func (r readmeResource) Handler(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      "docs://readme",
			MIMEType: "text/markdown",
			Text:     "readme content",
		},
	}, nil
}

var _ goneMcp.IResource = (*readmeResource)(nil)
```

</details>

#### Functional Define，基于`github.com/mark3labs/mcp-go`中的函数式的定义

**例子，代码所在位置: [examples/mcp/stdio/server/functional_define](../examples/mcp/stdio/server/functional_add/mcp_server_define.go)**

### 服务端配置

基于Gone强大的配置能力，可以让服务快速支持本地文件配置和远程配置，更多关于配置的内容参考：[viper本地配置](../viper/README_CN.md)、[viper远程配置](../viper/remote/README_CN.md)
。下面讲解对应Key的意义

| key                                      | 说明                         | 对应`mcp-go`选项                          | 默认值      |
| ---------------------------------------- | ---------------------------- | ----------------------------------------- | ----------- |
| mcp.name                                 | 服务器名称                   | 无                                        | `MCP服务器` |
| mcp.version                              | 服务器版本                   | 无                                        | `1.0.0`     |
| mcp.withRecovery                         | 是否启用恢复机制             | `server.WithRecovery`                     | `false`     |
| mcp.withPromptCapabilities               | 是否启用提示词能力           | `server.WithPromptCapabilities`           | `false`     |
| mcp.withToolCapabilities                 | 是否启用工具能力             | `server.WithToolCapabilities`             | `false`     |
| mcp.withLogging                          | 是否启用日志                 | `server.WithLogging`                      | `false`     |
| mcp.withInstructions                     | 自定义指令内容               | `server.WithInstructions`                 | `""`        |
| mcp.transportType                        | 传输类型，支持`stdio`和`sse` | `server.WithTransportType`                | `stdio`     |
| mcp.sse.withBaseURL                      | SSE传输类型的BaseURL         | `server.WithSSEBaseURL`                   | `""`        |
| mcp.sse.withBasePath                     | SSE传输类型的基础路径        | `server.WithBasePath`                     | `""`        |
| mcp.sse.withMessageEndpoint              | SSE消息端点                  | `server.WithMessageEndpoint`              | `""`        |
| mcp.sse.withUseFullURLForMessageEndpoint | 是否使用完整URL作为消息端点  | `server.WithUseFullURLForMessageEndpoint` | `false`     |
| mcp.sse.withSSEEndpoint                  | SSE端点                      | `server.WithSSEEndpoint`                  | `""`        |
| mcp.sse.withKeepAlive                    | 是否启用保持连接             | `server.WithKeepAlive`                    | `false`     |
| mcp.sse.withKeepAliveInterval            | 保持连接的时间间隔           | `server.WithKeepAliveInterval`            | `0`         |
| mcp.sse.address                          | SSE服务器监听地址            | 无                                        | `:8080`     |

### 支持Stdio和SSE两种通道
#### Stdio
将`mcp.transportType`设置为`stdio`，MCP服务器将使用Stdio作为通信通道。

#### SSE
将`mcp.transportType`设置为`sse`，MCP服务器将使用SSE作为通信通道。在SSE通道模式，可以配置服务相关参数:`mcp.sse.*`，具体含义参考 [服务端配置](#服务端配置)。


### 支持加载`*server.Hooks`到gone框架，MCP 组件会使用相关的Hook实现功能拦截

MCP服务器支持通过自定义Hooks来实现功能拦截，可以在工具调用前后执行自定义逻辑。以下是使用示例：

```go
// 在应用中注册Hooks
func main() {
    gone.
		NewApp(goneMcp.Load).
        Load(g.NamedThirdComponentLoadFunc("mcp.hooks", &server.Hooks{
			OnBeforeCallTool: []server.OnBeforeCallToolFunc{
				func(ctx context.Context, id any, message *mcp.CallToolRequest) {
					//todo 处理逻辑
				},
			},
		})).
        Run()
}
```

说明：
1. `g.NamedThirdComponentLoadFunc`的作用是将`*server.Hooks`转换为Gone框架能够识别的加载函数；
2. 这里使用的名字为`mcp.hooks`；如果是多实例的情况，前缀`mcp`应该替换为具体实例的前缀；
3. Hook函数参考 [`github.com/mark3labs/mcp-go/mcp`的文档](http://github.com/mark3labs/mcp-go/blob/main/mcp)

### 支持加载 `server.StdioContextFunc`(使用stdio通道时) 或 `server.SseContextFunc`（使用sse通道时），MCP 组件会使用相关的ContextFunc设置Context

下面以使用stdio通道为例，其他通道的使用方式类似。

```go
func main() {
    gone.
		NewApp(goneMcp.Load).
        Load(
			// g.NamedThirdComponentLoadFunc 将 server.StdioContextFunc 转换为 Gone框架能够识别的加载函数
			g.NamedThirdComponentLoadFunc(
				"mcp.stdio.context", // 前缀为`mcp`，多实例前缀依赖具体注入的定义
				// 需要将定义函数转换为 server.StdioContextFunc， MCP组件才能识别
				server.StdioContextFunc(
					func(ctx context.Context) context.Context {
						return context.WithValue(ctx, "test", "test")
					},
				),
			)
		).
        Run()
}
```

说明：
1. server.StdioContextFunc 加载的名字为 `mcp.stdio.context`；
2. server.SseContextFunc 加载的名字为 `mcp.sse.context`；
3. 多实例的情况，前缀`mcp`应该替换为具体实例的前缀；

### 多实例支持
对于 **Functional Define**，可以使用`gone`标签来指定实例名称，例如：

```go
package service
import (
	"github.com/gone-io/gone/v2"
	"github.com/mark3labs/mcp-go/server"
)

type MyService struct {
	gone.Flag
	s1 *server.MCPServer `gone:"*,customer1"` // 实例名称为customer1"`, 注意前面的`*,`，星号可以省略，逗号不能省略
	s2 *server.MCPServer `gone:",customer2"`
}
```
说明：
 1. 对于上面定义的`MyService`，`s1`和`s2`都是`*server.MCPServer`类型的实例，分别对应`customer1`和`customer2`实例。
 2. gone框架将启动两个MCP服务器，分别对应`customer1`和`customer2`实例。
 3. 他们的配置读取，配置的前缀应该使用对应实例的名字，例如服务名字分别是`customer1.name`和`customer2.name`
 4. `*server.Hooks`、`server.StdioContextFunc`、`server.SseContextFunc`，也支持多实例，注入的名字前缀需要将`mcp`替换为具体实例的名称。



## 客户端使用

### 注入客户端

```go
package service
import (
	"github.com/gone-io/gone/v2"
	"github.com/mark3labs/mcp-go/client"
)
type MyService struct {
gone.Flag
	// 使用Stdio客户端
	StdioClient *client.Client `gone:"*,type=stdio,param=./mcp-server --param=value"`

	// 使用配置文件中的定义的Stdio客户端，配置的key为`mcp.stdio`【使用配置】
	StdioClientFromConfig *client.Client `gone:"*,type=stdio,configKey=mcp.stdio"`

	// 使用SSE客户端
	SseClient *client.Client `gone:"*,type=sse,param=http://localhost:8080/sse"`

	// 使用配置文件中定义的SSE客户端,配置的key为`mcp.sse`【使用配置】
	SseClientFromConfig *client.Client `gone:"*,type=sse,configKey=mcp.sse"`

	// 使用自定义Transport客户端
	CustomTransportClient *client.Client `gone:"*,type=transport,param=myTransport"`
}
```

### 配置
#### Stdio客户端配置
| key                  | 说明     | 默认值 |
| -------------------- | -------- | ------ |
| ${configKey}.command | 命令     | 无     |
| ${configKey}.env     | 环境变量 | 无     | 无 |
| ${configKey}.args    | 命令参数 | 无     | 无 |

**例子：**

```yaml
# Stdio客户端配置
mcp.stdio:
  command: "./mcp-server"  # 命令路径
  env: # 环境变量
    - "ENV_VAR=value"
  args: # 命令参数
    - "--param=value"
```

#### SSE客户端配置
| key                  | 说明         | 默认值 |
| -------------------- | ------------ | ------ |
| ${configKey}.baseUrl | SSE服务器URL | 无     |
| ${configKey}.header  | 请求头       | 无     |
**例子：**

```yaml
# SSE客户端配置
mcp.sse:
  baseUrl: "http://localhost:8080/sse"  # SSE服务器URL
  header: # 请求头
    Authorization: "Bearer token"
```

### 使用客户端
具体如何使用客户端，请参考 [github.com/mark3labs/mcp-go/client](https://github.com/mark3labs/mcp-go)

### 示例

完整示例可以参考 [examples/mcp](https://github.com/gone-io/goner/tree/main/examples/mcp) 目录。

### 高级用法
#### 自定义Transport
1. 定义Transport
```go
type xTransport struct {
	gone.Flag
}

func (x xTransport) Start(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (x xTransport) SendRequest(ctx context.Context, request transport.JSONRPCRequest) (*transport.JSONRPCResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (x xTransport) SendNotification(ctx context.Context, notification mcp.JSONRPCNotification) error {
	//TODO implement me
	panic("implement me")
}

func (x xTransport) SetNotificationHandler(handler func(notification mcp.JSONRPCNotification)) {
	//TODO implement me
	panic("implement me")
}

func (x xTransport) Close() error {
	//TODO implement me
	panic("implement me")
}

var _ transport.Interface = (*xTransport)(nil)
```
2. 加载自定义Transport

```go
func main() {
	gone.
		NewApp(goneMcp.Load).
		Load(
			g.NamedThirdComponentLoadFunc("myTransport", &xTransport{}),
		).
		Run()
}
```
3. 在客户端中使用自定义Transport
```go
package service
import (
	"github.com/gone-io/gone/v2"
	"github.com/mark3labs/mcp-go/client"
)
type MyService struct {
	gone.Flag
	// 使用自定义Transport客户端；注意：参数为自定义Transport的名称，即在Load函数中定义的名称
	CustomTransportClient *client.Client `gone:"*,type=transport,param=myTransport"`
}
```
