package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	mcpserver "github.com/m2tx/mcpserver_example/internal/mcp"
	"github.com/m2tx/mcpserver_example/internal/tools"
)

func getHTTPPort() string {
	port := os.Getenv("HTTP_PORT")
	if port != "" {
		return port
	}

	return "9000"
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	toolsList := []mcpserver.Tool{
		&tools.Add{},
		&tools.Greet{},
	}

	server := mcpserver.New("mcp-server", "1.0.0", toolsList...)

	if err := server.Run(ctx, fmt.Sprintf(":%s", getHTTPPort())); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
