# mcpserver_example

A minimal [Model Context Protocol (MCP)](https://modelcontextprotocol.io) server written in Go using the [go-sdk](https://github.com/modelcontextprotocol/go-sdk).

## Project Structure

```
mcpserver_example/
├── cmd/
│   └── server/
│       └── main.go        # Entry point
├── internal/
│   ├── mcp/
│   │   ├── handler.go     # Tool interface and Handler
│   │   └── server.go      # MCP server wrapper
│   └── tools/
│       ├── add.go         # "add" tool
│       └── greet.go       # "greet" tool
├── go.mod
└── go.sum
```

## Tools

| Tool    | Description              | Arguments              |
|---------|--------------------------|------------------------|
| `add`   | Add two numbers together | `a` (float), `b` (float) |
| `greet` | Greet someone by name    | `name` (string)        |

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

The server communicates over stdio using the MCP protocol.

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

2. Register it in `cmd/server/main.go`:

```go
h := mcpserver.NewHandler(&tools.Add{}, &tools.Greet{}, &tools.MyTool{})
```
