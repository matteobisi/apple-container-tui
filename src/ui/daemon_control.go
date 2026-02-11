package ui

import (
	"errors"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"container-tui/src/models"
	"container-tui/src/services"
)

type daemonStatusMsg struct {
	status models.DaemonStatus
	err    error
}

type daemonActionMsg struct {
	result models.Result
	err    error
}

// DaemonControlScreen manages daemon start/stop.
type DaemonControlScreen struct {
	executor services.CommandExecutor
	status   models.DaemonStatus
	loading  bool
	errorMsg string
	result   *models.Result
	confirm  *YesNoConfirmModal
}

// NewDaemonControlScreen creates the daemon control screen.
func NewDaemonControlScreen(executor services.CommandExecutor) DaemonControlScreen {
	return DaemonControlScreen{executor: executor}
}

// Init loads daemon status.
func (m DaemonControlScreen) Init() tea.Cmd {
	return m.fetchStatusCmd()
}

// Update handles daemon control input.
func (m DaemonControlScreen) Update(msg tea.Msg) (DaemonControlScreen, tea.Cmd) {
	switch message := msg.(type) {
	case daemonStatusMsg:
		m.loading = false
		if message.err != nil {
			m.errorMsg = services.FormatError(message.err, "")
			return m, nil
		}
		m.errorMsg = ""
		m.status = message.status
		return m, nil
	case daemonActionMsg:
		m.loading = false
		if message.err != nil {
			m.errorMsg = services.FormatError(message.err, message.result.Stderr)
		}
		m.result = &message.result
		return m, m.fetchStatusCmd()
	case tea.KeyMsg:
		if m.confirm != nil {
			confirmed, canceled := m.confirm.Handle(message)
			if confirmed {
				command := m.confirm.Command
				m.confirm = nil
				m.loading = true
				return m, m.executeCommandCmd(command)
			}
			if canceled {
				m.confirm = nil
				return m, nil
			}
			return m, nil
		}

		switch message.String() {
		case "esc":
			return m, func() tea.Msg { return screenChangeMsg{target: ScreenContainerList} }
		case "?":
			return m, func() tea.Msg { return screenChangeMsg{target: ScreenHelp} }
		case "r":
			m.loading = true
			return m, m.fetchStatusCmd()
		case "s":
			updated, cmd := m.buildAndConfirmStart()
			return updated, cmd
		case "t":
			updated, cmd := m.buildAndConfirmStop()
			return updated, cmd
		}
	}

	return m, nil
}

// View renders daemon controls.
func (m DaemonControlScreen) View() string {
	builder := strings.Builder{}
	builder.WriteString(RenderTitle("Daemon Control") + "\n\n")
	status := "stopped"
	if m.status.Running {
		status = "running"
	}
	statusLine := "Status: " + status
	if m.status.Running {
		statusLine = RenderSuccess(statusLine)
	} else {
		statusLine = RenderWarning(statusLine)
	}
	builder.WriteString(statusLine + "\n")
	if m.loading {
		builder.WriteString("\n" + RenderMuted("Loading...") + "\n")
	}
	if m.errorMsg != "" {
		builder.WriteString("\n" + RenderError("Error: "+m.errorMsg) + "\n")
	}
	builder.WriteString("\nActions:\n")
	builder.WriteString("  s - start daemon\n")
	builder.WriteString("  t - stop daemon (!)\n")
	builder.WriteString("\n" + RenderMuted("Keys: s=start, t=stop, r=refresh, ?=help, esc=back") + "\n")

	if m.confirm != nil {
		builder.WriteString("\n")
		builder.WriteString(m.confirm.View())
	}
	if m.result != nil {
		builder.WriteString("\n\n")
		builder.WriteString(RenderResult(*m.result))
	}
	return builder.String()
}

func (m DaemonControlScreen) buildAndConfirmStart() (DaemonControlScreen, tea.Cmd) {
	builder := services.StartDaemonBuilder{}
	cmd, err := builder.Build()
	if err != nil {
		m.errorMsg = err.Error()
		return m, nil
	}
	m.confirm = &YesNoConfirmModal{
		Title:   "Start Daemon",
		Body:    "Start container services?",
		Command: cmd,
		Warning: false,
	}
	return m, nil
}

func (m DaemonControlScreen) buildAndConfirmStop() (DaemonControlScreen, tea.Cmd) {
	builder := services.StopDaemonBuilder{}
	cmd, err := builder.Build()
	if err != nil {
		m.errorMsg = err.Error()
		return m, nil
	}
	m.confirm = &YesNoConfirmModal{
		Title:   "Stop Daemon",
		Body:    "Stop container services and running containers?",
		Command: cmd,
		Warning: true,
	}
	return m, nil
}

func (m DaemonControlScreen) fetchStatusCmd() tea.Cmd {
	return func() tea.Msg {
		builder := services.CheckDaemonStatusBuilder{}
		cmd, err := builder.Build()
		if err != nil {
			return daemonStatusMsg{err: err}
		}
		result, err := m.executor.Execute(cmd)
		if err != nil {
			return daemonStatusMsg{err: errors.New(services.FormatError(err, result.Stderr))}
		}
		status := services.ParseDaemonStatus(result.Stdout)
		return daemonStatusMsg{status: status}
	}
}

func (m DaemonControlScreen) executeCommandCmd(command models.Command) tea.Cmd {
	return func() tea.Msg {
		result, err := m.executor.Execute(command)
		return daemonActionMsg{result: result, err: err}
	}
}
