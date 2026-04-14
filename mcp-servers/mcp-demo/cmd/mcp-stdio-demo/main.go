package main

import (
	"context"
	"github.com/magomedcoder/gen/mcp-servers/mcp-demo/internal/mcpdemoserver"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"log"
)

func main() {
	srv := mcpdemoserver.NewServer()
	if err := srv.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatalf("mcp server: %v", err)
	}
}
