package services

import (
	"encoding/json"
	"strings"
	"time"

	"container-tui/src/models"
)

type daemonStatusPayload struct {
	Status           string `json:"status"`
	InstallRoot      string `json:"installRoot"`
	AppRoot          string `json:"appRoot"`
	APIServerVersion string `json:"apiServerVersion"`
}

// ParseDaemonStatus determines daemon state from structured JSON output.
func ParseDaemonStatus(output string) models.DaemonStatus {
	status := models.DaemonStatus{
		State:       models.DaemonStateUnknown,
		Running:     false,
		LastChecked: time.Now(),
	}

	var payload daemonStatusPayload
	if err := json.Unmarshal([]byte(strings.TrimSpace(output)), &payload); err != nil {
		status.MissingFields = []string{"status"}
		return status
	}

	status.Version = strings.TrimSpace(payload.APIServerVersion)
	status.InstallRoot = strings.TrimSpace(payload.InstallRoot)
	status.AppRoot = strings.TrimSpace(payload.AppRoot)

	classification := classifyDaemonState(strings.TrimSpace(payload.Status))
	if classification == models.DaemonStateUnknown {
		status.MissingFields = []string{"status"}
		return status
	}

	status.State = classification
	status.Running = classification == models.DaemonStateRunning
	return status
}

func classifyDaemonState(raw string) models.DaemonState {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "running":
		return models.DaemonStateRunning
	case "stopped", "not running", "stopping", "starting", "down":
		return models.DaemonStateStopped
	default:
		return models.DaemonStateUnknown
	}
}
