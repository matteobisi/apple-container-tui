package services

import "container-tui/src/models"

// DeleteContainerBuilder builds a delete container command.
type DeleteContainerBuilder struct {
	ContainerID string
}

// Validate ensures a container ID is provided.
func (b DeleteContainerBuilder) Validate() error {
	_, err := normalizeRequiredToken(b.ContainerID, "container id")
	if err != nil {
		return err
	}
	return nil
}

// Build returns the delete command.
func (b DeleteContainerBuilder) Build() (models.Command, error) {
	if err := b.Validate(); err != nil {
		return models.Command{}, err
	}
	containerID, _ := normalizeRequiredToken(b.ContainerID, "container id")
	return models.Command{Executable: "container", Args: []string{"delete", containerID}}, nil
}
