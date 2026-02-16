package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"container-tui/src/models"
	"container-tui/src/services"
)

type containerShellFinishedMsg struct {
	result models.Result
	err    error
}

// ContainerShellScreen detects shell and executes interactive container shell command.
type ContainerShellScreen struct {
	executor      services.CommandExecutor
	detector      *services.ShellDetector
	container     models.Container
	loading       bool
	errorMsg      string
	statusMessage string
}

func NewContainerShellScreen(executor services.CommandExecutor) ContainerShellScreen {
	return ContainerShellScreen{executor: executor, detector: services.NewShellDetector(executor)}
}

func (m ContainerShellScreen) SetContainer(container models.Container) ContainerShellScreen {
	m.container = container
	m.loading = true
	m.errorMsg = ""
	m.statusMessage = "Preparing shell session..."
	return m
}

func (m ContainerShellScreen) Init() tea.Cmd {
	if strings.TrimSpace(m.container.ID) == "" {
		return nil
	}
	return m.executeShellCmd()
}

func (m ContainerShellScreen) Update(msg tea.Msg) (ContainerShellScreen, tea.Cmd) {
	switch message := msg.(type) {
	case containerShellFinishedMsg:
		m.loading = false
		if message.err != nil {
			m.errorMsg = services.FormatError(message.err, message.result.Stderr)
			m.statusMessage = "Shell unavailable."
		} else {
			m.statusMessage = "Shell session completed."
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

func (m ContainerShellScreen) View() string {
	builder := strings.Builder{}
	builder.WriteString(RenderTitle("Container Shell") + "\n\n")
	builder.WriteString(RenderMuted("Container: "+m.container.Name) + "\n\n")
	if m.loading {
		builder.WriteString(RenderMuted("Opening shell...") + "\n")
	}
	if m.statusMessage != "" {
		builder.WriteString(m.statusMessage + "\n")
	}
	if m.errorMsg != "" {
		builder.WriteString(RenderError("Error: "+m.errorMsg) + "\n")
	}
	builder.WriteString("\n" + RenderMuted("Keys: esc=back") + "\n")
	return builder.String()
}

func (m ContainerShellScreen) executeShellCmd() tea.Cmd {
	return func() tea.Msg {
		shell, err := m.detector.DetectShell(m.container.ID)
		if err != nil {
			return containerShellFinishedMsg{err: err}
		}
		builder := services.ContainerExecBuilder{ContainerName: m.container.ID, Shell: shell}
		cmd, err := builder.Build()
		if err != nil {
			return containerShellFinishedMsg{err: err}
		}
		result, err := m.executor.Execute(cmd)
		return containerShellFinishedMsg{result: result, err: err}
	}
}
