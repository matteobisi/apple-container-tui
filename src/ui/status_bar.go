package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// RenderStatusBar renders a single-line status bar.
func RenderStatusBar(width int, left, right string) string {
	content := strings.TrimSpace(left)
	if right != "" {
		if content != "" {
			content += " | "
		}
		content += strings.TrimSpace(right)
	}

	if width > 0 {
		return lipgloss.NewStyle().Width(width).Render(content)
	}
	return content
}
