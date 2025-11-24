package main

import (
	"flag"
	"fmt"
	"os"

	"task-tracker-cli/internal/cli"
	"task-tracker-cli/internal/storage"
)

// These variables are set via build flags:
// go build -ldflags "-X main.version=1.0.0 -X main.buildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ) -X main.gitCommit=$(git rev-parse HEAD)"
var (
	version   = "dev"
	buildTime = "unknown"
	gitCommit = "unknown"
)

func main() {
	var (
		fileFlag    = flag.String("file", "", "Path to tasks JSON file (default: tasks.json in current directory)")
		versionFlag = flag.Bool("version", false, "Show version information")
	)

	// Parse flags but keep remaining args
	flag.Parse()

	// Handle version flag
	if *versionFlag {
		fmt.Printf("task-tracker-cli version %s\n", version)
		fmt.Printf("Build time: %s\n", buildTime)
		fmt.Printf("Git commit: %s\n", gitCommit)
		os.Exit(0)
	}

	// Determine file path: flag > environment variable > default
	filename := *fileFlag
	if filename == "" {
		filename = os.Getenv("TASK_TRACKER_FILE")
	}
	if filename == "" {
		filename = ""
	}

	stor := storage.New(filename)
	app, err := cli.New(stor)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Use flag.Args() which contains all non-flag arguments
	// Prepend program name to maintain expected format
	args := append([]string{os.Args[0]}, flag.Args()...)
	if err := app.Run(args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
