package services

import (
	"fmt"
	"time"

	"container-tui/src/models"
)

// DryRunExecutor echoes commands without executing.
type DryRunExecutor struct{}

// Execute returns a successful dry-run result.
func (DryRunExecutor) Execute(cmd models.Command) (models.Result, error) {
	start := time.Now()
	return models.Result{
		ExitCode: 0,
		Stdout:   fmt.Sprintf("dry-run: %s", cmd.String()),
		Duration: time.Since(start),
		Status:   models.ResultSuccess,
	}, nil
}
