package services

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"container-tui/src/models"
)

type machineListEntry struct {
	ID        string              `json:"id"`
	Image     string              `json:"image"`
	State     models.MachineState `json:"state"`
	Status    models.MachineState `json:"status"`
	IsDefault bool                `json:"default"`
	CPUs      int                 `json:"cpus"`
	Memory    json.RawMessage     `json:"memory"`
	HomeMount string              `json:"homeMount"`
}

// ParseMachineList parses `container machine list --format json` output.
func ParseMachineList(output string) ([]models.ContainerMachine, error) {
	trimmed := strings.TrimSpace(output)
	if trimmed == "" {
		return []models.ContainerMachine{}, nil
	}

	var entries []machineListEntry
	if err := json.Unmarshal([]byte(trimmed), &entries); err != nil {
		return nil, err
	}

	validated := make([]models.ContainerMachine, 0, len(entries))
	for _, entry := range entries {
		machine := models.ContainerMachine{
			ID:        entry.ID,
			Image:     entry.Image,
			State:     entry.State,
			IsDefault: entry.IsDefault,
			CPUs:      entry.CPUs,
			Memory:    formatMachineMemory(entry.Memory),
			HomeMount: entry.HomeMount,
		}
		if machine.State == "" {
			machine.State = entry.Status
		}
		if err := machine.Validate(); err != nil {
			continue
		}
		machine.State = machine.NormalizedState()
		machine.HomeMount = machine.NormalizedHomeMount()
		validated = append(validated, machine)
	}
	return validated, nil
}

func formatMachineMemory(raw json.RawMessage) string {
	if len(raw) == 0 || string(raw) == "null" {
		return ""
	}
	var text string
	if err := json.Unmarshal(raw, &text); err == nil {
		return strings.TrimSpace(text)
	}
	var bytes int64
	if err := json.Unmarshal(raw, &bytes); err != nil {
		return ""
	}
	const gib = int64(1024 * 1024 * 1024)
	if bytes > 0 && bytes%gib == 0 {
		return fmt.Sprintf("%dG", bytes/gib)
	}
	return strconv.FormatInt(bytes, 10)
}
