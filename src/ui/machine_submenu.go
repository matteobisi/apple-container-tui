package ui

import (
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"container-tui/src/models"
	"container-tui/src/services"
)

type machineSubmenuActionMsg struct {
	action string
	result models.Result
	err    error
}

type machineSubmenuOption struct {
	label  string
	action string
}

// MachineSubmenuScreen displays actions for the selected container machine.
type MachineSubmenuScreen struct {
	executor services.CommandExecutor
	machine  models.ContainerMachine
	options  []machineSubmenuOption
	cursor   int
	loading  bool
	errorMsg string
	result   *models.Result
	preview  *CommandPreviewModal
	confirm  *TypeToConfirmModal
	width    int
}

func NewMachineSubmenuScreen(executor services.CommandExecutor) MachineSubmenuScreen {
	return MachineSubmenuScreen{executor: executor}
}

func (m MachineSubmenuScreen) SetMachine(machine models.ContainerMachine) MachineSubmenuScreen {
	m.machine = machine
	m.cursor = 0
	m.loading = false
	m.errorMsg = ""
	m.result = nil
	m.preview = nil
	m.confirm = nil
	m.options = m.buildOptions()
	return m
}

func (m MachineSubmenuScreen) Init() tea.Cmd { return nil }

func (m MachineSubmenuScreen) Update(msg tea.Msg) (MachineSubmenuScreen, tea.Cmd) {
	switch message := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = message.Width
		return m, nil
	case machineSubmenuActionMsg:
		m.loading = false
		if message.err != nil {
			m.errorMsg = services.FormatError(message.err, message.result.Stderr)
		} else {
			m.errorMsg = ""
		}
		m.result = &message.result
		if message.err == nil && message.action == "delete" {
			return m, func() tea.Msg { return screenChangeMsg{target: ScreenMachineList} }
		}
		return m, nil
	case tea.KeyMsg:
		if m.confirm != nil {
			updatedConfirm, confirmed, canceled := m.confirm.Handle(message)
			m.confirm = &updatedConfirm
			if confirmed {
				command := updatedConfirm.Command
				m.confirm = nil
				m.loading = true
				return m, m.executeCommandCmd("delete", command)
			}
			if canceled {
				m.confirm = nil
				return m, nil
			}
			return m, nil
		}
		if m.preview != nil {
			switch strings.ToLower(message.String()) {
			case "y", "enter":
				previewed := m.preview.Command
				action := m.preview.Title
				m.preview = nil
				m.loading = true
				return m, m.executeCommandCmd(action, previewed)
			case "n", "esc":
				m.preview = nil
				return m, nil
			}
			return m, nil
		}

		switch message.String() {
		case "up", "k":
			m.cursor = max(0, m.cursor-1)
		case "down", "j":
			m.cursor = min(len(m.options)-1, m.cursor+1)
		case "esc":
			return m, func() tea.Msg { return screenChangeMsg{target: ScreenMachineList} }
		case "enter":
			if len(m.options) == 0 || m.cursor < 0 || m.cursor >= len(m.options) {
				return m, nil
			}
			return m.selectOption(m.options[m.cursor])
		}
	}
	return m, nil
}

