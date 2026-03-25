package huggingface

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/magomedcoder/llm-runner/service"
)

type Options struct {
	RepoID    string
	File      string
	Revision  string
	Token     string
	ListOnly  bool
	OutDir    string
	AsRef     string
	LogOutput io.Writer
	Context   context.Context
}

func logf(opts Options, format string, args ...interface{}) {
	w := opts.LogOutput
	if w == nil {
		w = os.Stdout
	}

	_, _ = fmt.Fprintf(w, format, args...)
}

func formatBytes(n int64) string {
	const unit = 1024
	if n < unit {
		return fmt.Sprintf("%d B", n)
	}

	div, exp := int64(unit), 0
	for v := n / unit; v >= unit; v /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf("%.1f %cB", float64(n)/float64(div), "KMGTPE"[exp])
}

func optsContext(opts Options) context.Context {
	if opts.Context != nil {
		return opts.Context
	}

	return context.Background()
}

func downloadFile(opts Options, client *Client, repoID, revision, filename, outDir string) error {
	ctx := optsContext(opts)
	dstPath := filepath.Join(outDir, filename)
	f, err := os.Create(dstPath)
	if err != nil {
		return err
	}

	defer f.Close()

	logf(opts, "Скачиваю %s ...\n", filename)
	n, err := client.Download(ctx, repoID, revision, filename, f)
	if err != nil {
		_ = os.Remove(dstPath)
		return err
	}

	logf(opts, "  %s - %s\n", filename, formatBytes(n))

	return nil
}

func Run(opts Options) error {
	if strings.TrimSpace(opts.RepoID) == "" {
		return fmt.Errorf("укажите repo")
	}

	if opts.Revision == "" {
		opts.Revision = "main"
	}

	if opts.OutDir == "" {
		opts.OutDir = "./models"
	}

	ctx := optsContext(opts)
	client := NewClient(opts.Token)
	info, err := client.ModelInfo(ctx, opts.RepoID)
	if err != nil {
		return fmt.Errorf("информация о модели: %w", err)
	}

	ggufFiles := info.GGFFFilenames()
	if len(ggufFiles) == 0 {
		return fmt.Errorf("в репозитории нет .gguf файлов")
	}

	if opts.ListOnly {
		logf(opts, "Доступные .gguf файлы:\n")
		for _, name := range ggufFiles {
			logf(opts, " %s\n", name)
		}

		return nil
	}

	toDownload := ggufFiles
	if opts.File != "" {
		found := false
		for _, f := range ggufFiles {
			if f == opts.File {
				found = true
				toDownload = []string{f}
				break
			}
		}

		if !found {
			return fmt.Errorf("файл %q не найден в репозитории (используйте -list)", opts.File)
		}
	}

	if err := os.MkdirAll(opts.OutDir, 0o755); err != nil {
		return fmt.Errorf("каталог %s: %w", opts.OutDir, err)
	}

	asTrim := strings.TrimSpace(opts.AsRef)
	if asTrim != "" && len(toDownload) != 1 {
		return fmt.Errorf("-as допустим только при одном файле (укажите -file)")
	}

	for _, name := range toDownload {
		if err := ctx.Err(); err != nil {
			return err
		}

		if err := downloadFile(opts, client, opts.RepoID, opts.Revision, name, opts.OutDir); err != nil {
			return fmt.Errorf("%s: %w", name, err)
		}

		if asTrim == "" {
			continue
		}

		base, tag := service.SplitModelRef(asTrim)
		if base == "" {
			return fmt.Errorf("-as: укажите имя, например phi-3:Q4_K_M")
		}

		var newName string
		if tag == "" || strings.EqualFold(tag, "latest") {
			newName = base + ".gguf"
		} else {
			newName = fmt.Sprintf("%s-%s.gguf", base, tag)
		}

		oldPath := filepath.Join(opts.OutDir, name)
		newPath := filepath.Join(opts.OutDir, newName)
		if filepath.Clean(oldPath) != filepath.Clean(newPath) {
			if err := os.Rename(oldPath, newPath); err != nil {
				return fmt.Errorf("переименование в %s: %w", newName, err)
			}

			logf(opts, "сохранено как %s\n", newName)
		}
	}

	return nil
}
