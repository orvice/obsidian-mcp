package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/mark3labs/mcp-go/server"
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

func NewMCPServer() *server.MCPServer {
	// Create MCP server
	s := server.NewMCPServer(
		"ObsidianMCP",
		version,
	)

	// Create Obsidian client
	client := obsidianrest.NewClient(os.Getenv("OBSIDIAN_BASE_URL"), os.Getenv("OBSIDIAN_API_KEY"),
		obsidianrest.WithInsecureSkipVerify(true))

	// Register tools
	obsidianToolServer := tools.NewObsidianToolServer(client)
	for _, tool := range obsidianToolServer.Tools() {
		s.AddTool(tool.Tool, tool.Handler)
	}

	// Register prompts
	obsidianPromptServer := prompts.NewObsidianPromptServer(client)
	for _, prompt := range obsidianPromptServer.Prompts() {
		s.AddPrompt(prompt.Prompt, prompt.Handler)
	}

	return s
}

func serveStdio() error {
	s := NewMCPServer()
	// Start the stdio server
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
	return nil
}

func serveSSEServer() error {
	// Create MCP server
	s := server.NewSSEServer(
		NewMCPServer(),
	)

	mux := http.NewServeMux()
	mux.HandleFunc("/sse", s.ServeHTTP)

	server := &http.Server{
		Addr:    ":8081",
		Handler: mux,
	}
	return server.ListenAndServe()
}
