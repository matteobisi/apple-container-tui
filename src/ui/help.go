package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// HelpScreen shows keyboard shortcuts and paths.
type HelpScreen struct {
	Version string
}

// Init returns no initial command.
func (h HelpScreen) Init() tea.Cmd {
	return nil
}

// Update handles help screen input.
func (h HelpScreen) Update(msg tea.Msg) (HelpScreen, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
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
	builder.WriteString("Navigation:\n")
	builder.WriteString("  up/down, j/k  - move selection\n")
	builder.WriteString("  enter         - toggle start/stop\n")
	builder.WriteString("  q             - quit\n")
	builder.WriteString("\nContainers:\n")
	builder.WriteString("  s             - start container\n")
	builder.WriteString("  t             - stop container\n")
	builder.WriteString("  d             - delete container (type-to-confirm)\n")
	builder.WriteString("  r             - refresh list\n")
	builder.WriteString("\nImages:\n")
	builder.WriteString("  p             - pull image\n")
	builder.WriteString("  b             - build image\n")
	builder.WriteString("\nDaemon:\n")
	builder.WriteString("  m             - manage daemon\n")
	builder.WriteString("\nPaths:\n")
	builder.WriteString("  Config: ~/.config/actui/config OR ~/Library/Application Support/actui/config\n")
	builder.WriteString("  Logs:   ~/Library/Application Support/actui/command.log\n")
	builder.WriteString("\n" + RenderMuted("Press ? or esc to return") + "\n")
	return builder.String()
}
