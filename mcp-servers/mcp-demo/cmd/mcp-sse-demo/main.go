package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/magomedcoder/gen/mcp-servers/mcp-demo/internal/mcpdemoserver"
	"github.com/magomedcoder/gen/pkg/mcpsafe"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	addr := flag.String("listen", "127.0.0.1:8765", "адрес HTTP (GET = SSE, POST = сообщения сессии)")
	flag.Parse()

	demo := mcpdemoserver.NewServer()
	h := mcp.NewSSEHandler(func(*http.Request) *mcp.Server {
		return demo
	}, nil)

	log.Printf("MCP SSE demo: transport=sse url=http://%s/", *addr)
	log.Fatal(http.ListenAndServe(*addr, mcpsafe.RecoverPanic("mcp-demo-sse", h)))
}
