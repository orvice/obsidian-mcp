# obsidian-mcp

A Model Context Protocol (MCP) server for Obsidian, written in Go.

## Features

- Interact with Obsidian through the MCP protocol
- Support for Obsidian REST API
- Provides tools and prompts functionality
- Supports both stdio and SSE server modes

## Prerequisites

This MCP server requires the [Obsidian Local REST API plugin](https://github.com/coddingtonbear/obsidian-local-rest-api) to be installed and enabled in Obsidian. This plugin provides a secure HTTPS interface that allows external tools to interact with your Obsidian notes through a REST API.

## Installation

### Using go install

```bash
go install github.com/orvice/obsidian-mcp/cmd/obsidian-mcp@latest
```

### Build from source

```bash
git clone https://github.com/orvice/obsidian-mcp.git
cd obsidian-mcp
go build -o obsidian-mcp ./cmd/obsidian-mcp
```

## Configuration

### Environment Variables

- `OBSIDIAN_BASE_URL`: Base URL for Obsidian REST API
- `OBSIDIAN_API_KEY`: API key for Obsidian REST API
- `SSE_SERVER`: Set to "true" to enable SSE server mode (optional)

### MCP Configuration Examples

#### Claude Desktop Configuration

Add to `~/Library/Application Support/Claude/claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "obsidian-mcp": {
      "command": "obsidian-mcp",
      "env": {
        "OBSIDIAN_BASE_URL": "http://localhost:27123",
        "OBSIDIAN_API_KEY": "your-obsidian-api-key"
      }
    }
  }
}
```

#### Continue Configuration

Add to `~/.continue/config.json`:

```json
{
  "experimental": {
    "modelContextProtocol": true
  },
  "mcpServers": {
    "obsidian-mcp": {
      "command": "obsidian-mcp",
      "args": [],
      "env": {
        "OBSIDIAN_BASE_URL": "http://localhost:27123",
        "OBSIDIAN_API_KEY": "your-obsidian-api-key"
      }
    }
  }
}
```

## Usage

1. Ensure Obsidian is installed with the REST API plugin enabled
2. Obtain an API key and configure environment variables
3. Start the MCP server
4. Use in MCP-compatible clients (such as Claude Desktop)

## Development

```bash
# Run the project
go run ./cmd/obsidian-mcp

# Run tests
go test ./...

# Build
make build
```

## License

[MIT License](LICENSE)