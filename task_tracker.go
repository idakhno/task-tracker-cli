package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"
)

// Task represents a task in the task tracker with an ID, description, and status.
type Task struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	Status      string    `json:"status"` // "not_done", "in_progress", "done"
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

var tasks []Task

func loadTasks(filename string) error {
	file, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	return json.Unmarshal(file, &tasks)
}

func saveTasks(filename string) error {
	data, err := json.MarshalIndent(tasks, "", " ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

func updateTask(id int, description string) error {
	for i, task := range tasks {
		if task.ID == id {
			tasks[i].Description = description
			tasks[i].UpdatedAt = time.Now()
			return saveTasks("tasks.json")
		}
	}
	return fmt.Errorf("task with ID %d not found", id)
}

func deleteTask(id int) error {
	for i, task := range tasks {
		if task.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			return saveTasks("tasks.json")
		}
	}
	return fmt.Errorf("task with ID %d not found", id)
}

func markInProgress(id int) error {
	for i, task := range tasks {
		if task.ID == id {
			tasks[i].Status = "in-progress"
			tasks[i].UpdatedAt = time.Now()
			return saveTasks("tasks.json")
		}
	}
	return fmt.Errorf("task with ID %d not found", id)
}

func markDone(id int) error {
	for i, task := range tasks {
		if task.ID == id {
			tasks[i].Status = "done"
			tasks[i].UpdatedAt = time.Now()
			return saveTasks("tasks.json")
		}
	}
	return fmt.Errorf("task with ID %d not found", id)
}

func listTasks(status string) {
	for _, task := range tasks {
		if status == "" || task.Status == status {
			fmt.Printf("ID: %d, Description: %s, Status: %s, CreatedAt: %s, UpdatedAt: %s\n",
				task.ID,
				task.Description,
				task.Status,
				task.CreatedAt.Format("2006-01-02 15:04:05"),
				task.UpdatedAt.Format("2006-01-02 15:04:05"))
		}
	}
}

func main() {
	const filename = "tasks.json"

	if err := loadTasks(filename); err != nil {
		fmt.Println("Error loading tasks:", err)
		return
	}

	fmt.Println("Welcome to Task Tracker!")

	if len(os.Args) < 2 {
		fmt.Println("Usage: task-tracker-cli [command] [args]")
		return
	}

	command := os.Args[1]

	switch command {
	case "add":
		if len(os.Args) < 3 {
			fmt.Println("Please specify a task description.")
			return
		}
		description := os.Args[2]

		newTask := Task{
			ID:          len(tasks) + 1,
			Description: description,
			Status:      "not_done",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		tasks = append(tasks, newTask)

		if err := saveTasks(filename); err != nil {
			fmt.Println("Error saving tasks:", err)
			return
		}
		fmt.Printf("Task added: ID: %d, Description: %s, Status: %s, CreatedAt: %s, UpdatedAt: %s\n",
			newTask.ID,
			newTask.Description,
			newTask.Status,
			newTask.CreatedAt.Format("2006-01-02 15:04:05"),
			newTask.UpdatedAt.Format("2006-01-02 15:04:05"))

	case "update":
		if len(os.Args) < 4 {
			fmt.Println("Usage: task-tracker-cli update [id] [description]")
			return
		}
		idStr := os.Args[2]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			fmt.Println("Invalid ID:", err)
			return
		}
		description := os.Args[3]
		if err := updateTask(id, description); err != nil {
			fmt.Println("Error updating task:", err)
			return
		}
		fmt.Println("Task updated successfully.")

	case "delete":
		if len(os.Args) < 3 {
			fmt.Println("Usage: task-tracker-cli delete [id]")
			return
		}
		idStr := os.Args[2]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			fmt.Println("Invalid ID:", err)
			return
		}
		if err := deleteTask(id); err != nil {
			fmt.Println("Error deleting task:", err)
			return
		}
		fmt.Println("Task deleted successfully.")

	case "mark-in-progress":
		if len(os.Args) < 3 {
			fmt.Println("Usage: task-tracker-cli mark-in-progress [id]")
			return
		}
		idStr := os.Args[2]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			fmt.Println("Invalid ID:", err)
			return
		}
		if err := markInProgress(id); err != nil {
			fmt.Println("Error marking task as in progress:", err)
			return
		}
		fmt.Println("Task marked as in progress.")

	case "mark-done":
		if len(os.Args) < 3 {
			fmt.Println("Usage: task-tracker-cli mark-done [id]")
			return
		}
		idStr := os.Args[2]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			fmt.Println("Invalid ID:", err)
			return
		}
		if err := markDone(id); err != nil {
			fmt.Println("Error marking task as done:", err)
			return
		}
		fmt.Println("Task marked as done.")

	case "list":
		var status string
		if len(os.Args) > 2 {
			status = os.Args[2]
		}
		listTasks(status)

	default:
		fmt.Println("Unknown command:", command)
	}
}
