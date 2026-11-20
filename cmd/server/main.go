package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"

	"github.com/m2tx/mcpserver_example/internal/client"
	"github.com/m2tx/mcpserver_example/internal/docs"
	"github.com/m2tx/mcpserver_example/internal/mcp"
	"github.com/m2tx/mcpserver_example/internal/prompts"
	"github.com/m2tx/mcpserver_example/internal/tools"
)

const (
	appName = "mcp-server"
	version = "1.0.0"
)

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}

	return def
}

func getHTTPPort() string {
	port := os.Getenv("HTTP_PORT")
	if port != "" {
		return port
	}

	return "9000"
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	ddg := client.NewDuckDuckGoClient()

	embedder := docs.NewEmbedder(logger)
	if err := embedder.Index(getEnv("DOCS_DIR", "docs")); err != nil {
		logger.Error("index docs", "err", err)
		os.Exit(1)
	}

	toolsList := []mcp.Tool{
		&tools.Add{},
		&tools.Greet{},
		tools.NewWebSearch(ddg),
		tools.NewSearchDocs(embedder),
	}

	promptsList := []mcp.Prompt{
		&prompts.CodeReview{},
	}

	port := getHTTPPort()

	logger.Info("starting server", "app", appName, "version", version, "port", port)

	for _, t := range toolsList {
		logger.Info("tool", "type", fmt.Sprintf("%T", t))
	}

	for _, p := range promptsList {
		logger.Info("prompt", "type", fmt.Sprintf("%T", p))
	}

	server := mcp.New(appName, version, toolsList, promptsList, logger)

	if err := server.Run(ctx, fmt.Sprintf(":%s", port)); err != nil {
		logger.Error("server error", "err", err)
		os.Exit(1)
	}
}
