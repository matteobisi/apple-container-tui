package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"container-tui/src/models"
	"container-tui/src/services"
)

type imageSubmenuActionMsg struct {
	result models.Result
	err    error
}

// ImageSubmenuScreen shows actions for selected image.
type ImageSubmenuScreen struct {
	executor services.CommandExecutor
	image    models.Image
	cursor   int
	errorMsg string
	confirm  *TypeToConfirmModal
	width    int
}

func NewImageSubmenuScreen(executor services.CommandExecutor) ImageSubmenuScreen {
	return ImageSubmenuScreen{executor: executor}
}

func (m ImageSubmenuScreen) SetImage(image models.Image) ImageSubmenuScreen {
	m.image = image
	m.cursor = 0
	m.errorMsg = ""
	m.confirm = nil
	return m
}

func (m ImageSubmenuScreen) Init() tea.Cmd { return nil }

func (m ImageSubmenuScreen) Update(msg tea.Msg) (ImageSubmenuScreen, tea.Cmd) {
	switch message := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = message.Width
		return m, nil
	case imageSubmenuActionMsg:
		if message.err != nil {
			m.errorMsg = services.FormatError(message.err, message.result.Stderr)
			return m, nil
		}
		return m, func() tea.Msg { return screenChangeMsg{target: ScreenImageList} }
	case tea.KeyMsg:
		if m.confirm != nil {
			updatedConfirm, confirmed, canceled := m.confirm.Handle(message)
			m.confirm = &updatedConfirm
			if confirmed {
				command := updatedConfirm.Command
				m.confirm = nil
				return m, m.executeCommandCmd(command)
			}
			if canceled {
				m.confirm = nil
				return m, nil
			}
			return m, nil
		}

		switch message.String() {
		case "up", "k":
			m.cursor = max(0, m.cursor-1)
		case "down", "j":
			m.cursor = min(2, m.cursor+1)
		case "esc":
			return m, func() tea.Msg { return screenChangeMsg{target: ScreenImageList} }
		case "enter":
			switch m.cursor {
			case 0:
				selected := m.image
				return m, func() tea.Msg { return screenChangeMsg{target: ScreenImageInspect, image: &selected, push: true} }
			case 1:
				cmd, err := (services.ImageDeleteBuilder{ImageReference: m.image.Reference()}).Build()
				if err != nil {
					m.errorMsg = err.Error()
					return m, nil
				}
				confirm := NewTypeToConfirmModal("Delete Image", "delete", cmd)
				m.confirm = &confirm
				return m, nil
			default:
				return m, func() tea.Msg { return screenChangeMsg{target: ScreenImageList} }
			}
		}
	}
	return m, nil
}

func (m ImageSubmenuScreen) View() string {
	options := []string{"Inspect image", "Delete image", "Back"}
	builder := strings.Builder{}
	builder.WriteString(RenderTitle("Image Actions") + "\n\n")

	// Image Details section with bold header
	headerStyle := lipgloss.NewStyle().Bold(true)
	builder.WriteString(headerStyle.Render("Image Details") + "\n")
	builder.WriteString("repository: " + m.image.Name + "\n")
	builder.WriteString("tag: " + m.image.Tag + "\n")
	builder.WriteString("digest: " + m.image.Digest + "\n")
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
	for i, option := range options {
		style := normalStyle
		if i == m.cursor {
			style = selectedStyle
		}
		builder.WriteString(style.Render(option) + "\n")
	}
	builder.WriteString("\n")

	// Horizontal separator after actions
	builder.WriteString(strings.Repeat("─", width) + "\n\n")

	if m.errorMsg != "" {
		builder.WriteString("\n" + RenderError("Error: "+m.errorMsg) + "\n")
	}
	if m.confirm != nil {
		builder.WriteString("\n" + m.confirm.View() + "\n")
	}
	builder.WriteString("\n" + RenderMuted("Keys: up/down=navigate, enter=select, esc=back") + "\n")
	return builder.String()
}

func (m ImageSubmenuScreen) executeCommandCmd(command models.Command) tea.Cmd {
	return func() tea.Msg {
		result, err := m.executor.Execute(command)
		return imageSubmenuActionMsg{result: result, err: err}
	}
}
