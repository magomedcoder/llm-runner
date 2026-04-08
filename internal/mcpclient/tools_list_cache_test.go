package mcpclient

import (
	"testing"
	"time"

	"github.com/magomedcoder/gen/internal/domain"
)

func TestServerConfigFingerprintChangesWithCommand(t *testing.T) {
	a := &domain.MCPServer{
		Transport:      "stdio",
		Command:        "npx",
		ArgsJSON:       `["x"]`,
		EnvJSON:        `{}`,
		URL:            "",
		HeadersJSON:    `{}`,
		TimeoutSeconds: 120,
	}

	fp1 := serverConfigFingerprint(a)
	a.Command = "node"
	fp2 := serverConfigFingerprint(a)
	if fp1 == fp2 {
		t.Fatal("fingerprint should change when command changes")
	}
}

func TestToolsListCacheInvalidateServerID(t *testing.T) {
	c := NewToolsListCache()
	key := toolsCacheKey{
		id: 1,
		fp: "abc",
	}

	c.mu.Lock()
	c.entries[key] = toolsCacheEntry{until: time.Now().Add(time.Hour)}
	c.mu.Unlock()
	c.InvalidateServerID(1)
	c.mu.RLock()
	_, ok := c.entries[key]
	c.mu.RUnlock()

	if ok {
		t.Fatal("ожидаемый ключ удален\n")
	}
}
