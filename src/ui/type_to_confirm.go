package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"container-tui/src/models"
	"container-tui/src/services"
)

// TypeToConfirmModal requires the user to type an exact value.
type TypeToConfirmModal struct {
	Title    string
	Expected string
	Command  models.Command
	input    textinput.Model
	ErrorMsg string
}

// NewTypeToConfirmModal creates a type-to-confirm modal.
func NewTypeToConfirmModal(title, expected string, command models.Command) TypeToConfirmModal {
	input := textinput.New()
	input.Prompt = "Type to confirm: "
	input.Placeholder = expected
	input.Focus()
	return TypeToConfirmModal{Title: title, Expected: expected, Command: command, input: input}
}

// Handle processes input and returns confirmation or cancel signals.
func (m TypeToConfirmModal) Handle(msg tea.Msg) (TypeToConfirmModal, bool, bool) {
	switch message := msg.(type) {
	case tea.KeyMsg:
		if message.String() == "esc" {
			return m, false, true
		}
		if message.String() == "enter" {
			if services.IsExactMatch(m.Expected, strings.TrimSpace(m.input.Value())) {
				m.ErrorMsg = ""
				return m, true, false
			}
			m.ErrorMsg = "input does not match"
			return m, false, false
		}
	}

	updated, cmd := m.input.Update(msg)
	m.input = updated
	_ = cmd
	return m, false, false
}

// View renders the modal.
func (m TypeToConfirmModal) View() string {
	builder := strings.Builder{}
	builder.WriteString(RenderTitle(m.Title) + "\n\n")
	builder.WriteString("Command: " + m.Command.String() + "\n")
	builder.WriteString(RenderWarning(fmt.Sprintf("Confirm by typing: %s", m.Expected)) + "\n\n")
	builder.WriteString(m.input.View() + "\n")
	if m.ErrorMsg != "" {
		builder.WriteString("\n" + RenderError("Error: "+m.ErrorMsg) + "\n")
	}
	builder.WriteString("\n" + RenderMuted("Press enter to confirm, esc to cancel") + "\n")
	return builder.String()
}
