package cli

import (
	"fmt"
	"strconv"
	"strings"

	"task-tracker-cli/internal/storage"
	"task-tracker-cli/internal/task"
)

// CLI handles all command-line interface operations.
// It maintains an in-memory copy of tasks and persists changes through storage.
type CLI struct {
	storage *storage.Storage
	tasks   []*task.Task
}

// New initializes a CLI instance by loading existing tasks from storage.
// Returns an error if tasks cannot be loaded.
func New(storage *storage.Storage) (*CLI, error) {
	tasks, err := storage.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load tasks: %w", err)
	}

	return &CLI{
		storage: storage,
		tasks:   tasks,
	}, nil
}

// Run parses command-line arguments and executes the appropriate command handler.
// Displays usage information if no command is provided.
func (c *CLI) Run(args []string) error {
	if len(args) < 2 {
		c.printUsage()
		return nil
	}

	command := args[1]

	switch command {
	case "add":
		return c.handleAdd(args)
	case "update":
		return c.handleUpdate(args)
	case "delete":
		return c.handleDelete(args)
	case "mark-in-progress":
		return c.handleMarkInProgress(args)
	case "mark-done":
		return c.handleMarkDone(args)
	case "list":
		c.handleList(args)
		return nil
	default:
		return fmt.Errorf("unknown command: %s", command)
	}
}

// handleAdd creates a new task with the provided description.
// The description is constructed by joining all arguments after the command.
func (c *CLI) handleAdd(args []string) error {
	if len(args) < 3 {
		return fmt.Errorf("usage: task-tracker-cli add <description>")
	}

	description := strings.TrimSpace(strings.Join(args[2:], " "))
	if description == "" {
		return fmt.Errorf("task description cannot be empty")
	}

	newTask, err := task.New(c.nextID(), description)
	if err != nil {
		return fmt.Errorf("failed to create task: %w", err)
	}

	c.tasks = append(c.tasks, newTask)

	if err := c.storage.Save(c.tasks); err != nil {
		return fmt.Errorf("failed to save task: %w", err)
	}

	c.printTask(newTask)
	return nil
}

func (c *CLI) handleUpdate(args []string) error {
	if len(args) < 4 {
		return fmt.Errorf("usage: task-tracker-cli update <id> <description>")
	}

	id, err := strconv.Atoi(args[2])
	if err != nil {
		return fmt.Errorf("invalid ID: %w", err)
	}

	if id <= 0 {
		return fmt.Errorf("task ID must be positive")
	}

	description := strings.TrimSpace(strings.Join(args[3:], " "))
	if description == "" {
		return fmt.Errorf("task description cannot be empty")
	}

	t, err := c.findTask(id)
	if err != nil {
		return err
	}

	if err := t.UpdateDescription(description); err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	if err := c.storage.Save(c.tasks); err != nil {
		return fmt.Errorf("failed to save task: %w", err)
	}

	fmt.Println("Task updated successfully.")
	return nil
}

func (c *CLI) handleDelete(args []string) error {
	if len(args) < 3 {
		return fmt.Errorf("usage: task-tracker-cli delete <id>")
	}

	id, err := strconv.Atoi(args[2])
	if err != nil {
		return fmt.Errorf("invalid ID: %w", err)
	}

	if id <= 0 {
		return fmt.Errorf("task ID must be positive")
	}

	idx := c.findTaskIndex(id)
	if idx == -1 {
		return fmt.Errorf("task with ID %d not found", id)
	}

	c.tasks = append(c.tasks[:idx], c.tasks[idx+1:]...)
	if err := c.storage.Save(c.tasks); err != nil {
		return fmt.Errorf("failed to save task: %w", err)
	}

	fmt.Println("Task deleted successfully.")
	return nil
}

func (c *CLI) handleMarkInProgress(args []string) error {
	if len(args) < 3 {
		return fmt.Errorf("usage: task-tracker-cli mark-in-progress <id>")
	}

	id, err := strconv.Atoi(args[2])
	if err != nil {
		return fmt.Errorf("invalid ID: %w", err)
	}

	if id <= 0 {
		return fmt.Errorf("task ID must be positive")
	}

	t, err := c.findTask(id)
	if err != nil {
		return err
	}

	t.MarkInProgress()
	if err := c.storage.Save(c.tasks); err != nil {
		return fmt.Errorf("failed to save task: %w", err)
	}

	fmt.Println("Task marked as in progress.")
	return nil
}

func (c *CLI) handleMarkDone(args []string) error {
	if len(args) < 3 {
		return fmt.Errorf("usage: task-tracker-cli mark-done <id>")
	}

	id, err := strconv.Atoi(args[2])
	if err != nil {
		return fmt.Errorf("invalid ID: %w", err)
	}

	if id <= 0 {
		return fmt.Errorf("task ID must be positive")
	}

	t, err := c.findTask(id)
	if err != nil {
		return err
	}

	t.MarkDone()
	if err := c.storage.Save(c.tasks); err != nil {
		return fmt.Errorf("failed to save task: %w", err)
	}

	fmt.Println("Task marked as done.")
	return nil
}

// handleList displays all tasks, optionally filtered by status.
// If no status is provided, all tasks are shown.
func (c *CLI) handleList(args []string) {
	var filterStatus string
	if len(args) > 2 {
		filterStatus = args[2]
		if !task.IsStatusValid(filterStatus) {
			fmt.Printf("Invalid status: %s\n", filterStatus)
			return
		}
	}

	for _, t := range c.tasks {
		if filterStatus == "" || t.Status == filterStatus {
			c.printTask(t)
		}
	}
}

func (c *CLI) findTask(id int) (*task.Task, error) {
	for _, t := range c.tasks {
		if t.ID == id {
			return t, nil
		}
	}
	return nil, fmt.Errorf("task with ID %d not found", id)
}

func (c *CLI) findTaskIndex(id int) int {
	for i, t := range c.tasks {
		if t.ID == id {
			return i
		}
	}
	return -1
}

// nextID generates the next available task ID by finding the maximum
// existing ID and incrementing it. This ensures unique IDs even after deletions.
func (c *CLI) nextID() int {
	maxID := 0
	for _, t := range c.tasks {
		if t.ID > maxID {
			maxID = t.ID
		}
	}
	return maxID + 1
}

func (c *CLI) printTask(t *task.Task) {
	fmt.Printf("ID: %d, Description: %s, Status: %s, CreatedAt: %s, UpdatedAt: %s\n",
		t.ID,
		t.Description,
		t.Status,
		t.CreatedAt.Format("2006-01-02 15:04:05"),
		t.UpdatedAt.Format("2006-01-02 15:04:05"))
}

func (c *CLI) printUsage() {
	fmt.Println("Usage: task-tracker-cli [command] [args]")
	fmt.Println("\nCommands:")
	fmt.Println("  add <description>              Add a new task")
	fmt.Println("  update <id> <description>     Update task description")
	fmt.Println("  delete <id>                    Delete a task")
	fmt.Println("  mark-in-progress <id>          Mark task as in progress")
	fmt.Println("  mark-done <id>                 Mark task as done")
	fmt.Println("  list [status]                  List all tasks or filter by status")
}

