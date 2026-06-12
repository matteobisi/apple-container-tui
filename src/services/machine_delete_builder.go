package services

import "container-tui/src/models"

// MachineDeleteBuilder builds `container machine delete <id>`.
type MachineDeleteBuilder struct {
	MachineID string
}

func (b MachineDeleteBuilder) Validate() error {
	_, err := normalizeRequiredToken(b.MachineID, "machine id")
	return err
}

func (b MachineDeleteBuilder) Build() (models.Command, error) {
	if err := b.Validate(); err != nil {
		return models.Command{}, err
	}
	machineID, _ := normalizeRequiredToken(b.MachineID, "machine id")
	return models.Command{Executable: "container", Args: []string{"machine", "delete", machineID}}, nil
}
