package services

import (
	"encoding/json"
	"strings"

	"container-tui/src/models"
)

// ParseRegistryList parses the runtime registry list output.
func ParseRegistryList(output string) ([]models.RegistryLogin, error) {
	trimmed := strings.TrimSpace(output)
	if trimmed == "" {
		return []models.RegistryLogin{}, nil
	}

	var entries []models.RegistryLogin
	if err := json.Unmarshal([]byte(trimmed), &entries); err != nil {
		return nil, err
	}
	validated := make([]models.RegistryLogin, 0, len(entries))
	for _, entry := range entries {
		if err := entry.Validate(); err != nil {
			continue
		}
		validated = append(validated, entry)
	}
	return validated, nil
}
