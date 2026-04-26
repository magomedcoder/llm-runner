package main

import (
	"context"
	"log"
	"os"

	"github.com/urfave/cli/v3"
)

func main() {
	app := &cli.Command{
		Name:  "gen-llm-runner",
		Usage: "Gen runner",
		Commands: []*cli.Command{
			cmdServe(),
			cmdCreate(),
			cmdShow(),
			cmdDownload(),
			cmdModels(),
			cmdRemote(),
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return cli.ShowAppHelp(cmd)
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
