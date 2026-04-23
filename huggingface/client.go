package huggingface

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path"
	"strings"
	"time"
)

type ModelInfo struct {
	ModelID  string `json:"modelId"`
	Private  bool   `json:"private"`
	Gated    bool   `json:"gated"`
	Siblings []File `json:"siblings"`
}

type File struct {
	Rfilename string `json:"rfilename"`
}

type Client struct {
	baseURL    string
	httpClient *http.Client
	token      string
}

func NewClient(token string) *Client {
	tr := &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DisableCompression:    true,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          10,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   30 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	return &Client{
		baseURL: "https://huggingface.co",
		httpClient: &http.Client{
			Transport: tr,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				if len(via) >= 10 {
					return fmt.Errorf("слишком много редиректов")
				}

				return nil
			},
		},
		token: strings.TrimSpace(token),
	}
}

func (c *Client) ModelInfo(ctx context.Context, repoID string) (*ModelInfo, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	url := fmt.Sprintf("%s/api/models/%s", c.baseURL, repoID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	c.setAuth(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return nil, fmt.Errorf("ответ api со статусом %d: %s", resp.StatusCode, string(body))
	}

	var info ModelInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, fmt.Errorf("разбор ответа: %w", err)
	}

	return &info, nil
}

func (info *ModelInfo) GGFFFilenames() []string {
	var out []string
	for _, s := range info.Siblings {
		if strings.ToLower(path.Ext(s.Rfilename)) == ".gguf" {
			out = append(out, s.Rfilename)
		}
	}

	return out
}

func (c *Client) Download(ctx context.Context, repoID, revision, filename string, w io.Writer, onProgress ProgressFunc) (written int64, contentLength int64, err error) {
	if ctx == nil {
		ctx = context.Background()
	}

	if revision == "" {
		revision = "main"
	}

	url := fmt.Sprintf("%s/%s/resolve/%s/%s", c.baseURL, repoID, revision, filename)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, 0, err
	}
	c.setAuth(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))

		return 0, 0, fmt.Errorf("загрузка вернула %d: %s", resp.StatusCode, string(body))
	}

	total := resp.ContentLength
	var body io.Reader = resp.Body
	if onProgress != nil {
		body = &progressReader{r: resp.Body, total: total, on: onProgress}
	}

	n, err := io.Copy(w, body)
	if err != nil {
		return n, total, err
	}

	return n, total, nil
}

func (c *Client) setAuth(req *http.Request) {
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
}
