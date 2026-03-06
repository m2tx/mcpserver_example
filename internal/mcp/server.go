package mcp

import (
	"context"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

type Server struct {
	s *sdkmcp.Server
}

func New(name, version string, h *Handler) *Server {
	s := sdkmcp.NewServer(&sdkmcp.Implementation{
		Name:    name,
		Version: version,
	}, nil)
	for _, t := range h.tools {
		t.Register(s)
	}
	return &Server{s: s}
}

func (srv *Server) Run(ctx context.Context) error {
	return srv.s.Run(ctx, &sdkmcp.StdioTransport{})
}
