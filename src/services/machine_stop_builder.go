package services

import "container-tui/src/models"

// MachineStopBuilder builds `container machine stop <id>`.
type MachineStopBuilder struct {
	MachineID string
}

func (b MachineStopBuilder) Validate() error {
	_, err := normalizeRequiredToken(b.MachineID, "machine id")
	return err
}

func (b MachineStopBuilder) Build() (models.Command, error) {
	if err := b.Validate(); err != nil {
		return models.Command{}, err
	}
	machineID, _ := normalizeRequiredToken(b.MachineID, "machine id")
	return models.Command{Executable: "container", Args: []string{"machine", "stop", machineID}}, nil
}
