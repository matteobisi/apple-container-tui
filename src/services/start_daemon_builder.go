package services

import "container-tui/src/models"

// StartDaemonBuilder builds a start daemon command.
type StartDaemonBuilder struct{}

// Validate returns nil because there are no inputs.
func (StartDaemonBuilder) Validate() error {
	return nil
}

// Build returns the start daemon command.
func (StartDaemonBuilder) Build() (models.Command, error) {
	return models.Command{Executable: "container", Args: []string{"system", "start"}}, nil
}
