package main

import (
	"log"

	"github.com/magomedcoder/gen/tools/tce-server/internal"
)

func main() {
	app, err := internal.NewFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	defer app.Close()

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
