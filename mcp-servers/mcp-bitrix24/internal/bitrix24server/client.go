package bitrix24server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
	"unicode/utf8"

	"golang.org/x/text/encoding/charmap"
)

type bitrixClient struct {
	baseURL    *url.URL
	httpClient *http.Client
}

func newBitrixClient(baseURL string, timeout time.Duration) (*bitrixClient, error) {
	trimmed := strings.TrimSpace(baseURL)
	if trimmed == "" {
		return nil, fmt.Errorf("B24_WEBHOOK_BASE is empty")
	}

	trimmed = strings.TrimSuffix(trimmed, "/")
	u, err := url.Parse(trimmed)
	if err != nil {
		return nil, fmt.Errorf("invalid B24_WEBHOOK_BASE: %w", err)
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return nil, fmt.Errorf("B24_WEBHOOK_BASE must start with http:// or https://")
	}

	if u.Host == "" {
		return nil, fmt.Errorf("B24_WEBHOOK_BASE host is empty")
	}

	return &bitrixClient{
		baseURL: u,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}, nil
}

func (c *bitrixClient) call(ctx context.Context, method string, payload any) (map[string]any, error) {
	method = strings.TrimSpace(method)
	if method == "" {
		return nil, fmt.Errorf("method is empty")
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal payload: %w", err)
	}

	endpoint := *c.baseURL
	endpoint.Path = strings.TrimSuffix(endpoint.Path, "/") + "/" + method
	log.Printf("[b24-mcp] http method=%q url=%s payload_bytes=%d", method, endpoint.String(), len(body))
	log.Printf("[b24-mcp] http method=%q request_body=%s", method, truncateForLog(prettyJSONString(body), 3000))

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint.String(), bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Printf("[b24-mcp] http method=%q request_err=%v", method, err)
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(io.LimitReader(resp.Body, 5<<20))
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}
	raw = normalizeResponseEncoding(raw)
	log.Printf("[b24-mcp] http method=%q status=%d response_bytes=%d", method, resp.StatusCode, len(raw))
	log.Printf("[b24-mcp] http method=%q response_body=%s", method, truncateForLog(prettyJSONString(raw), 3000))

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		log.Printf("[b24-mcp] http method=%q non_2xx status=%d body=%s", method, resp.StatusCode, truncateForLog(string(raw), 600))
		return nil, fmt.Errorf("bitrix http status %d: %s", resp.StatusCode, string(raw))
	}

	var response map[string]any
	if err := json.Unmarshal(raw, &response); err != nil {
		log.Printf("[b24-mcp] http method=%q decode_json_err=%v body=%s", method, err, truncateForLog(string(raw), 600))
		return nil, fmt.Errorf("decode json: %w", err)
	}

	if errObj, ok := response["error"]; ok {
		desc, _ := response["error_description"].(string)
		log.Printf("[b24-mcp] http method=%q api_error=%v desc=%q", method, errObj, desc)
		return nil, fmt.Errorf("bitrix api error: %v (%s)", errObj, desc)
	}

	log.Printf("[b24-mcp] http method=%q ok", method)

	return response, nil
}

func (c *bitrixClient) callTaskCommentItemGetList(ctx context.Context, taskID int, order map[string]any, filter map[string]any) (map[string]any, error) {
	payload := struct {
		TaskID int            `json:"TASKID"`
		Order  map[string]any `json:"ORDER,omitempty"`
		Filter map[string]any `json:"FILTER,omitempty"`
	}{
		TaskID: taskID,
		Order:  order,
		Filter: filter,
	}
	return c.call(ctx, "task.commentitem.getlist", payload)
}

func normalizeResponseEncoding(raw []byte) []byte {
	if len(raw) == 0 || utf8.Valid(raw) {
		return raw
	}

	decoded, err := charmap.Windows1251.NewDecoder().Bytes(raw)
	if err != nil || !utf8.Valid(decoded) || !json.Valid(decoded) {
		return raw
	}

	return decoded
}

func truncateForLog(s string, maxRunes int) string {
	if maxRunes <= 0 || utf8.RuneCountInString(s) <= maxRunes {
		return s
	}

	r := []rune(s)
	return string(r[:maxRunes]) + "...(truncated)"
}

func prettyJSONString(raw []byte) string {
	if len(raw) == 0 {
		return ""
	}

	var parsed any
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return string(raw)
	}

	formatted, err := json.MarshalIndent(parsed, "", "  ")
	if err != nil {
		return string(raw)
	}

	return string(formatted)
}
