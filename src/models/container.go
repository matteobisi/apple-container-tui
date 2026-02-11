package models

import (
	"errors"
	"strings"
)

// ContainerStatus captures the lifecycle state of a container.
type ContainerStatus string

const (
	// ContainerStatusRunning indicates a running container.
	ContainerStatusRunning ContainerStatus = "running"
	// ContainerStatusStopped indicates a stopped container.
	ContainerStatusStopped ContainerStatus = "stopped"
	// ContainerStatusPaused indicates a paused container.
	ContainerStatusPaused ContainerStatus = "paused"
	// ContainerStatusCreated indicates a container created but not started.
	ContainerStatusCreated ContainerStatus = "created"
	// ContainerStatusUnknown indicates an unrecognized status.
	ContainerStatusUnknown ContainerStatus = "unknown"
)

// Container represents a managed container instance.
type Container struct {
	ID      string
	Name    string
	Image   string
	Status  ContainerStatus
	Created string
	Ports   []PortMapping
}

// Validate ensures required fields are present and valid.
func (c Container) Validate() error {
	if strings.TrimSpace(c.ID) == "" {
		return errors.New("container id is required")
	}
	if strings.TrimSpace(c.Name) == "" {
		return errors.New("container name is required")
	}
	if !isValidContainerStatus(c.Status) {
		return errors.New("container status is invalid")
	}
	return nil
}

func isValidContainerStatus(status ContainerStatus) bool {
	switch status {
	case ContainerStatusRunning, ContainerStatusStopped, ContainerStatusPaused, ContainerStatusCreated, ContainerStatusUnknown:
		return true
	default:
		return false
	}
}
