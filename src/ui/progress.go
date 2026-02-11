package ui

import "github.com/charmbracelet/bubbles/progress"

// ProgressModel wraps a progress bar.
type ProgressModel struct {
	bar     progress.Model
	percent float64
}

// NewProgressModel returns a progress bar model.
func NewProgressModel() ProgressModel {
	return ProgressModel{bar: progress.New(progress.WithDefaultGradient())}
}

// SetPercent updates progress percent.
func (p *ProgressModel) SetPercent(percent float64) {
	p.percent = percent
}

// View renders the progress bar.
func (p ProgressModel) View(width int) string {
	if width <= 0 {
		return ""
	}
	p.bar.Width = width
	return p.bar.ViewAs(p.percent)
}
