package models

import (
	"errors"
	"strings"
)

// RegistryLogin represents one configured runtime registry login.
type RegistryLogin struct {
	Hostname     string
	Username     string
	CreatedDate  int64
	ModifiedDate int64
}

// Validate ensures the registry entry has the required identifying data.
func (r RegistryLogin) Validate() error {
	if strings.TrimSpace(r.Hostname) == "" {
		return errors.New("hostname is required")
	}
	return nil
}
