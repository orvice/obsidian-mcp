FROM golang:1.24 as builder

WORKDIR /app

# 复制 go.mod 和 go.sum 并下载依赖
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o /app/bin/obsidian-mcp ./cmd/obsidian-mcp/main.go

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/bin/obsidian-mcp /app/bin/obsidian-mcp

CMD ["/app/bin/obsidian-mcp"] 