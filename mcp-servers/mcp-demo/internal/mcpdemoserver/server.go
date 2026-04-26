package mcpdemoserver

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/magomedcoder/gen/pkg/mcpsafe"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func NewServer() *mcp.Server {
	srv := mcp.NewServer(&mcp.Implementation{
		Name:    "demo",
		Version: "0.1.0",
	}, nil)

	mcp.AddTool(srv, &mcp.Tool{
		Name:        "ping",
		Description: "Проверка связи; возвращает pong",
	}, func(_ context.Context, _ *mcp.CallToolRequest, _ any) (*mcp.CallToolResult, any, error) {
		return mcpsafe.SafeToolInvoke("mcp-demo", "ping", func() (*mcp.CallToolResult, any, error) {
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{
						Text: "pong",
					},
				},
			}, nil, nil
		})
	})

	type greetArgs struct {
		Name string `json:"name" jsonschema:"имя для приветствия"`
	}

	mcp.AddTool(srv, &mcp.Tool{
		Name:        "greet",
		Description: "Персональное приветствие",
	}, func(_ context.Context, _ *mcp.CallToolRequest, args greetArgs) (*mcp.CallToolResult, any, error) {
		return mcpsafe.SafeToolInvoke("mcp-demo", "greet", func() (*mcp.CallToolResult, any, error) {
			name := strings.TrimSpace(args.Name)
			if name == "" {
				name = "мир"
			}

			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{
						Text: "Привет, " + name + "!",
					},
				},
			}, nil, nil
		})
	})

	type addArgs struct {
		A int `json:"a" jsonschema:"первое слагаемое"`
		B int `json:"b" jsonschema:"второе слагаемое"`
	}
	mcp.AddTool(srv, &mcp.Tool{
		Name:        "add",
		Description: "Сложить два целых числа",
	}, func(_ context.Context, _ *mcp.CallToolRequest, args addArgs) (*mcp.CallToolResult, any, error) {
		return mcpsafe.SafeToolInvoke("mcp-demo", "add", func() (*mcp.CallToolResult, any, error) {
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{
						Text: fmt.Sprintf("%d", args.A+args.B),
					},
				},
			}, nil, nil
		})
	})

	type httpGetArgs struct {
		URL string `json:"url" jsonschema:"полный https URL для GET"`
	}
	mcp.AddTool(srv, &mcp.Tool{
		Name:        "http_get",
		Description: "GET по https; тело ответа до 64 KiB. В продакшене задайте allowlist хостов - иначе риск SSRF.",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, args httpGetArgs) (*mcp.CallToolResult, any, error) {
		return mcpsafe.SafeToolInvoke("mcp-demo", "http_get", func() (*mcp.CallToolResult, any, error) {
			u, err := url.Parse(strings.TrimSpace(args.URL))
			if err != nil || u.Scheme != "https" || u.Host == "" {
				return nil, nil, fmt.Errorf("нужен корректный https URL с хостом")
			}

			req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
			if err != nil {
				return nil, nil, err
			}

			client := &http.Client{Timeout: 15 * time.Second}
			resp, err := client.Do(req)
			if err != nil {
				return nil, nil, err
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(io.LimitReader(resp.Body, 64<<10))
			if err != nil {
				return nil, nil, err
			}

			text, err := json.Marshal(map[string]any{
				"status": resp.StatusCode,
				"body":   string(body),
			})

			if err != nil {
				return nil, nil, err
			}

			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: string(text)},
				},
			}, nil, nil
		})
	})

	return srv
}
