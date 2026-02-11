package services

import (
	"bytes"
	"os/exec"
	"time"

	"container-tui/src/models"
)

// RealExecutor executes commands using os/exec.Command.
type RealExecutor struct{}

// Execute runs the command and captures output.
func (RealExecutor) Execute(cmd models.Command) (models.Result, error) {
	start := time.Now()
	command := exec.Command(cmd.Executable, cmd.Args...)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr

	err := command.Run()
	result := models.Result{
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
		Duration: time.Since(start),
	}

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitErr.ExitCode()
		} else {
			result.ExitCode = -1
		}
		result.Status = models.ResultError
		return result, err
	}

	result.ExitCode = 0
	result.Status = models.ResultSuccess
	return result, nil
}
