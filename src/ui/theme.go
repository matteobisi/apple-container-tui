package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Theme defines UI styles.
type Theme struct {
	Title   lipgloss.Style
	Error   lipgloss.Style
	Warning lipgloss.Style
	Success lipgloss.Style
	Muted   lipgloss.Style
	Accent  lipgloss.Style
	Key     lipgloss.Style
}

var currentTheme = autoTheme()

// ApplyTheme sets the active theme mode.
func ApplyTheme(mode string) {
	currentTheme = ResolveTheme(mode)
}

// CurrentTheme returns the active theme.
func CurrentTheme() Theme {
	return currentTheme
}

// ResolveTheme maps a mode string to a theme.
func ResolveTheme(mode string) Theme {
	switch strings.ToLower(strings.TrimSpace(mode)) {
	case "light":
		return lightTheme()
	case "dark":
		return darkTheme()
	default:
		return autoTheme()
	}
}

func lightTheme() Theme {
	return Theme{
		Title:   lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("235")),
		Error:   lipgloss.NewStyle().Foreground(lipgloss.Color("160")),
		Warning: lipgloss.NewStyle().Foreground(lipgloss.Color("166")),
		Success: lipgloss.NewStyle().Foreground(lipgloss.Color("28")),
		Muted:   lipgloss.NewStyle().Foreground(lipgloss.Color("245")),
		Accent:  lipgloss.NewStyle().Foreground(lipgloss.Color("27")),
		Key:     lipgloss.NewStyle().Foreground(lipgloss.Color("62")),
	}
}

func darkTheme() Theme {
	return Theme{
		Title:   lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("255")),
		Error:   lipgloss.NewStyle().Foreground(lipgloss.Color("203")),
		Warning: lipgloss.NewStyle().Foreground(lipgloss.Color("214")),
		Success: lipgloss.NewStyle().Foreground(lipgloss.Color("77")),
		Muted:   lipgloss.NewStyle().Foreground(lipgloss.Color("244")),
		Accent:  lipgloss.NewStyle().Foreground(lipgloss.Color("81")),
		Key:     lipgloss.NewStyle().Foreground(lipgloss.Color("109")),
	}
}

func autoTheme() Theme {
	return Theme{
		Title:   lipgloss.NewStyle().Bold(true).Foreground(lipgloss.AdaptiveColor{Light: "235", Dark: "255"}),
		Error:   lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "160", Dark: "203"}),
		Warning: lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "166", Dark: "214"}),
		Success: lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "28", Dark: "77"}),
		Muted:   lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "245", Dark: "244"}),
		Accent:  lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "27", Dark: "81"}),
		Key:     lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "62", Dark: "109"}),
	}
}

// RenderTitle styles a title string.
func RenderTitle(text string) string {
	return currentTheme.Title.Render(text)
}

// RenderError styles an error string.
func RenderError(text string) string {
	return currentTheme.Error.Render(text)
}

// RenderWarning styles a warning string.
func RenderWarning(text string) string {
	return currentTheme.Warning.Render(text)
}

// RenderSuccess styles a success string.
func RenderSuccess(text string) string {
	return currentTheme.Success.Render(text)
}

// RenderMuted styles muted or secondary text.
func RenderMuted(text string) string {
	return currentTheme.Muted.Render(text)
}

// RenderAccent styles highlighted text.
func RenderAccent(text string) string {
	return currentTheme.Accent.Render(text)
}
