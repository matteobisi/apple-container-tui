package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"container-tui/src/models"
	"container-tui/src/services"
)

type kubernetesSubmenuActionMsg struct {
	action string
	result models.Result
	err    error
}

type kubernetesSubmenuOption struct{ label, action string }

// KubernetesClusterSubmenuScreen provides actions for a selected cluster.
type KubernetesClusterSubmenuScreen struct {
	executor services.CommandExecutor
	cluster  models.KubernetesCluster
	options  []kubernetesSubmenuOption
	cursor   int
	loading  bool
	errorMsg string
	result   *models.Result
	preview  *CommandPreviewModal
	confirm  *TypeToConfirmModal
	width    int
}

func NewKubernetesClusterSubmenuScreen(executor services.CommandExecutor) KubernetesClusterSubmenuScreen {
	return KubernetesClusterSubmenuScreen{executor: executor}
}

func (m KubernetesClusterSubmenuScreen) SetCluster(cluster models.KubernetesCluster) KubernetesClusterSubmenuScreen {
	m.cluster, m.cursor, m.loading, m.errorMsg, m.result, m.preview, m.confirm = cluster, 0, false, "", nil, nil, nil
	m.options = m.buildOptions()
	return m
}

func (m KubernetesClusterSubmenuScreen) Init() tea.Cmd { return nil }

func (m KubernetesClusterSubmenuScreen) Update(msg tea.Msg) (KubernetesClusterSubmenuScreen, tea.Cmd) {
	switch message := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = message.Width
	case kubernetesSubmenuActionMsg:
		m.loading = false
		m.result = &message.result
		if message.err != nil {
			m.errorMsg = services.FormatError(message.err, message.result.Stderr)
			return m, nil
		}
		m.errorMsg = ""
		if message.action == "delete" {
			return m, func() tea.Msg { return screenChangeMsg{target: ScreenKubernetesClusterList} }
		}
	case tea.KeyMsg:
		if m.confirm != nil {
			confirm, confirmed, canceled := m.confirm.Handle(message)
			m.confirm = &confirm
			if confirmed {
				m.confirm, m.loading = nil, true
				return m, m.executeCommandCmd("delete", confirm.Command)
			}
			if canceled {
				m.confirm = nil
			}
			return m, nil
		}
		if m.preview != nil {
			switch strings.ToLower(message.String()) {
			case "y", "enter":
				command, action := m.preview.Command, m.preview.Title
				m.preview, m.loading = nil, true
				return m, m.executeCommandCmd(action, command)
			case "n", "esc":
				m.preview = nil
			}
			return m, nil
		}
		switch message.String() {
		case "up", "k":
			m.cursor = max(0, m.cursor-1)
		case "down", "j":
			m.cursor = min(len(m.options)-1, m.cursor+1)
		case "esc":
			return m, func() tea.Msg { return screenChangeMsg{target: ScreenKubernetesClusterList} }
		case "enter":
			if m.cursor >= 0 && m.cursor < len(m.options) {
				return m.selectOption(m.options[m.cursor])
			}
		}
	}
	return m, nil
}

func (m KubernetesClusterSubmenuScreen) View() string {
	builder := strings.Builder{}
	builder.WriteString(RenderTitle("Kubernetes Cluster Actions") + "\n\n")
	header := lipgloss.NewStyle().Bold(true)
	builder.WriteString(header.Render("Cluster Details") + "\n")
	builder.WriteString("name: " + m.cluster.Name + "\nstate: " + string(m.cluster.State) + "\nnodes: " + intToDisplay(len(m.cluster.Nodes)) + "\naddress: " + m.cluster.Address + "\n")
	width := m.width
	if width == 0 {
		width = 80
	}
	builder.WriteString("\n" + strings.Repeat("─", width) + "\n\n" + header.Render("Available Actions") + "\n")
	for i, option := range m.options {
		style := lipgloss.NewStyle()
		if i == m.cursor {
			style = style.Reverse(true)
		}
		builder.WriteString(style.Render(option.label) + "\n")
	}
	if m.loading {
		builder.WriteString("\n" + RenderMuted("Running action...") + "\n")
	}
	if m.errorMsg != "" {
		builder.WriteString("\n" + RenderError("Error: "+m.errorMsg) + "\n")
	}
	if m.result != nil {
		builder.WriteString("\n" + RenderResult(*m.result) + "\n")
	}
	if m.preview != nil {
		builder.WriteString("\n" + m.preview.View() + "\n")
	}
	if m.confirm != nil {
		builder.WriteString("\n" + m.confirm.View() + "\n")
	}
	builder.WriteString("\n" + RenderMuted("Keys: up/down=navigate, enter=select, esc=back") + "\n")
	return builder.String()
}

func (m KubernetesClusterSubmenuScreen) buildOptions() []kubernetesSubmenuOption {
	options := []kubernetesSubmenuOption{{"Refresh clusters", "refresh"}}
	if m.cluster.State != models.KubernetesClusterStateRunning {
		options = append(options, kubernetesSubmenuOption{"Start cluster", "start"})
	}
	return append(options, kubernetesSubmenuOption{"Load image", "load-image"}, kubernetesSubmenuOption{"Write kubeconfig", "write-config"}, kubernetesSubmenuOption{"Delete cluster", "delete"}, kubernetesSubmenuOption{"Create cluster", "create"}, kubernetesSubmenuOption{"Back", "back"})
}

func (m KubernetesClusterSubmenuScreen) selectOption(option kubernetesSubmenuOption) (KubernetesClusterSubmenuScreen, tea.Cmd) {
	switch option.action {
	case "refresh", "back":
		return m, func() tea.Msg { return screenChangeMsg{target: ScreenKubernetesClusterList} }
	case "create":
		return m, func() tea.Msg { return screenChangeMsg{target: ScreenKubernetesCreate, push: true} }
	case "load-image":
		cluster := m.cluster
		return m, func() tea.Msg {
			return screenChangeMsg{target: ScreenKubernetesLoadImage, cluster: &cluster, push: true}
		}
	case "write-config":
		cluster := m.cluster
		return m, func() tea.Msg {
			return screenChangeMsg{target: ScreenKubernetesWriteConfig, cluster: &cluster, push: true}
		}
	case "start":
		command, err := (services.KubernetesStartBuilder{ClusterName: m.cluster.Name}).Build()
		if err != nil {
			m.errorMsg = err.Error()
			return m, nil
		}
		m.preview = &CommandPreviewModal{Title: "Start cluster", Command: command}
	case "delete":
		command, err := (services.KubernetesDeleteBuilder{ClusterName: m.cluster.Name}).Build()
		if err != nil {
			m.errorMsg = err.Error()
			return m, nil
		}
		confirm := NewTypeToConfirmModal("Delete Kubernetes Cluster", m.cluster.Name, command)
		m.confirm = &confirm
	}
	return m, nil
}

func (m KubernetesClusterSubmenuScreen) executeCommandCmd(action string, command models.Command) tea.Cmd {
	return func() tea.Msg {
		result, err := m.executor.Execute(command)
		return kubernetesSubmenuActionMsg{action: action, result: result, err: err}
	}
}
