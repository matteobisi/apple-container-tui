package services

import "container-tui/src/models"

// MachineSetDefaultBuilder builds `container machine set-default <id>`.
type MachineSetDefaultBuilder struct {
	MachineID string
}

func (b MachineSetDefaultBuilder) Validate() error {
	_, err := normalizeRequiredToken(b.MachineID, "machine id")
	return err
}

func (b MachineSetDefaultBuilder) Build() (models.Command, error) {
	if err := b.Validate(); err != nil {
		return models.Command{}, err
	}
	machineID, _ := normalizeRequiredToken(b.MachineID, "machine id")
	return models.Command{Executable: "container", Args: []string{"machine", "set-default", machineID}}, nil
}
