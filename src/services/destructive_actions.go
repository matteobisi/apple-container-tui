package services

// DestructiveAction identifies an action that needs extra confirmation.
type DestructiveAction string

const (
	// ActionDeleteContainer represents a container delete operation.
	ActionDeleteContainer DestructiveAction = "delete-container"
	// ActionStopDaemon represents a daemon stop operation.
	ActionStopDaemon DestructiveAction = "stop-daemon"
)

// ActionMetadata describes a destructive action in the UI.
type ActionMetadata struct {
	Label       string
	Description string
}

// DestructiveActionMetadata returns metadata for destructive actions.
func DestructiveActionMetadata() map[DestructiveAction]ActionMetadata {
	return map[DestructiveAction]ActionMetadata{
		ActionDeleteContainer: {
			Label:       "Delete Container",
			Description: "Permanently deletes a stopped container",
		},
		ActionStopDaemon: {
			Label:       "Stop Daemon",
			Description: "Stops container services and running containers",
		},
	}
}
