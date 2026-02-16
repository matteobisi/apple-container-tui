package ui

import (
	"errors"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"container-tui/src/models"
	"container-tui/src/services"
)

type containerListLoadedMsg struct {
	containers []models.Container
	err        error
}

type commandExecutedMsg struct {
	result models.Result
	err    error
}

// ContainerListScreen displays containers and actions.
type ContainerListScreen struct {
	executor   services.CommandExecutor
	containers []models.Container
	cursor     int
	loading    bool
	errorMsg   string
	result     *models.Result
	preview    *CommandPreviewModal
	confirm    *TypeToConfirmModal
	pendingCmd *models.Command
	hasLoaded  bool
}

// NewContainerListScreen creates the container list screen.
func NewContainerListScreen(executor services.CommandExecutor) ContainerListScreen {
	return ContainerListScreen{executor: executor}
}

// Init fetches the initial container list.
func (m ContainerListScreen) Init() tea.Cmd {
	return m.fetchContainersCmd(false)
}

// Update handles screen messages.
func (m ContainerListScreen) Update(msg tea.Msg) (ContainerListScreen, tea.Cmd) {
	switch message := msg.(type) {
	case containerListLoadedMsg:
		m.loading = false
		if message.err != nil {
			m.errorMsg = services.FormatError(message.err, "")
			return m, nil
		}
		m.errorMsg = ""
		m.containers = message.containers
		m.hasLoaded = true
		if m.cursor >= len(m.containers) {
			m.cursor = max(0, len(m.containers)-1)
		}
		return m, nil
	case commandExecutedMsg:
		m.loading = false
		m.pendingCmd = nil
		if message.err != nil {
			m.errorMsg = services.FormatError(message.err, message.result.Stderr)
		}
		m.result = &message.result
		return m, m.fetchContainersCmd(true)
	case tea.KeyMsg:
		if m.confirm != nil {
			updatedConfirm, confirmed, canceled := m.confirm.Handle(message)
			m.confirm = &updatedConfirm
			if confirmed {
				command := updatedConfirm.Command
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
		if m.preview != nil {
			switch strings.ToLower(message.String()) {
			case "y", "enter":
				previewed := m.preview.Command
				m.preview = nil
				m.pendingCmd = &previewed
				m.loading = true
				return m, m.executeCommandCmd(previewed)
			case "n", "esc":
				m.preview = nil
				m.pendingCmd = nil
				return m, nil
			}
			return m, nil
		}

		switch message.String() {
		case "up", "k":
			m.cursor = max(0, m.cursor-1)
		case "down", "j":
			m.cursor = min(len(m.containers)-1, m.cursor+1)
		case "r":
			m.loading = true
			return m, m.fetchContainersCmd(true)
		case "i":
			return m, func() tea.Msg { return screenChangeMsg{target: ScreenImageList, push: true} }
		case "m":
			return m, func() tea.Msg { return screenChangeMsg{target: ScreenDaemonControl} }
		case "?":
			return m, func() tea.Msg { return screenChangeMsg{target: ScreenHelp} }
		case "d":
			updated, cmd := m.buildAndConfirmDelete()
			return updated, cmd
		case "s":
			updated, cmd := m.buildAndPreviewStart()
			return updated, cmd
		case "t":
			updated, cmd := m.buildAndPreviewStop()
			return updated, cmd
		case "enter":
			selected, ok := m.selectedContainer()
			if !ok {
				return m, nil
			}
			containerCopy := selected
			return m, func() tea.Msg {
				return screenChangeMsg{target: ScreenContainerSubmenu, container: &containerCopy, push: true}
			}
		}
	}

	return m, nil
}

// View renders the screen content.
func (m ContainerListScreen) View() string {
	builder := strings.Builder{}
	builder.WriteString(RenderTitle("Containers") + "\n\n")
	if m.loading {
		builder.WriteString(RenderMuted("Loading...") + "\n")
	}
	if m.errorMsg != "" {
		builder.WriteString(RenderError("Error: "+m.errorMsg) + "\n\n")
	}

	if len(m.containers) == 0 {
		builder.WriteString(RenderMuted("No containers found.") + "\n")
	} else {
		for i, container := range m.containers {
			cursor := " "
			if i == m.cursor {
				cursor = ">"
			}
			line := fmt.Sprintf("%s %s [%s] %s", cursor, container.Name, container.Status, container.Image)
			if i == m.cursor {
				line = RenderAccent(line)
			}
			builder.WriteString(line + "\n")
		}
	}

	builder.WriteString("\n" + RenderMuted("Keys: up/down, enter=submenu, s=start, t=stop, d=delete(!), i=images, r=refresh, m=manage, ?=help, q=quit") + "\n")

	if m.preview != nil {
		builder.WriteString("\n")
		builder.WriteString(m.preview.View())
	}

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

func (m ContainerListScreen) buildAndPreviewStart() (ContainerListScreen, tea.Cmd) {
	selected, ok := m.selectedContainer()
	if !ok {
		return m, nil
	}
	builder := services.StartContainerBuilder{ContainerID: selected.ID}
	cmd, err := builder.Build()
	if err != nil {
		m.errorMsg = err.Error()
		return m, nil
	}
	m.preview = &CommandPreviewModal{Title: "Start Container", Command: cmd}
	return m, nil
}

func (m ContainerListScreen) buildAndPreviewStop() (ContainerListScreen, tea.Cmd) {
	selected, ok := m.selectedContainer()
	if !ok {
		return m, nil
	}
	builder := services.StopContainerBuilder{ContainerID: selected.ID}
	cmd, err := builder.Build()
	if err != nil {
		m.errorMsg = err.Error()
		return m, nil
	}
	m.preview = &CommandPreviewModal{Title: "Stop Container", Command: cmd}
	return m, nil
}

func (m ContainerListScreen) buildAndPreviewToggle() (ContainerListScreen, tea.Cmd) {
	selected, ok := m.selectedContainer()
	if !ok {
		return m, nil
	}
	if selected.Status == models.ContainerStatusRunning {
		return m.buildAndPreviewStop()
	}
	return m.buildAndPreviewStart()
}

func (m ContainerListScreen) buildAndConfirmDelete() (ContainerListScreen, tea.Cmd) {
	selected, ok := m.selectedContainer()
	if !ok {
		return m, nil
	}
	if selected.Status != models.ContainerStatusStopped {
		m.errorMsg = "container must be stopped to delete"
		return m, nil
	}

	builder := services.DeleteContainerBuilder{ContainerID: selected.ID}
	cmd, err := builder.Build()
	if err != nil {
		m.errorMsg = err.Error()
		return m, nil
	}

	expected := selected.Name
	if strings.TrimSpace(expected) == "" {
		expected = selected.ID
	}

	confirm := NewTypeToConfirmModal("Delete Container", expected, cmd)
	m.confirm = &confirm
	return m, nil
}

func (m ContainerListScreen) fetchContainersCmd(force bool) tea.Cmd {
	if m.hasLoaded && !force {
		return nil
	}
	return func() tea.Msg {
		builder := services.ListContainersBuilder{}
		cmd, err := builder.Build()
		if err != nil {
			return containerListLoadedMsg{err: err}
		}
		result, err := m.executor.Execute(cmd)
		if err != nil {
			return containerListLoadedMsg{err: errors.New(services.FormatError(err, result.Stderr))}
		}
		containers, parseErr := services.ParseContainerList(result.Stdout)
		if parseErr != nil {
			return containerListLoadedMsg{err: parseErr}
		}
		return containerListLoadedMsg{containers: containers}
	}
}

func (m ContainerListScreen) executeCommandCmd(command models.Command) tea.Cmd {
	return func() tea.Msg {
		result, err := m.executor.Execute(command)
		return commandExecutedMsg{result: result, err: err}
	}
}

func (m *ContainerListScreen) selectedContainer() (models.Container, bool) {
	if len(m.containers) == 0 || m.cursor < 0 || m.cursor >= len(m.containers) {
		m.errorMsg = "no container selected"
		return models.Container{}, false
	}
	return m.containers[m.cursor], true
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
