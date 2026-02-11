package services

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"

	"container-tui/src/models"
)

// LogEntry represents a command execution entry.
type LogEntry struct {
	Command    string    `json:"command"`
	DryRun     bool      `json:"dry_run"`
	ExitCode   int       `json:"exit_code"`
	Stdout     string    `json:"stdout,omitempty"`
	Stderr     string    `json:"stderr,omitempty"`
	StartTime  time.Time `json:"start_time"`
	DurationMs int64     `json:"duration_ms"`
	Status     string    `json:"status"`
}

// LogWriter appends command log entries to a JSON lines file.
type LogWriter struct {
	path          string
	retentionDays int
}

// NewLogWriter creates a log writer for the default log path.
func NewLogWriter(retentionDays int) (*LogWriter, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	logPath := filepath.Join(home, "Library", "Application Support", "apple-tui", "command.log")
	return &LogWriter{path: logPath, retentionDays: retentionDays}, nil
}

// Write appends a log entry to the log file.
func (w *LogWriter) Write(entry LogEntry) error {
	if w == nil {
		return errors.New("log writer is nil")
	}

	if err := w.rotateIfNeeded(); err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(w.path), 0o755); err != nil {
		return err
	}
	file, err := os.OpenFile(w.path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o600)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	encoder := json.NewEncoder(file)
	return encoder.Encode(entry)
}

func (w *LogWriter) rotateIfNeeded() error {
	if w.retentionDays <= 0 {
		return nil
	}
	if _, err := os.Stat(w.path); err != nil {
		return nil
	}

	cutoff := time.Now().AddDate(0, 0, -w.retentionDays)
	file, err := os.Open(w.path)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	kept := make([]LogEntry, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Bytes()
		var entry LogEntry
		if err := json.Unmarshal(line, &entry); err != nil {
			continue
		}
		if entry.StartTime.IsZero() || entry.StartTime.After(cutoff) {
			kept = append(kept, entry)
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	tempPath := w.path + ".tmp"
	out, err := os.OpenFile(tempPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		return err
	}
	encoder := json.NewEncoder(out)
	for _, entry := range kept {
		if err := encoder.Encode(entry); err != nil {
			_ = out.Close()
			return err
		}
	}
	if err := out.Close(); err != nil {
		return err
	}

	return os.Rename(tempPath, w.path)
}

// BuildLogEntry constructs a log entry from a command result.
func BuildLogEntry(command models.Command, result models.Result, dryRun bool) LogEntry {
	startTime := time.Now().Add(-result.Duration)
	return LogEntry{
		Command:    command.String(),
		DryRun:     dryRun,
		ExitCode:   result.ExitCode,
		Stdout:     result.Stdout,
		Stderr:     result.Stderr,
		StartTime:  startTime,
		DurationMs: result.Duration.Milliseconds(),
		Status:     string(result.Status),
	}
}
