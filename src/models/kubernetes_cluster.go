package models

import (
	"errors"
	"strings"
)

// KubernetesClusterState is the safe aggregate state of a local Kubernetes cluster.
type KubernetesClusterState string

const (
	KubernetesClusterStateRunning KubernetesClusterState = "running"
	KubernetesClusterStateStopped KubernetesClusterState = "stopped"
	KubernetesClusterStateUnknown KubernetesClusterState = "unknown"
)

// KubernetesNode is a node reported by `container k8s list`.
type KubernetesNode struct {
	Name    string
	Role    string
	State   KubernetesClusterState
	CPUs    string
	Memory  string
	Address string
	Ports   string
}

// KubernetesCluster groups the nodes reported for one local cluster.
type KubernetesCluster struct {
	Name    string
	Nodes   []KubernetesNode
	State   KubernetesClusterState
	CPUs    string
	Memory  string
	Address string
	Ports   string
}

func (c KubernetesCluster) Validate() error {
	if strings.TrimSpace(c.Name) == "" {
		return errors.New("cluster name is required")
	}
	return nil
}

// KubernetesCreateInput collects optional cluster-create settings.
type KubernetesCreateInput struct {
	Name         string
	CPUs         string
	Memory       string
	RemoveOnStop bool
	NodeImage    string
}

// KubernetesLoadImageInput collects image loading settings.
type KubernetesLoadImageInput struct {
	ClusterName string
	Image       string
	Platform    string
}

// KubernetesWriteConfigInput collects optional kubeconfig output settings.
type KubernetesWriteConfigInput struct {
	ClusterName    string
	KubeconfigPath string
}
