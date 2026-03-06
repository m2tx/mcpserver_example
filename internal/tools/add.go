package tools

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type AddArgs struct {
	A float64 `json:"a" jsonschema:"First number"`
	B float64 `json:"b" jsonschema:"Second number"`
}

type Add struct{}

func (a *Add) Register(s *mcp.Server) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "add",
		Description: "Add two numbers together",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args AddArgs) (*mcp.CallToolResult, any, error) {
		result := args.A + args.B
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("%.0f + %.0f = %.0f", args.A, args.B, result)},
			},
		}, nil, nil
	})
}
