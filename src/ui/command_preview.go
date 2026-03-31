package ui

import (
	"container-tui/src/models"
	"strings"
)

// CommandPreviewModal shows a confirmation prompt for a command.
type CommandPreviewModal struct {
	Title    string
	Command  models.Command
	Commands []models.Command
}

// View renders the command preview modal.
func (m CommandPreviewModal) View() string {
	title := m.Title
	if title == "" {
		title = "Command Preview"
	}
	commands := m.Commands
	if len(commands) == 0 {
		commands = []models.Command{m.Command}
	}
	lines := make([]string, 0, len(commands))
	for _, command := range commands {
		lines = append(lines, command.String())
	}
	return RenderTitle(title) + "\n\n" + strings.Join(lines, "\n") + "\n\n" + RenderMuted("Confirm (y/n)")
}
