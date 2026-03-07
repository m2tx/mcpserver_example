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

## Architecture

Each tool implements the `Tool` interface:

```go
type Tool interface {
    Register(s *mcp.Server)
}
```

Tools self-register with the MCP server via `mcp.AddTool`. The `Handler` collects tools and the `Server` wires everything together over stdio transport.

## Tools

| Tool    | Description              | Arguments                  |
|---------|--------------------------|----------------------------|
| `add`   | Add two numbers together | `a` (float), `b` (float)   |
| `greet` | Greet someone by name    | `name` (string)            |

## Getting Started

### Prerequisites

- Go 1.25+

### Build

```bash
go build ./cmd/server
```

### Run

**stdio mode** (default):

```bash
./server
```

**HTTP mode**:

```bash
HTTP_PORT=8080 ./server
```

When `HTTP_PORT` is set, the server listens for HTTP connections on that port using the streamable HTTP transport. Otherwise it communicates over stdio.

## Using with MCP Clients

### Claude Desktop

Add the following to your `claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "example": {
      "command": "/path/to/server"
    }
  }
}
```

### Claude Code

```bash
claude mcp add example /path/to/server
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

2. Register it in `cmd/server/main.go`:

```go
h := mcpserver.NewHandler(&tools.Add{}, &tools.Greet{}, &tools.MyTool{})
```
