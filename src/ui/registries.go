package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"container-tui/src/models"
	"container-tui/src/services"
)

type registriesLoadedMsg struct {
	entries []models.RegistryLogin
	err     error
}

// RegistriesScreen shows runtime-managed registries.
type RegistriesScreen struct {
	executor services.CommandExecutor
	entries  []models.RegistryLogin
	cursor   int
	width    int
	loading  bool
	errorMsg string
}

func NewRegistriesScreen(executor services.CommandExecutor) RegistriesScreen {
	return RegistriesScreen{executor: executor, entries: []models.RegistryLogin{}}
}

func (m RegistriesScreen) Init() tea.Cmd {
	m.loading = true
	return m.fetchRegistriesCmd()
}

func (m RegistriesScreen) Update(msg tea.Msg) (RegistriesScreen, tea.Cmd) {
	switch message := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = message.Width
		return m, nil
	case registriesLoadedMsg:
		m.loading = false
		if message.err != nil {
			m.errorMsg = message.err.Error()
			return m, nil
		}
		m.entries = message.entries
		if m.cursor >= len(m.entries) {
			m.cursor = max(0, len(m.entries)-1)
		}
		if len(m.entries) == 0 {
			m.errorMsg = "No registries found"
		} else {
			m.errorMsg = ""
		}
		return m, nil
	case tea.KeyMsg:
		switch message.String() {
		case "up", "k":
			m.cursor = max(0, m.cursor-1)
		case "down", "j":
			m.cursor = min(len(m.entries)-1, m.cursor+1)
		case "r":
			m.loading = true
			return m, m.fetchRegistriesCmd()
		case "?":
			return m, func() tea.Msg { return screenChangeMsg{target: ScreenHelp} }
		case "esc":
			return m, func() tea.Msg { return BackToListMsg{} }
		}
	}
	return m, nil
}

func (m RegistriesScreen) View() string {
	builder := strings.Builder{}
	builder.WriteString(RenderTitle("Registries") + "\n\n")
	if m.loading {
		builder.WriteString(RenderMuted("Loading registries...") + "\n")
	}

	table := NewTable([]TableColumn{
		{Header: "Hostname", MinWidth: 20, Priority: 1, Align: "left"},
		{Header: "Username", MinWidth: 16, Priority: 2, Align: "left"},
	})

	if len(m.entries) > 0 {
		rows := make([]TableRow, len(m.entries))
		for index, entry := range m.entries {
			username := entry.Username
			if strings.TrimSpace(username) == "" {
				username = "(no username)"
			}
			rows[index] = TableRow{Cells: []string{entry.Hostname, username}, Selected: index == m.cursor}
		}
		table.SetRows(rows)
	}

	tableWidth := m.width
	if tableWidth == 0 {
		tableWidth = 80
	}
	builder.WriteString(table.Render(tableWidth, m.cursor))
	builder.WriteString(strings.Repeat("─", tableWidth) + "\n")
	if m.errorMsg != "" {
		builder.WriteString("\n" + RenderMuted(m.errorMsg) + "\n")
	}
	builder.WriteString("\n" + RenderMuted("Keys: up/down=navigate, r=refresh, ?=help, esc=back") + "\n")
	return builder.String()
}

func (m RegistriesScreen) fetchRegistriesCmd() tea.Cmd {
	return func() tea.Msg {
		cmd, err := (services.RegistryListBuilder{}).Build()
		if err != nil {
			return registriesLoadedMsg{err: err}
		}
		result, err := m.executor.Execute(cmd)
		if err != nil {
			return registriesLoadedMsg{err: err}
		}
		entries, parseErr := services.ParseRegistryList(result.Stdout)
		if parseErr != nil {
			return registriesLoadedMsg{err: parseErr}
		}
		return registriesLoadedMsg{entries: entries}
	}
}
