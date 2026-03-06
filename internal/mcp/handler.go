package mcp

import sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"

type Tool interface {
	Register(s *sdkmcp.Server)
}

type Handler struct {
	tools []Tool
}

func NewHandler(tools ...Tool) *Handler {
	return &Handler{tools: tools}
}

