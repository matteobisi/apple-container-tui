package services

import (
	"fmt"
	"strconv"
	"strings"

	"container-tui/src/models"
)

// MachineSetBuilder builds `container machine set -n <id> cpus=<N> memory=<M> home-mount=<mode>`.
type MachineSetBuilder struct {
	MachineID string
	CPUs      string
	Memory    string
	HomeMount string
}

func (b MachineSetBuilder) Validate() error {
	if _, err := normalizeRequiredToken(b.MachineID, "machine id"); err != nil {
		return err
	}
	cpus, err := normalizeRequiredToken(b.CPUs, "cpus")
	if err != nil {
		return err
	}
	parsedCPUs, err := strconv.Atoi(cpus)
	if err != nil || parsedCPUs < 1 {
		return fmt.Errorf("cpus must be a positive integer")
	}
	if _, err := normalizeRequiredToken(b.Memory, "memory"); err != nil {
		return err
	}
	homeMount := strings.ToLower(strings.TrimSpace(b.HomeMount))
	if homeMount != "rw" && homeMount != "ro" && homeMount != "none" {
		return fmt.Errorf("home-mount must be one of rw, ro, none")
	}
	return nil
}

func (b MachineSetBuilder) Build() (models.Command, error) {
	if err := b.Validate(); err != nil {
		return models.Command{}, err
	}
	machineID, _ := normalizeRequiredToken(b.MachineID, "machine id")
	cpus, _ := normalizeRequiredToken(b.CPUs, "cpus")
	memory, _ := normalizeRequiredToken(b.Memory, "memory")
	homeMount := strings.ToLower(strings.TrimSpace(b.HomeMount))
	return models.Command{Executable: "container", Args: []string{"machine", "set", "-n", machineID, "cpus=" + cpus, "memory=" + memory, "home-mount=" + homeMount}}, nil
}
