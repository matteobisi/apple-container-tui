package services

import "container-tui/src/models"

// ImagePruneBuilder builds `container image prune`.
type ImagePruneBuilder struct{}

func (b ImagePruneBuilder) Validate() error { return nil }

func (b ImagePruneBuilder) Build() (models.Command, error) {
	return models.Command{Executable: "container", Args: []string{"image", "prune"}}, nil
}
