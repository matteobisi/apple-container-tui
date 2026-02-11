package models

import "time"

// DaemonStatus represents the container daemon state.
type DaemonStatus struct {
	Running     bool
	Version     string
	LastChecked time.Time
}
