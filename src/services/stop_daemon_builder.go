package services

import "container-tui/src/models"

// StopDaemonBuilder builds a stop daemon command.
type StopDaemonBuilder struct{}

// Validate returns nil because there are no inputs.
func (StopDaemonBuilder) Validate() error {
	return nil
}

// Build returns the stop daemon command.
func (StopDaemonBuilder) Build() (models.Command, error) {
	return models.Command{Executable: "container", Args: []string{"system", "stop"}}, nil
}
