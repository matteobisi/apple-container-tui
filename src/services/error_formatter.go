package services

import (
	"strings"
)

// FormatError returns a user-friendly error message based on stderr and err.
func FormatError(err error, stderr string) string {
	message := strings.TrimSpace(stderr)
	if message == "" && err != nil {
		message = err.Error()
	}
	return formatErrorMessage(message)
}

func formatErrorMessage(message string) string {
	if strings.TrimSpace(message) == "" {
		return "unknown error"
	}

	lower := strings.ToLower(message)
	switch {
	case strings.Contains(lower, "not found") || strings.Contains(lower, "no such container") || strings.Contains(lower, "manifest unknown"):
		return "resource not found; check the name or reference"
	case strings.Contains(lower, "already running"):
		return "container is already running"
	case strings.Contains(lower, "already stopped") || strings.Contains(lower, "not running"):
		return "container is already stopped"
	case strings.Contains(lower, "permission") || strings.Contains(lower, "sudo"):
		return "permission denied; check your privileges"
	case strings.Contains(lower, "daemon") || strings.Contains(lower, "connection"):
		return "container daemon is not running; start it from the Daemon screen"
	case strings.Contains(lower, "unauthorized") || strings.Contains(lower, "authentication"):
		return "authentication required; check registry credentials"
	case strings.Contains(lower, "builder") || strings.Contains(lower, "buildkit"):
		return "builder is not running; start it and retry"
	case strings.Contains(lower, "invalid reference format"):
		return "invalid image reference"
	case strings.Contains(lower, "no such file") || strings.Contains(lower, "file not found"):
		return "build file or context path not found"
	case strings.Contains(lower, "not a directory") && strings.Contains(lower, "context"):
		return "build context must be a directory"
	default:
		return message
	}
}
