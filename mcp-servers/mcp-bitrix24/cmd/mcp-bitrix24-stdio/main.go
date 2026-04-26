package main

import (
	"context"
	"log"
	"os"

	"github.com/magomedcoder/gen/mcp-servers/mcp-bitrix24/internal/bitrix24server"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	webhookBase := os.Getenv("B24_WEBHOOK_BASE")
	log.Printf("MCP Bitrix24 stdio: starting webhook_base_set=%t", webhookBase != "")
	srv, err := bitrix24server.NewServer(bitrix24server.Config{
		WebhookBase: webhookBase,
	})
	if err != nil {
		log.Fatalf("init bitrix24 server: %v", err)
	}

	log.Printf("MCP Bitrix24 stdio: transport=stdio ready")
	if err := srv.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatalf("mcp server: %v", err)
	}
}
