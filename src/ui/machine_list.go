package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"container-tui/src/models"
	"container-tui/src/services"
)

type machineListLoadedMsg struct {
	machines []models.ContainerMachine
	err      error
}

// MachineListScreen displays container machines.
type MachineListScreen struct {
	executor  services.CommandExecutor
	machines  []models.ContainerMachine
	cursor    int
	width     int
	loading   bool
	errorMsg  string
	hasLoaded bool
}

func NewMachineListScreen(executor services.CommandExecutor) MachineListScreen {
	return MachineListScreen{executor: executor}
}

func (m MachineListScreen) Init() tea.Cmd {
	return m.fetchMachinesCmd()
}

func (m MachineListScreen) Update(msg tea.Msg) (MachineListScreen, tea.Cmd) {
	switch message := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = message.Width
		return m, nil
	case machineListLoadedMsg:
		m.loading = false
		m.hasLoaded = true
		if message.err != nil {
			m.errorMsg = services.FormatError(message.err, "")
			return m, nil
		}
		m.errorMsg = ""
		m.machines = message.machines
		if m.cursor >= len(m.machines) {
			m.cursor = max(0, len(m.machines)-1)
		}
		return m, nil
	case tea.KeyMsg:
		switch message.String() {
		case "up", "k":
			m.cursor = max(0, m.cursor-1)
		case "down", "j":
			m.cursor = min(len(m.machines)-1, m.cursor+1)
		case "r":
			m.loading = true
			return m, m.fetchMachinesCmd()
		case "c":
			return m, func() tea.Msg { return screenChangeMsg{target: ScreenMachineCreate, push: true} }
		case "esc":
			return m, func() tea.Msg { return screenChangeMsg{target: ScreenContainerList} }
		case "enter":
			selected, ok := m.selectedMachine()
			if !ok {
				return m, nil
			}
			machineCopy := selected
			return m, func() tea.Msg {
				return screenChangeMsg{target: ScreenMachineSubmenu, machine: &machineCopy, push: true}
			}
		}
	}
	return m, nil
}

func (m MachineListScreen) View() string {
	builder := strings.Builder{}
	builder.WriteString(RenderTitle("Container Machines") + "\n\n")
	if m.loading {
		builder.WriteString(RenderMuted("Loading machines...") + "\n")
	}
	if m.errorMsg != "" {
		builder.WriteString(RenderError("Error: "+m.errorMsg) + "\n\n")
	}

	table := NewTable([]TableColumn{
		{Header: "Name", MinWidth: 12, Priority: 1, Align: "left"},
		{Header: "State", MinWidth: 8, Priority: 2, Align: "left"},
		{Header: "Image", MinWidth: 16, Priority: 3, Align: "left"},
		{Header: "Default", MinWidth: 8, Priority: 4, Align: "left"},
	})

	if len(m.machines) > 0 {
		rows := make([]TableRow, len(m.machines))
		for i, machine := range m.machines {
			defaultValue := ""
			if machine.IsDefault {
				defaultValue = "yes"
			}
			rows[i] = TableRow{
				Cells:    []string{machine.ID, string(machine.State), machine.Image, defaultValue},
				Selected: i == m.cursor,
				Data:     &machine,
			}
		}
		table.SetRows(rows)
	}

	tableWidth := m.width
	if tableWidth == 0 {
		tableWidth = 80
	}
	if len(m.machines) == 0 && m.hasLoaded && m.errorMsg == "" {
		builder.WriteString(RenderMuted("No container machines found.") + "\n")
	} else {
		builder.WriteString(table.Render(tableWidth, m.cursor))
	}
	builder.WriteString(strings.Repeat("─", tableWidth) + "\n")
	builder.WriteString("\n" + RenderMuted("Keys: up/down, enter=submenu, c=create, r=refresh, esc=containers, q=quit") + "\n")
	return builder.String()
}

func (m MachineListScreen) selectedMachine() (models.ContainerMachine, bool) {
	if len(m.machines) == 0 || m.cursor < 0 || m.cursor >= len(m.machines) {
		return models.ContainerMachine{}, false
	}
	return m.machines[m.cursor], true
}

func (m MachineListScreen) fetchMachinesCmd() tea.Cmd {
	return func() tea.Msg {
		cmd, err := (services.MachineListBuilder{}).Build()
		if err != nil {
			return machineListLoadedMsg{err: err}
		}
		result, err := m.executor.Execute(cmd)
		if err != nil {
			return machineListLoadedMsg{err: err, machines: nil}
		}
		machines, err := services.ParseMachineList(result.Stdout)
		return machineListLoadedMsg{machines: machines, err: err}
	}
}
