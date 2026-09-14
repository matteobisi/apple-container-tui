package services

import (
	"strings"

	"container-tui/src/models"
)

// ParseKubernetesList parses the whitespace table printed by `container k8s list`.
func ParseKubernetesList(output string) ([]models.KubernetesCluster, error) {
	clusters := map[string]*models.KubernetesCluster{}
	order := []string{}
	currentName := ""
	for _, line := range strings.Split(output, "\n") {
		fields := strings.Fields(line)
		if len(fields) == 0 || fields[0] == "CLUSTER" || fields[0] == "PORTS" {
			continue
		}
		name, offset := currentName, 0
		if len(fields) >= 7 && isKubernetesRole(fields[2]) {
			name, currentName, offset = fields[0], fields[0], 1
		}
		if name == "" && len(fields) >= 6 && isKubernetesRole(fields[1]) && strings.Contains(fields[1], "control-plane") {
			name, currentName = clusterNameFromControlPlane(fields[0]), clusterNameFromControlPlane(fields[0])
		}
		if name == "" || len(fields)-offset < 5 {
			continue
		}
		node := models.KubernetesNode{Name: fields[offset], Role: fields[offset+1], State: normalizeKubernetesState(fields[offset+2]), CPUs: fields[offset+3], Memory: fields[offset+4]}
		remainder := fields[offset+5:]
		if len(remainder) > 0 && isMemoryUnit(remainder[0]) {
			node.Memory += " " + remainder[0]
			remainder = remainder[1:]
		}
		if len(remainder) == 1 && strings.Contains(remainder[0], "->") {
			node.Ports = remainder[0]
		} else if len(remainder) > 0 {
			node.Address = remainder[0]
			node.Ports = strings.Join(remainder[1:], " ")
		}
		cluster := clusters[name]
		if cluster == nil {
			cluster = &models.KubernetesCluster{Name: name}
			clusters[name] = cluster
			order = append(order, name)
		}
		cluster.Nodes = append(cluster.Nodes, node)
		if cluster.CPUs == "" {
			cluster.CPUs, cluster.Memory, cluster.Address, cluster.Ports = node.CPUs, node.Memory, node.Address, node.Ports
		}
	}
	result := make([]models.KubernetesCluster, 0, len(order))
	for _, name := range order {
		cluster := clusters[name]
		cluster.State = aggregateKubernetesState(cluster.Nodes)
		result = append(result, *cluster)
	}
	return result, nil
}

func isKubernetesRole(value string) bool {
	return strings.Contains(value, "control-plane") || value == "worker"
}

func clusterNameFromControlPlane(nodeName string) string {
	return strings.TrimSuffix(nodeName, "-control-plane")
}

func isMemoryUnit(value string) bool {
	return value == "MB" || value == "GB" || value == "MiB" || value == "GiB"
}

func normalizeKubernetesState(value string) models.KubernetesClusterState {
	switch strings.ToLower(value) {
	case "running":
		return models.KubernetesClusterStateRunning
	case "stopped":
		return models.KubernetesClusterStateStopped
	default:
		return models.KubernetesClusterStateUnknown
	}
}
func aggregateKubernetesState(nodes []models.KubernetesNode) models.KubernetesClusterState {
	if len(nodes) == 0 {
		return models.KubernetesClusterStateUnknown
	}
	state := nodes[0].State
	for _, node := range nodes[1:] {
		if node.State != state {
			return models.KubernetesClusterStateUnknown
		}
	}
	return state
}
