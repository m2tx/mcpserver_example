package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/m2tx/mcpserver_example/internal/client"
	"github.com/m2tx/mcpserver_example/internal/mcp"
	"github.com/m2tx/mcpserver_example/internal/prompts"
	"github.com/m2tx/mcpserver_example/internal/tools"
)

const (
	appName = "mcp-server"
	version = "1.0.0"
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

	ddg := client.NewDuckDuckGoClient()

	toolsList := []mcp.Tool{
		&tools.Add{},
		&tools.Greet{},
		tools.NewWebSearch(ddg),
	}

	promptsList := []mcp.Prompt{
		&prompts.CodeReview{},
	}

	port := getHTTPPort()

	log.Printf("Starting %s v%s on :%s", appName, version, port)

	log.Printf("Tools loaded (%d):", len(toolsList))
	for _, t := range toolsList {
		log.Printf("  - %T", t)
	}

	log.Printf("Prompts loaded (%d):", len(promptsList))
	for _, p := range promptsList {
		log.Printf("  - %T", p)
	}

	server := mcp.New(appName, version, toolsList, promptsList)

	if err := server.Run(ctx, fmt.Sprintf(":%s", port)); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
