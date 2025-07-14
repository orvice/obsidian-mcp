package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/orvice/obsidian-mcp/internal/prompts"
	"github.com/orvice/obsidian-mcp/internal/tools"
	obsidianrest "github.com/orvice/obsidian-mcp/pkg/obsidian-rest"
)

const (
	version = "0.0.1"
)

func main() {
	go func() {
		if err := serveStdio(); err != nil {
			fmt.Printf("Server error: %v\n", err)
		}
	}()

	if os.Getenv("SSE_SERVER") == "true" {
		if err := serveSSEServer(); err != nil {
			fmt.Printf("Server error: %v\n", err)
		}
	}
	select {}
}

func NewMCPServer() *mcp.Server {
	// Create MCP server using official SDK
	implementation := &mcp.Implementation{
		Name:    "ObsidianMCP",
		Version: version,
	}

	server := mcp.NewServer(implementation, nil)

	// Create Obsidian client
	client := obsidianrest.NewClient(os.Getenv("OBSIDIAN_REST_URL"), os.Getenv("OBSIDIAN_API_KEY"))

	// Register tools
	tools.RegisterTools(server, client)

	// Register prompts
	prompts.RegisterPrompts(server)

	return server
}

func serveStdio() error {
	ctx := context.Background()
	server := NewMCPServer()

	// Create stdio transport
	transport := mcp.NewStdioTransport()

	// Run server using stdio
	return server.Run(ctx, transport)
}

func serveSSEServer() error {
	mux := http.NewServeMux()
	// Add SSE endpoint here when needed
	return http.ListenAndServe(":8000", mux)
}
