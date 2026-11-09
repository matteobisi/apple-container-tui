package services

import (
	"fmt"
	"strconv"
	"strings"

	"container-tui/src/models"
)

type KubernetesListBuilder struct{ Nodes bool }

func (KubernetesListBuilder) Validate() error { return nil }

func (b KubernetesListBuilder) Build() (models.Command, error) {
	args := []string{"cluster", "list"}
	if b.Nodes {
		args = append(args, "--nodes")
	}
	return models.Command{Executable: "container", Args: args}, nil
}

type KubernetesCreateBuilder struct {
	Name             string
	CPUs             string
	Memory           string
	RemoveOnStop     bool
	NodeImage        string
	Workers          string
	CNI              string
	SharedImageStore bool
}

func (b KubernetesCreateBuilder) Validate() error {
	if b.Name != "" {
		if _, err := normalizeRequiredToken(b.Name, "cluster name"); err != nil {
			return err
		}
	}
	if b.CPUs != "" {
		value, err := strconv.Atoi(strings.TrimSpace(b.CPUs))
		if err != nil || value < 1 {
			return fmt.Errorf("cpus must be a positive integer")
		}
	}
	if b.Memory != "" {
		if _, err := normalizeRequiredToken(b.Memory, "memory"); err != nil {
			return err
		}
	}
	if b.NodeImage != "" {
		if _, err := normalizeRequiredToken(b.NodeImage, "node image"); err != nil {
			return err
		}
	}
	if strings.TrimSpace(b.Workers) != "" {
		if _, err := splitWorkersSpecs(b.Workers); err != nil {
			return err
		}
	}
	return nil
}

// splitWorkersSpecs parses the Workers field into one --workers spec per
// group. Each group is either a bare non-negative count ("2") or a
// count=<n>,cpus=<n>,memory=<size> spec; groups are separated by whitespace
// or semicolons (commas are intra-spec separators).
func splitWorkersSpecs(raw string) ([]string, error) {
	normalized := strings.ReplaceAll(raw, ";", " ")
	var specs []string
	for _, token := range strings.Fields(normalized) {
		if !strings.Contains(token, "=") {
			for _, piece := range strings.Split(token, ",") {
				piece = strings.TrimSpace(piece)
				if piece == "" {
					continue
				}
				value, err := strconv.Atoi(piece)
				if err != nil || value < 0 {
					return nil, fmt.Errorf("workers must be a non-negative integer or count=<n>,cpus=<n>,memory=<size>")
				}
				if value == 0 {
					continue
				}
				specs = append(specs, piece)
			}
			continue
		}
		count := -1
		for _, segment := range strings.Split(token, ",") {
			segment = strings.TrimSpace(segment)
			kv := strings.SplitN(segment, "=", 2)
			if len(kv) != 2 || strings.TrimSpace(kv[0]) == "" || strings.TrimSpace(kv[1]) == "" {
				return nil, fmt.Errorf("workers must be a non-negative integer or count=<n>,cpus=<n>,memory=<size>")
			}
			key, value := strings.TrimSpace(kv[0]), strings.TrimSpace(kv[1])
			switch key {
			case "count":
				parsed, err := strconv.Atoi(value)
				if err != nil || parsed < 0 {
					return nil, fmt.Errorf("workers count must be a non-negative integer")
				}
				count = parsed
			case "cpus":
				parsed, err := strconv.Atoi(value)
				if err != nil || parsed < 1 {
					return nil, fmt.Errorf("workers cpus must be a positive integer")
				}
			case "memory":
				if _, err := normalizeRequiredToken(value, "workers memory"); err != nil {
					return nil, err
				}
			default:
				return nil, fmt.Errorf("unsupported workers setting %q (expected count, cpus, or memory)", key)
			}
		}
		if count < 0 {
			return nil, fmt.Errorf("workers spec %q must include count=<n>", token)
		}
		if count == 0 {
			continue
		}
		specs = append(specs, token)
	}
	return specs, nil
}

func (b KubernetesCreateBuilder) Build() (models.Command, error) {
	if err := b.Validate(); err != nil {
		return models.Command{}, err
	}
	args := []string{"cluster", "create"}
	if b.Name != "" {
		value, _ := normalizeRequiredToken(b.Name, "cluster name")
		args = append(args, "--name", value)
	}
	if b.CPUs != "" {
		args = append(args, "--cpus", strings.TrimSpace(b.CPUs))
	}
	if b.Memory != "" {
		value, _ := normalizeRequiredToken(b.Memory, "memory")
		args = append(args, "--memory", value)
	}
	if b.RemoveOnStop {
		args = append(args, "--rm")
	}
	if b.NodeImage != "" {
		value, _ := normalizeRequiredToken(b.NodeImage, "node image")
		args = append(args, "--node-image", value)
	}
	if strings.TrimSpace(b.Workers) != "" {
		specs, err := splitWorkersSpecs(b.Workers)
		if err != nil {
			return models.Command{}, err
		}
		for _, spec := range specs {
			args = append(args, "--workers", spec)
		}
	}
	if strings.TrimSpace(b.CNI) != "" {
		args = append(args, "--cni", strings.TrimSpace(b.CNI))
	}
	if b.SharedImageStore {
		args = append(args, "--shared-image-store")
	}
	return models.Command{Executable: "container", Args: args}, nil
}

type KubernetesDeleteBuilder struct{ ClusterName string }

func (b KubernetesDeleteBuilder) Build() (models.Command, error) {
	name, err := normalizeRequiredToken(b.ClusterName, "cluster name")
	if err != nil {
		return models.Command{}, err
	}
	return models.Command{Executable: "container", Args: []string{"cluster", "delete", "--name", name}}, nil
}

type KubernetesWriteConfigBuilder struct{ ClusterName, KubeconfigPath string }

func (b KubernetesWriteConfigBuilder) Build() (models.Command, error) {
	name, err := normalizeRequiredToken(b.ClusterName, "cluster name")
	if err != nil {
		return models.Command{}, err
	}
	args := []string{"cluster", "write-config", "--name", name}
	if strings.TrimSpace(b.KubeconfigPath) != "" {
		args = append(args, "--kubeconfig", strings.TrimSpace(b.KubeconfigPath))
	}
	return models.Command{Executable: "container", Args: args}, nil
}
