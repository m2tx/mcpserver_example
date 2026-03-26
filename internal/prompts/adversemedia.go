package prompts

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"text/template"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

//go:embed adversemedia.md
var adverseMediaTemplate string

type AdverseMedia struct{}

func (p *AdverseMedia) Register(s *mcp.Server) {
	tmpl := template.Must(template.New("adverse_media").Parse(adverseMediaTemplate))

	s.AddPrompt(&mcp.Prompt{
		Name:        "adverse_media",
		Title:       "Analyze Media Search Results for Adverse Media",
		Description: "Guide to analyze online media search results and identify adverse media about a person",
		Arguments: []*mcp.PromptArgument{
			{Name: "name", Required: true, Description: "Full name of the person being investigated"},
			{Name: "search_results", Required: true, Description: "Raw search results from the media_search tool"},
			{Name: "context", Required: false, Description: "Optional context such as country or organization"},
		},
	}, func(ctx context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		name := req.Params.Arguments["name"]
		searchResults := req.Params.Arguments["search_results"]
		additionalContext := req.Params.Arguments["context"]

		var header bytes.Buffer
		fmt.Fprintf(&header, "You are a compliance analyst conducting an adverse media screening for **%s**.", name)
		if additionalContext != "" {
			fmt.Fprintf(&header, " Context: %s.", additionalContext)
		}

		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, map[string]string{
			"Header":        header.String(),
			"SearchResults": searchResults,
			"Name":          name,
		}); err != nil {
			return nil, err
		}

		return &mcp.GetPromptResult{
			Description: "Adverse media analysis prompt for compliance screening",
			Messages: []*mcp.PromptMessage{
				{Role: "user", Content: &mcp.TextContent{Text: buf.String()}},
			},
		}, nil
	})
}
