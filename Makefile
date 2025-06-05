.PHONY: build
build:
	go build -o ${GOBIN}/obsidian-mcp ./cmd/obsidian-mcp/main.go