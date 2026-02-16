package services

import "container-tui/src/models"

// ImageDeleteBuilder builds `container image rm <ref>`.
type ImageDeleteBuilder struct {
	ImageReference string
}

func (b ImageDeleteBuilder) Validate() error {
	_, err := normalizeRequiredToken(b.ImageReference, "image reference")
	return err
}

func (b ImageDeleteBuilder) Build() (models.Command, error) {
	if err := b.Validate(); err != nil {
		return models.Command{}, err
	}
	ref, _ := normalizeRequiredToken(b.ImageReference, "image reference")
	return models.Command{Executable: "container", Args: []string{"image", "rm", ref}}, nil
}
