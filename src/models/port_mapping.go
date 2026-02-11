package models

import "errors"

// PortMapping represents a host-to-container port binding.
type PortMapping struct {
	HostPort      int
	ContainerPort int
	Protocol      string
}

// Validate ensures port mapping values are in range.
func (p PortMapping) Validate() error {
	if p.HostPort < 1 || p.HostPort > 65535 {
		return errors.New("host port out of range")
	}
	if p.ContainerPort < 1 || p.ContainerPort > 65535 {
		return errors.New("container port out of range")
	}
	if p.Protocol == "" {
		return errors.New("protocol is required")
	}
	return nil
}
