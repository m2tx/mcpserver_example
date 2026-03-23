package prompts

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type CodeReview struct{}

func (p *CodeReview) Register(s *mcp.Server) {
	s.AddPrompt(&mcp.Prompt{
		Name:        "code_review",
		Title:       "Code Review Prompt",
		Description: "Gera um prompt para revisão de código",
		Arguments: []*mcp.PromptArgument{
			{Name: "code", Required: true},
			{Name: "language", Required: false},
		},
	}, func(ctx context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		code := req.Params.Arguments["code"]
		lang := req.Params.Arguments["language"]
		var text string
		if lang != "" {
			text = fmt.Sprintf("Por favor, revise o seguinte código %s:\n\n%s", lang, code)
		} else {
			text = fmt.Sprintf("Por favor, revise o seguinte código:\n\n%s", code)
		}
		return &mcp.GetPromptResult{
			Description: "Prompt de revisão de código",
			Messages: []*mcp.PromptMessage{
				{Role: "user", Content: &mcp.TextContent{Text: text}},
			},
		}, nil
	})
}
