package functional_add

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gone-io/gone/v2"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type mcpServerDefine struct {
	gone.Flag

	// use MCPServer pointer Receive Injected Value
	s *server.MCPServer `gone:"*"`
}

func (m *mcpServerDefine) Init() {
	// add resource
	// Dynamic resource example - user profiles by ID
	template := mcp.NewResourceTemplate(
		"users://{id}/profile",
		"User Profile",
		mcp.WithTemplateDescription("Returns user profile information"),
		mcp.WithTemplateMIMEType("application/json"),
	)

	// Add template with its handler
	m.s.AddResourceTemplate(template, func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		// Extract ID from the URI using regex matching
		// The server automatically matches URIs to templates
		//userID := extractIDFromURI(request.Params.URI)

		//profile, err := getUserProfile(userID) // Your DB/API call here
		//if err != nil {
		//	return nil, err
		//}

		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      request.Params.URI,
				MIMEType: "application/json",
				Text:     `{"id": 10, "name":"Jim"}`,
			},
		}, nil
	})

	// add prompt
	// Code review prompt with embedded resource
	m.s.AddPrompt(mcp.NewPrompt("code_review",
		mcp.WithPromptDescription("Code review assistance"),
		mcp.WithArgument("pr_number",
			mcp.ArgumentDescription("Pull request number to review"),
			mcp.RequiredArgument(),
		),
	), func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		prNumber := request.Params.Arguments["pr_number"]
		if prNumber == "" {
			return nil, fmt.Errorf("pr_number is required")
		}

		contents := mcp.TextResourceContents{
			URI:      fmt.Sprintf("git://pulls/%s/diff", prNumber),
			MIMEType: "text/x-diff",
		}

		marshal, _ := json.Marshal(&contents)

		return mcp.NewGetPromptResult(
			"Code review assistance",
			[]mcp.PromptMessage{
				mcp.NewPromptMessage(
					mcp.RoleUser,
					mcp.NewTextContent("You are a helpful code reviewer. Review the changes and provide constructive feedback."),
				),
				mcp.NewPromptMessage(
					mcp.RoleAssistant,
					mcp.NewTextContent(string(marshal)),
				),
			},
		), nil
	})

	calculatorTool := mcp.NewTool("calculate",
		mcp.WithDescription("Perform basic arithmetic calculations"),
		mcp.WithString("operation",
			mcp.Required(),
			mcp.Description("The arithmetic operation to perform"),
			mcp.Enum("add", "subtract", "multiply", "divide"),
		),
		mcp.WithNumber("x",
			mcp.Required(),
			mcp.Description("First number"),
		),
		mcp.WithNumber("y",
			mcp.Required(),
			mcp.Description("Second number"),
		),
	)
	// add tool
	m.s.AddTool(calculatorTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		op := request.Params.Arguments["operation"].(string)
		x := request.Params.Arguments["x"].(float64)
		y := request.Params.Arguments["y"].(float64)

		var result float64
		switch op {
		case "add":
			result = x + y
		case "subtract":
			result = x - y
		case "multiply":
			result = x * y
		case "divide":
			if y == 0 {
				return mcp.NewToolResultError("cannot divide by zero"), nil
			}
			result = x / y
		}

		return mcp.FormatNumberResult(result), nil
	})

}
