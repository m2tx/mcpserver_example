# mcpserver_example

A minimal [Model Context Protocol (MCP)](https://modelcontextprotocol.io) server written in Go using the [go-sdk](https://github.com/modelcontextprotocol/go-sdk).

## Project Structure

```
mcpserver_example/
├── cmd/
│   └── server/
│       └── main.go          # Entry point
├── internal/
│   ├── client/
│   │   └── duckduckgo.go    # DuckDuckGo HTTP client
│   ├── docs/
│   │   └── embedder.go      # Local hash-embedding document index
│   ├── mcp/
│   │   ├── handler.go       # Tool and Prompt interfaces
│   │   └── server.go        # MCP server wrapper
│   ├── prompts/
│   │   └── codereview.go    # "code_review" prompt
│   └── tools/
│       ├── add.go           # "add" tool
│       ├── greet.go         # "greet" tool
│       ├── searchdocs.go    # "search_docs" tool
│       └── websearch.go     # "web_search" tool
├── docs/                    # Documents indexed by search_docs (contents git-ignored)
├── go.mod
└── go.sum
```

## Architecture

Each tool implements the `Tool` interface:

```go
type Tool interface {
    Register(s *mcp.Server)
}
```

Each prompt implements the `Prompt` interface:

```go
type Prompt interface {
    Register(s *mcp.Server)
}
```

Tools and prompts self-register with the MCP server. `mcp.New` accepts both slices, registers them all, and `Run` starts the streamable HTTP transport on the given address.

The `DuckDuckGoClient` ([internal/client/duckduckgo.go](internal/client/duckduckgo.go)) scrapes DuckDuckGo's HTML endpoint with built-in rate limiting (1 request/second minimum interval).

The `Embedder` ([internal/docs/embedder.go](internal/docs/embedder.go)) indexes every `.txt`, `.md` and `.pdf` file in `DOCS_DIR` at startup, split into ~800-character chunks, using local feature-hashing embeddings (no external API). `search_docs` returns the 3 chunks with the highest cosine similarity to the query.

## Tools

| Tool         | Description                                      | Arguments               |
|--------------|--------------------------------------------------|-------------------------|
| `add`        | Add two numbers together                         | `a` (float), `b` (float)|
| `greet`      | Greet someone by name                            | `name` (string)         |
| `web_search` | Search the web using DuckDuckGo Instant Answer API | `query` (string)      |
| `search_docs` | Search the local document library (`DOCS_DIR`)   | `query` (string)        |

## Prompts

| Prompt        | Description                      | Arguments                                    |
|---------------|----------------------------------|----------------------------------------------|
| `code_review` | Structured code review guidance  | `code` (string), `language` (string, optional)|

## Getting Started

### Prerequisites

- Go 1.25+

### Build

```bash
go build ./cmd/server
```

### Run

```bash
./server
```

The server defaults to HTTP mode on port **9000**. Set `HTTP_PORT` to override:

```bash
HTTP_PORT=8080 ./server
```

| Variable | Default | Description |
|---|---|---|
| `HTTP_PORT` | `9000` | HTTP listen port |
| `DOCS_DIR` | `docs` | Directory indexed for `search_docs` (missing/empty → no results) |

## Using with MCP Clients

### Claude Desktop / Claude Code (HTTP)

Start the server, then register it using the streamable HTTP transport:

```bash
./server  # listens on :9000
```

```bash
claude mcp add --transport http example http://localhost:9000/mcp
```

## Adding a New Tool

1. Create a new file in `internal/tools/`, e.g. `internal/tools/mytool.go`:

```go
package tools

import (
    "context"
    "fmt"

    "github.com/modelcontextprotocol/go-sdk/mcp"
)

type MyToolArgs struct {
    Input string `json:"input" jsonschema:"The input value"`
}

type MyTool struct{}

func (t *MyTool) Register(s *mcp.Server) {
    mcp.AddTool(s, &mcp.Tool{
        Name:        "mytool",
        Description: "Does something useful",
    }, func(ctx context.Context, req *mcp.CallToolRequest, args MyToolArgs) (*mcp.CallToolResult, any, error) {
        return &mcp.CallToolResult{
            Content: []mcp.Content{
                &mcp.TextContent{Text: fmt.Sprintf("Input: %s", args.Input)},
            },
        }, nil, nil
    })
}
```

2. Register it in `cmd/server/main.go` by adding it to the `toolsList` slice:

```go
toolsList := []mcp.Tool{
    &tools.Add{},
    &tools.Greet{},
    tools.NewWebSearch(ddg),
    tools.NewSearchDocs(embedder),
    &tools.MyTool{},
}
```

## Adding a New Prompt

1. Create a new file in `internal/prompts/`, e.g. `internal/prompts/myprompt.go`.
2. Register it in `cmd/server/main.go` by adding it to the `promptsList` slice:

```go
promptsList := []mcp.Prompt{
    &prompts.CodeReview{},
    &prompts.MyPrompt{}, // your new prompt
}
```