func (m MachineSubmenuScreen) View() string {
	builder := strings.Builder{}
	builder.WriteString(RenderTitle("Machine Actions") + "\n\n")

	headerStyle := lipgloss.NewStyle().Bold(true)
	builder.WriteString(headerStyle.Render("Machine Details") + "\n")
	builder.WriteString("name: " + m.machine.ID + "\n")
	builder.WriteString("state: " + string(m.machine.State) + "\n")
	builder.WriteString("image: " + m.machine.Image + "\n")
	builder.WriteString("cpus: " + intToDisplay(m.machine.CPUs) + "\n")
	builder.WriteString("memory: " + m.machine.Memory + "\n")
	builder.WriteString("home-mount: " + m.machine.HomeMount + "\n")
	if m.machine.IsDefault {
		builder.WriteString("default: yes\n")
	}
	builder.WriteString("\n")

	width := m.width
	if width == 0 {
		width = 80
	}
	builder.WriteString(strings.Repeat("─", width) + "\n\n")
	builder.WriteString(headerStyle.Render("Available Actions") + "\n")

	normalStyle := lipgloss.NewStyle()
	selectedStyle := lipgloss.NewStyle().Reverse(true)
	for i, option := range m.options {
		style := normalStyle
		if i == m.cursor {
			style = selectedStyle
		}
		builder.WriteString(style.Render(option.label) + "\n")
	}
	builder.WriteString("\n")
	builder.WriteString(strings.Repeat("─", width) + "\n")

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

func (m MachineSubmenuScreen) buildOptions() []machineSubmenuOption {
	options := make([]machineSubmenuOption, 0, 8)
	options = append(options, machineSubmenuOption{label: "Inspect machine", action: "inspect"})
	options = append(options, machineSubmenuOption{label: "View machine logs", action: "logs"})
	switch m.machine.NormalizedState() {
	case models.MachineStateRunning:
		options = append(options, machineSubmenuOption{label: "Stop machine", action: "stop"})
	case models.MachineStateStopped:
		options = append(options, machineSubmenuOption{label: "Start machine", action: "start"})
	}
	options = append(options, machineSubmenuOption{label: "Edit resources", action: "edit"})
	if !m.machine.IsDefault {
		options = append(options, machineSubmenuOption{label: "Set as default", action: "set-default"})
	}
	options = append(options, machineSubmenuOption{label: "Delete machine", action: "delete"})
	options = append(options, machineSubmenuOption{label: "Back", action: "back"})
	return options
}

func (m MachineSubmenuScreen) selectOption(option machineSubmenuOption) (MachineSubmenuScreen, tea.Cmd) {
	switch option.action {
	case "inspect":
		machineCopy := m.machine
		return m, func() tea.Msg {
			return screenChangeMsg{target: ScreenMachineInspect, machine: &machineCopy, push: true}
		}
	case "logs":
		machineCopy := m.machine
		return m, func() tea.Msg { return screenChangeMsg{target: ScreenMachineLogs, machine: &machineCopy, push: true} }
	case "edit":
		machineCopy := m.machine
		return m, func() tea.Msg {
			return screenChangeMsg{target: ScreenMachineEditResources, machine: &machineCopy, push: true}
		}
	case "start":
		cmd, err := (services.MachineStartBuilder{MachineID: m.machine.ID}).Build()
		if err != nil {
			m.errorMsg = err.Error()
			return m, nil
		}
		m.preview = &CommandPreviewModal{Title: "start", Command: cmd}
	case "stop":
		cmd, err := (services.MachineStopBuilder{MachineID: m.machine.ID}).Build()
		if err != nil {
			m.errorMsg = err.Error()
			return m, nil
		}
		m.preview = &CommandPreviewModal{Title: "stop", Command: cmd}
	case "set-default":
		cmd, err := (services.MachineSetDefaultBuilder{MachineID: m.machine.ID}).Build()
		if err != nil {
			m.errorMsg = err.Error()
			return m, nil
		}
		m.preview = &CommandPreviewModal{Title: "set-default", Command: cmd}
	case "delete":
		cmd, err := (services.MachineDeleteBuilder{MachineID: m.machine.ID}).Build()
		if err != nil {
			m.errorMsg = err.Error()
			return m, nil
		}
		confirm := NewTypeToConfirmModal("Delete Machine", m.machine.ID, cmd)
		m.confirm = &confirm
	case "back":
		return m, func() tea.Msg { return screenChangeMsg{target: ScreenMachineList} }
	}
	return m, nil
}

func (m MachineSubmenuScreen) executeCommandCmd(action string, command models.Command) tea.Cmd {
	return func() tea.Msg {
		result, err := m.executor.Execute(command)
		return machineSubmenuActionMsg{action: action, result: result, err: err}
	}
}

func intToDisplay(value int) string {
	if value == 0 {
		return ""
	}
	return strconv.Itoa(value)
}
