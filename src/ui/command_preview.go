package ui

import "container-tui/src/models"

// CommandPreviewModal shows a confirmation prompt for a command.
type CommandPreviewModal struct {
	Title   string
	Command models.Command
}

// View renders the command preview modal.
func (m CommandPreviewModal) View() string {
	title := m.Title
	if title == "" {
		title = "Command Preview"
	}
	return RenderTitle(title) + "\n\n" + m.Command.String() + "\n\n" + RenderMuted("Confirm (y/n)")
}
