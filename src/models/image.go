package models

import (
	"errors"
	"strings"
)

// Image represents a local container image.
type Image struct {
	Name   string
	Tag    string
	Digest string
}

// Validate checks image fields.
func (i Image) Validate() error {
	if strings.TrimSpace(i.Name) == "" {
		return errors.New("image name is required")
	}
	if strings.TrimSpace(i.Tag) == "" {
		return errors.New("image tag is required")
	}
	if strings.TrimSpace(i.Digest) != "" && !strings.HasPrefix(strings.TrimSpace(i.Digest), "sha256:") {
		return errors.New("image digest must start with sha256:")
	}
	return nil
}

// Reference returns NAME:TAG or NAME@DIGEST if tag is missing.
func (i Image) Reference() string {
	name := strings.TrimSpace(i.Name)
	tag := strings.TrimSpace(i.Tag)
	digest := strings.TrimSpace(i.Digest)
	if tag != "" && tag != "<none>" {
		return name + ":" + tag
	}
	if digest != "" {
		return name + "@" + digest
	}
	return name
}
