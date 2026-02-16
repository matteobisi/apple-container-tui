package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/filepicker"
	tea "github.com/charmbracelet/bubbletea"

	"container-tui/src/services"
)

// FilePickerScreen allows selecting a build file.
type FilePickerScreen struct {
	executor     services.CommandExecutor
	picker       filepicker.Model
	returnTarget ActiveScreen
	errorMsg     string
	width        int
	height       int
}

// NewFilePickerScreen creates a file picker screen.
func NewFilePickerScreen(executor services.CommandExecutor) FilePickerScreen {
	picker := filepicker.New()
	picker.AllowedTypes = []string{"Containerfile", "Dockerfile"}
	picker.CurrentDirectory = "."
	return FilePickerScreen{executor: executor, picker: picker, returnTarget: ScreenContainerList}
}

// SetReturnTarget sets the screen to return to after flow completion/cancel.
func (m FilePickerScreen) SetReturnTarget(target ActiveScreen) FilePickerScreen {
	m.returnTarget = target
	return m
}

// Init starts the file picker.
func (m FilePickerScreen) Init() tea.Cmd {
	return m.picker.Init()
}

// Update handles file picker events.
func (m FilePickerScreen) Update(msg tea.Msg) (FilePickerScreen, tea.Cmd) {
	switch message := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = message.Width
		m.height = message.Height
		m.picker.Height = message.Height - 6
	case tea.KeyMsg:
		switch message.String() {
		case "esc":
			return m, func() tea.Msg { return BackToListMsg{} }
		case "?":
			return m, func() tea.Msg { return screenChangeMsg{target: ScreenHelp} }
		}
	}

	updatedPicker, cmd := m.picker.Update(msg)
	m.picker = updatedPicker
	if didSelect, path := m.picker.DidSelectFile(msg); didSelect {
		returnTarget := m.returnTarget
		return m, func() tea.Msg { return buildFileSelectedMsg{path: path, returnTarget: returnTarget} }
	}
	if didSelect, path := m.picker.DidSelectDisabledFile(msg); didSelect {
		m.errorMsg = "unsupported file: " + path
	}
	return m, cmd
}

// View renders the file picker screen.
func (m FilePickerScreen) View() string {
	builder := strings.Builder{}
	builder.WriteString(RenderTitle("Select Build File") + "\n\n")
	builder.WriteString(m.picker.View())
	if m.errorMsg != "" {
		builder.WriteString("\n" + RenderError("Error: "+m.errorMsg) + "\n")
	}
	builder.WriteString("\n" + RenderMuted("Keys: enter=select, ?=help, esc=back") + "\n")
	return builder.String()
}
