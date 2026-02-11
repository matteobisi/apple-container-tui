package services

import (
	"fmt"
	"strings"
)

// normalizeRequiredToken trims input and rejects whitespace for single-token args.
func normalizeRequiredToken(value, field string) (string, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "", fmt.Errorf("%s is required", field)
	}
	if strings.ContainsAny(trimmed, " \t\n") {
		return "", fmt.Errorf("%s must not contain whitespace", field)
	}
	return trimmed, nil
}
