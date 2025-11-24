package logger

import (
	"fmt"
	"io"
	"os"
	"time"
)

// Level represents the logging level
type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

var (
	currentLevel = LevelWarn
	output       io.Writer = os.Stderr
)

// SetLevel sets the minimum logging level
func SetLevel(level Level) {
	currentLevel = level
}

// SetOutput sets the output writer for logs
func SetOutput(w io.Writer) {
	output = w
}

// Debug logs a debug message
func Debug(format string, args ...interface{}) {
	log(LevelDebug, "DEBUG", format, args...)
}

// Info logs an info message
func Info(format string, args ...interface{}) {
	log(LevelInfo, "INFO", format, args...)
}

// Warn logs a warning message
func Warn(format string, args ...interface{}) {
	log(LevelWarn, "WARN", format, args...)
}

// Error logs an error message
func Error(format string, args ...interface{}) {
	log(LevelError, "ERROR", format, args...)
}

func log(level Level, levelStr, format string, args ...interface{}) {
	if level < currentLevel {
		return
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	message := fmt.Sprintf(format, args...)
	fmt.Fprintf(output, "[%s] [%s] %s\n", timestamp, levelStr, message)
}

