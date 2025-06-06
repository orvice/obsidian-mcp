package prompts

import (
	"context"
	"errors"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	obsidianrest "github.com/orvice/obsidian-mcp/pkg/obsidian-rest"
)

type Prompt struct {
	Name        string
	Description string
	Prompt      mcp.Prompt
	Handler     server.PromptHandlerFunc
}

type ObsidianPromptServer struct {
	client *obsidianrest.Client
}

func NewObsidianPromptServer(client *obsidianrest.Client) *ObsidianPromptServer {
	return &ObsidianPromptServer{
		client: client,
	}
}

func (s *ObsidianPromptServer) Prompts() []Prompt {
	return []Prompt{
		{
			Prompt: mcp.NewPrompt("obsidian_note_summarizer",
				mcp.WithPromptDescription("Summarize the content of an Obsidian note"),
				mcp.WithArgument("path",
					mcp.ArgumentDescription("path to the note to summarize"),
					mcp.RequiredArgument(),
				),
			),
			Handler: s.NoteSummarizer,
		},
		{
			Prompt: mcp.NewPrompt("obsidian_note_analyzer",
				mcp.WithPromptDescription("Analyze the structure and content of an Obsidian note"),
				mcp.WithArgument("path",
					mcp.ArgumentDescription("path to the note to analyze"),
					mcp.RequiredArgument(),
				),
				mcp.WithArgument("analysis_type",
					mcp.ArgumentDescription("type of analysis (structure, content, links, tags)"),
					mcp.RequiredArgument(),
				),
			),
			Handler: s.NoteAnalyzer,
		},
		{
			Prompt: mcp.NewPrompt("obsidian_vault_overview",
				mcp.WithPromptDescription("Generate an overview of the Obsidian vault structure"),
			),
			Handler: s.VaultOverview,
		},
	}
}

func (s *ObsidianPromptServer) NoteSummarizer(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	args := request.Params.Arguments
	path, ok := args["path"]
	if !ok {
		return nil, errors.New("path argument is required")
	}

	note, err := s.client.GetVaultFile(path)
	if err != nil {
		return nil, err
	}

	promptText := "Please summarize the following Obsidian note:\n\n" +
		"**File Path:** " + path + "\n\n" +
		"**Content:**\n" + note.Content + "\n\n" +
		"Please provide a concise summary highlighting the main points, key concepts, and any important links or references."

	return mcp.NewGetPromptResult(
		"Summarize the content of note: "+path,
		[]mcp.PromptMessage{
			mcp.NewPromptMessage(
				mcp.RoleUser,
				mcp.NewTextContent(promptText),
			),
		},
	), nil
}

func (s *ObsidianPromptServer) NoteAnalyzer(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	args := request.Params.Arguments
	path, ok := args["path"]
	if !ok {
		return nil, errors.New("path argument is required")
	}

	analysisType, ok := args["analysis_type"]
	if !ok {
		return nil, errors.New("analysis_type argument is required")
	}

	note, err := s.client.GetVaultFile(path)
	if err != nil {
		return nil, err
	}

	var promptText string
	switch analysisType {
	case "structure":
		promptText = "Please analyze the structure of the following Obsidian note:\n\n" +
			"**File Path:** " + path + "\n\n" +
			"**Content:**\n" + note.Content + "\n\n" +
			"Focus on: headings hierarchy, sections organization, and overall document structure."
	case "content":
		promptText = "Please analyze the content of the following Obsidian note:\n\n" +
			"**File Path:** " + path + "\n\n" +
			"**Content:**\n" + note.Content + "\n\n" +
			"Focus on: main themes, key concepts, arguments, and conclusions."
	case "links":
		promptText = "Please analyze the links and references in the following Obsidian note:\n\n" +
			"**File Path:** " + path + "\n\n" +
			"**Content:**\n" + note.Content + "\n\n" +
			"Focus on: internal links [[]], external links, backlinks potential, and connection patterns."
	case "tags":
		promptText = "Please analyze the tags and metadata in the following Obsidian note:\n\n" +
			"**File Path:** " + path + "\n\n" +
			"**Content:**\n" + note.Content + "\n\n" +
			"Focus on: existing tags, suggested tags, metadata, and categorization."
	default:
		return nil, errors.New("invalid analysis_type: must be one of structure, content, links, tags")
	}

	return mcp.NewGetPromptResult(
		"Analyze "+analysisType+" of note: "+path,
		[]mcp.PromptMessage{
			mcp.NewPromptMessage(
				mcp.RoleUser,
				mcp.NewTextContent(promptText),
			),
		},
	), nil
}

func (s *ObsidianPromptServer) VaultOverview(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	// For now, we'll provide a general prompt for vault overview
	// In a future enhancement, we could add an API endpoint to get vault structure
	promptText := "Please provide an overview of the current Obsidian vault structure. " +
		"Analyze the organization, main categories, note relationships, and suggest improvements for " +
		"better knowledge management. Consider the following aspects:\n\n" +
		"1. Folder structure and organization\n" +
		"2. Note naming conventions\n" +
		"3. Tag usage patterns\n" +
		"4. Link density and connection quality\n" +
		"5. Content categories and themes\n" +
		"6. Suggestions for improvement"

	return mcp.NewGetPromptResult(
		"Generate an overview of the Obsidian vault structure",
		[]mcp.PromptMessage{
			mcp.NewPromptMessage(
				mcp.RoleUser,
				mcp.NewTextContent(promptText),
			),
		},
	), nil
}
