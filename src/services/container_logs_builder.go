package services

import "container-tui/src/models"

// ContainerLogsBuilder builds `container logs -f <containerName>`.
type ContainerLogsBuilder struct {
	ContainerName string
}

func (b ContainerLogsBuilder) Validate() error {
	_, err := normalizeRequiredToken(b.ContainerName, "container name")
	return err
}

func (b ContainerLogsBuilder) Build() (models.Command, error) {
	if err := b.Validate(); err != nil {
		return models.Command{}, err
	}
	containerName, _ := normalizeRequiredToken(b.ContainerName, "container name")
	return models.Command{Executable: "container", Args: []string{"logs", "-f", containerName}}, nil
}
