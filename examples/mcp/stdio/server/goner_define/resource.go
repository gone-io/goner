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
