package ui

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

// SpinnerModel wraps a Bubbles spinner.
type SpinnerModel struct {
	model  spinner.Model
	active bool
}

// NewSpinnerModel creates a new spinner model.
func NewSpinnerModel() SpinnerModel {
	sp := spinner.New()
	sp.Spinner = spinner.Line
	return SpinnerModel{model: sp}
}

// SetActive enables or disables the spinner.
func (s *SpinnerModel) SetActive(active bool) {
	s.active = active
}

// Update advances the spinner animation.
func (s SpinnerModel) Update(msg tea.Msg) (SpinnerModel, tea.Cmd) {
	if !s.active {
		return s, nil
	}
	updated, cmd := s.model.Update(msg)
	s.model = updated
	return s, cmd
}

// View renders the spinner.
func (s SpinnerModel) View() string {
	if !s.active {
		return ""
	}
	return s.model.View()
}
