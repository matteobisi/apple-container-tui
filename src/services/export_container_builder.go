package services

import "container-tui/src/models"

// ExportContainerBuilder builds `container export --image <ref> <container-id>`.
type ExportContainerBuilder struct {
	ContainerID    string
	ImageReference string
}

func (b ExportContainerBuilder) Validate() error {
	if _, err := normalizeRequiredToken(b.ContainerID, "container id"); err != nil {
		return err
	}
	_, err := normalizeRequiredToken(b.ImageReference, "image reference")
	return err
}

func (b ExportContainerBuilder) Build() (models.Command, error) {
	if err := b.Validate(); err != nil {
		return models.Command{}, err
	}
	containerID, _ := normalizeRequiredToken(b.ContainerID, "container id")
	imageRef, _ := normalizeRequiredToken(b.ImageReference, "image reference")
	return models.Command{Executable: "container", Args: []string{"export", "--image", imageRef, containerID}}, nil
}
