package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"container-tui/src/models"
	"container-tui/src/services"
)

type kubernetesLoadImageResultMsg struct {
	result models.Result
	err    error
}
type KubernetesLoadImageScreen struct {
	executor services.CommandExecutor
	cluster  models.KubernetesCluster
	inputs   []textinput.Model
	focus    int
	preview  *CommandPreviewModal
	loading  bool
	errorMsg string
	result   *models.Result
}

func NewKubernetesLoadImageScreen(executor services.CommandExecutor) KubernetesLoadImageScreen {
	return KubernetesLoadImageScreen{executor: executor, inputs: newKubernetesLoadImageInputs()}
}
func (m KubernetesLoadImageScreen) SetCluster(cluster models.KubernetesCluster) KubernetesLoadImageScreen {
	m.cluster, m.inputs, m.focus, m.preview, m.loading, m.errorMsg, m.result = cluster, newKubernetesLoadImageInputs(), 0, nil, false, "", nil
	return m
}
func (m KubernetesLoadImageScreen) Init() tea.Cmd { return textinput.Blink }
func (m KubernetesLoadImageScreen) Update(msg tea.Msg) (KubernetesLoadImageScreen, tea.Cmd) {
	switch message := msg.(type) {
	case kubernetesLoadImageResultMsg:
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
		case "tab", "down":
			m.focus = (m.focus + 1) % len(m.inputs)
			m.updateFocus()
			return m, nil
		case "shift+tab", "up":
			m.focus = (m.focus - 1 + len(m.inputs)) % len(m.inputs)
			m.updateFocus()
			return m, nil
		case "enter":
			command, err := (services.KubernetesLoadImageBuilder{ClusterName: m.cluster.Name, Image: m.inputs[0].Value(), Platform: m.inputs[1].Value()}).Build()
			if err != nil {
				m.errorMsg = err.Error()
				return m, nil
			}
			m.preview = &CommandPreviewModal{Title: "Load image", Command: command}
			return m, nil
		}
	}
	for i := range m.inputs {
		updated, cmd := m.inputs[i].Update(msg)
		m.inputs[i] = updated
		if cmd != nil {
			return m, cmd
		}
	}
	return m, nil
}
func (m KubernetesLoadImageScreen) View() string {
	builder := strings.Builder{}
	builder.WriteString(RenderTitle("Load Kubernetes Image") + "\n\n" + RenderMuted("Cluster: "+m.cluster.Name) + "\n\nImage\n" + m.inputs[0].View() + "\n\nPlatform (optional)\n" + m.inputs[1].View() + "\n")
	if m.loading {
		builder.WriteString(RenderMuted("Loading image...") + "\n")
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
	builder.WriteString(RenderMuted("Keys: tab/up/down=field, enter=preview, esc=back") + "\n")
	return builder.String()
}
func (m *KubernetesLoadImageScreen) updateFocus() {
	for i := range m.inputs {
		if i == m.focus {
			m.inputs[i].Focus()
		} else {
			m.inputs[i].Blur()
		}
	}
}
func (m KubernetesLoadImageScreen) executeCommandCmd(command models.Command) tea.Cmd {
	return func() tea.Msg {
		result, err := m.executor.Execute(command)
		return kubernetesLoadImageResultMsg{result: result, err: err}
	}
}
func newKubernetesLoadImageInputs() []textinput.Model {
	image, platform := textinput.New(), textinput.New()
	image.Prompt, image.Placeholder = "", "nginx:latest"
	image.Focus()
	platform.Prompt, platform.Placeholder = "", "linux/arm64"
	return []textinput.Model{image, platform}
}
