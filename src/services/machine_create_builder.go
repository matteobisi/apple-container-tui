package services

import (
	"fmt"
	"strconv"
	"strings"

	"container-tui/src/models"
)

// MachineCreateBuilder builds `container machine create <image> [--name <name>] [--cpus <n>] [--memory <size>] [--home-mount <mode>]`.
type MachineCreateBuilder struct {
	Image     string
	Name      string
	CPUs      string
	Memory    string
	HomeMount string
}

func (b MachineCreateBuilder) Validate() error {
	_, err := normalizeRequiredToken(b.Image, "image")
	if err != nil {
		return err
	}
	if b.Name != "" {
		_, err = normalizeRequiredToken(b.Name, "machine name")
	}
	if err != nil {
		return err
	}
	if strings.TrimSpace(b.CPUs) != "" {
		value, parseErr := strconv.Atoi(strings.TrimSpace(b.CPUs))
		if parseErr != nil || value < 1 {
			return fmt.Errorf("cpus must be a positive integer")
		}
	}
	if strings.TrimSpace(b.Memory) != "" {
		if _, memErr := normalizeRequiredToken(b.Memory, "memory"); memErr != nil {
			return memErr
		}
	}
	if strings.TrimSpace(b.HomeMount) != "" {
		homeMount := strings.ToLower(strings.TrimSpace(b.HomeMount))
		if homeMount != "rw" && homeMount != "ro" && homeMount != "none" {
			return fmt.Errorf("home-mount must be one of rw, ro, none")
		}
	}
	return nil
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
	if strings.TrimSpace(b.CPUs) != "" {
		args = append(args, "--cpus", strings.TrimSpace(b.CPUs))
	}
	if strings.TrimSpace(b.Memory) != "" {
		memory, _ := normalizeRequiredToken(b.Memory, "memory")
		args = append(args, "--memory", memory)
	}
	if strings.TrimSpace(b.HomeMount) != "" {
		args = append(args, "--home-mount", strings.ToLower(strings.TrimSpace(b.HomeMount)))
	}
	return models.Command{Executable: "container", Args: args}, nil
}
