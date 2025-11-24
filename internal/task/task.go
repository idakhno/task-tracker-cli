package task

import (
	"errors"
	"strings"
	"time"
)

// Task status constants define the valid states a task can be in.
const (
	StatusNotDone    = "not_done"
	StatusInProgress = "in_progress"
	StatusDone       = "done"
)

const (
	MaxDescriptionLength = 1000
	MinDescriptionLength = 1
)

var (
	ErrEmptyDescription   = errors.New("task description cannot be empty")
	ErrDescriptionTooLong = errors.New("task description exceeds maximum length")
	ErrInvalidID          = errors.New("task ID must be positive")
)

// Task represents a single task with its metadata.
// Tasks are persisted to JSON and track their lifecycle through timestamps.
type Task struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// New creates a new task with the given description.
// Returns an error if the description is invalid.
func New(id int, description string) (*Task, error) {
	if err := ValidateDescription(description); err != nil {
		return nil, err
	}
	if id <= 0 {
		return nil, ErrInvalidID
	}

	now := time.Now()
	return &Task{
		ID:          id,
		Description: description,
		Status:      StatusNotDone,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

// ValidateDescription checks if a task description is valid.
func ValidateDescription(description string) error {
	trimmed := strings.TrimSpace(description)
	if len(trimmed) < MinDescriptionLength {
		return ErrEmptyDescription
	}
	if len(trimmed) > MaxDescriptionLength {
		return ErrDescriptionTooLong
	}
	return nil
}

// UpdateDescription updates the task description and automatically
// updates the UpdatedAt timestamp to track when the change occurred.
// Returns an error if the description is invalid.
func (t *Task) UpdateDescription(description string) error {
	if err := ValidateDescription(description); err != nil {
		return err
	}
	t.Description = strings.TrimSpace(description)
	t.UpdatedAt = time.Now()
	return nil
}

// MarkInProgress transitions the task to in-progress status and
// updates the timestamp to reflect the state change.
func (t *Task) MarkInProgress() {
	t.Status = StatusInProgress
	t.UpdatedAt = time.Now()
}

// MarkDone transitions the task to completed status and
// updates the timestamp to reflect the state change.
func (t *Task) MarkDone() {
	t.Status = StatusDone
	t.UpdatedAt = time.Now()
}

// IsStatusValid checks if the status is valid
func IsStatusValid(status string) bool {
	return status == StatusNotDone || status == StatusInProgress || status == StatusDone
}
