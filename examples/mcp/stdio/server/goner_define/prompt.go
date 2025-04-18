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
