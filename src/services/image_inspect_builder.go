package services

import "container-tui/src/models"

// ImageInspectBuilder builds `container image inspect <ref> | jq` (jq can be applied by shell wrapper).
type ImageInspectBuilder struct {
	ImageReference string
}

func (b ImageInspectBuilder) Validate() error {
	_, err := normalizeRequiredToken(b.ImageReference, "image reference")
	return err
}

func (b ImageInspectBuilder) Build() (models.Command, error) {
	if err := b.Validate(); err != nil {
		return models.Command{}, err
	}
	ref, _ := normalizeRequiredToken(b.ImageReference, "image reference")
	return models.Command{Executable: "container", Args: []string{"image", "inspect", ref}}, nil
}
