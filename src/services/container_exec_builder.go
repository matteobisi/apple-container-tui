package services

import "container-tui/src/models"

// ContainerExecBuilder builds `container exec -it <containerName> <shell>`.
type ContainerExecBuilder struct {
	ContainerName string
	Shell         string
}

func (b ContainerExecBuilder) Validate() error {
	if _, err := normalizeRequiredToken(b.ContainerName, "container name"); err != nil {
		return err
	}
	if _, err := normalizeRequiredToken(b.Shell, "shell"); err != nil {
		return err
	}
	return nil
}

func (b ContainerExecBuilder) Build() (models.Command, error) {
	if err := b.Validate(); err != nil {
		return models.Command{}, err
	}
	containerName, _ := normalizeRequiredToken(b.ContainerName, "container name")
	shell, _ := normalizeRequiredToken(b.Shell, "shell")
	return models.Command{Executable: "container", Args: []string{"exec", "-it", containerName, shell}}, nil
}
