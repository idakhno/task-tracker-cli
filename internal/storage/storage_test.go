package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"task-tracker-cli/internal/task"
)

func TestStorage_Load_Save(t *testing.T) {
	tmpDir := t.TempDir()
	filename := filepath.Join(tmpDir, "test_tasks.json")
	stor := New(filename)

	// Test loading non-existent file
	tasks, err := stor.Load()
	if err != nil {
		t.Fatalf("Load() error = %v, want nil", err)
	}
	if len(tasks) != 0 {
		t.Errorf("Load() len = %d, want 0", len(tasks))
	}

	// Create and save tasks
	task1, _ := task.New(1, "Task 1")
	task2, _ := task.New(2, "Task 2")
	tasksToSave := []*task.Task{task1, task2}

	if err := stor.Save(tasksToSave); err != nil {
		t.Fatalf("Save() error = %v, want nil", err)
	}

	// Load tasks back
	loadedTasks, err := stor.Load()
	if err != nil {
		t.Fatalf("Load() error = %v, want nil", err)
	}
	if len(loadedTasks) != 2 {
		t.Errorf("Load() len = %d, want 2", len(loadedTasks))
	}
}

func TestStorage_Load_InvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	filename := filepath.Join(tmpDir, "invalid.json")
	stor := New(filename)

	// Write invalid JSON
	if err := os.WriteFile(filename, []byte("{invalid json}"), 0644); err != nil {
		t.Fatalf("Failed to write invalid JSON: %v", err)
	}

	_, err := stor.Load()
	if err == nil {
		t.Error("Load() error = nil, want error")
	}
}

func TestStorage_Load_DuplicateIDs(t *testing.T) {
	tmpDir := t.TempDir()
	filename := filepath.Join(tmpDir, "duplicate.json")
	stor := New(filename)

	// Create tasks with duplicate IDs
	task1, _ := task.New(1, "Task 1")
	task2, _ := task.New(1, "Task 2") // Same ID
	tasks := []*task.Task{task1, task2}

	// Manually write JSON with duplicate IDs
	data, _ := json.Marshal(tasks)
	if err := os.WriteFile(filename, data, 0644); err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	_, err := stor.Load()
	if err == nil {
		t.Error("Load() error = nil, want duplicate ID error")
	}
}

func TestStorage_Save_Validation(t *testing.T) {
	tmpDir := t.TempDir()
	filename := filepath.Join(tmpDir, "validation.json")
	stor := New(filename)

	// Try to save tasks with duplicate IDs
	task1, _ := task.New(1, "Task 1")
	task2, _ := task.New(1, "Task 2") // Same ID
	tasks := []*task.Task{task1, task2}

	err := stor.Save(tasks)
	if err == nil {
		t.Error("Save() error = nil, want validation error")
	}
}

func TestStorage_Save_Atomic(t *testing.T) {
	tmpDir := t.TempDir()
	filename := filepath.Join(tmpDir, "atomic.json")
	stor := New(filename)

	task1, _ := task.New(1, "Task 1")
	tasks := []*task.Task{task1}

	if err := stor.Save(tasks); err != nil {
		t.Fatalf("Save() error = %v, want nil", err)
	}

	// Check that temp file doesn't exist
	tmpFile := filename + ".tmp"
	if _, err := os.Stat(tmpFile); err == nil {
		t.Error("Temporary file still exists after save")
	}
}

