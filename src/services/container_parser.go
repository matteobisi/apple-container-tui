package services

import (
	"fmt"
	"strconv"
	"strings"

	"container-tui/src/models"
)

var containerListHeaders = []string{
	"CONTAINER ID",
	"IMAGE",
	"COMMAND",
	"CREATED",
	"STATUS",
	"PORTS",
}

// ParseContainerList parses `container list --all` output into containers.
func ParseContainerList(output string) ([]models.Container, error) {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) == 0 || strings.TrimSpace(lines[0]) == "" {
		return nil, nil
	}

	headerLine := lines[0]
	if strings.Contains(headerLine, "CONTAINER ID") {
		indices, err := headerIndices(headerLine, containerListHeaders)
		if err != nil {
			return nil, err
		}

		containers := make([]models.Container, 0)
		for _, line := range lines[1:] {
			if strings.TrimSpace(line) == "" {
				continue
			}
			columns := sliceColumns(line, indices)
			if len(columns) < len(containerListHeaders) {
				continue
			}

			id := columns[0]
			image := columns[1]
			command := columns[2]
			created := columns[3]
			status := parseContainerStatus(columns[4])
			ports := parsePortMappings(columns[5])

			name := command
			if strings.TrimSpace(name) == "" {
				name = id
			}

			containers = append(containers, models.Container{
				ID:      id,
				Name:    name,
				Image:   image,
				Status:  status,
				Created: created,
				Ports:   ports,
			})
		}

		return containers, nil
	}

	return parseContainerListFields(lines)
}

func parseContainerListFields(lines []string) ([]models.Container, error) {
	headers := strings.Fields(lines[0])
	index := make(map[string]int, len(headers))
	for i, header := range headers {
		index[strings.ToUpper(header)] = i
	}

	idIndex, okID := index["ID"]
	imageIndex, okImage := index["IMAGE"]
	statusIndex, okStatus := index["STATUS"]
	if !okStatus {
		statusIndex, okStatus = index["STATE"]
	}
	if !okID || !okImage || !okStatus {
		return nil, fmt.Errorf("missing header in container list output")
	}

	containers := make([]models.Container, 0)
	for _, line := range lines[1:] {
		if strings.TrimSpace(line) == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) <= statusIndex || len(fields) <= imageIndex || len(fields) <= idIndex {
			continue
		}

		id := fields[idIndex]
		image := fields[imageIndex]
		status := parseContainerStatus(fields[statusIndex])
		name := id

		containers = append(containers, models.Container{
			ID:     id,
			Name:   name,
			Image:  image,
			Status: status,
		})
	}

	return containers, nil
}

func headerIndices(line string, headers []string) ([]int, error) {
	indices := make([]int, 0, len(headers))
	for _, header := range headers {
		idx := strings.Index(line, header)
		if idx == -1 {
			return nil, fmt.Errorf("missing header %q", header)
		}
		indices = append(indices, idx)
	}
	return indices, nil
}

func sliceColumns(line string, indices []int) []string {
	columns := make([]string, len(indices))
	for i, idx := range indices {
		end := len(line)
		if i+1 < len(indices) {
			end = indices[i+1]
			if end > len(line) {
				end = len(line)
			}
		}
		if idx >= len(line) {
			columns[i] = ""
			continue
		}
		columns[i] = strings.TrimSpace(line[idx:end])
	}
	return columns
}

func parseContainerStatus(value string) models.ContainerStatus {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "running":
		return models.ContainerStatusRunning
	case "stopped", "exited":
		return models.ContainerStatusStopped
	case "paused":
		return models.ContainerStatusPaused
	case "created":
		return models.ContainerStatusCreated
	default:
		return models.ContainerStatusUnknown
	}
}

func parsePortMappings(value string) []models.PortMapping {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}

	entries := strings.Split(value, ",")
	ports := make([]models.PortMapping, 0, len(entries))
	for _, entry := range entries {
		entry = strings.TrimSpace(entry)
		if entry == "" {
			continue
		}

		parts := strings.Split(entry, "->")
		if len(parts) != 2 {
			continue
		}
		hostPart := strings.TrimSpace(parts[0])
		containerPart := strings.TrimSpace(parts[1])

		colonIndex := strings.LastIndex(hostPart, ":")
		if colonIndex == -1 {
			continue
		}
		hostPortStr := hostPart[colonIndex+1:]

		containerSegments := strings.SplitN(containerPart, "/", 2)
		containerPortStr := strings.TrimSpace(containerSegments[0])
		protocol := ""
		if len(containerSegments) == 2 {
			protocol = strings.TrimSpace(containerSegments[1])
		}

		hostPort, err := strconv.Atoi(hostPortStr)
		if err != nil {
			continue
		}
		containerPort, err := strconv.Atoi(containerPortStr)
		if err != nil {
			continue
		}

		ports = append(ports, models.PortMapping{
			HostPort:      hostPort,
			ContainerPort: containerPort,
			Protocol:      protocol,
		})
	}

	return ports
}

// ParseImageList parses `container image list` output into images.
func ParseImageList(output string) ([]models.Image, error) {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) == 0 || strings.TrimSpace(lines[0]) == "" {
		return []models.Image{}, nil
	}

	headers := strings.Fields(lines[0])
	indices := map[string]int{}
	for i, header := range headers {
		indices[strings.ToUpper(header)] = i
	}

	nameIndex, okName := indices["NAME"]
	tagIndex, okTag := indices["TAG"]
	digestIndex, okDigest := indices["DIGEST"]
	if !okName || !okTag {
		return nil, fmt.Errorf("missing required image headers")
	}

	images := make([]models.Image, 0)
	for _, line := range lines[1:] {
		if strings.TrimSpace(line) == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) <= tagIndex || len(fields) <= nameIndex {
			continue
		}

		image := models.Image{
			Name: fields[nameIndex],
			Tag:  fields[tagIndex],
		}
		if okDigest && len(fields) > digestIndex {
			image.Digest = fields[digestIndex]
		}
		if strings.TrimSpace(image.Tag) == "" {
			image.Tag = "<none>"
		}
		images = append(images, image)
	}

	return images, nil
}
