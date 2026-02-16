package ui

import "github.com/charmbracelet/bubbles/key"

// KeyMap stores key bindings used across the app.
type KeyMap struct {
	Quit   key.Binding
	Images key.Binding
}

// DefaultKeyMap returns the default key bindings.
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
		Images: key.NewBinding(
			key.WithKeys("i"),
			key.WithHelp("i", "images"),
		),
	}
}
