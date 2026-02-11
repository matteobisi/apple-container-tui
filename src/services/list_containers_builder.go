package services

import (
	"container-tui/src/models"
)

// ListContainersBuilder builds the list containers command.
type ListContainersBuilder struct{}

// Validate returns nil because there are no inputs.
func (ListContainersBuilder) Validate() error {
	return nil
}

// Build returns the list command.
func (ListContainersBuilder) Build() (models.Command, error) {
	return models.Command{Executable: "container", Args: []string{"list", "--all"}}, nil
}
