package tools

import (
	"context"
	"errors"

	"github.com/modelcontextprotocol/go-sdk/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	obsidianrest "github.com/orvice/obsidian-mcp/pkg/obsidian-rest"
)

type ObsidianToolServer struct {
	client *obsidianrest.Client
}

func NewObsidianToolServer(client *obsidianrest.Client) *ObsidianToolServer {
	return &ObsidianToolServer{
		client: client,
	}
}

// RegisterTools registers all Obsidian tools with the MCP server
func RegisterTools(server *mcp.Server, client *obsidianrest.Client) {
	toolServer := NewObsidianToolServer(client)

	// Register obsidian_get_note tool
	server.AddTool(&mcp.Tool{
		Name:        "obsidian_get_note",
		Description: "Get Obsidian Note",
		InputSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"path": {
					Type:        "string",
					Description: "path to the note",
				},
			},
			Required: []string{"path"},
		},
	}, toolServer.GetNote)

	// Register obsidian_update_note tool
	server.AddTool(&mcp.Tool{
		Name:        "obsidian_update_note",
		Description: "Update Obsidian Note Content",
		InputSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"path": {
					Type:        "string",
					Description: "path to the note",
				},
				"content": {
					Type:        "string",
					Description: "new content for the note",
				},
			},
			Required: []string{"path", "content"},
		},
	}, toolServer.UpdateNote)
}

func (s *ObsidianToolServer) GetNote(ctx context.Context, session *mcp.ServerSession, params *mcp.CallToolParamsFor[map[string]any]) (*mcp.CallToolResultFor[any], error) {
	path, ok := params.Arguments["path"].(string)
	if !ok {
		return nil, errors.New("path must be a string")
	}

	note, err := s.client.GetVaultFile(path)
	if err != nil {
		return &mcp.CallToolResultFor[any]{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{Text: "Error: " + err.Error()},
			},
		}, nil
	}

	return &mcp.CallToolResultFor[any]{
		IsError: false,
		Content: []mcp.Content{
			&mcp.TextContent{Text: note.Content},
		},
	}, nil
}

func (s *ObsidianToolServer) UpdateNote(ctx context.Context, session *mcp.ServerSession, params *mcp.CallToolParamsFor[map[string]any]) (*mcp.CallToolResultFor[any], error) {
	path, ok := params.Arguments["path"].(string)
	if !ok {
		return nil, errors.New("path must be a string")
	}

	content, ok := params.Arguments["content"].(string)
	if !ok {
		return nil, errors.New("content must be a string")
	}

	err := s.client.UpdateVaultFile(path, content)
	if err != nil {
		return &mcp.CallToolResultFor[any]{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{Text: "Error: " + err.Error()},
			},
		}, nil
	}

	return &mcp.CallToolResultFor[any]{
		IsError: false,
		Content: []mcp.Content{
			&mcp.TextContent{Text: "Note updated successfully"},
		},
	}, nil
}
