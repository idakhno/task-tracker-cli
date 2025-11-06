package main

import (
	"fmt"
	"os"

	"task-tracker-cli/internal/cli"
	"task-tracker-cli/internal/storage"
)

func main() {
	stor := storage.New("")
	app, err := cli.New(stor)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

