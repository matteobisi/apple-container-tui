package services

import "container-tui/src/models"

// MachineCreateBuilder builds `container machine create <image> [--name <name>]`.
type MachineCreateBuilder struct {
	Image string
	Name  string
}

func (b MachineCreateBuilder) Validate() error {
	_, err := normalizeRequiredToken(b.Image, "image")
	if err != nil {
		return err
	}
	if b.Name != "" {
		_, err = normalizeRequiredToken(b.Name, "machine name")
	}
	return err
}

func (b MachineCreateBuilder) Build() (models.Command, error) {
	if err := b.Validate(); err != nil {
		return models.Command{}, err
	}
	image, _ := normalizeRequiredToken(b.Image, "image")
	args := []string{"machine", "create", image}
	if b.Name != "" {
		name, _ := normalizeRequiredToken(b.Name, "machine name")
		args = append(args, "--name", name)
	}
	return models.Command{Executable: "container", Args: args}, nil
}
