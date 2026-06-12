package services

import "container-tui/src/models"

// MachineLogsBuilder builds `container machine logs <id>`.
type MachineLogsBuilder struct {
	MachineID string
}

func (b MachineLogsBuilder) Validate() error {
	_, err := normalizeRequiredToken(b.MachineID, "machine id")
	return err
}

func (b MachineLogsBuilder) Build() (models.Command, error) {
	if err := b.Validate(); err != nil {
		return models.Command{}, err
	}
	machineID, _ := normalizeRequiredToken(b.MachineID, "machine id")
	return models.Command{Executable: "container", Args: []string{"machine", "logs", machineID}}, nil
}
