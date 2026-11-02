package services

import (
	"strings"

	"container-tui/src/models"
)

// CommitContainerBuilder builds `container commit <id> <reference>`.
type CommitContainerBuilder struct {
	ContainerID string
	Reference   string
}

func (b CommitContainerBuilder) Validate() error {
	if _, err := normalizeRequiredToken(b.ContainerID, "container id"); err != nil {
		return err
	}
	_, err := normalizeRequiredToken(b.Reference, "image reference")
	return err
}

func (b CommitContainerBuilder) Build() (models.Command, error) {
	if err := b.Validate(); err != nil {
		return models.Command{}, err
	}
	containerID, _ := normalizeRequiredToken(b.ContainerID, "container id")
	return models.Command{Executable: "container", Args: []string{"commit", containerID, strings.TrimSpace(b.Reference)}}, nil
}
