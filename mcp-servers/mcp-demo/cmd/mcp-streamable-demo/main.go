package main

import (
	"flag"
	"github.com/magomedcoder/gen/mcp-servers/mcp-demo/internal/mcpdemoserver"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"log"
	"net/http"
)

func main() {
	addr := flag.String("listen", "127.0.0.1:8766", "адрес HTTP (POST JSON-RPC, при необходимости GET SSE)")
	flag.Parse()

	demo := mcpdemoserver.NewServer()
	h := mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server {
		return demo
	}, nil)

	log.Printf("MCP streamable demo: transport=streamable url=http://%s/", *addr)
	log.Fatal(http.ListenAndServe(*addr, h))
}
