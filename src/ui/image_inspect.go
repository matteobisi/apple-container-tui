package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"

	"container-tui/src/models"
	"container-tui/src/services"
)

type imageInspectLoadedMsg struct {
	content string
	err     error
}

// ImageInspectScreen shows image inspect output.
type ImageInspectScreen struct {
	executor services.CommandExecutor
	image    models.Image
	viewport viewport.Model
	loading  bool
	errorMsg string
	width    int
	height   int
}

func NewImageInspectScreen(executor services.CommandExecutor) ImageInspectScreen {
	return ImageInspectScreen{executor: executor, viewport: viewport.New(0, 0)}
}

func (m ImageInspectScreen) SetImage(image models.Image) ImageInspectScreen {
	m.image = image
	m.loading = true
	m.errorMsg = ""
	m.viewport.SetContent("")
	return m
}

func (m ImageInspectScreen) Init() tea.Cmd {
	if strings.TrimSpace(m.image.Name) == "" {
		return nil
	}
	return m.loadInspectCmd()
}

func (m ImageInspectScreen) Update(msg tea.Msg) (ImageInspectScreen, tea.Cmd) {
	switch message := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = message.Width
		m.height = message.Height
		m.viewport.Width = max(20, message.Width-4)
		m.viewport.Height = max(5, message.Height-8)
	case imageInspectLoadedMsg:
		m.loading = false
		if message.err != nil {
			m.errorMsg = message.err.Error()
			return m, nil
		}
		m.viewport.SetContent(message.content)
	case tea.KeyMsg:
		switch message.String() {
		case "up", "down", "pgup", "pgdown", "home", "end":
			updatedViewport, _ := m.viewport.Update(message)
			m.viewport = updatedViewport
			return m, nil
		case "esc":
			selected := m.image
			return m, func() tea.Msg { return screenChangeMsg{target: ScreenImageSubmenu, image: &selected} }
		}
	}
	return m, nil
}

func (m ImageInspectScreen) View() string {
	builder := strings.Builder{}
	builder.WriteString(RenderTitle("Image Inspection") + "\n\n")
	builder.WriteString(RenderMuted("Image: "+m.image.Reference()) + "\n\n")
	if m.loading {
		builder.WriteString(RenderMuted("Loading inspection...") + "\n")
	}
	if m.errorMsg != "" {
		builder.WriteString(RenderError("Error: "+m.errorMsg) + "\n")
	}
	builder.WriteString(m.viewport.View() + "\n")
	builder.WriteString("\n" + RenderMuted("Keys: up/down/pgup/pgdn/home/end=scroll, esc=back") + "\n")
	return builder.String()
}

func (m ImageInspectScreen) loadInspectCmd() tea.Cmd {
	return func() tea.Msg {
		cmd, err := (services.ImageInspectBuilder{ImageReference: m.image.Reference()}).Build()
		if err != nil {
			return imageInspectLoadedMsg{err: err}
		}
		result, err := m.executor.Execute(cmd)
		if err != nil {
			return imageInspectLoadedMsg{err: err}
		}
		content := strings.TrimSpace(result.Stdout)
		if content == "" {
			content = result.Stderr
		}
		if strings.TrimSpace(content) == "" {
			content = "No inspection output."
		}
		return imageInspectLoadedMsg{content: content}
	}
}
