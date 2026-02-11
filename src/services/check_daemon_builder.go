package services

import "container-tui/src/models"

// CheckDaemonStatusBuilder builds a daemon status command.
type CheckDaemonStatusBuilder struct{}

// Validate returns nil because there are no inputs.
func (CheckDaemonStatusBuilder) Validate() error {
	return nil
}

// Build returns the daemon status command.
func (CheckDaemonStatusBuilder) Build() (models.Command, error) {
	return models.Command{Executable: "container", Args: []string{"system", "status"}}, nil
}
