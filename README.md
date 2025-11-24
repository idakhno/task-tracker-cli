# Task Tracker CLI

A simple command-line utility for managing tasks with file-based storage.

## Features

- ✅ Create, update, and delete tasks
- ✅ Track task status (not_done, in_progress, done)
- ✅ Filter tasks by status
- ✅ File-based storage with JSON format
- ✅ Thread-safe operations with file locking
- ✅ Input validation
- ✅ Comprehensive error handling
- ✅ Logging support
- ✅ Version information

## Installation

### From Source

```bash
go build -o task-tracker ./cmd/task-tracker
```

### Build with Version Information

```bash
go build -ldflags "-X main.version=1.0.0 -X main.buildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ) -X main.gitCommit=$(git rev-parse HEAD)" -o task-tracker ./cmd/task-tracker
```

## Usage

### Basic Commands

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

### Options

```bash
# Specify custom tasks file
task-tracker --file /path/to/tasks.json add "New task"

# Or use environment variable
export TASK_TRACKER_FILE=/path/to/tasks.json
task-tracker add "New task"

# Show version information
task-tracker --version
```

## Configuration

### Environment Variables

- `TASK_TRACKER_FILE`: Path to the tasks JSON file (default: `tasks.json` in current directory)

### File Format

Tasks are stored in JSON format:

```json
[
  {
    "id": 1,
    "description": "Buy groceries",
    "status": "not_done",
    "createdAt": "2024-01-01T10:00:00Z",
    "updatedAt": "2024-01-01T10:00:00Z"
  }
]
```

## Development

### Running Tests

```bash
go test ./...
```

### Running Tests with Coverage

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Code Quality

```bash
# Format code
gofmt -s -w .

# Run linter
golangci-lint run
```

## Security

See [SECURITY.md](SECURITY.md) for security policy and reporting vulnerabilities.

## License

This project is open source and available under the MIT License.
