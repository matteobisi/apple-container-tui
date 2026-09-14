package ui

import (
	"errors"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"container-tui/src/models"
	"container-tui/src/services"
)

type kubernetesListLoadedMsg struct {
	clusters []models.KubernetesCluster
	err      error
}

// KubernetesClusterListScreen displays local Apple Container Kubernetes clusters.
type KubernetesClusterListScreen struct {
	executor  services.CommandExecutor
	clusters  []models.KubernetesCluster
	cursor    int
	width     int
	loading   bool
	errorMsg  string
	hasLoaded bool
}

func NewKubernetesClusterListScreen(executor services.CommandExecutor) KubernetesClusterListScreen {
	return KubernetesClusterListScreen{executor: executor}
}

func (m KubernetesClusterListScreen) Init() tea.Cmd { return m.fetchClustersCmd() }

func (m KubernetesClusterListScreen) Update(msg tea.Msg) (KubernetesClusterListScreen, tea.Cmd) {
	switch message := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = message.Width
	case kubernetesListLoadedMsg:
		m.loading = false
		m.hasLoaded = true
		if message.err != nil {
			m.errorMsg = formatKubernetesListError(message.err)
			return m, nil
		}
		m.errorMsg = ""
		m.clusters = message.clusters
		if m.cursor >= len(m.clusters) {
			m.cursor = max(0, len(m.clusters)-1)
		}
	case tea.KeyMsg:
		switch message.String() {
		case "up", "k":
			m.cursor = max(0, m.cursor-1)
		case "down", "j":
			m.cursor = min(len(m.clusters)-1, m.cursor+1)
		case "r":
			m.loading = true
			return m, m.fetchClustersCmd()
		case "c":
			return m, func() tea.Msg { return screenChangeMsg{target: ScreenKubernetesCreate, push: true} }
		case "esc":
			return m, func() tea.Msg { return screenChangeMsg{target: ScreenContainerList} }
		case "enter":
			cluster, ok := m.selectedCluster()
			if !ok {
				return m, nil
			}
			return m, func() tea.Msg {
				return screenChangeMsg{target: ScreenKubernetesClusterSubmenu, cluster: &cluster, push: true}
			}
		}
	}
	return m, nil
}

func formatKubernetesListError(err error) string {
	message := strings.ToLower(err.Error())
	if strings.Contains(message, "unknown command") || strings.Contains(message, "not found") || errors.Is(err, exec.ErrNotFound) {
		return "Kubernetes support is unavailable in the installed Container CLI"
	}
	return services.FormatError(err, "")
}

func (m KubernetesClusterListScreen) View() string {
	builder := strings.Builder{}
	builder.WriteString(RenderTitle("Kubernetes Clusters") + "\n\n")
	if m.loading {
		builder.WriteString(RenderMuted("Loading clusters...") + "\n")
	}
	if m.errorMsg != "" {
		builder.WriteString(RenderError("Error: "+m.errorMsg) + "\n\n")
	}
	table := NewTable([]TableColumn{{Header: "Name", MinWidth: 16, Priority: 1, Align: "left"}, {Header: "State", MinWidth: 8, Priority: 2, Align: "left"}, {Header: "Nodes", MinWidth: 6, Priority: 3, Align: "right"}, {Header: "Address", MinWidth: 16, Priority: 4, Align: "left"}})
	rows := make([]TableRow, len(m.clusters))
	for i, cluster := range m.clusters {
		rows[i] = TableRow{Cells: []string{cluster.Name, string(cluster.State), intToDisplay(len(cluster.Nodes)), cluster.Address}, Selected: i == m.cursor, Data: &cluster}
	}
	table.SetRows(rows)
	width := m.width
	if width == 0 {
		width = 80
	}
	if len(m.clusters) == 0 && m.hasLoaded && m.errorMsg == "" {
		builder.WriteString(RenderMuted("No Kubernetes clusters found.") + "\n")
	} else {
		builder.WriteString(table.Render(width, m.cursor))
	}
	builder.WriteString(strings.Repeat("─", width) + "\n")
	builder.WriteString("\n" + RenderMuted("Keys: up/down, enter=submenu, c=create, r=refresh, esc=containers, q=quit") + "\n")
	return builder.String()
}

func (m KubernetesClusterListScreen) selectedCluster() (models.KubernetesCluster, bool) {
	if len(m.clusters) == 0 || m.cursor < 0 || m.cursor >= len(m.clusters) {
		return models.KubernetesCluster{}, false
	}
	return m.clusters[m.cursor], true
}

func (m KubernetesClusterListScreen) fetchClustersCmd() tea.Cmd {
	return func() tea.Msg {
		command, err := (services.KubernetesListBuilder{}).Build()
		if err != nil {
			return kubernetesListLoadedMsg{err: err}
		}
		result, err := m.executor.Execute(command)
		if err != nil {
			return kubernetesListLoadedMsg{err: err}
		}
		clusters, err := services.ParseKubernetesList(result.Stdout)
		return kubernetesListLoadedMsg{clusters: clusters, err: err}
	}
}
