<p>
    English&nbsp ｜&nbsp <a href="README_CN.md">中文</a>
</p>

# goner/mcp Component

- [Gone MCP Component Usage Guide](#gone-mcp-component-usage-guide)
  - [Introduction](#introduction)
  - [Quick Start](#quick-start)
  - [Gone MCP Component Features](#gone-mcp-component-features)
  - [Server-side Usage](#server-side-usage)
    - [Basic Concepts](#basic-concepts)
    - [Gone MCP Component Supports Two Ways to Define **Tool**, **Prompt**, and **Resource**](#gone-mcp-component-supports-two-ways-to-define-tool-prompt-and-resource)
      - [Goner Define: Based on Implementing Interfaces `goneMcp.ITool`, `goneMcp.IPrompt`, `goneMcp.IResource`](#goner-define-based-on-implementing-interfaces-gonemcpitool-gonemcpiprompt-gonemcpiresource)
      - [Functional Define: Based on Function-style Definitions in `github.com/mark3labs/mcp-go`](#functional-define-based-on-function-style-definitions-in-githubcommark3labsmcp-go)
    - [Server Configuration](#server-configuration)
    - [Support for Stdio and SSE Channels](#support-for-stdio-and-sse-channels)
      - [Stdio](#stdio)
      - [SSE](#sse)
    - [Support for Loading `*server.Hooks` into Gone Framework](#support-for-loading-serverhooks-into-gone-framework)
    - [Support for Loading `server.StdioContextFunc` or `server.SseContextFunc`](#support-for-loading-serverstdiocontextfunc-or-serverssecontextfunc)
    - [Multi-instance Support](#multi-instance-support)
  - [Client Usage](#client-usage)
    - [Client Injection](#client-injection)
    - [Configuration](#configuration)
      - [Stdio Client Configuration](#stdio-client-configuration)
      - [SSE Client Configuration](#sse-client-configuration)
    - [Using the Client](#using-the-client)
    - [Examples](#examples)
    - [Advanced Usage](#advanced-usage)
      - [Custom Transport](#custom-transport)


## Introduction

MCP (Model Context Protocol) is an open protocol led by Anthropic (the developer of Claude model). Its main goal is to provide a standardized framework for AI models to easily interact with external data sources, tools, and services.

The **goner/mcp** component is a toolkit wrapped around `github.com/mark3labs/mcp-go`, helping developers quickly build MCP server and client applications. By using the Gone MCP component, you can easily integrate AI models with your business systems.

## Quick Start

Complete code reference: [mcp/quick_start](../examples/mcp/quick_start)

<details>
<summary>1. Create project and install dependencies</summary>

```bash
go mod init quickstart
go get github.com/gone-io/goner/mcp
```

</details>

<details>
<summary>2. Create server based on stdio, create file `server/main.go`</summary>

```go
// This is a server-side example program based on MCP (Model-Controller-Provider)
// It demonstrates how to create a simple MCP tool that accepts a name parameter and returns a greeting

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

// main function initializes and starts the MCP server
// Loads necessary components and registers the helloTool
func main() {
	_ = os.Setenv("GONE_MCP", `{"name":"quick-start", "version": "0.0.1"}`)
	gone.
		Loads(goneMcp.ServerLoad). // Load MCP server component
		Load(&helloTool{}). // Load helloTool
		Serve() // Start service
}

// helloTool implements a simple greeting tool
// Inherits gone.Flag to support basic MCP tool functionality
type helloTool struct {
	gone.Flag
}

// Define defines the tool's name, description, and parameters
// Return value: goneMcp.Tool - Tool definition information
func (h helloTool) Define() goneMcp.Tool {
	return mcp.NewTool("hello_world",
		mcp.WithDescription("Say hello to someone"), // Tool description
		mcp.WithString("name", // Define string type name parameter
			mcp.Required(),                                 // Set name parameter as required
			mcp.Description("Name of the person to greet"), // name parameter description
		),
	)
}

// Handler handles tool call requests
// Parameters:
//   - ctx: Context information
//   - request: Call request containing parameter information
//
// Return values:
//   - *mcp.CallToolResult: Call result
//   - error: Error information
func (h helloTool) Handler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get and validate name parameter
	name, ok := request.Params.Arguments["name"].(string)
	if !ok {
		return nil, errors.New("name must be a string") // Return error if parameter type is incorrect
	}
	// Return formatted greeting
	return mcp.NewToolResultText(fmt.Sprintf("Hello, %s!", name)), nil
}

// Compile-time type check to ensure helloTool implements goneMcp.ITool interface
var _ goneMcp.ITool = (*helloTool)(nil)
```

</details>

<details>
<summary>3. Create client based on stdio, create file `client/main.go`</summary>

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

// main function is the program entry point
// This is a client example program based on MCP (Message Communication Protocol)
// It demonstrates how to initialize MCP client, connect to server, list available tools, and call tools
func main() {
	gone.
		NewApp(goneMcp.ClientLoad). // Create a new Gone application using MCP client loader
		Run(func(in struct {
			// Inject MCP client, using stdio type communication, and run through local server with command: go run ./server
			client *client.Client `gone:"*,type=stdio,param=go run ./server"`
		}) {
			c := in.client
			// Create context with 30-second timeout
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			// Start client
			err := c.Start(ctx)
			if err != nil {
				log.Fatalf("Failed to start client: %v", err)
			}

			// Prepare initialization request
			initRequest := mcp.InitializeRequest{}
			// Set protocol version and client information
			initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
			initRequest.Params.ClientInfo = mcp.Implementation{
				Name:    "example-client",
				Version: "1.0.0",
			}

			// Send initialization request to server
			initResult, err := c.Initialize(ctx, initRequest)
			if err != nil {
				log.Fatalf("Failed to initialize: %v", err)
			}
			// Print server information
			fmt.Printf(
				"Initialized with server: %s %s\n\n",
				initResult.ServerInfo.Name,
				initResult.ServerInfo.Version,
			)

			// List all available tools provided by the server
			fmt.Println("Listing available tools...")
			toolsRequest := mcp.ListToolsRequest{}
			tools, err := c.ListTools(ctx, toolsRequest)
			if err != nil {
				log.Fatalf("Failed to list tools: %v", err)
			}
			// Iterate and print each tool's name and description
			for _, tool := range tools.Tools {
				fmt.Printf("- %s: %s\n", tool.Name, tool.Description)
			}
			fmt.Println()

			// Call hello_world tool example
			fmt.Println("Calling `hello_world`")
			request := mcp.CallToolRequest{}
			// Set tool name and parameters to call
			request.Params.Name = "hello_world"
			request.Params.Arguments = map[string]any{
				"name": "John",
			}
			// Execute tool call
			tool, err := c.CallTool(ctx, request)
			if err != nil {
				log.Fatalf("Failed to call tool: %v", err)
			}
			// Print tool execution result
			printToolResult(tool)
		})
}

// printToolResult helper function for printing tool execution results
// Supports formatting text content and other types of content
func printToolResult(result *mcp.CallToolResult) {
	for _, content := range result.Content {
		// If it's text content, print directly
		if textContent, ok := content.(mcp.TextContent); ok {
			fmt.Println(textContent.Text)
		} else {
			// For other types of content, convert to JSON format and print
			jsonBytes, _ := json.MarshalIndent(content, "", "  ")
			fmt.Println(string(jsonBytes))
		}
	}
}
```

</details>

<details>
<summary>4. Run: `go run client/main.go`</summary>

```bash
go mod tidy
go run client/main.go
```

Execution result:

```
Initialized with server: quick-start 0.0.1

Listing available tools...
- hello_world: Say hello to someone

Calling `hello_world`
Hello, John!
```

</details>

## Gone MCP Component Features

- Server-side
    - Supports `Goner Define` and `Functional Define` for defining MCP Tool, Prompt, Resource
    - Supports configuration files
    - Supports Stdio and SSE channels
    - Supports Hook injection
    - Supports ContextFunc injection
    - Supports multiple instances
- Client-side
    - Supports multiple instances based on different `gone` tag configurations
    - Supports parameter reading from configuration
    - Supports Stdio and SSE
    - SSE supports header setting
    - Supports custom `transport.Interface`

## Server-side Usage

### Basic Concepts

MCP (Model Control Protocol) is a protocol for controlling AI models, and the Gone MCP component provides support for this protocol. The server-side mainly includes the following core concepts:

- **Tool**: Defines tools that can be called by AI models
- **Prompt**: Defines prompt templates
- **Resource**: Defines resources that can be accessed by AI models

### Gone MCP Component Supports Two Ways to Define **Tool**, **Prompt**, and **Resource**

#### Goner Define: Based on Implementing Interfaces `goneMcp.ITool`, `goneMcp.IPrompt`, `goneMcp.IResource`

**Examples, code location: [examples/mcp/stdio/server/goner_define](../examples/mcp/stdio/server/goner_define):**

<details>
<summary>1. Define Tool</summary>

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
<summary>2. Define Prompt</summary>

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
<summary>3. Define Resource</summary>

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

#### Functional Define: Based on Function-style Definitions in `github.com/mark3labs/mcp-go`

**Examples, code location: [examples/mcp/stdio/server/functional_define](../examples/mcp/stdio/server/functional_add/mcp_server_define.go)**

### Server Configuration

Based on Gone's powerful configuration capabilities, services can quickly support local file configuration and remote configuration. For more about configuration, refer to: [viper local configuration](../viper/README.md), [viper remote configuration](../viper/remote/README.md).
Below explains the meaning of corresponding keys:

| key                                      | Description                    | Corresponding `mcp-go` option            | Default     |
| ---------------------------------------- | ------------------------------ | --------------------------------------- | ----------- |
| mcp.name                                 | Server name                    | None                                    | `MCP Server`|
| mcp.version                              | Server version                 | None                                    | `1.0.0`     |
| mcp.withRecovery                         | Enable recovery mechanism      | `server.WithRecovery`                   | `false`     |
| mcp.withPromptCapabilities               | Enable prompt capabilities     | `server.WithPromptCapabilities`         | `false`     |
| mcp.withToolCapabilities                 | Enable tool capabilities       | `server.WithToolCapabilities`           | `false`     |
| mcp.withLogging                          | Enable logging                 | `server.WithLogging`                    | `false`     |
| mcp.withInstructions                     | Custom instruction content     | `server.WithInstructions`               | `""`      |
| mcp.transportType                        | Transport type, supports `stdio` and `sse` | `server.WithTransportType`      | `stdio`     |
| mcp.sse.withBaseURL                      | SSE transport type BaseURL     | `server.WithSSEBaseURL`                 | `""`      |
| mcp.sse.withBasePath                     | SSE transport base path        | `server.WithBasePath`                   | `""`      |
| mcp.sse.withMessageEndpoint              | SSE message endpoint           | `server.WithMessageEndpoint`            | `""`      |
| mcp.sse.withUseFullURLForMessageEndpoint | Use full URL as message endpoint | `server.WithUseFullURLForMessageEndpoint` | `false`   |
| mcp.sse.withSSEEndpoint                  | SSE endpoint                   | `server.WithSSEEndpoint`                | `""`      |
| mcp.sse.withKeepAlive                    | Enable keep-alive              | `server.WithKeepAlive`                  | `false`     |
| mcp.sse.withKeepAliveInterval            | Keep-alive interval            | `server.WithKeepAliveInterval`          | `0`         |
| mcp.sse.address                          | SSE server listening address   | None                                    | `:8080`     |

### Support for Stdio and SSE Channels
#### Stdio
Set `mcp.transportType` to `stdio`, and the MCP server will use Stdio as the communication channel.

#### SSE
Set `mcp.transportType` to `sse`, and the MCP server will use SSE as the communication channel. In SSE channel mode, you can configure service-related parameters: `mcp.sse.*`, refer to [Server Configuration](#server-configuration) for specific meanings.

### Support for Loading `*server.Hooks` into Gone Framework

MCP server supports implementing function interception through custom Hooks, allowing execution of custom logic before and after tool calls. Here's an example:

```go
// Register Hooks in the application
func main() {
    gone.
		NewApp(goneMcp.Load).
        Load(g.NamedThirdComponentLoadFunc("mcp.hooks", &server.Hooks{
			OnBeforeCallTool: []server.OnBeforeCallToolFunc{
				func(ctx context.Context, id any, message *mcp.CallToolRequest) {
					//todo processing logic
				},
			},
		})).
        Run()
}
```

Notes:
1. `g.NamedThirdComponentLoadFunc` converts `*server.Hooks` into a load function that Gone framework can recognize;
2. The name used here is `mcp.hooks`; for multi-instance cases, the prefix `mcp` should be replaced with the specific instance prefix;
3. Hook functions refer to [`github.com/mark3labs/mcp-go/mcp` documentation](http://github.com/mark3labs/mcp-go/blob/main/mcp)

### Support for Loading `server.StdioContextFunc` or `server.SseContextFunc`

Below is an example using stdio channel, other channels are used similarly.

```go
func main() {
    gone.
		NewApp(goneMcp.Load).
        Load(
			// g.NamedThirdComponentLoadFunc converts server.StdioContextFunc to a load function that Gone framework can recognize
			g.NamedThirdComponentLoadFunc(
				"mcp.stdio.context", // prefix is `mcp`, multi-instance prefix depends on specific injection definition
				// Need to convert the defined function to server.StdioContextFunc for MCP component to recognize
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

Notes:
1. server.StdioContextFunc is loaded with the name `mcp.stdio.context`;
2. server.SseContextFunc is loaded with the name `mcp.sse.context`;
3. For multi-instance cases, the prefix `mcp` should be replaced with the specific instance prefix;

### Multi-instance Support
For **Functional Define**, you can use the `gone` tag to specify instance names, for example:

```go
package service
import (
	"github.com/gone-io/gone/v2"
	"github.com/mark3labs/mcp-go/server"
)

type MyService struct {
	gone.Flag
	s1 *server.MCPServer `gone:"*,customer1"` // instance name is "customer1", note the `*,` prefix, asterisk can be omitted, comma cannot
	s2 *server.MCPServer `gone:",customer2"`
}
```
Notes:
 1. For the `MyService` defined above, both `s1` and `s2` are instances of `*server.MCPServer`, corresponding to `customer1` and `customer2` instances.
 2. The Gone framework will start two MCP servers, corresponding to `customer1` and `customer2` instances.
 3. Their configuration reading should use the corresponding instance names as prefixes, for example, server names are `customer1.name` and `customer2.name` respectively
 4. `*server.Hooks`, `server.StdioContextFunc`, `server.SseContextFunc` also support multiple instances, the injection name prefix needs to replace `mcp` with the specific instance name.

## Client Usage

### Client Injection

```go
package service
import (
	"github.com/gone-io/gone/v2"
	"github.com/mark3labs/mcp-go/client"
)
type MyService struct {
gone.Flag
	// Use Stdio client
	StdioClient *client.Client `gone:"*,type=stdio,param=./mcp-server --param=value"`

	// Use Stdio client defined in configuration file, config key is `mcp.stdio` [Using Configuration]
	StdioClientFromConfig *client.Client `gone:"*,type=stdio,configKey=mcp.stdio"`

	// Use SSE client
	SseClient *client.Client `gone:"*,type=sse,param=http://localhost:8080/sse"`

	// Use SSE client defined in configuration file, config key is `mcp.sse` [Using Configuration]
	SseClientFromConfig *client.Client `gone:"*,type=sse,configKey=mcp.sse"`

	// Use custom Transport client
	CustomTransportClient *client.Client `gone:"*,type=transport,param=myTransport"`
}
```

### Configuration
#### Stdio Client Configuration
| key                  | Description      | Default |
| -------------------- | ---------------- | ------- |
| ${configKey}.command | Command          | None    |
| ${configKey}.env     | Environment vars | None    |
| ${configKey}.args    | Command args     | None    |

**Example:**

```yaml
# Stdio client configuration
mcp.stdio:
  command: "./mcp-server"  # Command path
  env: # Environment variables
    - "ENV_VAR=value"
  args: # Command arguments
    - "--param=value"
```

#### SSE Client Configuration
| key                  | Description      | Default |
| -------------------- | ---------------- | ------- |
| ${configKey}.baseUrl | SSE server URL   | None    |
| ${configKey}.header  | Request headers  | None    |

**Example:**

```yaml
# SSE client configuration
mcp.sse:
  baseUrl: "http://localhost:8080/sse"  # SSE server URL
  header: # Request headers
    Authorization: "Bearer token"
```

### Using the Client
For details on how to use the client, please refer to [github.com/mark3labs/mcp-go/client](https://github.com/mark3labs/mcp-go)

### Examples

For complete examples, refer to the [examples/mcp](https://github.com/gone-io/goner/tree/main/examples/mcp) directory.

### Advanced Usage
#### Custom Transport
1. Define Transport
```go
type xTransport struct {
	gone.Flag
}

func (x xTransport) Start(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (x xTransport) SendRequest(ctx context.Context, request *mcp.Request) (*mcp.Response, error) {
	//TODO implement me
	panic("implement me")
}

func (x xTransport) Close() error {
	//TODO implement me
	panic("implement me")
}
```