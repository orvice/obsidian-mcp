package tools

import (
	"context"
	"errors"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	obsidianrest "github.com/orvice/obsidian-mcp/pkg/obsidian-rest"
)

type Tool struct {
	Name        string
	Description string
	// Handler     any

	Tool    mcp.Tool
	Handler server.ToolHandlerFunc
}

type ObsidianToolServer struct {
	client *obsidianrest.Client
}

func NewObsidianToolServer(client *obsidianrest.Client) *ObsidianToolServer {
	return &ObsidianToolServer{
		client: client,
	}
}

func (s *ObsidianToolServer) Tools() []Tool {
	return []Tool{
		{
			Tool: mcp.NewTool("obsidian_get_note",
				mcp.WithDescription("Get Obsidian Note"),
				mcp.WithString("path",
					mcp.Description("path to the note"),
					mcp.Required(),
				),
			),
			Handler: s.GetNote,
		},
		{
			Tool: mcp.NewTool("obsidian_update_note",
				mcp.WithDescription("Update Obsidian Note Content"),
				mcp.WithString("path",
					mcp.Description("path to the note"),
					mcp.Required(),
				),
				mcp.WithString("content",
					mcp.Description("new content for the note"),
					mcp.Required(),
				),
			),
			Handler: s.UpdateNote,
		},
	}
}

func (s *ObsidianToolServer) GetNote(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	path, ok := args["path"].(string)
	if !ok {
		return nil, errors.New("path must be a string")
	}

	note, err := s.client.GetVaultFile(path)
	if err != nil {
		return nil, err
	}
	return mcp.NewToolResultText(note.Content), nil
}

func (s *ObsidianToolServer) UpdateNote(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	path, ok := args["path"].(string)
	if !ok {
		return nil, errors.New("path must be a string")
	}

	content, ok := args["content"].(string)
	if !ok {
		return nil, errors.New("content must be a string")
	}

	err := s.client.UpdateVaultFile(path, content)
	if err != nil {
		return nil, err
	}
	return mcp.NewToolResultText("Note updated successfully"), nil
}
