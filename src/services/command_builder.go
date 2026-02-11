package services

import "container-tui/src/models"

// CommandBuilder validates and builds a command.
type CommandBuilder interface {
	Validate() error
	Build() (models.Command, error)
}
