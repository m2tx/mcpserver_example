package mcp

import (
	"context"
	"net/http"

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

// Run starts the server. If addr is non-empty it listens for HTTP connections
// on that address; otherwise it communicates over stdio.
func (srv *Server) Run(ctx context.Context, httpPort string) error {
	if httpPort != "" {
		handler := sdkmcp.NewStreamableHTTPHandler(func(r *http.Request) *sdkmcp.Server {
			return srv.s
		}, nil)
		httpSrv := &http.Server{Addr: httpPort, Handler: handler}
		go func() {
			<-ctx.Done()
			httpSrv.Shutdown(context.Background())
		}()
		return httpSrv.ListenAndServe()
	}

	return srv.s.Run(ctx, &sdkmcp.StdioTransport{})
}
