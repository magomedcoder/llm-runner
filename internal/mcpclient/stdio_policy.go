package mcpclient

import "github.com/magomedcoder/gen/internal/domain"

var stdioServerValidator func(*domain.MCPServer) error

func SetStdioServerValidator(v func(*domain.MCPServer) error) {
	stdioServerValidator = v
}

func validateStdioPolicy(srv *domain.MCPServer) error {
	if stdioServerValidator == nil {
		return nil
	}

	return stdioServerValidator(srv)
}
