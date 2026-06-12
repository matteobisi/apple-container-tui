package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"container-tui/src/models"
	"container-tui/src/services"
)

type machineEditActionMsg struct {
	result models.Result
	err    error
}

// MachineEditResourcesScreen edits machine CPU, memory, and home mount settings.
type MachineEditResourcesScreen struct {
	executor services.CommandExecutor
	machine  models.ContainerMachine
	inputs   []textinput.Model
	focus    int
	loading  bool
	errorMsg string
	result   *models.Result
	preview  *CommandPreviewModal
}

func NewMachineEditResourcesScreen(executor services.CommandExecutor) MachineEditResourcesScreen {
	return MachineEditResourcesScreen{executor: executor, inputs: newMachineEditInputs(models.ContainerMachine{})}
}

func (m MachineEditResourcesScreen) SetMachine(machine models.ContainerMachine) MachineEditResourcesScreen {
	m.machine = machine
	m.inputs = newMachineEditInputs(machine)
	m.focus = 0
	m.loading = false
	m.errorMsg = ""
	m.result = nil
	m.preview = nil
	return m
}

func (m MachineEditResourcesScreen) Init() tea.Cmd { return nil }

func (m MachineEditResourcesScreen) Update(msg tea.Msg) (MachineEditResourcesScreen, tea.Cmd) {
	switch message := msg.(type) {
	case machineEditActionMsg:
		m.loading = false
		if message.err != nil {
			m.errorMsg = services.FormatError(message.err, message.result.Stderr)
		} else {
			m.errorMsg = ""
		}
		m.result = &message.result
		return m, nil
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
		case "esc", "q":
			machineCopy := m.machine
			return m, func() tea.Msg { return screenChangeMsg{target: ScreenMachineSubmenu, machine: &machineCopy} }
		case "tab", "down":
			m.focus = (m.focus + 1) % len(m.inputs)
			m.updateFocus()
			return m, nil
		case "shift+tab", "up":
			m.focus = (m.focus - 1 + len(m.inputs)) % len(m.inputs)
			m.updateFocus()
			return m, nil
		case "enter":
			cmd, err := (services.MachineSetBuilder{MachineID: m.machine.ID, CPUs: m.inputs[0].Value(), Memory: m.inputs[1].Value(), HomeMount: m.inputs[2].Value()}).Build()
			if err != nil {
				m.errorMsg = err.Error()
				return m, nil
			}
			m.preview = &CommandPreviewModal{Title: "Edit Machine Resources", Command: cmd}
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

func (m MachineEditResourcesScreen) View() string {
	builder := strings.Builder{}
	builder.WriteString(RenderTitle("Edit Machine Resources") + "\n\n")
	builder.WriteString(RenderMuted("Machine: "+m.machine.ID) + "\n\n")
	builder.WriteString("CPUs\n" + m.inputs[0].View() + "\n\n")
	builder.WriteString("Memory\n" + m.inputs[1].View() + "\n\n")
	builder.WriteString("Home mount\n" + m.inputs[2].View() + "\n")
	builder.WriteString("\n" + RenderMuted("Changes take effect after next stop and restart") + "\n")
	if m.loading {
		builder.WriteString("\n" + RenderMuted("Updating resources...") + "\n")
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

func (m *MachineEditResourcesScreen) updateFocus() {
	for i := range m.inputs {
		if i == m.focus {
			m.inputs[i].Focus()
		} else {
			m.inputs[i].Blur()
		}
	}
}

func (m MachineEditResourcesScreen) executeCommandCmd(command models.Command) tea.Cmd {
	return func() tea.Msg {
		result, err := m.executor.Execute(command)
		return machineEditActionMsg{result: result, err: err}
	}
}

func newMachineEditInputs(machine models.ContainerMachine) []textinput.Model {
	cpus := textinput.New()
	cpus.Prompt = ""
	cpus.Placeholder = "2"
	if machine.CPUs > 0 {
		cpus.SetValue(fmt.Sprintf("%d", machine.CPUs))
	}
	cpus.Focus()

	memory := textinput.New()
	memory.Prompt = ""
	memory.Placeholder = "4G"
	memory.SetValue(machine.Memory)

	homeMount := textinput.New()
	homeMount.Prompt = ""
	homeMount.Placeholder = "rw"
	homeMount.SetValue(machine.NormalizedHomeMount())

	return []textinput.Model{cpus, memory, homeMount}
}
