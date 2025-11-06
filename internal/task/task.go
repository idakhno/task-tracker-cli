package task

import "time"

// Task status constants define the valid states a task can be in.
const (
	StatusNotDone    = "not_done"
	StatusInProgress = "in_progress"
	StatusDone       = "done"
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

// New creates a new task with the given description
func New(id int, description string) *Task {
	now := time.Now()
	return &Task{
		ID:          id,
		Description: description,
		Status:      StatusNotDone,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// UpdateDescription updates the task description and automatically
// updates the UpdatedAt timestamp to track when the change occurred.
func (t *Task) UpdateDescription(description string) {
	t.Description = description
	t.UpdatedAt = time.Now()
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

