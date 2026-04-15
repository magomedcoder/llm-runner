package mcpclient

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/magomedcoder/gen/internal/domain"
)

const DefaultToolsListCacheTTL = 2 * time.Minute

type listCacheKey struct {
	id int64
	fp string
}

type toolsCacheEntry struct {
	tools []DeclaredTool
	until time.Time
}

type resourcesCacheEntry struct {
	items []DeclaredResource
	until time.Time
}

type promptsCacheEntry struct {
	items []DeclaredPrompt
	until time.Time
}

type ToolsListCache struct {
	mu             sync.RWMutex
	toolEntries    map[listCacheKey]toolsCacheEntry
	resEntries     map[listCacheKey]resourcesCacheEntry
	promptsEntries map[listCacheKey]promptsCacheEntry
}

func NewToolsListCache() *ToolsListCache {
	return &ToolsListCache{
		toolEntries:    make(map[listCacheKey]toolsCacheEntry),
		resEntries:     make(map[listCacheKey]resourcesCacheEntry),
		promptsEntries: make(map[listCacheKey]promptsCacheEntry),
	}
}

var toolsListNotifyDefault atomic.Pointer[ToolsListCache]

func SetToolsListCacheForNotifications(c *ToolsListCache) {
	toolsListNotifyDefault.Store(c)
}

func notifyForListChangedHandlers(explicit *ToolsListCache) *ToolsListCache {
	if explicit != nil {
		return explicit
	}
	return toolsListNotifyDefault.Load()
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

func cloneDeclaredResources(in []DeclaredResource) []DeclaredResource {
	if len(in) == 0 {
		return nil
	}
	out := make([]DeclaredResource, len(in))
	copy(out, in)
	return out
}

func cloneDeclaredPrompts(in []DeclaredPrompt) []DeclaredPrompt {
	if len(in) == 0 {
		return nil
	}
	out := make([]DeclaredPrompt, len(in))
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
		return listTools(ctx, srv, nil)
	}

	if ttl <= 0 {
		ttl = DefaultToolsListCacheTTL
	}

	fp := serverConfigFingerprint(srv)
	key := listCacheKey{
		id: srv.ID,
		fp: fp,
	}
	now := time.Now()

	c.mu.RLock()
	e, ok := c.toolEntries[key]
	c.mu.RUnlock()
	if ok && now.Before(e.until) {
		recordListCacheHit()
		return cloneDeclaredTools(e.tools), nil
	}

	recordListCacheMiss()
	tools, err := listTools(ctx, srv, c)
	if err != nil {
		return nil, err
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	if c.toolEntries == nil {
		c.toolEntries = make(map[listCacheKey]toolsCacheEntry)
	}

	for k, v := range c.toolEntries {
		if !now.Before(v.until) {
			delete(c.toolEntries, k)
		}
	}

	for k := range c.toolEntries {
		if k.id == srv.ID && k.fp != fp {
			delete(c.toolEntries, k)
		}
	}

	c.toolEntries[key] = toolsCacheEntry{
		tools: cloneDeclaredTools(tools),
		until: now.Add(ttl),
	}

	return cloneDeclaredTools(tools), nil
}

func (c *ToolsListCache) ListResourcesCached(ctx context.Context, srv *domain.MCPServer, ttl time.Duration) ([]DeclaredResource, error) {
	if c == nil {
		return ListResources(ctx, srv)
	}
	if srv == nil {
		return nil, errors.New("nil mcp server")
	}
	if srv.ID <= 0 {
		return listResources(ctx, srv, nil)
	}
	if ttl <= 0 {
		ttl = DefaultToolsListCacheTTL
	}
	fp := serverConfigFingerprint(srv)
	key := listCacheKey{id: srv.ID, fp: fp}
	now := time.Now()

	c.mu.RLock()
	e, ok := c.resEntries[key]
	c.mu.RUnlock()
	if ok && now.Before(e.until) {
		recordListCacheHit()
		return cloneDeclaredResources(e.items), nil
	}

	recordListCacheMiss()
	items, err := listResources(ctx, srv, c)
	if err != nil {
		return nil, err
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	if c.resEntries == nil {
		c.resEntries = make(map[listCacheKey]resourcesCacheEntry)
	}
	for k, v := range c.resEntries {
		if !now.Before(v.until) {
			delete(c.resEntries, k)
		}
	}
	for k := range c.resEntries {
		if k.id == srv.ID && k.fp != fp {
			delete(c.resEntries, k)
		}
	}
	c.resEntries[key] = resourcesCacheEntry{
		items: cloneDeclaredResources(items),
		until: now.Add(ttl),
	}
	return cloneDeclaredResources(items), nil
}

func (c *ToolsListCache) ListPromptsCached(ctx context.Context, srv *domain.MCPServer, ttl time.Duration) ([]DeclaredPrompt, error) {
	if c == nil {
		return ListPrompts(ctx, srv)
	}
	if srv == nil {
		return nil, errors.New("nil mcp server")
	}
	if srv.ID <= 0 {
		return listPrompts(ctx, srv, nil)
	}
	if ttl <= 0 {
		ttl = DefaultToolsListCacheTTL
	}
	fp := serverConfigFingerprint(srv)
	key := listCacheKey{id: srv.ID, fp: fp}
	now := time.Now()

	c.mu.RLock()
	e, ok := c.promptsEntries[key]
	c.mu.RUnlock()
	if ok && now.Before(e.until) {
		recordListCacheHit()
		return cloneDeclaredPrompts(e.items), nil
	}

	recordListCacheMiss()
	items, err := listPrompts(ctx, srv, c)
	if err != nil {
		return nil, err
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	if c.promptsEntries == nil {
		c.promptsEntries = make(map[listCacheKey]promptsCacheEntry)
	}
	for k, v := range c.promptsEntries {
		if !now.Before(v.until) {
			delete(c.promptsEntries, k)
		}
	}
	for k := range c.promptsEntries {
		if k.id == srv.ID && k.fp != fp {
			delete(c.promptsEntries, k)
		}
	}
	c.promptsEntries[key] = promptsCacheEntry{
		items: cloneDeclaredPrompts(items),
		until: now.Add(ttl),
	}
	return cloneDeclaredPrompts(items), nil
}

func (c *ToolsListCache) InvalidateServerTools(id int64) {
	if c == nil || id <= 0 {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	for k := range c.toolEntries {
		if k.id == id {
			delete(c.toolEntries, k)
		}
	}
}

func (c *ToolsListCache) InvalidateServerResources(id int64) {
	if c == nil || id <= 0 {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	for k := range c.resEntries {
		if k.id == id {
			delete(c.resEntries, k)
		}
	}
}

func (c *ToolsListCache) InvalidateServerPrompts(id int64) {
	if c == nil || id <= 0 {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	for k := range c.promptsEntries {
		if k.id == id {
			delete(c.promptsEntries, k)
		}
	}
}

func (c *ToolsListCache) InvalidateServerID(id int64) {
	if id <= 0 {
		return
	}
	closePooledHTTPSession(id)
	if c == nil {
		return
	}
	c.InvalidateServerTools(id)
	c.InvalidateServerResources(id)
	c.InvalidateServerPrompts(id)
}
