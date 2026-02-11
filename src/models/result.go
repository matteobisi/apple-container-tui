package models

import "time"

// ResultStatus describes command execution state.
type ResultStatus string

const (
	// ResultPending indicates a command has not started yet.
	ResultPending ResultStatus = "pending"
	// ResultRunning indicates a command is currently executing.
	ResultRunning ResultStatus = "running"
	// ResultSuccess indicates a command completed successfully.
	ResultSuccess ResultStatus = "success"
	// ResultError indicates a command failed.
	ResultError ResultStatus = "error"
	// ResultCanceled indicates a command was canceled.
	ResultCanceled ResultStatus = "cancelled"
)

// Result captures command execution output.
type Result struct {
	ExitCode int
	Stdout   string
	Stderr   string
	Duration time.Duration
	Status   ResultStatus
}
