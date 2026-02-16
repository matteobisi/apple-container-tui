package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"container-tui/src/models"
	"container-tui/src/services"
)

type containerLogsLoadedMsg struct {
	result models.Result
	err    error
}

// ContainerLogsScreen shows container log output.
type ContainerLogsScreen struct {
	executor  services.CommandExecutor
	container models.Container
	loading   bool
	errorMsg  string
	lines     []string
}

func NewContainerLogsScreen(executor services.CommandExecutor) ContainerLogsScreen {
	return ContainerLogsScreen{executor: executor, lines: []string{}}
}

func (m ContainerLogsScreen) SetContainer(container models.Container) ContainerLogsScreen {
	m.container = container
	m.errorMsg = ""
	m.lines = []string{}
	m.loading = true
	return m
}

func (m ContainerLogsScreen) Init() tea.Cmd {
	if strings.TrimSpace(m.container.ID) == "" {
		return nil
	}
	return m.loadLogsCmd()
}

func (m ContainerLogsScreen) Update(msg tea.Msg) (ContainerLogsScreen, tea.Cmd) {
	switch message := msg.(type) {
	case containerLogsLoadedMsg:
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
		return m, nil
	case tea.KeyMsg:
		switch message.String() {
		case "esc":
			containerCopy := m.container
			return m, func() tea.Msg { return screenChangeMsg{target: ScreenContainerSubmenu, container: &containerCopy} }
		}
	}
	return m, nil
}

func (m ContainerLogsScreen) View() string {
	builder := strings.Builder{}
	builder.WriteString(RenderTitle("Container Logs") + "\n\n")
	builder.WriteString(RenderMuted("Container: "+m.container.Name) + "\n\n")
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

func (m ContainerLogsScreen) loadLogsCmd() tea.Cmd {
	return func() tea.Msg {
		builder := services.ContainerLogsBuilder{ContainerName: m.container.ID}
		cmd, err := builder.Build()
		if err != nil {
			return containerLogsLoadedMsg{err: err}
		}
		result, err := m.executor.Execute(cmd)
		if err != nil {
			return containerLogsLoadedMsg{result: result, err: err}
		}
		return containerLogsLoadedMsg{result: result}
	}
}
