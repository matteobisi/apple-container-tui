package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"container-tui/src/models"
	"container-tui/src/services"
)

type imagePullResultMsg struct {
	result models.Result
	err    error
}

// ImagePullScreen collects image reference and runs pull.
type ImagePullScreen struct {
	executor     services.CommandExecutor
	returnTarget ActiveScreen
	input        textinput.Model
	preview      *CommandPreviewModal
	loading      bool
	errorMsg     string
	result       *models.Result
	progress     ProgressModel
	width        int
}

// NewImagePullScreen creates the pull screen.
func NewImagePullScreen(executor services.CommandExecutor) ImagePullScreen {
	input := textinput.New()
	input.Placeholder = "nginx:latest"
	input.Prompt = "Image reference: "
	input.Focus()
	return ImagePullScreen{
		executor:     executor,
		returnTarget: ScreenContainerList,
		input:        input,
		progress:     NewProgressModel(),
	}
}

// SetReturnTarget sets the screen to return to for Esc and post-success flow.
func (m ImagePullScreen) SetReturnTarget(target ActiveScreen) ImagePullScreen {
	m.returnTarget = target
	return m
}

// Init starts the text input.
func (m ImagePullScreen) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles input and command execution.
func (m ImagePullScreen) Update(msg tea.Msg) (ImagePullScreen, tea.Cmd) {
	switch message := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = message.Width
	case imagePullResultMsg:
		m.loading = false
		if message.err != nil {
			m.errorMsg = services.FormatError(message.err, message.result.Stderr)
			m.result = &message.result
			m.progress.SetPercent(1)
			return m, nil
		}
		m.result = &message.result
		m.progress.SetPercent(1)
		if m.returnTarget != ScreenContainerList {
			target := m.returnTarget
			return m, func() tea.Msg { return screenChangeMsg{target: target} }
		}
		return m, nil
	case tea.KeyMsg:
		if m.preview != nil {
			switch strings.ToLower(message.String()) {
			case "y", "enter":
				previewed := m.preview.Command
				m.preview = nil
				m.loading = true
				m.progress.SetPercent(0)
				return m, m.executeCommandCmd(previewed)
			case "n", "esc":
				m.preview = nil
				return m, nil
			}
			return m, nil
		}

		switch message.String() {
		case "esc":
			return m, func() tea.Msg { return BackToListMsg{} }
		case "?":
			return m, func() tea.Msg { return screenChangeMsg{target: ScreenHelp} }
		case "enter":
			builder := services.PullImageBuilder{Reference: strings.TrimSpace(m.input.Value())}
			cmd, err := builder.Build()
			if err != nil {
				m.errorMsg = err.Error()
				return m, nil
			}
			m.preview = &CommandPreviewModal{Title: "Pull Image", Command: cmd}
			return m, nil
		}
	}

	updatedInput, cmd := m.input.Update(msg)
	m.input = updatedInput
	return m, cmd
}

// View renders the pull screen.
func (m ImagePullScreen) View() string {
	builder := strings.Builder{}
	builder.WriteString(RenderTitle("Pull Image") + "\n\n")
	builder.WriteString(m.input.View() + "\n")
	if m.loading {
		builder.WriteString("\n" + RenderMuted("Pulling...") + "\n")
	}
	if m.width > 0 {
		builder.WriteString(m.progress.View(m.width-4) + "\n")
	}
	if m.errorMsg != "" {
		builder.WriteString("\n" + RenderError("Error: "+m.errorMsg) + "\n")
	}
	if m.preview != nil {
		builder.WriteString("\n")
		builder.WriteString(m.preview.View())
	}
	if m.result != nil {
		builder.WriteString("\n\n")
		builder.WriteString(RenderResult(*m.result))
	}
	builder.WriteString("\n" + RenderMuted("Keys: enter=preview, ?=help, esc=back") + "\n")
	return builder.String()
}

func (m ImagePullScreen) executeCommandCmd(command models.Command) tea.Cmd {
	return func() tea.Msg {
		result, err := m.executor.Execute(command)
		return imagePullResultMsg{result: result, err: err}
	}
}
