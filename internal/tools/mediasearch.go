package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/m2tx/mcpserver_example/internal/client"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// MediaSearchArgs defines the input for media_search.
// AGENT INSTRUCTIONS: Always call this tool multiple times with different name
// variations (e.g. full name, name without middle name, common aliases) to
// maximise recall. Do not stop after a single call — iterate until you have
// tried all plausible variants for the person being searched.
type MediaSearchArgs struct {
	Name    string `json:"name"    jsonschema:"Full or partial name of the person to search. Call the tool multiple times with different name variations (full name, without middle name, aliases) to maximise results."`
	Context string `json:"context,omitempty" jsonschema:"Optional context such as country or organization to narrow results"`
}

type MediaSearch struct {
	client *client.DuckDuckGoClient
}

func NewMediaSearch(c *client.DuckDuckGoClient) *MediaSearch {
	return &MediaSearch{client: c}
}

func (m *MediaSearch) Register(s *mcp.Server) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "media_search",
		Title:       "Media Search",
		Description: "Search online media about a person using multiple targeted queries covering news, legal issues, fraud, and sanctions",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args MediaSearchArgs) (*mcp.CallToolResult, any, error) {
		name := strings.TrimSpace(args.Name)
		if name == "" {
			return nil, nil, fmt.Errorf("name must not be empty")
		}

		base := name
		if c := strings.TrimSpace(args.Context); c != "" {
			base = name + " " + c
		}

		queries := []string{
			base,
			base + " fraud",
			base + " corruption",
			base + " criminal",
			base + " lawsuit",
			base + " scandal",
			base + " sanction",
			base + " arrest",
			base + " fraude",
			base + " corrupção",
			base + " crime",
			base + " processo",
			base + " escândalo",
			base + " sanção",
			base + " prisão",
			base + " homicídio",
			base + " latrocínio",
			base + " sequestro",
			base + " estupro",
			base + " tráfico",
			base + " terrorismo",
			base + " tortura",
			base + " genocídio",
		}

		var sb strings.Builder
		fmt.Fprintf(&sb, "# Media Search Results for: %s\n\n", name)

		for _, q := range queries {
			results, err := m.client.Search(ctx, q, 5)
			if err != nil || len(results) == 0 {
				fmt.Printf("Error searching for %q: %v\n", q, err)
				continue
			}
			fmt.Fprintf(&sb, "## Query: %q\n", q)
			for i, r := range results {
				fmt.Fprintf(&sb, "%d. %s\n   %s\n", i+1, r.Title, r.URL)
				if r.Snippet != "" {
					fmt.Fprintf(&sb, "   %s\n", r.Snippet)
				}
			}
			sb.WriteString("\n")
		}

		text := sb.String()
		if strings.Count(text, "\n") <= 2 {
			text = fmt.Sprintf("No media results found for: %s", name)
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: text},
			},
		}, nil, nil
	})
}
