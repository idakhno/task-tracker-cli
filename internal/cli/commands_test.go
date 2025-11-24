package cli

import (
	"path/filepath"
	"testing"

	"task-tracker-cli/internal/storage"
)

func TestCLI_handleAdd(t *testing.T) {
	tmpDir := t.TempDir()
	filename := filepath.Join(tmpDir, "test.json")
	stor := storage.New(filename)
	app, _ := New(stor)

	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{"valid add", []string{"", "add", "Test task"}, false},
		{"empty description", []string{"", "add", ""}, true},
		{"missing args", []string{"", "add"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := app.handleAdd(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("handleAdd() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCLI_handleUpdate(t *testing.T) {
	tmpDir := t.TempDir()
	filename := filepath.Join(tmpDir, "test.json")
	stor := storage.New(filename)
	app, _ := New(stor)

	// Add a task first
	app.handleAdd([]string{"", "add", "Original task"})

	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{"valid update", []string{"", "update", "1", "Updated task"}, false},
		{"empty description", []string{"", "update", "1", ""}, true},
		{"invalid ID", []string{"", "update", "0", "Task"}, true},
		{"non-existent ID", []string{"", "update", "999", "Task"}, true},
		{"missing args", []string{"", "update"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := app.handleUpdate(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("handleUpdate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCLI_handleDelete(t *testing.T) {
	tmpDir := t.TempDir()
	filename := filepath.Join(tmpDir, "test.json")
	stor := storage.New(filename)
	app, _ := New(stor)

	// Add a task first
	app.handleAdd([]string{"", "add", "Task to delete"})

	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{"valid delete", []string{"", "delete", "1"}, false},
		{"invalid ID", []string{"", "delete", "0"}, true},
		{"non-existent ID", []string{"", "delete", "999"}, true},
		{"missing args", []string{"", "delete"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := app.handleDelete(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("handleDelete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCLI_nextID(t *testing.T) {
	tmpDir := t.TempDir()
	filename := filepath.Join(tmpDir, "test.json")
	stor := storage.New(filename)
	app, _ := New(stor)

	// Test initial ID
	if id := app.nextID(); id != 1 {
		t.Errorf("nextID() = %d, want 1", id)
	}

	// Add tasks and check ID generation
	app.handleAdd([]string{"", "add", "Task 1"})
	app.handleAdd([]string{"", "add", "Task 2"})

	if id := app.nextID(); id != 3 {
		t.Errorf("nextID() = %d, want 3", id)
	}
}

func TestCLI_findTask(t *testing.T) {
	tmpDir := t.TempDir()
	filename := filepath.Join(tmpDir, "test.json")
	stor := storage.New(filename)
	app, _ := New(stor)

	app.handleAdd([]string{"", "add", "Task 1"})

	task, err := app.findTask(1)
	if err != nil {
		t.Fatalf("findTask() error = %v, want nil", err)
	}
	if task.ID != 1 {
		t.Errorf("findTask() ID = %d, want 1", task.ID)
	}

	_, err = app.findTask(999)
	if err == nil {
		t.Error("findTask() error = nil, want error")
	}
}

