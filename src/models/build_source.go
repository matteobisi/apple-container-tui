package models

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

// BuildFileType describes a supported container build file.
type BuildFileType string

const (
	// BuildFileTypeContainerfile represents a Containerfile build script.
	BuildFileTypeContainerfile BuildFileType = "Containerfile"
	// BuildFileTypeDockerfile represents a Dockerfile build script.
	BuildFileTypeDockerfile BuildFileType = "Dockerfile"
)

// BuildSource represents a build file and its context.
type BuildSource struct {
	FilePath         string
	FileType         BuildFileType
	WorkingDirectory string
	Exists           bool
}

// Validate ensures required fields and file existence.
func (b BuildSource) Validate() error {
	if strings.TrimSpace(b.FilePath) == "" {
		return errors.New("file path is required")
	}
	if strings.TrimSpace(b.WorkingDirectory) == "" {
		return errors.New("working directory is required")
	}
	if b.FileType != BuildFileTypeContainerfile && b.FileType != BuildFileTypeDockerfile {
		return errors.New("file type is invalid")
	}
	info, err := os.Stat(b.WorkingDirectory)
	if err != nil || !info.IsDir() {
		return errors.New("working directory is not a directory")
	}
	if _, err := os.Stat(b.FilePath); err != nil {
		return errors.New("build file does not exist")
	}
	return nil
}

// WithComputedExists returns a copy with Exists filled.
func (b BuildSource) WithComputedExists() BuildSource {
	if strings.TrimSpace(b.FilePath) == "" {
		return b
	}
	if _, err := os.Stat(b.FilePath); err == nil {
		b.Exists = true
	}
	if b.WorkingDirectory == "" {
		b.WorkingDirectory = filepath.Dir(b.FilePath)
	}
	return b
}
