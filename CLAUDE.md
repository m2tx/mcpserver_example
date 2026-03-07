# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Build
go build ./cmd/server

# Run (stdio mode)
./server

# Run (HTTP mode)
HTTP_PORT=8080 ./server

# Run tests
go test ./...

# Run a single test
go test ./internal/tools/ -run TestName

# Tidy dependencies
go mod tidy
```

## Architecture

This is a minimal [MCP](https://modelcontextprotocol.io) server in Go using the official `github.com/modelcontextprotocol/go-sdk` SDK.

### Key interfaces and flow

- **`Tool` interface** ([internal/mcp/handler.go](internal/mcp/handler.go)): Every tool implements `Register(s *sdkmcp.Server)`. Inside `Register`, call `mcp.AddTool(s, ...)` with a typed args struct — the SDK uses reflection + JSON schema tags to generate the tool schema automatically.

- **`Server`** ([internal/mcp/server.go](internal/mcp/server.go)): Wraps the SDK server. `New(name, version, tools...)` registers all tools. `Run(ctx, addr)` selects transport: empty string → stdio, non-empty → streamable HTTP on that address with graceful shutdown on context cancellation.

- **Entry point** ([cmd/server/main.go](cmd/server/main.go)): Sets up signal-aware context (`signal.NotifyContext`), reads `HTTP_PORT` env var, and calls `Run`.

### Adding a new tool

1. Create `internal/tools/mytool.go` with a struct implementing `Register`.
2. Add `&tools.MyTool{}` to the `toolsList` slice in `cmd/server/main.go`.

### Transport

- **stdio** (default): for use with Claude Desktop / `claude mcp add`.
- **HTTP** (`HTTP_PORT=<port>`): streamable HTTP transport, one SDK server instance shared across requests.
