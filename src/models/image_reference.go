package models

import (
	"errors"
	"strings"
)

// ImageReference represents a container image reference.
type ImageReference struct {
	Registry   string
	Repository string
	Tag        string
	Digest     string
}

// Validate ensures required fields are present.
func (r ImageReference) Validate() error {
	if strings.TrimSpace(r.Repository) == "" {
		return errors.New("repository is required")
	}
	return nil
}

// String formats the reference for CLI usage.
func (r ImageReference) String() string {
	reference := strings.TrimSpace(r.Repository)
	if r.Registry != "" {
		reference = strings.TrimSpace(r.Registry) + "/" + reference
	}
	if r.Digest != "" {
		return reference + "@" + strings.TrimSpace(r.Digest)
	}
	if r.Tag != "" {
		return reference + ":" + strings.TrimSpace(r.Tag)
	}
	return reference
}
