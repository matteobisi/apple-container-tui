package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"container-tui/src/models"
	"container-tui/src/services"
)

type kubernetesCreateResultMsg struct {
	result models.Result
	err    error
}

type KubernetesCreateScreen struct {
	executor     services.CommandExecutor
	returnTarget ActiveScreen
	inputs       []textinput.Model
	focus        int
	preview      *CommandPreviewModal
	loading      bool
	errorMsg     string
	result       *models.Result
}

func NewKubernetesCreateScreen(executor services.CommandExecutor) KubernetesCreateScreen {
	return KubernetesCreateScreen{executor: executor, returnTarget: ScreenKubernetesClusterList, inputs: newKubernetesCreateInputs()}
}

func (m KubernetesCreateScreen) SetReturnTarget(target ActiveScreen) KubernetesCreateScreen {
	m.returnTarget = target
	return m
}
func (m KubernetesCreateScreen) Init() tea.Cmd { return textinput.Blink }
func (m KubernetesCreateScreen) Reset() KubernetesCreateScreen {
	m.inputs, m.focus, m.preview, m.loading, m.errorMsg, m.result = newKubernetesCreateInputs(), 0, nil, false, "", nil
	return m
}
func (m KubernetesCreateScreen) Update(msg tea.Msg) (KubernetesCreateScreen, tea.Cmd) {
	switch message := msg.(type) {
	case kubernetesCreateResultMsg:
		m.loading, m.result = false, &message.result
		if message.err != nil {
			m.errorMsg = services.FormatError(message.err, message.result.Stderr)
			return m, nil
		}
		return m, func() tea.Msg { return screenChangeMsg{target: ScreenKubernetesClusterList} }
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
			return m, func() tea.Msg { return screenChangeMsg{target: m.returnTarget} }
		case "tab", "down":
			m.focus = (m.focus + 1) % len(m.inputs)
			m.updateFocus()
			return m, nil
		case "shift+tab", "up":
			m.focus = (m.focus - 1 + len(m.inputs)) % len(m.inputs)
			m.updateFocus()
			return m, nil
		case "enter":
			command, err := (services.KubernetesCreateBuilder{Name: m.inputs[0].Value(), CPUs: m.inputs[1].Value(), Memory: m.inputs[2].Value(), RemoveOnStop: strings.EqualFold(m.inputs[3].Value(), "y"), NodeImage: m.inputs[4].Value()}).Build()
			if err != nil {
				m.errorMsg = err.Error()
				return m, nil
			}
			m.preview = &CommandPreviewModal{Title: "Create cluster", Command: command}
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
func (m KubernetesCreateScreen) View() string {
	labels := []string{"Name (optional)", "CPUs (optional)", "Memory (optional)", "Remove on stop (y/N)", "Node image (optional)"}
	builder := strings.Builder{}
	builder.WriteString(RenderTitle("Create Kubernetes Cluster") + "\n\n")
	for i, label := range labels {
		builder.WriteString(label + "\n" + m.inputs[i].View() + "\n\n")
	}
	if m.loading {
		builder.WriteString(RenderMuted("Creating cluster...") + "\n")
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
func (m *KubernetesCreateScreen) updateFocus() {
	for i := range m.inputs {
		if i == m.focus {
			m.inputs[i].Focus()
		} else {
			m.inputs[i].Blur()
		}
	}
}
func (m KubernetesCreateScreen) executeCommandCmd(command models.Command) tea.Cmd {
	return func() tea.Msg {
		result, err := m.executor.Execute(command)
		return kubernetesCreateResultMsg{result: result, err: err}
	}
}
func newKubernetesCreateInputs() []textinput.Model {
	placeholders := []string{"k8s-dev", "4", "8G", "N", "ghcr.io/apple/container-kubernetes-node:latest"}
	inputs := make([]textinput.Model, len(placeholders))
	for i, placeholder := range placeholders {
		inputs[i] = textinput.New()
		inputs[i].Prompt, inputs[i].Placeholder = "", placeholder
	}
	inputs[0].Focus()
	return inputs
}
