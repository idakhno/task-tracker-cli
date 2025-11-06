# Task Tracker CLI

A simple command-line utility for managing tasks.

## Installation

```bash
go build -o task-tracker ./cmd/task-tracker
```

## Usage

```bash
# Add a new task
task-tracker add "Buy groceries"

# Update a task
task-tracker update 1 "Buy groceries and cook dinner"

# Delete a task
task-tracker delete 1

# Mark task as in progress
task-tracker mark-in-progress 1

# Mark task as done
task-tracker mark-done 1

# List all tasks
task-tracker list

# List tasks by status
task-tracker list not_done
task-tracker list in_progress
task-tracker list done
```

Tasks are saved to `tasks.json` in the current directory.
