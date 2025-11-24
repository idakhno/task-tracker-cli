package task

import (
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name        string
		id          int
		description string
		wantErr     bool
	}{
		{"valid task", 1, "Test task", false},
		{"empty description", 1, "", true},
		{"whitespace only", 1, "   ", true},
		{"too long description", 1, string(make([]byte, MaxDescriptionLength+1)), true},
		{"zero ID", 0, "Test task", true},
		{"negative ID", -1, "Test task", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task, err := New(tt.id, tt.description)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && task == nil {
				t.Error("New() returned nil task without error")
			}
			if !tt.wantErr {
				if task.ID != tt.id {
					t.Errorf("New() ID = %d, want %d", task.ID, tt.id)
				}
				if task.Status != StatusNotDone {
					t.Errorf("New() Status = %s, want %s", task.Status, StatusNotDone)
				}
			}
		})
	}
}

func TestValidateDescription(t *testing.T) {
	tests := []struct {
		name        string
		description string
		wantErr     bool
	}{
		{"valid", "Valid task description", false},
		{"empty", "", true},
		{"whitespace", "   ", true},
		{"too long", string(make([]byte, MaxDescriptionLength+1)), true},
		{"max length", string(make([]byte, MaxDescriptionLength)), false},
		{"min length", "a", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateDescription(tt.description)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateDescription() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdateDescription(t *testing.T) {
	task, _ := New(1, "Original description")
	
	tests := []struct {
		name        string
		description string
		wantErr     bool
	}{
		{"valid update", "Updated description", false},
		{"empty", "", true},
		{"whitespace", "   ", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := task.UpdateDescription(tt.description)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateDescription() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && task.Description != strings.TrimSpace(tt.description) {
				t.Errorf("UpdateDescription() description = %s, want %s", task.Description, strings.TrimSpace(tt.description))
			}
		})
	}
}

func TestMarkInProgress(t *testing.T) {
	task, _ := New(1, "Test task")
	task.MarkInProgress()
	if task.Status != StatusInProgress {
		t.Errorf("MarkInProgress() status = %s, want %s", task.Status, StatusInProgress)
	}
}

func TestMarkDone(t *testing.T) {
	task, _ := New(1, "Test task")
	task.MarkDone()
	if task.Status != StatusDone {
		t.Errorf("MarkDone() status = %s, want %s", task.Status, StatusDone)
	}
}

func TestIsStatusValid(t *testing.T) {
	tests := []struct {
		status string
		want   bool
	}{
		{StatusNotDone, true},
		{StatusInProgress, true},
		{StatusDone, true},
		{"invalid", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.status, func(t *testing.T) {
			if got := IsStatusValid(tt.status); got != tt.want {
				t.Errorf("IsStatusValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

