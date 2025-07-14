package prompts

import (
	"context"
	"errors"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	obsidianrest "github.com/orvice/obsidian-mcp/pkg/obsidian-rest"
)

type ObsidianPromptServer struct {
	client *obsidianrest.Client
}

func NewObsidianPromptServer(client *obsidianrest.Client) *ObsidianPromptServer {
	return &ObsidianPromptServer{
		client: client,
	}
}

// RegisterPrompts registers all Obsidian prompts with the MCP server
func RegisterPrompts(server *mcp.Server) {
	// Register obsidian_note_summarizer prompt
	server.AddPrompt(&mcp.Prompt{
		Name:        "obsidian_note_summarizer",
		Description: "Summarize the content of an Obsidian note",
		Arguments: []*mcp.PromptArgument{
			{
				Name:        "path",
				Description: "path to the note to summarize",
				Required:    true,
			},
		},
	}, NoteSummarizerHandler)

	// Register obsidian_note_analyzer prompt
	server.AddPrompt(&mcp.Prompt{
		Name:        "obsidian_note_analyzer",
		Description: "Analyze the structure and content of an Obsidian note",
		Arguments: []*mcp.PromptArgument{
			{
				Name:        "path",
				Description: "path to the note to analyze",
				Required:    true,
			},
			{
				Name:        "analysis_type",
				Description: "type of analysis (structure, content, links, tags)",
				Required:    true,
			},
		},
	}, NoteAnalyzerHandler)

	// Register obsidian_vault_overview prompt
	server.AddPrompt(&mcp.Prompt{
		Name:        "obsidian_vault_overview",
		Description: "Generate an overview of the Obsidian vault structure",
	}, VaultOverviewHandler)
}

func NoteSummarizerHandler(ctx context.Context, session *mcp.ServerSession, params *mcp.GetPromptParams) (*mcp.GetPromptResult, error) {
	path, ok := params.Arguments["path"]
	if !ok {
		return nil, errors.New("path argument is required")
	}

	// For now, create a static prompt since we don't have client access here
	// In a real implementation, you might want to pass the client or fetch content
	promptText := "Please summarize the following Obsidian note:\n\n" +
		"**File Path:** " + path + "\n\n" +
		"Please provide a concise summary highlighting the main points, key concepts, and any important links or references."

	return &mcp.GetPromptResult{
		Description: "Summarize the content of note: " + path,
		Messages: []*mcp.PromptMessage{
			{
				Role:    mcp.Role("user"),
				Content: &mcp.TextContent{Text: promptText},
			},
		},
	}, nil
}

func NoteAnalyzerHandler(ctx context.Context, session *mcp.ServerSession, params *mcp.GetPromptParams) (*mcp.GetPromptResult, error) {
	path, ok := params.Arguments["path"]
	if !ok {
		return nil, errors.New("path argument is required")
	}

	analysisType, ok := params.Arguments["analysis_type"]
	if !ok {
		return nil, errors.New("analysis_type argument is required")
	}

	var promptText string
	switch analysisType {
	case "structure":
		promptText = "Please analyze the structure of the following Obsidian note:\n\n" +
			"**File Path:** " + path + "\n\n" +
			"Focus on: headings hierarchy, sections organization, and overall document structure."
	case "content":
		promptText = "Please analyze the content of the following Obsidian note:\n\n" +
			"**File Path:** " + path + "\n\n" +
			"Focus on: main themes, key concepts, arguments, and conclusions."
	case "links":
		promptText = "Please analyze the links and references in the following Obsidian note:\n\n" +
			"**File Path:** " + path + "\n\n" +
			"Focus on: internal links [[]], external links, backlinks potential, and connection patterns."
	case "tags":
		promptText = "Please analyze the tags and metadata in the following Obsidian note:\n\n" +
			"**File Path:** " + path + "\n\n" +
			"Focus on: existing tags, suggested tags, metadata, and categorization."
	default:
		return nil, errors.New("invalid analysis_type: must be one of structure, content, links, tags")
	}

	return &mcp.GetPromptResult{
		Description: "Analyze " + analysisType + " of note: " + path,
		Messages: []*mcp.PromptMessage{
			{
				Role:    mcp.Role("user"),
				Content: &mcp.TextContent{Text: promptText},
			},
		},
	}, nil
}

func VaultOverviewHandler(ctx context.Context, session *mcp.ServerSession, params *mcp.GetPromptParams) (*mcp.GetPromptResult, error) {
	promptText := "Please provide an overview of the current Obsidian vault structure. " +
		"Analyze the organization, main categories, note relationships, and suggest improvements for " +
		"better knowledge management. Consider the following aspects:\n\n" +
		"1. Folder structure and organization\n" +
		"2. Note naming conventions\n" +
		"3. Tag usage patterns\n" +
		"4. Link density and connection quality\n" +
		"5. Content categories and themes\n" +
		"6. Suggestions for improvement"

	return &mcp.GetPromptResult{
		Description: "Generate an overview of the Obsidian vault structure",
		Messages: []*mcp.PromptMessage{
			{
				Role:    mcp.Role("user"),
				Content: &mcp.TextContent{Text: promptText},
			},
		},
	}, nil
}
