package mcpclient

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/magomedcoder/gen/internal/domain"
)

const DefaultToolsListCacheTTL = 2 * time.Minute

type toolsCacheKey struct {
	id int64
	fp string
}

type toolsCacheEntry struct {
	tools []DeclaredTool
	until time.Time
}

type ToolsListCache struct {
	mu      sync.RWMutex
	entries map[toolsCacheKey]toolsCacheEntry
}

func NewToolsListCache() *ToolsListCache {
	return &ToolsListCache{
		entries: make(map[toolsCacheKey]toolsCacheEntry),
	}
}

func serverConfigFingerprint(s *domain.MCPServer) string {
	if s == nil {
		return ""
	}

	uid := ""
	if s.UserID != nil {
		uid = fmt.Sprintf("%d", *s.UserID)
	}

	payload := fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%d", uid, s.Transport, s.Command, s.ArgsJSON, s.EnvJSON, s.URL, s.HeadersJSON, s.TimeoutSeconds)
	sum := sha256.Sum256([]byte(payload))

	return hex.EncodeToString(sum[:12])
}

func cloneDeclaredTools(in []DeclaredTool) []DeclaredTool {
	if len(in) == 0 {
		return nil
	}

	out := make([]DeclaredTool, len(in))
	copy(out, in)
	return out
}

func (c *ToolsListCache) ListToolsCached(ctx context.Context, srv *domain.MCPServer, ttl time.Duration) ([]DeclaredTool, error) {
	if c == nil {
		return ListTools(ctx, srv)
	}

	if srv == nil {
		return nil, errors.New("nil mcp server")
	}

	if srv.ID <= 0 {
		return ListTools(ctx, srv)
	}

	if ttl <= 0 {
		ttl = DefaultToolsListCacheTTL
	}

	fp := serverConfigFingerprint(srv)
	key := toolsCacheKey{
		id: srv.ID,
		fp: fp,
	}
	now := time.Now()

	c.mu.RLock()
	e, ok := c.entries[key]
	c.mu.RUnlock()
	if ok && now.Before(e.until) {
		return cloneDeclaredTools(e.tools), nil
	}

	tools, err := ListTools(ctx, srv)
	if err != nil {
		return nil, err
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	if c.entries == nil {
		c.entries = make(map[toolsCacheKey]toolsCacheEntry)
	}

	for k, v := range c.entries {
		if !now.Before(v.until) {
			delete(c.entries, k)
		}
	}

	for k := range c.entries {
		if k.id == srv.ID && k.fp != fp {
			delete(c.entries, k)
		}
	}

	c.entries[key] = toolsCacheEntry{
		tools: cloneDeclaredTools(tools),
		until: now.Add(ttl),
	}

	return cloneDeclaredTools(tools), nil
}

func (c *ToolsListCache) InvalidateServerID(id int64) {
	if c == nil || id <= 0 {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	for k := range c.entries {
		if k.id == id {
			delete(c.entries, k)
		}
	}
}
