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
	return models.Command{Executable: "container", Args: []string{"k8s", "list"}}, nil
}

type KubernetesCreateBuilder struct {
	Name         string
	CPUs         string
	Memory       string
	RemoveOnStop bool
	NodeImage    string
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
	return nil
}

func (b KubernetesCreateBuilder) Build() (models.Command, error) {
	if err := b.Validate(); err != nil {
		return models.Command{}, err
	}
	args := []string{"k8s", "create"}
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
	return models.Command{Executable: "container", Args: args}, nil
}

type KubernetesStartBuilder struct{ ClusterName string }

func (b KubernetesStartBuilder) Build() (models.Command, error) {
	name, err := normalizeRequiredToken(b.ClusterName, "cluster name")
	if err != nil {
		return models.Command{}, err
	}
	return models.Command{Executable: "container", Args: []string{"k8s", "start", "--name", name}}, nil
}

type KubernetesDeleteBuilder struct{ ClusterName string }

func (b KubernetesDeleteBuilder) Build() (models.Command, error) {
	name, err := normalizeRequiredToken(b.ClusterName, "cluster name")
	if err != nil {
		return models.Command{}, err
	}
	return models.Command{Executable: "container", Args: []string{"k8s", "delete", "--name", name}}, nil
}

type KubernetesLoadImageBuilder struct{ ClusterName, Image, Platform string }

func (b KubernetesLoadImageBuilder) Build() (models.Command, error) {
	image, err := normalizeRequiredToken(b.Image, "image")
	if err != nil {
		return models.Command{}, err
	}
	args := []string{"k8s", "load-image"}
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
	return models.Command{Executable: "container", Args: args}, nil
}

type KubernetesWriteConfigBuilder struct{ ClusterName, KubeconfigPath string }

func (b KubernetesWriteConfigBuilder) Build() (models.Command, error) {
	name, err := normalizeRequiredToken(b.ClusterName, "cluster name")
	if err != nil {
		return models.Command{}, err
	}
	args := []string{"k8s", "write-config", "--name", name}
	if strings.TrimSpace(b.KubeconfigPath) != "" {
		args = append(args, "--kubeconfig", strings.TrimSpace(b.KubeconfigPath))
	}
	return models.Command{Executable: "container", Args: args}, nil
}
