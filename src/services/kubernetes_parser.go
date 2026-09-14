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
		if len(fields) >= 7 && !isNodeName(fields[0]) {
			name, currentName, offset = fields[0], fields[0], 1
		}
		if name == "" || len(fields)-offset < 5 {
			continue
		}
		node := models.KubernetesNode{Name: fields[offset], Role: fields[offset+1], State: normalizeKubernetesState(fields[offset+2]), CPUs: fields[offset+3], Memory: fields[offset+4]}
		if len(fields)-offset > 5 {
			node.Address = fields[offset+5]
		}
		if len(fields)-offset > 6 {
			node.Ports = strings.Join(fields[offset+6:], " ")
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

func isNodeName(value string) bool { return strings.Contains(value, "-") }
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
