package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gone-io/gone/v2"
	goneMcp "github.com/gone-io/goner/mcp"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
	"log"
	"time"
)

func main() {
	gone.
		NewApp(goneMcp.ClientLoad).
		Run(func(in struct {
			client *client.Client `gone:"*,type=sse,param=http://localhost:8082/sse"`
		}) {
			c := in.client
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			err := c.Start(ctx)
			if err != nil {
				log.Fatalf("Failed to start client: %v", err)
			}

			initRequest := mcp.InitializeRequest{}
			initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
			initRequest.Params.ClientInfo = mcp.Implementation{
				Name:    "example-client",
				Version: "1.0.0",
			}

			initResult, err := c.Initialize(ctx, initRequest)
			if err != nil {
				log.Fatalf("Failed to initialize: %v", err)
			}
			fmt.Printf(
				"Initialized with server: %s %s\n\n",
				initResult.ServerInfo.Name,
				initResult.ServerInfo.Version,
			)

			// List Tools
			fmt.Println("Listing available tools...")
			toolsRequest := mcp.ListToolsRequest{}
			tools, err := c.ListTools(ctx, toolsRequest)
			if err != nil {
				log.Fatalf("Failed to list tools: %v", err)
			}
			for _, tool := range tools.Tools {
				fmt.Printf("- %s: %s\n", tool.Name, tool.Description)
			}
			fmt.Println()

			fmt.Println("Listing available prompts...")
			prompts, err := c.ListPrompts(ctx, mcp.ListPromptsRequest{})
			if err != nil {
				log.Fatalf("Failed to list prompts: %v", err)
			}
			for _, prompt := range prompts.Prompts {
				fmt.Printf("- %s: %s\n", prompt.Name, prompt.Description)
			}
			fmt.Println()

			fmt.Println("Listing available resources...")
			resources, err := c.ListResources(ctx, mcp.ListResourcesRequest{})
			if err != nil {
				log.Fatalf("Failed to list resources: %v", err)
			}
			for _, resource := range resources.Resources {
				fmt.Printf("- %s: %s\n", resource.Name, resource.Description)
			}
			fmt.Println()

			fmt.Println("Listing available resource templates...")
			templates, err := c.ListResourceTemplates(ctx, mcp.ListResourceTemplatesRequest{})
			if err != nil {
				log.Fatalf("Failed to list resource templates: %v", err)
			}
			for _, template := range templates.ResourceTemplates {
				fmt.Printf("- %s: %s\n", template.Name, template.Description)
			}
			fmt.Println()

			fmt.Println("Calling `hello_world`")
			request := mcp.CallToolRequest{}
			request.Params.Name = "hello_world"
			request.Params.Arguments = map[string]any{
				"name": "John",
			}
			tool, err := c.CallTool(ctx, request)
			if err != nil {
				log.Fatalf("Failed to call tool: %v", err)
			}
			printToolResult(tool)

			fmt.Println()

			fmt.Println("Reading resource...")
			resourceRequest := mcp.ReadResourceRequest{}
			resourceRequest.Params.URI = "docs://readme"
			resource, err := c.ReadResource(ctx, resourceRequest)
			if err != nil {
				log.Fatalf("Failed to read resource: %v", err)
			}
			fmt.Printf("%#v\n", resource.Contents[0])
			fmt.Println()

			resourceRequest.Params.URI = "users://10/profile"
			resource, err = c.ReadResource(ctx, resourceRequest)
			if err != nil {
				log.Fatalf("Failed to read resource: %v", err)
			}
			fmt.Printf("%#v", resource.Contents[0])

			fmt.Println()

			fmt.Println("Get Prompt...")
			promptRequest := mcp.GetPromptRequest{}
			promptRequest.Params.Name = "code_review"
			promptRequest.Params.Arguments = map[string]string{
				"pr_number": "123",
			}
			prompt, err := c.GetPrompt(ctx, promptRequest)
			if err != nil {
				log.Fatalf("Failed to get prompt: %v", err)
			}
			for _, msg := range prompt.Messages {
				fmt.Printf("%s: %s\n", msg.Role, msg.Content.(mcp.TextContent).Text)
			}
			fmt.Println()

		})
}

// Helper function to print tool results
func printToolResult(result *mcp.CallToolResult) {
	for _, content := range result.Content {
		if textContent, ok := content.(mcp.TextContent); ok {
			fmt.Println(textContent.Text)
		} else {
			jsonBytes, _ := json.MarshalIndent(content, "", "  ")
			fmt.Println(string(jsonBytes))
		}
	}
}
