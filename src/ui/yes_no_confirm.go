package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"container-tui/src/models"
)

// YesNoConfirmModal asks for a yes/no confirmation.
type YesNoConfirmModal struct {
	Title   string
	Body    string
	Command models.Command
	Warning bool
}

// Handle processes key input and returns confirmation or cancel signals.
func (m YesNoConfirmModal) Handle(msg tea.Msg) (bool, bool) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch strings.ToLower(keyMsg.String()) {
		case "y", "enter":
			return true, false
		case "n", "esc":
			return false, true
		}
	}
	return false, false
}

// View renders the modal.
func (m YesNoConfirmModal) View() string {
	builder := strings.Builder{}
	builder.WriteString(RenderTitle(m.Title) + "\n\n")
	if m.Body != "" {
		builder.WriteString(m.Body + "\n\n")
	}
	builder.WriteString("Command: " + m.Command.String() + "\n")
	builder.WriteString(RenderMuted("Confirm (y/n)") + "\n")
	if m.Warning {
		builder.WriteString("\n" + RenderWarning("Warning: destructive action!") + "\n")
	}
	return builder.String()
}
