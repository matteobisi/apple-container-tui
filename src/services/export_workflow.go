package services

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"container-tui/src/models"
)

// ContainerExportPlan captures the previewable export workflow state.
type ContainerExportPlan struct {
	Container            models.Container
	DestinationDirectory string
	GeneratedImageRef    string
	ArchivePath          string
	Commands             []models.Command
	CleanupCommand       models.Command
}

// ExportWorkflowResult captures the final export result.
type ExportWorkflowResult struct {
	Result      models.Result
	ArchivePath string
}

// ExportWorkflowService plans and executes the container export workflow.
type ExportWorkflowService struct {
	Executor CommandExecutor
	Now      func() time.Time
}

func NewExportWorkflowService(executor CommandExecutor) ExportWorkflowService {
	return ExportWorkflowService{Executor: executor, Now: time.Now}
}

// Plan creates the command sequence for exporting a stopped container.
func (s ExportWorkflowService) Plan(container models.Container, destinationDirectory string) (ContainerExportPlan, error) {
	if container.Status != models.ContainerStatusStopped {
		return ContainerExportPlan{}, errors.New("only stopped containers can be exported")
	}
	destination := strings.TrimSpace(destinationDirectory)
	if destination == "" {
		return ContainerExportPlan{}, errors.New("destination directory is required")
	}
	info, err := os.Stat(destination)
	if err != nil {
		return ContainerExportPlan{}, fmt.Errorf("destination directory is not available: %w", err)
	}
	if !info.IsDir() {
		return ContainerExportPlan{}, errors.New("destination must be a directory")
	}

	now := s.Now()
	imageRef := models.BuildExportImageReference(container.Name, container.ID, now)
	archivePath := filepath.Join(destination, models.BuildExportArchiveName(container.Name, container.ID, now))
	if _, err := os.Stat(archivePath); err == nil {
		return ContainerExportPlan{}, errors.New("generated export archive already exists")
	} else if !errors.Is(err, os.ErrNotExist) {
		return ContainerExportPlan{}, fmt.Errorf("cannot validate export archive path: %w", err)
	}

	exportCmd, err := (ExportContainerBuilder{ContainerID: container.ID, ImageReference: imageRef}).Build()
	if err != nil {
		return ContainerExportPlan{}, err
	}
	saveCmd, err := (ImageSaveBuilder{OutputPath: archivePath, ImageReference: imageRef}).Build()
	if err != nil {
		return ContainerExportPlan{}, err
	}
	cleanupCmd, err := (ImageDeleteBuilder{ImageReference: imageRef}).Build()
	if err != nil {
		return ContainerExportPlan{}, err
	}

	return ContainerExportPlan{
		Container:            container,
		DestinationDirectory: destination,
		GeneratedImageRef:    imageRef,
		ArchivePath:          archivePath,
		Commands:             []models.Command{exportCmd, saveCmd},
		CleanupCommand:       cleanupCmd,
	}, nil
}

// Execute runs the export and save steps in order.
func (s ExportWorkflowService) Execute(plan ContainerExportPlan) (ExportWorkflowResult, error) {
	if s.Executor == nil {
		return ExportWorkflowResult{}, errors.New("executor is required")
	}
	if len(plan.Commands) < 2 {
		return ExportWorkflowResult{}, errors.New("export plan requires at least export and save commands")
	}

	var stdoutParts []string
	var stderrParts []string
	stepNames := []string{"export", "save", "cleanup"}

	for index, command := range plan.Commands[:2] {
		result, err := s.Executor.Execute(command)
		if strings.TrimSpace(result.Stdout) != "" {
			stdoutParts = append(stdoutParts, fmt.Sprintf("[%s]\n%s", stepNames[index], strings.TrimSpace(result.Stdout)))
		}
		if strings.TrimSpace(result.Stderr) != "" {
			stderrParts = append(stderrParts, fmt.Sprintf("[%s]\n%s", stepNames[index], strings.TrimSpace(result.Stderr)))
		}
		if err != nil {
			combined := models.Result{
				Status: models.ResultError,
				Stdout: strings.Join(stdoutParts, "\n\n"),
				Stderr: strings.Join(append(stderrParts, err.Error()), "\n\n"),
			}
			return ExportWorkflowResult{Result: combined, ArchivePath: plan.ArchivePath}, err
		}
	}

	finalResult := models.Result{
		Status: models.ResultSuccess,
		Stdout: strings.TrimSpace(strings.Join(append([]string{fmt.Sprintf("Exported OCI archive: %s", plan.ArchivePath)}, stdoutParts...), "\n\n")),
	}

	return ExportWorkflowResult{Result: finalResult, ArchivePath: plan.ArchivePath}, nil
}
