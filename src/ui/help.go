package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// HelpScreen shows keyboard shortcuts and paths.
type HelpScreen struct {
	Version string
	width   int
}

// Init returns no initial command.
func (h HelpScreen) Init() tea.Cmd {
	return nil
}

// Update handles help screen input.
func (h HelpScreen) Update(msg tea.Msg) (HelpScreen, tea.Cmd) {
	switch message := msg.(type) {
	case tea.WindowSizeMsg:
		h.width = message.Width
		return h, nil
	case tea.KeyMsg:
		switch message.String() {
		case "esc", "?":
			return h, func() tea.Msg { return screenChangeMsg{target: ScreenContainerList} }
		}
	}
	return h, nil
}

// View renders help content.
func (h HelpScreen) View() string {
	builder := strings.Builder{}
	builder.WriteString(RenderTitle("Help") + "\n\n")

	if strings.TrimSpace(h.Version) != "" {
		builder.WriteString("Version: " + h.Version + "\n\n")
	}

	headerStyle := lipgloss.NewStyle().Bold(true)
	width := h.width
	if width == 0 {
		width = 80
	}

	// Section 1: Navigation
	builder.WriteString(headerStyle.Render("Navigation") + "\n")
	builder.WriteString("up/down, j/k       Navigate lists\n")
	builder.WriteString("enter              Open submenu/select\n")
	builder.WriteString("esc                Back/cancel\n")
	builder.WriteString("\n")
	builder.WriteString(strings.Repeat("─", width) + "\n\n")

	// Section 2: Container Actions
	builder.WriteString(headerStyle.Render("Container Actions") + "\n")
	builder.WriteString("s                  Start container\n")
	builder.WriteString("t                  Stop container\n")
	builder.WriteString("d                  Delete container\n")
	builder.WriteString("enter              Open container submenu\n")
	builder.WriteString("r                  Refresh list\n")
	builder.WriteString("\n")
	builder.WriteString(strings.Repeat("─", width) + "\n\n")

	// Section 3: Image Actions
	builder.WriteString(headerStyle.Render("Image Actions") + "\n")
	builder.WriteString("i                  Switch to images view\n")
	builder.WriteString("p                  Pull image\n")
	builder.WriteString("b                  Build image\n")
	builder.WriteString("n                  Prune images\n")
	builder.WriteString("enter              Open image submenu\n")
	builder.WriteString("\n")
	builder.WriteString(strings.Repeat("─", width) + "\n\n")

	// Section 4: General
	builder.WriteString(headerStyle.Render("General") + "\n")
	builder.WriteString("m                  Manage daemon\n")
	builder.WriteString("?                  Show this help\n")
	builder.WriteString("q                  Quit application\n")
	builder.WriteString("\n")
	builder.WriteString("Config: ~/.config/actui/config OR ~/Library/Application Support/actui/config\n")
	builder.WriteString("Logs:   ~/Library/Application Support/actui/command.log\n")
	builder.WriteString("\n")
	builder.WriteString(strings.Repeat("─", width) + "\n\n")

	builder.WriteString(RenderMuted("Press any key to return") + "\n")
	return builder.String()
}
