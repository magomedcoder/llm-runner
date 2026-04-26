package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/magomedcoder/gen/mcp-servers/mcp-bitrix24/internal/bitrix24server"
	"github.com/magomedcoder/gen/pkg/mcpsafe"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	addr := flag.String("listen", "127.0.0.1:8786", "адрес HTTP (POST JSON-RPC, при необходимости GET SSE)")
	flag.Parse()
	webhookBase := os.Getenv("B24_WEBHOOK_BASE")
	log.Printf("MCP Bitrix24 streamable: starting webhook_base_set=%t listen=%s", webhookBase != "", *addr)

	srv, err := bitrix24server.NewServer(bitrix24server.Config{
		WebhookBase: webhookBase,
	})
	if err != nil {
		log.Fatalf("init bitrix24 server: %v", err)
	}

	h := mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server {
		return srv
	}, nil)

	log.Printf("MCP Bitrix24 streamable: transport=streamable url=http://%s/", *addr)
	log.Fatal(http.ListenAndServe(*addr, mcpsafe.RecoverPanic("mcp-bitrix24-streamable", h)))
}
