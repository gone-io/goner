package main

import (
	"context"
	"github.com/gone-io/gone/v2"
	goneMcp "github.com/gone-io/goner/mcp"
	"github.com/mark3labs/mcp-go/mcp"
)

type toolDefine struct {
	gone.Flag
}

func (t toolDefine) Define() goneMcp.Tool {
	return mcp.NewTool("test", mcp.WithDescription("test tool"))
}

func (t toolDefine) Handler(ctx context.Context, request goneMcp.CallToolRequest) (*goneMcp.CallToolResult, error) {
	return mcp.NewToolResultText("test"), nil
}

var _ goneMcp.ITool = (*toolDefine)(nil)

func main() {
	gone.
		NewApp(goneMcp.ServerLoad).
		Load(&toolDefine{}).
		Serve()
}
