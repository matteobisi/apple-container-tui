package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"

	"container-tui/src/models"
	"container-tui/src/services"
)

type machineInspectLoadedMsg struct {
	content string
	err     error
}

// MachineInspectScreen shows machine inspect output.
type MachineInspectScreen struct {
	executor services.CommandExecutor
	machine  models.ContainerMachine
	viewport viewport.Model
	loading  bool
	errorMsg string
	width    int
	height   int
}

func NewMachineInspectScreen(executor services.CommandExecutor) MachineInspectScreen {
	return MachineInspectScreen{executor: executor, viewport: viewport.New(0, 0)}
}

func (m MachineInspectScreen) SetMachine(machine models.ContainerMachine) MachineInspectScreen {
	m.machine = machine
	m.loading = true
	m.errorMsg = ""
	m.viewport.SetContent("")
	return m
}

func (m MachineInspectScreen) Init() tea.Cmd {
	if strings.TrimSpace(m.machine.ID) == "" {
		return nil
	}
	return m.loadInspectCmd()
}

func (m MachineInspectScreen) Update(msg tea.Msg) (MachineInspectScreen, tea.Cmd) {
	switch message := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = message.Width
		m.height = message.Height
		m.viewport.Width = max(20, message.Width-4)
		m.viewport.Height = max(5, message.Height-8)
	case machineInspectLoadedMsg:
		m.loading = false
		if message.err != nil {
			m.errorMsg = message.err.Error()
			return m, nil
		}
		m.viewport.SetContent(message.content)
	case tea.KeyMsg:
		switch message.String() {
		case "up", "down", "pgup", "pgdown", "home", "end":
			updatedViewport, _ := m.viewport.Update(message)
			m.viewport = updatedViewport
			return m, nil
		case "esc", "q":
			machineCopy := m.machine
			return m, func() tea.Msg { return screenChangeMsg{target: ScreenMachineSubmenu, machine: &machineCopy} }
		}
	}
	return m, nil
}

func (m MachineInspectScreen) View() string {
	builder := strings.Builder{}
	builder.WriteString(RenderTitle("Machine Inspection") + "\n\n")
	builder.WriteString(RenderMuted("Machine: "+m.machine.ID) + "\n\n")
	if m.loading {
		builder.WriteString(RenderMuted("Loading inspection...") + "\n")
	}
	if m.errorMsg != "" {
		builder.WriteString(RenderError("Error: "+m.errorMsg) + "\n")
	}
	builder.WriteString(m.viewport.View() + "\n")
	builder.WriteString("\n" + RenderMuted("Keys: up/down/pgup/pgdn/home/end=scroll, esc=back") + "\n")
	return builder.String()
}

func (m MachineInspectScreen) loadInspectCmd() tea.Cmd {
	return func() tea.Msg {
		cmd, err := (services.MachineInspectBuilder{MachineID: m.machine.ID}).Build()
		if err != nil {
			return machineInspectLoadedMsg{err: err}
		}
		result, err := m.executor.Execute(cmd)
		if err != nil {
			return machineInspectLoadedMsg{err: err}
		}
		content := strings.TrimSpace(result.Stdout)
		if content == "" {
			content = result.Stderr
		}
		if strings.TrimSpace(content) == "" {
			content = "No inspection output."
		}
		return machineInspectLoadedMsg{content: content}
	}
}
