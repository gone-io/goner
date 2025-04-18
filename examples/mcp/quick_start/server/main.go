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
		Load(&helloTool{}).        // 加载helloTool工具
		Serve()                    // 启动服务
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
			mcp.Required(), // 将name参数设置为必填
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
