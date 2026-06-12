package services

import "container-tui/src/models"

// MachineListBuilder builds `container machine list --format json`.
type MachineListBuilder struct{}

func (MachineListBuilder) Validate() error { return nil }

func (MachineListBuilder) Build() (models.Command, error) {
	return models.Command{Executable: "container", Args: []string{"machine", "list", "--format", "json"}}, nil
}
