package services

import "container-tui/src/models"

// CommandExecutor runs a command and returns a result.
type CommandExecutor interface {
	Execute(cmd models.Command) (models.Result, error)
}
