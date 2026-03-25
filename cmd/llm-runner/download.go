package main

import (
	"context"

	"github.com/magomedcoder/llm-runner/huggingface"
	"github.com/urfave/cli/v3"
)

func cmdDownload() *cli.Command {
	return &cli.Command{
		Name:    "download",
		Aliases: []string{"dl"},
		Usage:   "Скачать .gguf с Hugging Face локально (без запущенного раннера)",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "repo",
				Usage: "id репозитория Hugging Face",
			},
			&cli.StringFlag{
				Name:  "file",
				Usage: "конкретный .gguf в репо",
			},
			&cli.StringFlag{
				Name:  "revision",
				Value: "main",
				Usage: "ветка / ревизия",
			},
			&cli.StringFlag{
				Name:  "token",
				Usage: "HF token"},
			&cli.BoolFlag{
				Name:    "list",
				Aliases: []string{"l"},
				Usage:   "только список .gguf",
			},
			&cli.StringFlag{
				Name:  "out",
				Value: "./models",
				Usage: "каталог для файлов",
			},
			&cli.StringFlag{
				Name:  "as",
				Usage: "сохранить как base:tag → base-tag.gguf",
			},
		},
		Action: runDownload,
	}
}

func runDownload(_ context.Context, cmd *cli.Command) error {
	return huggingface.Run(huggingface.Options{
		RepoID:   cmd.String("repo"),
		File:     cmd.String("file"),
		Revision: cmd.String("revision"),
		Token:    cmd.String("token"),
		ListOnly: cmd.Bool("list"),
		OutDir:   cmd.String("out"),
		AsRef:    cmd.String("as"),
	})
}
