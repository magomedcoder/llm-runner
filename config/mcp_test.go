package config

import "testing"

func TestMCPHTTPHostAllowedLoopback(t *testing.T) {
	c := &Config{MCP: MCPConfig{}}
	for _, h := range []string{"127.0.0.1", "::1"} {
		if !c.MCPHTTPHostAllowed(h) {
			t.Fatalf("loopback %q should be allowed", h)
		}
	}
}

func TestMCPHTTPHostAllowedSuffix(t *testing.T) {
	c := &Config{MCP: MCPConfig{HTTPAllowHosts: []string{"example.com"}}}
	if !c.MCPHTTPHostAllowed("api.example.com") {
		t.Fatal("suffix match expected")
	}

	if c.MCPHTTPHostAllowed("evil.com") {
		t.Fatal("unrelated host must be denied")
	}
}

func TestMCPHTTPHostAllowedAny(t *testing.T) {
	c := &Config{MCP: MCPConfig{HTTPAllowAny: true}}
	if !c.MCPHTTPHostAllowed("anywhere.example") {
		t.Fatal("allow any")
	}
}

func TestMCPHTTPHostAllowedExactIP(t *testing.T) {
	c := &Config{MCP: MCPConfig{HTTPAllowHosts: []string{"10.1.2.3"}}}
	if !c.MCPHTTPHostAllowed("10.1.2.3") {
		t.Fatal("exact IP")
	}

	if c.MCPHTTPHostAllowed("10.1.2.4") {
		t.Fatal("other IP denied")
	}
}
