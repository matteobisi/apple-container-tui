package services

import (
	"fmt"
	"strconv"
	"strings"

	"container-tui/src/models"
)

type KubernetesListBuilder struct{}

func (KubernetesListBuilder) Validate() error { return nil }

func (KubernetesListBuilder) Build() (models.Command, error) {
	return models.Command{Executable: "container", Args: []string{"cluster", "list"}}, nil
}

type KubernetesCreateBuilder struct {
	Name         string
	CPUs         string
	Memory       string
	RemoveOnStop bool
	NodeImage    string
	Workers      string
	CNI          string
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
		value, err := strconv.Atoi(strings.TrimSpace(b.Workers))
		if err != nil || value < 0 {
			return fmt.Errorf("workers must be a non-negative integer")
		}
	}
	return nil
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
	if strings.TrimSpace(b.Workers) != "" && strings.TrimSpace(b.Workers) != "0" {
		args = append(args, "--workers", strings.TrimSpace(b.Workers))
	}
	if strings.TrimSpace(b.CNI) != "" {
		args = append(args, "--cni", strings.TrimSpace(b.CNI))
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

type KubernetesLoadImageBuilder struct {
	ClusterName string
	Image       string
	Platform    string
	Nodes       []string
}

func (b KubernetesLoadImageBuilder) Build() (models.Command, error) {
	image, err := normalizeRequiredToken(b.Image, "image")
	if err != nil {
		return models.Command{}, err
	}
	args := []string{"cluster", "load-image"}
	if b.ClusterName != "" {
		name, err := normalizeRequiredToken(b.ClusterName, "cluster name")
		if err != nil {
			return models.Command{}, err
		}
		args = append(args, "--name", name)
	}
	args = append(args, image)
	if b.Platform != "" {
		platform, err := normalizeRequiredToken(b.Platform, "platform")
		if err != nil {
			return models.Command{}, err
		}
		args = append(args, "--platform", platform)
	}
	for _, node := range b.Nodes {
		node = strings.TrimSpace(node)
		if node == "" {
			continue
		}
		if strings.ContainsAny(node, " \t") {
			return models.Command{}, fmt.Errorf("node name must not contain whitespace")
		}
		args = append(args, "--node", node)
	}
	return models.Command{Executable: "container", Args: args}, nil
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
