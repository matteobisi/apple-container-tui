package services

import (
	"errors"
	"strings"

	"container-tui/src/models"
)

// ImageSaveBuilder builds `container image save --output <path> <ref>`.
type ImageSaveBuilder struct {
	OutputPath     string
	ImageReference string
}

func (b ImageSaveBuilder) Validate() error {
	if strings.TrimSpace(b.OutputPath) == "" {
		return errors.New("output path is required")
	}
	_, err := normalizeRequiredToken(b.ImageReference, "image reference")
	return err
}

func (b ImageSaveBuilder) Build() (models.Command, error) {
	if err := b.Validate(); err != nil {
		return models.Command{}, err
	}
	imageRef, _ := normalizeRequiredToken(b.ImageReference, "image reference")
	return models.Command{Executable: "container", Args: []string{"image", "save", "--output", strings.TrimSpace(b.OutputPath), imageRef}}, nil
}
