package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"container-tui/src/models"
	"container-tui/src/services"
)

type machineLogsLoadedMsg struct {
	result models.Result
	err    error
}

// MachineLogsScreen shows container machine log output.
type MachineLogsScreen struct {
	executor services.CommandExecutor
	machine  models.ContainerMachine
	loading  bool
	errorMsg string
	lines    []string
}

func NewMachineLogsScreen(executor services.CommandExecutor) MachineLogsScreen {
	return MachineLogsScreen{executor: executor, lines: []string{}}
}

func (m MachineLogsScreen) SetMachine(machine models.ContainerMachine) MachineLogsScreen {
	m.machine = machine
	m.errorMsg = ""
	m.lines = []string{}
	m.loading = true
	return m
}

func (m MachineLogsScreen) Init() tea.Cmd {
	if strings.TrimSpace(m.machine.ID) == "" {
		return nil
	}
	return m.loadLogsCmd()
}

func (m MachineLogsScreen) Update(msg tea.Msg) (MachineLogsScreen, tea.Cmd) {
	switch message := msg.(type) {
	case machineLogsLoadedMsg:
		m.loading = false
		if message.err != nil {
			m.errorMsg = services.FormatError(message.err, message.result.Stderr)
			return m, nil
		}
		content := strings.TrimSpace(message.result.Stdout)
		if content == "" {
			m.lines = []string{"No logs available."}
		} else {
			m.lines = strings.Split(content, "\n")
		}
	case tea.KeyMsg:
		switch message.String() {
		case "esc", "q":
			machineCopy := m.machine
			return m, func() tea.Msg { return screenChangeMsg{target: ScreenMachineSubmenu, machine: &machineCopy} }
		}
	}
	return m, nil
}

func (m MachineLogsScreen) View() string {
	builder := strings.Builder{}
	builder.WriteString(RenderTitle("Machine Logs") + "\n\n")
	builder.WriteString(RenderMuted("Machine: "+m.machine.ID) + "\n\n")
	if m.loading {
		builder.WriteString(RenderMuted("Loading logs...") + "\n")
	}
	if m.errorMsg != "" {
		builder.WriteString(RenderError("Error: "+m.errorMsg) + "\n")
	}
	for _, line := range m.lines {
		builder.WriteString(line + "\n")
	}
	builder.WriteString("\n" + RenderMuted("Keys: esc=back") + "\n")
	return builder.String()
}

func (m MachineLogsScreen) loadLogsCmd() tea.Cmd {
	return func() tea.Msg {
		cmd, err := (services.MachineLogsBuilder{MachineID: m.machine.ID}).Build()
		if err != nil {
			return machineLogsLoadedMsg{err: err}
		}
		result, err := m.executor.Execute(cmd)
		return machineLogsLoadedMsg{result: result, err: err}
	}
}
