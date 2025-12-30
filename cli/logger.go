package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Logger handles dual output to console and file
type Logger struct {
	filePath string
	file     *os.File
	enabled  bool
}

// NewLogger creates a new logger instance
func NewLogger() *Logger {
	return &Logger{
		enabled: false,
	}
}

// EnableFileLogging starts logging to a file
func (l *Logger) EnableFileLogging(moduleName string) error {
	// Create logs directory if it doesn't exist
	logsDir := "./logs"
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		return fmt.Errorf("failed to create logs directory: %v", err)
	}

	// Create log file with timestamp
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	l.filePath = filepath.Join(logsDir, fmt.Sprintf("%s_%s.log", moduleName, timestamp))

	file, err := os.Create(l.filePath)
	if err != nil {
		return fmt.Errorf("failed to create log file: %v", err)
	}

	l.file = file
	l.enabled = true

	// Write header
	l.writeToFile(fmt.Sprintf("=== Module Execution Log ===\n"))
	l.writeToFile(fmt.Sprintf("Module: %s\n", moduleName))
	l.writeToFile(fmt.Sprintf("Started: %s\n", time.Now().Format(time.RFC3339)))
	l.writeToFile(fmt.Sprintf("================================\n\n"))

	return nil
}

// Close closes the log file
func (l *Logger) Close() error {
	if l.file != nil {
		l.writeToFile(fmt.Sprintf("\n================================\n"))
		l.writeToFile(fmt.Sprintf("Ended: %s\n", time.Now().Format(time.RFC3339)))
		return l.file.Close()
	}
	return nil
}

// Log writes to both console and file
func (l *Logger) Log(message string) {
	fmt.Println(message)
	if l.enabled {
		l.writeToFile(message + "\n")
	}
}

// Logf writes formatted output to both console and file
func (l *Logger) Logf(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	fmt.Println(message)
	if l.enabled {
		l.writeToFile(message + "\n")
	}
}

// LogSection writes a section header
func (l *Logger) LogSection(title string) {
	l.Log(fmt.Sprintf("\n[%s]", title))
	if l.enabled {
		l.writeToFile(fmt.Sprintf("%s\n", "─────────────────────────────"))
	}
}

// GetFilePath returns the log file path
func (l *Logger) GetFilePath() string {
	return l.filePath
}

// writeToFile writes directly to file
func (l *Logger) writeToFile(content string) {
	if l.file != nil {
		l.file.WriteString(content)
	}
}
