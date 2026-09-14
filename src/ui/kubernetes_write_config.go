package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"container-tui/src/models"
	"container-tui/src/services"
)

type kubernetesWriteConfigResultMsg struct {
	result models.Result
	err    error
}
type KubernetesWriteConfigScreen struct {
	executor services.CommandExecutor
	cluster  models.KubernetesCluster
	input    textinput.Model
	preview  *CommandPreviewModal
	loading  bool
	errorMsg string
	result   *models.Result
}

func NewKubernetesWriteConfigScreen(executor services.CommandExecutor) KubernetesWriteConfigScreen {
	return KubernetesWriteConfigScreen{executor: executor, input: newKubernetesWriteConfigInput()}
}
func (m KubernetesWriteConfigScreen) SetCluster(cluster models.KubernetesCluster) KubernetesWriteConfigScreen {
	m.cluster, m.input, m.preview, m.loading, m.errorMsg, m.result = cluster, newKubernetesWriteConfigInput(), nil, false, "", nil
	return m
}
func (m KubernetesWriteConfigScreen) Init() tea.Cmd { return textinput.Blink }
func (m KubernetesWriteConfigScreen) Update(msg tea.Msg) (KubernetesWriteConfigScreen, tea.Cmd) {
	switch message := msg.(type) {
	case kubernetesWriteConfigResultMsg:
		m.loading, m.result = false, &message.result
		if message.err != nil {
			m.errorMsg = services.FormatError(message.err, message.result.Stderr)
			return m, nil
		}
		return m, func() tea.Msg { return screenChangeMsg{target: ScreenKubernetesClusterSubmenu} }
	case tea.KeyMsg:
		if m.preview != nil {
			switch strings.ToLower(message.String()) {
			case "y", "enter":
				command := m.preview.Command
				m.preview, m.loading = nil, true
				return m, m.executeCommandCmd(command)
			case "n", "esc":
				m.preview = nil
			}
			return m, nil
		}
		switch message.String() {
		case "esc":
			return m, func() tea.Msg { return screenChangeMsg{target: ScreenKubernetesClusterSubmenu} }
		case "enter":
			command, err := (services.KubernetesWriteConfigBuilder{ClusterName: m.cluster.Name, KubeconfigPath: m.input.Value()}).Build()
			if err != nil {
				m.errorMsg = err.Error()
				return m, nil
			}
			m.preview = &CommandPreviewModal{Title: "Write kubeconfig", Command: command}
			return m, nil
		}
	}
	updated, cmd := m.input.Update(msg)
	m.input = updated
	return m, cmd
}
func (m KubernetesWriteConfigScreen) View() string {
	builder := strings.Builder{}
	builder.WriteString(RenderTitle("Write Kubernetes Kubeconfig") + "\n\n" + RenderMuted("Cluster: "+m.cluster.Name) + "\n\nKubeconfig path (optional)\n" + m.input.View() + "\n")
	if m.loading {
		builder.WriteString(RenderMuted("Writing kubeconfig...") + "\n")
	}
	if m.errorMsg != "" {
		builder.WriteString(RenderError("Error: "+m.errorMsg) + "\n")
	}
	if m.result != nil {
		builder.WriteString(RenderResult(*m.result) + "\n")
	}
	if m.preview != nil {
		builder.WriteString(m.preview.View() + "\n")
	}
	builder.WriteString(RenderMuted("Keys: enter=preview, esc=back") + "\n")
	return builder.String()
}
func (m KubernetesWriteConfigScreen) executeCommandCmd(command models.Command) tea.Cmd {
	return func() tea.Msg {
		result, err := m.executor.Execute(command)
		return kubernetesWriteConfigResultMsg{result: result, err: err}
	}
}
func newKubernetesWriteConfigInput() textinput.Model {
	input := textinput.New()
	input.Prompt, input.Placeholder = "", "~/.kube/config"
	input.Focus()
	return input
}
