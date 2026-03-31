package models

import "time"

// DaemonState describes the classified daemon state.
type DaemonState string

const (
	// DaemonStateRunning indicates the daemon is running.
	DaemonStateRunning DaemonState = "running"
	// DaemonStateStopped indicates the daemon is stopped.
	DaemonStateStopped DaemonState = "stopped"
	// DaemonStateUnknown indicates the daemon state could not be classified safely.
	DaemonStateUnknown DaemonState = "unknown"
)

// DaemonStatus represents the container daemon state.
type DaemonStatus struct {
	State         DaemonState
	Running       bool
	Version       string
	InstallRoot   string
	AppRoot       string
	LastChecked   time.Time
	MissingFields []string
}
