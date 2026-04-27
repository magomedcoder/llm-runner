package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/magomedcoder/gen/internal/domain"
)

type ServerConfig struct {
	Port              string `yaml:"port"`
	Host              string `yaml:"host"`
	MetricsListenAddr string `yaml:"metrics_listen"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

type JWTConfig struct {
	AccessSecret  string        `yaml:"access_secret"`
	RefreshSecret string        `yaml:"refresh_secret"`
	AccessTTL     time.Duration `yaml:"access_ttl"`
	RefreshTTL    time.Duration `yaml:"refresh_ttl"`
}

type MCPConfig struct {
	HTTPAllowAny                        bool     `yaml:"http_allow_any"`
	HTTPAllowHosts                      []string `yaml:"http_allow_hosts"`
	Roots                               []string `yaml:"roots"`
	SamplingEnabled                     bool     `yaml:"sampling_enabled"`
	LogServerMessages                   bool     `yaml:"log_server_messages"`
	HTTPReuseSessions                   bool     `yaml:"http_reuse_sessions"`
	HTTPSessionMaxIdleSeconds           int      `yaml:"http_session_max_idle_seconds"`
	StdioDisabled                       bool     `yaml:"stdio_disabled"`
	StdioCommandPrefixes                []string `yaml:"stdio_command_prefixes"`
	MaxMCPServersPerUser                int      `yaml:"max_mcp_servers_per_user"`
	MaxTrackedServerIDsForCallStats     int      `yaml:"max_tracked_server_ids_for_call_stats"`
	ToolsAllowlistWhenUserImageAttached []string `yaml:"tools_allowlist_when_user_image_attached"`
}

type RAGConfig struct {
	ChunkSizeRunes                int     `yaml:"chunk_size_runes"`
	ChunkOverlapRunes             int     `yaml:"chunk_overlap_runes"`
	EmbedBatchSize                int     `yaml:"embed_batch_size"`
	MaxChunkEmbedRunes            int     `yaml:"max_chunk_embed_runes"`
	BackgroundIndexTimeoutSeconds int     `yaml:"background_index_timeout_seconds"`
	LLMContextFallbackTokens      int     `yaml:"llm_context_fallback_tokens"`
	MaxExtractedRunesOnUpload     int     `yaml:"max_extracted_runes_on_upload"`
	NeighborChunkWindow           int     `yaml:"neighbor_chunk_window"`
	PreferFullDocumentWhenFits    *bool   `yaml:"prefer_full_document_when_fits"`
	QueryRewriteEnabled           *bool   `yaml:"query_rewrite_enabled"`
	QueryRewriteMaxTokens         int     `yaml:"query_rewrite_max_tokens"`
	QueryRewriteTimeoutSeconds    int     `yaml:"query_rewrite_timeout_seconds"`
	HydeEnabled                   *bool   `yaml:"hyde_enabled"`
	HydeMaxTokens                 int     `yaml:"hyde_max_tokens"`
	HydeTimeoutSeconds            int     `yaml:"hyde_timeout_seconds"`
	AdaptiveKEnabled              *bool   `yaml:"adaptive_k_enabled"`
	AdaptiveKMultiplier           int     `yaml:"adaptive_k_multiplier"`
	MinSimilarityScore            float64 `yaml:"min_similarity_score"`
	DeepRAGEnabled                *bool   `yaml:"deep_rag_enabled"`
	DeepRAGMaxMapCalls            int     `yaml:"deep_rag_max_map_calls"`
	DeepRAGChunksPerMap           int     `yaml:"deep_rag_chunks_per_map"`
	DeepRAGMapMaxTokens           int     `yaml:"deep_rag_map_max_tokens"`
	DeepRAGMapTimeoutSeconds      int     `yaml:"deep_rag_map_timeout_seconds"`
	DeepRAGMaxMapOutputRunes      int     `yaml:"deep_rag_max_map_output_runes"`
	RerankEnabled                 *bool   `yaml:"rerank_enabled"`
	RerankMaxCandidates           int     `yaml:"rerank_max_candidates"`
	RerankMaxTokens               int     `yaml:"rerank_max_tokens"`
	RerankTimeoutSeconds          int     `yaml:"rerank_timeout_seconds"`
	RerankPassageMaxRunes         int     `yaml:"rerank_passage_max_runes"`
}

const (
	ragDefaultChunkSizeRunes           = 2048
	ragDefaultChunkOverlapRunes        = 256
	ragDefaultEmbedBatchSize           = 64
	ragDefaultMaxChunkEmbedRunes       = 16384
	ragDefaultBackgroundIndexSeconds   = 7200
	ragDefaultLLMContextFallbackTokens = 131072
	ragDefaultNeighborChunkWindow      = 8
	ragDefaultQueryRewriteMaxTokens    = 512
	ragDefaultQueryRewriteTimeoutSec   = 300
	ragDefaultHydeMaxTokens            = 768
	ragDefaultHydeTimeoutSec           = 300
	ragDefaultAdaptiveKMultiplier      = 6
	ragDefaultDeepRAGMaxMapCalls       = 16
	ragDefaultDeepRAGChunksPerMap      = 12
	ragDefaultDeepRAGMapMaxTokens      = 2048
	ragDefaultDeepRAGMapTimeoutSec     = 600
	ragDefaultDeepRAGMaxMapOutputRunes = 50000
	ragDefaultRerankMaxCandidates      = 32
	ragDefaultRerankMaxTokens          = 1024
	ragDefaultRerankTimeoutSec         = 300
	ragDefaultRerankPassageMaxRunes    = 4000
)

func (r RAGConfig) EffectiveChunkSizeRunes() int {
	if r.ChunkSizeRunes <= 0 {
		return ragDefaultChunkSizeRunes
	}
	return r.ChunkSizeRunes
}

func (r RAGConfig) EffectiveChunkOverlapRunes() int {
	if r.ChunkOverlapRunes < 0 {
		return 0
	}
	if r.ChunkOverlapRunes == 0 {
		return ragDefaultChunkOverlapRunes
	}
	return r.ChunkOverlapRunes
}

func (r RAGConfig) EffectiveEmbedBatchSize() int {
	if r.EmbedBatchSize <= 0 {
		return ragDefaultEmbedBatchSize
	}
	return r.EmbedBatchSize
}

func (r RAGConfig) EffectiveMaxChunkEmbedRunes() int {
	if r.MaxChunkEmbedRunes <= 0 {
		return ragDefaultMaxChunkEmbedRunes
	}
	return r.MaxChunkEmbedRunes
}

func (r RAGConfig) BackgroundIndexTimeout() time.Duration {
	if r.BackgroundIndexTimeoutSeconds <= 0 {
		return time.Duration(ragDefaultBackgroundIndexSeconds) * time.Second
	}

	return time.Duration(r.BackgroundIndexTimeoutSeconds) * time.Second
}

func (r RAGConfig) EffectiveLLMContextFallbackTokens() int {
	if r.LLMContextFallbackTokens < 0 {
		return 0
	}

	if r.LLMContextFallbackTokens == 0 {
		return ragDefaultLLMContextFallbackTokens
	}

	return r.LLMContextFallbackTokens
}

func (r RAGConfig) EffectiveMaxExtractedRunesOnUpload() int {
	if r.MaxExtractedRunesOnUpload <= 0 {
		return 0
	}
	return r.MaxExtractedRunesOnUpload
}

func (r RAGConfig) EffectiveNeighborChunkWindow() int {
	if r.NeighborChunkWindow < 0 {
		return 0
	}

	if r.NeighborChunkWindow == 0 {
		return ragDefaultNeighborChunkWindow
	}

	if r.NeighborChunkWindow > 8 {
		return 8
	}

	return r.NeighborChunkWindow
}

func (r RAGConfig) EffectivePreferFullDocumentWhenFits() bool {
	if r.PreferFullDocumentWhenFits == nil {
		return true
	}

	return *r.PreferFullDocumentWhenFits
}

func (r RAGConfig) EffectiveQueryRewriteEnabled() bool {
	if r.QueryRewriteEnabled == nil {
		return true
	}

	return *r.QueryRewriteEnabled
}

func (r RAGConfig) EffectiveQueryRewriteMaxTokens() int32 {
	t := r.QueryRewriteMaxTokens
	if t <= 0 {
		return int32(ragDefaultQueryRewriteMaxTokens)
	}

	if t > 512 {
		return 512
	}

	if t < 16 {
		return 16
	}

	return int32(t)
}

func (r RAGConfig) EffectiveQueryRewriteTimeoutSeconds() int32 {
	if r.QueryRewriteTimeoutSeconds <= 0 {
		return int32(ragDefaultQueryRewriteTimeoutSec)
	}

	if r.QueryRewriteTimeoutSeconds > 300 {
		return 300
	}

	return int32(r.QueryRewriteTimeoutSeconds)
}

func (r RAGConfig) EffectiveHydeEnabled() bool {
	if r.HydeEnabled == nil {
		return true
	}

	return *r.HydeEnabled
}

func (r RAGConfig) EffectiveHydeMaxTokens() int32 {
	t := r.HydeMaxTokens
	if t <= 0 {
		return int32(ragDefaultHydeMaxTokens)
	}

	if t > 768 {
		return 768
	}

	if t < 32 {
		return 32
	}

	return int32(t)
}

func (r RAGConfig) EffectiveHydeTimeoutSeconds() int32 {
	if r.HydeTimeoutSeconds <= 0 {
		return int32(ragDefaultHydeTimeoutSec)
	}

	if r.HydeTimeoutSeconds > 300 {
		return 300
	}

	return int32(r.HydeTimeoutSeconds)
}

func (r RAGConfig) EffectiveAdaptiveKEnabled() bool {
	if r.AdaptiveKEnabled == nil {
		return true
	}

	return *r.AdaptiveKEnabled
}

func (r RAGConfig) EffectiveAdaptiveKMultiplier() int {
	m := r.AdaptiveKMultiplier
	if m <= 0 {
		return ragDefaultAdaptiveKMultiplier
	}

	if m < 1 {
		return 1
	}

	if m > 6 {
		return 6
	}

	return m
}

func (r RAGConfig) EffectiveMinSimilarityScore() float64 {
	if r.MinSimilarityScore == 0 {
		return -1
	}

	if r.MinSimilarityScore <= -1 {
		return -1
	}

	if r.MinSimilarityScore > 1 {
		return 1
	}

	return r.MinSimilarityScore
}

func (r RAGConfig) EffectiveDeepRAGEnabled() bool {
	if r.DeepRAGEnabled == nil {
		return true
	}

	return *r.DeepRAGEnabled
}

func (r RAGConfig) EffectiveDeepRAGMaxMapCalls() int {
	n := r.DeepRAGMaxMapCalls
	if n <= 0 {
		return ragDefaultDeepRAGMaxMapCalls
	}

	if n > 16 {
		return 16
	}

	return n
}

func (r RAGConfig) EffectiveDeepRAGChunksPerMap() int {
	n := r.DeepRAGChunksPerMap
	if n <= 0 {
		return ragDefaultDeepRAGChunksPerMap
	}

	if n > 12 {
		return 12
	}

	return n
}

func (r RAGConfig) EffectiveDeepRAGMapMaxTokens() int32 {
	t := r.DeepRAGMapMaxTokens
	if t <= 0 {
		return int32(ragDefaultDeepRAGMapMaxTokens)
	}

	if t > 2048 {
		return 2048
	}

	if t < 64 {
		return 64
	}

	return int32(t)
}

func (r RAGConfig) EffectiveDeepRAGMapTimeoutSeconds() int32 {
	if r.DeepRAGMapTimeoutSeconds <= 0 {
		return int32(ragDefaultDeepRAGMapTimeoutSec)
	}

	if r.DeepRAGMapTimeoutSeconds > 600 {
		return 600
	}

	return int32(r.DeepRAGMapTimeoutSeconds)
}

func (r RAGConfig) EffectiveDeepRAGMaxMapOutputRunes() int {
	n := r.DeepRAGMaxMapOutputRunes
	if n <= 0 {
		return ragDefaultDeepRAGMaxMapOutputRunes
	}

	if n < 400 {
		return 400
	}

	if n > 50000 {
		return 50000
	}

	return n
}

func (r RAGConfig) EffectiveRerankEnabled() bool {
	if r.RerankEnabled == nil {
		return true
	}

	return *r.RerankEnabled
}

func (r RAGConfig) EffectiveRerankMaxCandidates() int {
	n := r.RerankMaxCandidates
	if n <= 0 {
		return ragDefaultRerankMaxCandidates
	}

	if n < 2 {
		return 2
	}

	if n > 32 {
		return 32
	}

	return n
}

func (r RAGConfig) EffectiveRerankMaxTokens() int32 {
	t := r.RerankMaxTokens
	if t <= 0 {
		return int32(ragDefaultRerankMaxTokens)
	}

	if t > 1024 {
		return 1024
	}

	if t < 32 {
		return 32
	}

	return int32(t)
}

func (r RAGConfig) EffectiveRerankTimeoutSeconds() int32 {
	if r.RerankTimeoutSeconds <= 0 {
		return int32(ragDefaultRerankTimeoutSec)
	}

	if r.RerankTimeoutSeconds > 300 {
		return 300
	}

	return int32(r.RerankTimeoutSeconds)
}

func (r RAGConfig) EffectiveRerankPassageMaxRunes() int {
	n := r.RerankPassageMaxRunes
	if n <= 0 {
		return ragDefaultRerankPassageMaxRunes
	}

	if n < 200 {
		return 200
	}

	if n > 4000 {
		return 4000
	}

	return n
}

type Config struct {
	Server                       ServerConfig
	Database                     DatabaseConfig
	JWT                          JWTConfig
	MCP                          MCPConfig
	RAG                          RAGConfig `yaml:"rag"`
	DataDir                      string    `yaml:"data_dir"`
	AttachmentHydrateParallelism int       `yaml:"attachment_hydrate_parallelism"`
	LogLevel                     string    `yaml:"log_level"`
	MinClientBuild               int32
}

func LoadFrom(path string) (*Config, error) {
	var conf Config
	path = strings.TrimSpace(path)
	if path == "" {
		return nil, fmt.Errorf("путь к файлу конфигурации пустой")
	}

	content, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	if err := yaml.Unmarshal(content, &conf); err != nil {
		log.Println(err)
		panic(fmt.Sprintf("Ошибка при разборе: %v", err))
	}

	//conf.MinClientBuild = 1

	return &conf, nil
}

func (c *Config) MCPHTTPHostAllowed(host string) bool {
	if c == nil {
		return false
	}

	if c.MCP.HTTPAllowAny {
		return true
	}

	h := strings.TrimSpace(host)
	if h == "" {
		return false
	}

	if ip := net.ParseIP(h); ip != nil {
		if ip.IsLoopback() {
			return true
		}

		for _, e := range c.MCP.HTTPAllowHosts {
			e = strings.TrimSpace(e)
			if e == "" {
				continue
			}

			if ip2 := net.ParseIP(e); ip2 != nil && ip.Equal(ip2) {
				return true
			}
		}

		return false
	}

	h = strings.ToLower(h)
	if h == "localhost" {
		return true
	}

	for _, s := range c.MCP.HTTPAllowHosts {
		s = strings.ToLower(strings.TrimSpace(s))
		if s == "" || net.ParseIP(s) != nil {
			continue
		}

		if h == s || strings.HasSuffix(h, "."+s) {
			return true
		}
	}

	return false
}

func (c *Config) ValidateMCPServerHTTP(s *domain.MCPServer) error {
	if c == nil || s == nil {
		return nil
	}

	tr := strings.ToLower(strings.TrimSpace(s.Transport))
	if tr != "sse" && tr != "streamable" {
		return nil
	}

	raw := strings.TrimSpace(s.URL)
	if raw == "" {
		return fmt.Errorf("для транспорта %s нужен непустой url", tr)
	}

	u, err := url.Parse(raw)
	if err != nil {
		return fmt.Errorf("url: %w", err)
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("url: ожидается http или https")
	}

	host := u.Hostname()
	if host == "" {
		return fmt.Errorf("url: пустой хост")
	}

	if !c.MCPHTTPHostAllowed(host) {
		return fmt.Errorf("хост %q не разрешён политикой MCP", host)
	}

	return nil
}

func StdioCommandAllowed(command string, prefixes []string) bool {
	cmd := strings.TrimSpace(command)
	if cmd == "" {
		return false
	}

	if len(prefixes) == 0 {
		return true
	}

	cleaned := filepath.Clean(cmd)
	for _, p := range prefixes {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}

		prefix := filepath.Clean(p)
		if cleaned == prefix {
			return true
		}

		sep := string(filepath.Separator)
		if strings.HasPrefix(cleaned, prefix+sep) {
			return true
		}
	}

	return false
}

func (c *Config) ValidateMCPServerStdio(s *domain.MCPServer) error {
	if c == nil || s == nil {
		return nil
	}

	tr := strings.ToLower(strings.TrimSpace(s.Transport))
	if tr != "stdio" {
		return nil
	}

	if c.MCP.StdioDisabled {
		return fmt.Errorf("транспорт stdio отключён в конфигурации (mcp.stdio_disabled)")
	}

	cmd := strings.TrimSpace(s.Command)
	if cmd == "" {
		return fmt.Errorf("stdio: пустая команда")
	}

	if !StdioCommandAllowed(cmd, c.MCP.StdioCommandPrefixes) {
		return fmt.Errorf("команда stdio не разрешена политикой mcp.stdio_command_prefixes")
	}

	return nil
}

func (c *Config) ValidateMCPServer(s *domain.MCPServer) error {
	if err := c.ValidateMCPServerHTTP(s); err != nil {
		return err
	}

	if err := c.ValidateMCPServerStdio(s); err != nil {
		return err
	}
	return nil
}
