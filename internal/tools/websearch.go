package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/m2tx/mcpserver_example/internal/client"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type WebSearchArgs struct {
	Query string `json:"query" jsonschema:"Search query"`
}

type WebSearch struct {
	client *client.DuckDuckGoClient
}

func NewWebSearch(c *client.DuckDuckGoClient) *WebSearch {
	return &WebSearch{client: c}
}

func (w *WebSearch) Register(s *mcp.Server) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "web_search",
		Title:       "Web Search",
		Description: "Search the web using DuckDuckGo Instant Answer API",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args WebSearchArgs) (*mcp.CallToolResult, any, error) {
		query := args.Query
		if query == "" {
			return nil, nil, fmt.Errorf("query must not be empty")
		}

		results, err := w.client.Search(ctx, query, 5)
		if err != nil {
			return nil, nil, err
		}

		if len(results) == 0 {
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: fmt.Sprintf("No results found for: %s", query)},
				},
			}, nil, nil
		}

		lines := []string{fmt.Sprintf("Results for: %s (via DuckDuckGo)", query)}
		for i, r := range results {
			lines = append(lines, fmt.Sprintf("%d. %s\n   %s", i+1, r.Title, r.URL))
			if r.Snippet != "" {
				lines = append(lines, fmt.Sprintf("   %s", r.Snippet))
			}
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: strings.Join(lines, "\n")},
			},
		}, nil, nil
	})
}
