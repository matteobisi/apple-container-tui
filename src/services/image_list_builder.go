package services

import "container-tui/src/models"

// ImageListBuilder builds `container image list`.
type ImageListBuilder struct{}

func (b ImageListBuilder) Validate() error { return nil }

func (b ImageListBuilder) Build() (models.Command, error) {
	return models.Command{Executable: "container", Args: []string{"image", "list"}}, nil
}
