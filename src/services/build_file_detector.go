package services

import (
	"errors"
	"os"
	"path/filepath"

	"container-tui/src/models"
)

// DetectBuildFile finds a Containerfile or Dockerfile in the directory.
func DetectBuildFile(workingDir string) (models.BuildSource, error) {
	if workingDir == "" {
		return models.BuildSource{}, errors.New("working directory is required")
	}

	containerfilePath := filepath.Join(workingDir, "Containerfile")
	if _, err := os.Stat(containerfilePath); err == nil {
		return models.BuildSource{
			FilePath:         containerfilePath,
			FileType:         models.BuildFileTypeContainerfile,
			WorkingDirectory: workingDir,
			Exists:           true,
		}, nil
	}

	dockerfilePath := filepath.Join(workingDir, "Dockerfile")
	if _, err := os.Stat(dockerfilePath); err == nil {
		return models.BuildSource{
			FilePath:         dockerfilePath,
			FileType:         models.BuildFileTypeDockerfile,
			WorkingDirectory: workingDir,
			Exists:           true,
		}, nil
	}

	return models.BuildSource{}, errors.New("no Containerfile or Dockerfile found")
}
