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
