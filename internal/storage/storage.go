package storage

import (
	"encoding/json"
	"fmt"
	"os"

	"task-tracker-cli/internal/task"
)

const defaultFilename = "tasks.json"

// Storage provides persistence layer for tasks using JSON file storage.
// It handles loading and saving tasks with proper error handling.
type Storage struct {
	filename string
}

// New creates a new storage instance with the specified filename.
// If filename is empty, it defaults to "tasks.json".
func New(filename string) *Storage {
	if filename == "" {
		filename = defaultFilename
	}
	return &Storage{filename: filename}
}

// Load reads tasks from the JSON file and returns them as a slice.
// Returns an empty slice if the file doesn't exist (not an error).
// Returns an error if the file exists but cannot be read or parsed.
func (s *Storage) Load() ([]*task.Task, error) {
	data, err := os.ReadFile(s.filename)
	if err != nil {
		if os.IsNotExist(err) {
			return []*task.Task{}, nil
		}
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var tasks []*task.Task
	if len(data) == 0 {
		return []*task.Task{}, nil
	}

	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, fmt.Errorf("failed to unmarshal tasks: %w", err)
	}

	return tasks, nil
}

// Save writes tasks to the JSON file with indented formatting for readability.
// The file is created if it doesn't exist, or overwritten if it does.
func (s *Storage) Save(tasks []*task.Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal tasks: %w", err)
	}

	if err := os.WriteFile(s.filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

