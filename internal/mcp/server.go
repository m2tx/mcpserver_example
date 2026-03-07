package mcp

import (
	"context"
	"errors"
	"net/http"
	"time"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

type Server struct {
	s *sdkmcp.Server
}

func New(name, version string, tools ...Tool) *Server {
	s := sdkmcp.NewServer(&sdkmcp.Implementation{
		Name:    name,
		Version: version,
	}, nil)
	for _, t := range tools {
		t.Register(s)
	}
	return &Server{s: s}
}

// Run starts the server. If addr is non-empty it listens for HTTP connections
// on that address; otherwise it communicates over stdio.
func (srv *Server) Run(ctx context.Context, addr string) error {
	if addr != "" {
		handler := sdkmcp.NewStreamableHTTPHandler(func(r *http.Request) *sdkmcp.Server {
			return srv.s
		}, nil)
		httpSrv := &http.Server{Addr: addr, Handler: handler}
		go func() {
			<-ctx.Done()
			shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			httpSrv.Shutdown(shutdownCtx)
		}()
		if err := httpSrv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	}

	return srv.s.Run(ctx, &sdkmcp.StdioTransport{})
}
