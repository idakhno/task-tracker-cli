package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"

	"task-tracker-cli/internal/logger"
	"task-tracker-cli/internal/task"
)

const defaultFilename = "tasks.json"

var (
	ErrDuplicateID   = errors.New("duplicate task ID found")
	ErrInvalidJSON   = errors.New("invalid JSON format")
	ErrCorruptedData = errors.New("corrupted task data")
)

// Storage provides persistence layer for tasks using JSON file storage.
// It handles loading and saving tasks with proper error handling and file locking.
type Storage struct {
	filename string
	mu       sync.Mutex
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
// Validates task data and checks for duplicate IDs.
func (s *Storage) Load() ([]*task.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	logger.Debug("Loading tasks from file: %s", s.filename)
	data, err := os.ReadFile(s.filename)
	if err != nil {
		if os.IsNotExist(err) {
			logger.Debug("Tasks file does not exist, returning empty list")
			return []*task.Task{}, nil
		}
		logger.Error("Failed to read tasks file: %v", err)
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var tasks []*task.Task
	if len(data) == 0 {
		return []*task.Task{}, nil
	}

	// Try to unmarshal JSON
	if err := json.Unmarshal(data, &tasks); err != nil {
		var syntaxErr *json.SyntaxError
		var unmarshalErr *json.UnmarshalTypeError
		if errors.As(err, &syntaxErr) {
			logger.Error("Invalid JSON syntax in tasks file: %v", err)
			return nil, fmt.Errorf("%w: %v", ErrInvalidJSON, err)
		}
		if errors.As(err, &unmarshalErr) {
			logger.Error("Invalid data type in tasks file: %v", err)
			return nil, fmt.Errorf("%w: %v", ErrCorruptedData, err)
		}
		logger.Error("Failed to unmarshal tasks: %v", err)
		return nil, fmt.Errorf("failed to unmarshal tasks: %w", err)
	}

	// Validate tasks and check for duplicate IDs
	if err := s.validateTasks(tasks); err != nil {
		logger.Error("Task validation failed: %v", err)
		return nil, err
	}

	logger.Info("Loaded %d tasks from file", len(tasks))
	return tasks, nil
}

// validateTasks checks for duplicate IDs and validates task data.
func (s *Storage) validateTasks(tasks []*task.Task) error {
	seenIDs := make(map[int]bool)
	for _, t := range tasks {
		if t == nil {
			return fmt.Errorf("%w: nil task found", ErrCorruptedData)
		}
		if t.ID <= 0 {
			return fmt.Errorf("%w: invalid task ID %d", ErrCorruptedData, t.ID)
		}
		if seenIDs[t.ID] {
			return fmt.Errorf("%w: task ID %d appears multiple times", ErrDuplicateID, t.ID)
		}
		seenIDs[t.ID] = true
		if !task.IsStatusValid(t.Status) {
			return fmt.Errorf("%w: invalid status '%s' for task %d", ErrCorruptedData, t.Status, t.ID)
		}
	}
	return nil
}

// Save writes tasks to the JSON file with indented formatting for readability.
// The file is created if it doesn't exist, or overwritten if it does.
// Uses file locking to prevent race conditions.
func (s *Storage) Save(tasks []*task.Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	logger.Debug("Saving %d tasks to file: %s", len(tasks), s.filename)

	// Validate before saving
	if err := s.validateTasks(tasks); err != nil {
		logger.Error("Validation failed before save: %v", err)
		return fmt.Errorf("validation failed before save: %w", err)
	}

	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		logger.Error("Failed to marshal tasks: %v", err)
		return fmt.Errorf("failed to marshal tasks: %w", err)
	}

	// Write to temporary file first, then rename (atomic operation)
	tmpFile := s.filename + ".tmp"
	if err := os.WriteFile(tmpFile, data, 0644); err != nil {
		logger.Error("Failed to write temporary file: %v", err)
		return fmt.Errorf("failed to write temporary file: %w", err)
	}

	// Atomic rename
	if err := os.Rename(tmpFile, s.filename); err != nil {
		os.Remove(tmpFile) // Clean up on error
		logger.Error("Failed to rename file: %v", err)
		return fmt.Errorf("failed to rename file: %w", err)
	}

	logger.Info("Successfully saved %d tasks to file", len(tasks))
	return nil
}
