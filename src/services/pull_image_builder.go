package services

import "container-tui/src/models"

// PullImageBuilder builds an image pull command.
type PullImageBuilder struct {
	Reference string
}

// Validate ensures reference is provided.
func (b PullImageBuilder) Validate() error {
	_, err := normalizeRequiredToken(b.Reference, "reference")
	if err != nil {
		return err
	}
	return nil
}

// Build returns the pull command.
func (b PullImageBuilder) Build() (models.Command, error) {
	if err := b.Validate(); err != nil {
		return models.Command{}, err
	}
	reference, _ := normalizeRequiredToken(b.Reference, "reference")
	return models.Command{Executable: "container", Args: []string{"image", "pull", reference}}, nil
}
