package mcpclient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/magomedcoder/gen/internal/domain"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type DeclaredTool struct {
	Name           string
	Description    string
	ParametersJSON string
}

func parseStringSliceJSON(raw string) ([]string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, nil
	}

	var out []string
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		return nil, err
	}

	return out, nil
}

func parseStringMapJSON(raw string) (map[string]string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return map[string]string{}, nil
	}

	var out map[string]string
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		return nil, err
	}

	if out == nil {
		return map[string]string{}, nil
	}

	return out, nil
}

func mergeEnv(base []string, extra map[string]string) []string {
	if len(extra) == 0 {
		return base
	}

	out := append([]string(nil), base...)
	for k, v := range extra {
		if strings.TrimSpace(k) == "" {
			continue
		}

		out = append(out, fmt.Sprintf("%s=%s", k, v))
	}

	return out
}

type headerRoundTripper struct {
	next http.RoundTripper
	h    map[string]string
}

func (t *headerRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	for k, v := range t.h {
		if strings.TrimSpace(k) != "" && v != "" {
			req.Header.Set(k, v)
		}
	}

	return t.next.RoundTrip(req)
}

func httpClientWithHeaders(headers map[string]string) *http.Client {
	if len(headers) == 0 {
		return nil
	}

	base := http.DefaultTransport
	if base == nil {
		base = &http.Transport{}
	}

	return &http.Client{
		Transport: &headerRoundTripper{next: base, h: headers},
		Timeout:   0,
	}
}

func transportFor(ctx context.Context, srv *domain.MCPServer) (mcp.Transport, error) {
	if srv == nil {
		return nil, errors.New("nil mcp server")
	}

	tr := strings.ToLower(strings.TrimSpace(srv.Transport))
	headers, err := parseStringMapJSON(srv.HeadersJSON)
	if err != nil {
		return nil, fmt.Errorf("headers_json: %w", err)
	}

	hc := httpClientWithHeaders(headers)

	switch tr {
	case "stdio":
		cmdPath := strings.TrimSpace(srv.Command)
		if cmdPath == "" {
			return nil, errors.New("stdio: пустая команда")
		}

		args, err := parseStringSliceJSON(srv.ArgsJSON)
		if err != nil {
			return nil, fmt.Errorf("args_json: %w", err)
		}

		envExtra, err := parseStringMapJSON(srv.EnvJSON)
		if err != nil {
			return nil, fmt.Errorf("env_json: %w", err)
		}

		cmd := exec.CommandContext(ctx, cmdPath, args...)
		cmd.Env = mergeEnv(os.Environ(), envExtra)
		return &mcp.CommandTransport{Command: cmd}, nil

	case "sse":
		raw := strings.TrimSpace(srv.URL)
		if raw == "" {
			return nil, errors.New("sse: пустой url")
		}

		if err := validateHTTPMCPURL("sse", raw); err != nil {
			return nil, err
		}

		t := &mcp.SSEClientTransport{Endpoint: raw}
		if hc != nil {
			t.HTTPClient = hc
		}

		return t, nil

	case "streamable":
		raw := strings.TrimSpace(srv.URL)
		if raw == "" {
			return nil, errors.New("streamable: пустой url")
		}

		if err := validateHTTPMCPURL("streamable", raw); err != nil {
			return nil, err
		}

		t := &mcp.StreamableClientTransport{
			Endpoint:             raw,
			DisableStandaloneSSE: true,
		}

		if hc != nil {
			t.HTTPClient = hc
		}
		return t, nil

	default:
		return nil, fmt.Errorf("неизвестный transport %q", tr)
	}
}

func timeoutFor(srv *domain.MCPServer) time.Duration {
	sec := int64(srv.TimeoutSeconds)
	if sec <= 0 {
		sec = 120
	}

	if sec > 600 {
		sec = 600
	}

	return time.Duration(sec) * time.Second
}

func withSession(ctx context.Context, srv *domain.MCPServer, fn func(context.Context, *mcp.ClientSession) error) error {
	transport, err := transportFor(ctx, srv)
	if err != nil {
		return err
	}

	tctx, cancel := context.WithTimeout(ctx, timeoutFor(srv))
	defer cancel()

	cli := mcp.NewClient(&mcp.Implementation{Name: "gen", Version: "1.0.0"}, &mcp.ClientOptions{
		Capabilities: &mcp.ClientCapabilities{},
	})
	session, err := cli.Connect(tctx, transport, nil)
	if err != nil {
		return err
	}
	defer func() {
		_ = session.Close()
	}()

	return fn(tctx, session)
}

func inputSchemaToParametersJSON(schema any) string {
	if schema == nil {
		return `{"type":"object","properties":{}}`
	}

	b, err := json.Marshal(schema)
	if err != nil || len(b) == 0 || string(b) == "null" {
		return `{"type":"object","properties":{}}`
	}

	return string(b)
}

func ListTools(ctx context.Context, srv *domain.MCPServer) ([]DeclaredTool, error) {
	var out []DeclaredTool
	err := withSession(ctx, srv, func(cctx context.Context, session *mcp.ClientSession) error {
		res, err := session.ListTools(cctx, &mcp.ListToolsParams{})
		if err != nil {
			return err
		}

		for _, t := range res.Tools {
			if t == nil || strings.TrimSpace(t.Name) == "" {
				continue
			}

			out = append(out, DeclaredTool{
				Name:           t.Name,
				Description:    strings.TrimSpace(t.Description),
				ParametersJSON: inputSchemaToParametersJSON(t.InputSchema),
			})
		}

		return nil
	})

	return out, err
}

func callToolResultString(r *mcp.CallToolResult) string {
	if r == nil {
		return ""
	}

	var parts []string
	for _, c := range r.Content {
		if c == nil {
			continue
		}

		switch x := c.(type) {
		case *mcp.TextContent:
			if strings.TrimSpace(x.Text) != "" {
				parts = append(parts, x.Text)
			}
		default:
			b, err := json.Marshal(c)
			if err == nil && len(b) > 0 {
				parts = append(parts, string(b))
			}
		}
	}

	if r.StructuredContent != nil {
		b, _ := json.Marshal(r.StructuredContent)
		if len(b) > 0 {
			parts = append(parts, string(b))
		}
	}

	if len(parts) == 0 && r.IsError {
		return "ошибка инструмента (пустой ответ)"
	}

	return strings.Join(parts, "\n")
}

func CallTool(ctx context.Context, srv *domain.MCPServer, mcpToolName string, arguments json.RawMessage) (string, error) {
	if strings.TrimSpace(mcpToolName) == "" {
		return "", errors.New("пустое имя инструмента MCP")
	}

	var result string
	var callErr error
	err := withSession(ctx, srv, func(cctx context.Context, session *mcp.ClientSession) error {
		var args any
		if len(arguments) > 0 {
			if err := json.Unmarshal(arguments, &args); err != nil {
				return fmt.Errorf("аргументы инструмента: %w", err)
			}
		}

		if args == nil {
			args = map[string]any{}
		}

		res, err := session.CallTool(cctx, &mcp.CallToolParams{
			Name:      mcpToolName,
			Arguments: args,
		})
		if err != nil {
			return err
		}

		result = callToolResultString(res)
		if res != nil && res.IsError {
			callErr = errors.New(strings.TrimSpace(result))
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	if callErr != nil {
		return result, callErr
	}

	return result, nil
}
