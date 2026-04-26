package main

import (
	"flag"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/magomedcoder/gen/mcp-servers/mcp-wiki/internal/mcpwikiserver"
	"github.com/magomedcoder/gen/pkg/mcpsafe"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	addr := flag.String("listen", "127.0.0.1:8772", "адрес HTTP (POST JSON-RPC, при необходимости GET SSE)")
	wikiDir := flag.String("wiki-dir", "", "обязательный каталог wiki для index_wiki_folder (единственный корень индексации)")
	flag.Parse()
	if strings.TrimSpace(*wikiDir) == "" {
		log.Fatal("mcp-wiki-streamable: обязателен флаг -wiki-dir (каталог wiki)")
	}

	var (
		mu      sync.Mutex
		servers = map[string]*mcp.Server{}
	)
	getServer := func(dir string) *mcp.Server {
		key := strings.TrimSpace(dir)
		mu.Lock()
		defer mu.Unlock()

		if srv, ok := servers[key]; ok {
			return srv
		}

		srv := mcpwikiserver.NewServerWithOptions(mcpwikiserver.Options{
			DefaultDirectory: key,
		})

		servers[key] = srv

		return srv
	}

	handler := mcp.NewStreamableHTTPHandler(func(_ *http.Request) *mcp.Server {
		return getServer(strings.TrimSpace(*wikiDir))
	}, nil)

	log.Printf("MCP wiki server (streamable): transport=streamable url=http://%s/ default_wiki_dir=%q", *addr, strings.TrimSpace(*wikiDir))
	log.Fatal(http.ListenAndServe(*addr, mcpsafe.RecoverPanic("mcp-wiki-streamable", handler)))
}
