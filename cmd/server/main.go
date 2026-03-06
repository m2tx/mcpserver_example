package main

import (
	"context"
	"log"

	mcpserver "github.com/m2tx/mcpserver_example/internal/mcp"
	"github.com/m2tx/mcpserver_example/internal/tools"
)

func main() {
	h := mcpserver.NewHandler(&tools.Add{}, &tools.Greet{})
	s := mcpserver.New("example-mcp-server", "1.0.0", h)

	if err := s.Run(context.Background()); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
