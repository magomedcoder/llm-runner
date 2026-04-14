package main

import (
	"flag"
	"github.com/magomedcoder/gen/mcp-servers/mcp-demo/internal/mcpdemoserver"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"log"
	"net/http"
)

func main() {
	addr := flag.String("listen", "127.0.0.1:8765", "адрес HTTP (GET = SSE, POST = сообщения сессии)")
	flag.Parse()

	demo := mcpdemoserver.NewServer()
	h := mcp.NewSSEHandler(func(*http.Request) *mcp.Server {
		return demo
	}, nil)

	log.Printf("MCP SSE demo: transport=sse url=http://%s/", *addr)
	log.Fatal(http.ListenAndServe(*addr, h))
}
