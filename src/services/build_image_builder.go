package services

import (
	"errors"
	"strings"

	"container-tui/src/models"
)

// BuildImageBuilder builds a container build command.
type BuildImageBuilder struct {
	Tag         string
	FilePath    string
	ContextPath string
}

// Validate ensures required inputs are provided.
func (b BuildImageBuilder) Validate() error {
	if _, err := normalizeRequiredToken(b.Tag, "tag"); err != nil {
		return err
	}
	if strings.TrimSpace(b.FilePath) == "" {
		return errors.New("file path is required")
	}
	if strings.TrimSpace(b.ContextPath) == "" {
		return errors.New("context path is required")
	}
	return nil
}

// Build returns the build command.
func (b BuildImageBuilder) Build() (models.Command, error) {
	if err := b.Validate(); err != nil {
		return models.Command{}, err
	}
	tag, _ := normalizeRequiredToken(b.Tag, "tag")
	filePath := strings.TrimSpace(b.FilePath)
	contextPath := strings.TrimSpace(b.ContextPath)
	return models.Command{
		Executable: "container",
		Args:       []string{"build", "-t", tag, "-f", filePath, contextPath},
	}, nil
}
