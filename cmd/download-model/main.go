package main

import (
	"flag"
	"fmt"
	"github.com/magomedcoder/llm-runner/huggingface"
	"os"
	"path/filepath"
)

func main() {
	repoId := flag.String("repo", "", "")
	file := flag.String("file", "", "")
	revision := flag.String("revision", "main", "")
	token := flag.String("token", "", "")
	listOnly := flag.Bool("list", false, "")
	flag.Parse()

	client := huggingface.NewClient(*token)
	info, err := client.ModelInfo(*repoId)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка получения информации о модели: %v\n", err)
		os.Exit(1)
	}

	ggufFiles := info.GGFFFilenames()
	if len(ggufFiles) == 0 {
		fmt.Fprintln(os.Stderr, "В репозитории нет .gguf файлов")
		os.Exit(1)
	}

	if *listOnly {
		fmt.Println("Доступные .gguf файлы:")
		for _, name := range ggufFiles {
			fmt.Println(" ", name)
		}
		return
	}

	toDownload := ggufFiles
	if *file != "" {
		found := false
		for _, f := range ggufFiles {
			if f == *file {
				found = true
				toDownload = []string{f}
				break
			}
		}
		if !found {
			fmt.Fprintf(os.Stderr, "Файл %q не найден в репозитории. Используйте -list для списка.\n", *file)
			os.Exit(1)
		}
	}

	if err := os.MkdirAll("./models", 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка создания каталога %s: %v\n", "./models", err)
		os.Exit(1)
	}

	for _, name := range toDownload {
		if err := downloadFile(client, *repoId, *revision, name, "./models"); err != nil {
			fmt.Fprintf(os.Stderr, "Ошибка загрузки %s: %v\n", name, err)
			os.Exit(1)
		}
	}
}

func downloadFile(client *huggingface.Client, repoId, revision, filename, outDir string) error {
	dstPath := filepath.Join(outDir, filename)
	f, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer f.Close()

	fmt.Printf("Скачиваю %s ...\n", filename)
	n, err := client.Download(repoId, revision, filename, f)
	if err != nil {
		os.Remove(dstPath)
		return err
	}
	fmt.Printf("  %s - %s\n", filename, formatBytes(n))

	return nil
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
