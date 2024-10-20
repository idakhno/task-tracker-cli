# Task Tracker CLI
A simple command-line interface (CLI) application to track and manage your tasks. This project helps you practice programming skills, including file system interaction, user input handling, and building a CLI application.

## 📚 Table of Contents
- [🚀 Features](#-features)
- [⚙️ Task Properties](#-task-properties)
- [🛠️ Requirements](#-requirements)
- [Usage](#usage)
    - [✔️ Adding a New Task](#adding-a-new-task)
    - [✔️ Updating and Deleting Tasks](#updating-and-deleting-tasks)
    - [✔️ Marking a Task as In Progress or Done](#marking-a-task-as-in-progress-or-done)
    - [✔️ Listing All Tasks](#listing-all-tasks)
    - [✔️ Listing Tasks by Status](#listing-tasks-by-status)
- [💾 Installation](#installation)
- [👋 Contributing](#contributing)
- [©️ License](#license)

## 🚀 Features
- Add, Update, and Delete Tasks: Manage your tasks with ease.
- Mark Tasks: Mark tasks as in progress or done.
- List Tasks: View all tasks or filter by their status.

## ⚙️ Task Properties
Each task has the following properties:
- id: A unique identifier for the task;
- description: A short description of the task;
- status: The status of the task (not_done, in_progress, done);
- createdAt: The date and time when the task was created;
- updatedAt: The date and time when the task was last updated.

## 🛠️ Requirements
- The application should run from the command line and accept user actions and inputs as arguments.
- Tasks should be stored in a JSON file in the current directory.
- The JSON file should be created if it does not exist.
- Use the native file system module for file interaction.
- Handle errors and edge cases gracefully.

## Usage
### ✔️ Adding a New Task
```
task-tracker-cli add "Buy groceries"
```
Output: Task added successfully (ID: 1)

### ✔️ Updating and Deleting Tasks
```
task-tracker-cli update 1 "Buy groceries and cook dinner"
task-tracker-cli delete 1
```

### ✔️ Marking a Task as In Progress or Done
```
task-tracker-cli mark-in-progress 1
task-tracker-cli mark-done 1
```

### ✔️ Listing All Tasks
```
task-tracker-cli list
```

### ✔️ Listing Tasks by Status
```
task-tracker-cli list done
task-tracker-cli list not_done
task-tracker-cli list in_progress
```

## 💾 Installation

Clone the repository:
```
git clone https://github.com/idakhno/task-tracker-cli.git
cd task-tracker-cli
```

Run the application:
```
go run task_tracker.go
```

## 👋 Contributing
Feel free to submit issues or pull requests.

## ©️ License
This project is licensed under the MIT License.