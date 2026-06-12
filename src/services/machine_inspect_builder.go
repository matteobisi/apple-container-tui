package services

import "container-tui/src/models"

// MachineInspectBuilder builds `container machine inspect <id>`.
type MachineInspectBuilder struct {
	MachineID string
}

func (b MachineInspectBuilder) Validate() error {
	_, err := normalizeRequiredToken(b.MachineID, "machine id")
	return err
}

func (b MachineInspectBuilder) Build() (models.Command, error) {
	if err := b.Validate(); err != nil {
		return models.Command{}, err
	}
	machineID, _ := normalizeRequiredToken(b.MachineID, "machine id")
	return models.Command{Executable: "container", Args: []string{"machine", "inspect", machineID}}, nil
}
