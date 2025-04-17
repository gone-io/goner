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

func (h helloTool) Process() func(ctx context.Context, request goneMcp.CallToolRequest) (*goneMcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		name, ok := request.Params.Arguments["name"].(string)
		if !ok {
			return nil, errors.New("name must be a string")
		}
		return mcp.NewToolResultText(fmt.Sprintf("Hello, %s!", name)), nil
	}
}

var _ goneMcp.ITool = (*helloTool)(nil)
