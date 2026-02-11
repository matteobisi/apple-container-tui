package services

import "container-tui/src/models"

// StartContainerBuilder builds a start command.
type StartContainerBuilder struct {
	ContainerID string
}

// Validate ensures a container ID is provided.
func (b StartContainerBuilder) Validate() error {
	_, err := normalizeRequiredToken(b.ContainerID, "container id")
	if err != nil {
		return err
	}
	return nil
}

// Build returns the start command.
func (b StartContainerBuilder) Build() (models.Command, error) {
	if err := b.Validate(); err != nil {
		return models.Command{}, err
	}
	containerID, _ := normalizeRequiredToken(b.ContainerID, "container id")
	return models.Command{Executable: "container", Args: []string{"start", containerID}}, nil
}
