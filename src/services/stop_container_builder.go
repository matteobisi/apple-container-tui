package services

import "container-tui/src/models"

// StopContainerBuilder builds a stop command.
type StopContainerBuilder struct {
	ContainerID string
}

// Validate ensures a container ID is provided.
func (b StopContainerBuilder) Validate() error {
	_, err := normalizeRequiredToken(b.ContainerID, "container id")
	if err != nil {
		return err
	}
	return nil
}

// Build returns the stop command.
func (b StopContainerBuilder) Build() (models.Command, error) {
	if err := b.Validate(); err != nil {
		return models.Command{}, err
	}
	containerID, _ := normalizeRequiredToken(b.ContainerID, "container id")
	return models.Command{Executable: "container", Args: []string{"stop", containerID}}, nil
}
