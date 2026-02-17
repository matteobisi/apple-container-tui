package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"container-tui/src/models"
	"container-tui/src/services"
)

type containerSubmenuActionMsg struct {
	result models.Result
	err    error
}

type containerSubmenuOption struct {
	label  string
	action string
}

// ContainerSubmenuScreen displays actions for the selected container.
type ContainerSubmenuScreen struct {
	executor  services.CommandExecutor
	container models.Container
	options   []containerSubmenuOption
	cursor    int
	loading   bool
	errorMsg  string
	result    *models.Result
	preview   *CommandPreviewModal
	width     int
}

func NewContainerSubmenuScreen(executor services.CommandExecutor) ContainerSubmenuScreen {
	return ContainerSubmenuScreen{executor: executor}
}

func (m ContainerSubmenuScreen) SetContainer(container models.Container) ContainerSubmenuScreen {
	m.container = container
	m.cursor = 0
	m.errorMsg = ""
	m.result = nil
	m.preview = nil
	m.options = m.buildOptions()
	return m
}

func (m ContainerSubmenuScreen) Init() tea.Cmd { return nil }

func (m ContainerSubmenuScreen) Update(msg tea.Msg) (ContainerSubmenuScreen, tea.Cmd) {
	switch message := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = message.Width
		return m, nil
	case containerSubmenuActionMsg:
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
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.options)-1 {
				m.cursor++
			}
		case "esc":
			return m, func() tea.Msg { return screenChangeMsg{target: ScreenContainerList} }
		case "enter":
			if len(m.options) == 0 || m.cursor < 0 || m.cursor >= len(m.options) {
				return m, nil
			}
			selected := m.options[m.cursor]
			switch selected.action {
			case "start":
				cmd, err := (services.StartContainerBuilder{ContainerID: m.container.ID}).Build()
				if err != nil {
					m.errorMsg = err.Error()
					return m, nil
				}
				m.preview = &CommandPreviewModal{Title: "Start Container", Command: cmd}
				return m, nil
			case "stop":
				cmd, err := (services.StopContainerBuilder{ContainerID: m.container.ID}).Build()
				if err != nil {
					m.errorMsg = err.Error()
					return m, nil
				}
				m.preview = &CommandPreviewModal{Title: "Stop Container", Command: cmd}
				return m, nil
			case "logs":
				containerCopy := m.container
				return m, func() tea.Msg {
					return screenChangeMsg{target: ScreenContainerLogs, container: &containerCopy, push: true}
				}
			case "shell":
				containerCopy := m.container
				return m, func() tea.Msg {
					return screenChangeMsg{target: ScreenContainerShell, container: &containerCopy, push: true}
				}
			default:
				return m, func() tea.Msg { return screenChangeMsg{target: ScreenContainerList} }
			}
		}
	}

	return m, nil
}

func (m ContainerSubmenuScreen) View() string {
	builder := strings.Builder{}
	builder.WriteString(RenderTitle("Container Actions") + "\n\n")

	// Container Details section with bold header
	headerStyle := lipgloss.NewStyle().Bold(true)
	builder.WriteString(headerStyle.Render("Container Details") + "\n")
	builder.WriteString("name: " + m.container.Name + "\n")
	builder.WriteString("status: " + string(m.container.Status) + "\n")
	builder.WriteString("image: " + m.container.Image + "\n")
	builder.WriteString("\n")

	// Horizontal separator
	width := m.width
	if width == 0 {
		width = 80
	}
	builder.WriteString(strings.Repeat("─", width) + "\n\n")

	// Available Actions section with bold header
	builder.WriteString(headerStyle.Render("Available Actions") + "\n")

	// Action items with inverse video selection
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

	// Horizontal separator after actions
	builder.WriteString(strings.Repeat("─", width) + "\n\n")

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
	builder.WriteString("\n" + RenderMuted("Keys: up/down=navigate, enter=select, esc=back") + "\n")
	return builder.String()
}

func (m ContainerSubmenuScreen) buildOptions() []containerSubmenuOption {
	options := make([]containerSubmenuOption, 0, 4)
	if m.container.Status == models.ContainerStatusRunning {
		options = append(options, containerSubmenuOption{label: "Stop container", action: "stop"})
	} else {
		options = append(options, containerSubmenuOption{label: "Start container", action: "start"})
	}
	options = append(options, containerSubmenuOption{label: "Tail container log", action: "logs"})
	if m.container.Status == models.ContainerStatusRunning {
		options = append(options, containerSubmenuOption{label: "Enter container", action: "shell"})
	}
	options = append(options, containerSubmenuOption{label: "Back", action: "back"})
	return options
}

func (m ContainerSubmenuScreen) executeCommandCmd(command models.Command) tea.Cmd {
	return func() tea.Msg {
		result, err := m.executor.Execute(command)
		return containerSubmenuActionMsg{result: result, err: err}
	}
}
