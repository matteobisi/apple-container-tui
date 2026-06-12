package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"container-tui/src/models"
	"container-tui/src/services"
)

type machineCreateResultMsg struct {
	result models.Result
	err    error
}

// MachineCreateScreen creates a new container machine from an image reference.
type MachineCreateScreen struct {
	executor services.CommandExecutor
	inputs   []textinput.Model
	focus    int
	preview  *CommandPreviewModal
	loading  bool
	errorMsg string
	result   *models.Result
}

func NewMachineCreateScreen(executor services.CommandExecutor) MachineCreateScreen {
	return MachineCreateScreen{executor: executor, inputs: newMachineCreateInputs()}
}

func (m MachineCreateScreen) Init() tea.Cmd { return textinput.Blink }

func (m MachineCreateScreen) Reset() MachineCreateScreen {
	m.inputs = newMachineCreateInputs()
	m.focus = 0
	m.preview = nil
	m.loading = false
	m.errorMsg = ""
	m.result = nil
	return m
}

func (m MachineCreateScreen) Update(msg tea.Msg) (MachineCreateScreen, tea.Cmd) {
	switch message := msg.(type) {
	case machineCreateResultMsg:
		m.loading = false
		if message.err != nil {
			m.errorMsg = services.FormatError(message.err, message.result.Stderr)
			m.result = &message.result
			return m, nil
		}
		m.result = &message.result
		return m, func() tea.Msg { return screenChangeMsg{target: ScreenMachineList} }
	case tea.KeyMsg:
		if m.preview != nil {
			switch strings.ToLower(message.String()) {
			case "y", "enter":
				previewed := m.preview.Command
				m.preview = nil
				m.loading = true
				return m, m.executeCommandCmd(previewed)
			case "n", "esc":
				m.preview = nil
				return m, nil
			}
			return m, nil
		}

		switch message.String() {
		case "esc":
			return m, func() tea.Msg { return screenChangeMsg{target: ScreenMachineList} }
		case "tab", "down":
			m.focus = (m.focus + 1) % len(m.inputs)
			m.updateFocus()
			return m, nil
		case "shift+tab", "up":
			m.focus = (m.focus - 1 + len(m.inputs)) % len(m.inputs)
			m.updateFocus()
			return m, nil
		case "enter":
			cmd, err := (services.MachineCreateBuilder{Image: m.inputs[0].Value(), Name: m.inputs[1].Value()}).Build()
			if err != nil {
				m.errorMsg = err.Error()
				return m, nil
			}
			m.preview = &CommandPreviewModal{Title: "Create Machine", Command: cmd}
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

func (m MachineCreateScreen) View() string {
	builder := strings.Builder{}
	builder.WriteString(RenderTitle("Create Machine") + "\n\n")
	builder.WriteString("Image\n" + m.inputs[0].View() + "\n\n")
	builder.WriteString("Name (optional)\n" + m.inputs[1].View() + "\n")
	if m.loading {
		builder.WriteString("\n" + RenderMuted("Creating machine...") + "\n")
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
	builder.WriteString("\n" + RenderMuted("Keys: tab/up/down=field, enter=preview, esc=back") + "\n")
	return builder.String()
}

func (m *MachineCreateScreen) updateFocus() {
	for i := range m.inputs {
		if i == m.focus {
			m.inputs[i].Focus()
		} else {
			m.inputs[i].Blur()
		}
	}
}

func (m MachineCreateScreen) executeCommandCmd(command models.Command) tea.Cmd {
	return func() tea.Msg {
		result, err := m.executor.Execute(command)
		return machineCreateResultMsg{result: result, err: err}
	}
}

func newMachineCreateInputs() []textinput.Model {
	image := textinput.New()
	image.Prompt = ""
	image.Placeholder = "alpine:latest"
	image.Focus()

	name := textinput.New()
	name.Prompt = ""
	name.Placeholder = "dev"

	return []textinput.Model{image, name}
}
