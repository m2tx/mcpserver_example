package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	mcpserver "github.com/m2tx/mcpserver_example/internal/mcp"
	"github.com/m2tx/mcpserver_example/internal/tools"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	addr := ""
	if port := os.Getenv("HTTP_PORT"); port != "" {
		addr = ":" + port
	}

	toolsList := []mcpserver.Tool{
		&tools.Add{},
		&tools.Greet{},
	}

	server := mcpserver.New("example-mcp-server", "1.0.0", toolsList...)

	if err := server.Run(ctx, addr); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
