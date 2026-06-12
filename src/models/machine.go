package models

import (
	"errors"
	"strings"
)

// MachineState identifies the current runtime state of a container machine.
type MachineState string

const (
	MachineStateRunning MachineState = "running"
	MachineStateStopped MachineState = "stopped"
	MachineStateUnknown MachineState = "unknown"
)

// ContainerMachine represents one Apple Container machine.
type ContainerMachine struct {
	ID        string       `json:"id"`
	Image     string       `json:"image"`
	State     MachineState `json:"state"`
	IsDefault bool         `json:"default"`
	CPUs      int          `json:"cpus"`
	Memory    string       `json:"memory"`
	HomeMount string       `json:"homeMount"`
}

// Validate ensures the machine has the identifying data required for commands.
func (m ContainerMachine) Validate() error {
	if strings.TrimSpace(m.ID) == "" {
		return errors.New("machine id is required")
	}
	return nil
}

// NormalizedState returns a safe state value for UI decisions.
func (m ContainerMachine) NormalizedState() MachineState {
	switch MachineState(strings.ToLower(strings.TrimSpace(string(m.State)))) {
	case MachineStateRunning:
		return MachineStateRunning
	case MachineStateStopped:
		return MachineStateStopped
	default:
		return MachineStateUnknown
	}
}

// NormalizedHomeMount returns a supported home mount value.
func (m ContainerMachine) NormalizedHomeMount() string {
	switch strings.ToLower(strings.TrimSpace(m.HomeMount)) {
	case "ro", "none":
		return strings.ToLower(strings.TrimSpace(m.HomeMount))
	default:
		return "rw"
	}
}

// MachineEditInput captures resource edit form values.
type MachineEditInput struct {
	Name      string
	CPUs      string
	Memory    string
	HomeMount string
}
