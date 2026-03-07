package main

import (
	"context"
	"fmt"
	"log"
	"os"

	mcpserver "github.com/m2tx/mcpserver_example/internal/mcp"
	"github.com/m2tx/mcpserver_example/internal/tools"
)

func getHttpPort() string {
	if port := os.Getenv("HTTP_PORT"); port != "" {
		return fmt.Sprintf(":%s", port)
	}

	return ""
}

func main() {
	ctx := context.Background()

	httpPort := getHttpPort()

	handler := mcpserver.NewHandler(&tools.Add{}, &tools.Greet{})
	server := mcpserver.New("example-mcp-server", "1.0.0", handler)

	if err := server.Run(ctx, httpPort); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
