package huggingface

import (
	"path/filepath"
	"strings"
)

func sanitizePathSegment(s string) string {
	s = strings.TrimSpace(s)
	if s == "" || s == "." || s == ".." {
		return ""
	}

	repl := strings.NewReplacer(
		`<`, "_",
		`>`, "_",
		`:`, "-",
		`"`, "_",
		`|`, "_",
		`?`, "_",
		`*`, "_",
		`\`, "_",
		`/`, "_",
	)
	s = repl.Replace(s)
	s = strings.Trim(s, " .")

	if s == "" || s == "." || s == ".." {
		return ""
	}

	return s
}

func RepoDownloadSubdir(repoID string) string {
	repoID = strings.Trim(repoID, "/")
	if repoID == "" {
		return "unknown-repo"
	}

	parts := strings.Split(repoID, "/")
	var seg []string
	for _, p := range parts {
		c := sanitizePathSegment(p)
		if c != "" {
			seg = append(seg, c)
		}
	}

	if len(seg) == 0 {
		return "unknown-repo"
	}

	return filepath.Join(seg...)
}
