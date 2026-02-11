package services

import (
	"strings"
	"time"

	"container-tui/src/models"
)

// ParseDaemonStatus determines daemon running state from output.
func ParseDaemonStatus(output string) models.DaemonStatus {
	lower := strings.ToLower(output)
	running := strings.Contains(lower, "running") && !strings.Contains(lower, "not running")
	return models.DaemonStatus{
		Running:     running,
		Version:     "",
		LastChecked: time.Now(),
	}
}
