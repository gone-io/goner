package goneMcp

import (
	"context"
	"github.com/mark3labs/mcp-go/mcp"
)

type Tool = mcp.Tool
type Prompt = mcp.Prompt
type Resource = mcp.Resource

type CallToolRequest = mcp.CallToolRequest
type CallToolResult = mcp.CallToolResult
type GetPromptRequest = mcp.GetPromptRequest
type GetPromptResult = mcp.GetPromptResult
type ReadResourceRequest = mcp.ReadResourceRequest
type ResourceContents = mcp.ResourceContents

type ITool interface {
	Define() Tool
	Process() func(ctx context.Context, request CallToolRequest) (*CallToolResult, error)
}

type IPrompt interface {
	Define() Prompt
	Process() func(ctx context.Context, request GetPromptRequest) (*GetPromptResult, error)
}

type IResource interface {
	Define() Resource
	Process() func(ctx context.Context, request ReadResourceRequest) ([]ResourceContents, error)
}
