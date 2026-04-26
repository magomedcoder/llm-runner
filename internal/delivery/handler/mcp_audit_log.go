package handler

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"

	"github.com/magomedcoder/gen/internal/domain"
	"github.com/magomedcoder/gen/pkg/logger"
)

func logMCPServerStdioAudit(op string, s *domain.MCPServer, actorUserID int) {
	if s == nil || strings.ToLower(strings.TrimSpace(s.Transport)) != "stdio" {
		return
	}

	sum := sha256.Sum256([]byte(strings.TrimSpace(s.Command)))
	short := hex.EncodeToString(sum[:8])

	switch {
	case actorUserID > 0:
		logger.I("MCP audit: %s stdio server_id=%d user_id=%d command_sha256_8=%s", op, s.ID, actorUserID, short)
	default:
		logger.I("MCP audit: %s stdio server_id=%d command_sha256_8=%s (admin)", op, s.ID, short)
	}
}
