package services

import "container-tui/src/models"

// MachineStartBuilder builds `container machine run -n <id>`.
type MachineStartBuilder struct {
	MachineID string
}

func (b MachineStartBuilder) Validate() error {
	_, err := normalizeRequiredToken(b.MachineID, "machine id")
	return err
}

func (b MachineStartBuilder) Build() (models.Command, error) {
	if err := b.Validate(); err != nil {
		return models.Command{}, err
	}
	machineID, _ := normalizeRequiredToken(b.MachineID, "machine id")
	return models.Command{Executable: "container", Args: []string{"machine", "run", "-n", machineID}}, nil
}
