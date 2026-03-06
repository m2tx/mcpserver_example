package tools

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type GreetArgs struct {
	Name string `json:"name" jsonschema:"The person's name"`
}

type Greet struct{}

func (g *Greet) Register(s *mcp.Server) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "greet",
		Description: "Greet someone by name",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args GreetArgs) (*mcp.CallToolResult, any, error) {
		greeting := fmt.Sprintf("Hello, %s! Welcome to the MCP server.", args.Name)
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: greeting},
			},
		}, nil, nil
	})
}
