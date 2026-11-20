package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/m2tx/mcpserver_example/internal/docs"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const searchDocsTopK = 3

type SearchDocsArgs struct {
	Query string `json:"query" jsonschema:"The search query describing what information you need"`
}

type SearchDocs struct {
	embedder *docs.Embedder
}

func NewSearchDocs(e *docs.Embedder) *SearchDocs {
	return &SearchDocs{embedder: e}
}

func (sd *SearchDocs) Register(s *mcp.Server) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "search_docs",
		Title:       "Search Docs",
		Description: "Searches the internal document library for information relevant to the query. Use this whenever the user asks about topics that might be covered in internal documentation.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args SearchDocsArgs) (*mcp.CallToolResult, any, error) {
		query := strings.TrimSpace(args.Query)
		if query == "" {
			return nil, nil, fmt.Errorf("query must not be empty")
		}

		results, err := sd.embedder.Search(query, searchDocsTopK)
		if err != nil {
			return nil, nil, err
		}

		if len(results) == 0 {
			return textResult(fmt.Sprintf("No documents found for: %s", query)), nil, nil
		}

		blocks := make([]string, 0, len(results))
		for i, r := range results {
			blocks = append(blocks, fmt.Sprintf("[%d] %s\n%s", i+1, r.Filename, r.Text))
		}

		return textResult(strings.Join(blocks, "\n\n---\n\n")), nil, nil
	})
}

func textResult(text string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: text}},
	}
}
