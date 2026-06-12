package models

import (
	"errors"
	"strings"
	"time"
)

// RegistryLogin represents one configured runtime registry login.
type RegistryLogin struct {
	Hostname     string    `json:"name"`
	Username     string    `json:"username"`
	CreatedDate  time.Time `json:"creationDate"`
	ModifiedDate time.Time `json:"modificationDate"`
}

// Validate ensures the registry entry has the required identifying data.
func (r RegistryLogin) Validate() error {
	if strings.TrimSpace(r.Hostname) == "" {
		return errors.New("hostname is required")
	}
	return nil
}
