package ui

import (
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"

	"container-tui/src/models"
	"container-tui/src/services"
)

type buildResultMsg struct {
	result models.Result
	err    error
}

// BuildScreen builds an image from a selected file.
type BuildScreen struct {
	executor services.CommandExecutor
	filePath string
	context  string
	input    textinput.Model
	preview  *CommandPreviewModal
	loading  bool
	errorMsg string
	result   *models.Result
	viewport viewport.Model
	progress ProgressModel
	width    int
	height   int
}

// NewBuildScreen creates the build screen.
func NewBuildScreen(executor services.CommandExecutor, filePath string) BuildScreen {
	input := textinput.New()
	input.Placeholder = "my-image:latest"
	input.Prompt = "Tag: "
	input.Focus()

	viewportModel := viewport.New(0, 0)
	return BuildScreen{
		executor: executor,
		filePath: filePath,
		context:  filepath.Dir(filePath),
		input:    input,
		viewport: viewportModel,
		progress: NewProgressModel(),
	}
}

// Init starts the tag input.
func (m BuildScreen) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles build flow.
func (m BuildScreen) Update(msg tea.Msg) (BuildScreen, tea.Cmd) {
	switch message := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = message.Width
		m.height = message.Height
		m.viewport.Width = message.Width - 4
		m.viewport.Height = max(3, message.Height-12)
	case buildResultMsg:
		m.loading = false
		if message.err != nil {
			m.errorMsg = services.FormatError(message.err, message.result.Stderr)
		}
		m.result = &message.result
		m.viewport.SetContent(RenderResult(message.result))
		m.progress.SetPercent(1)
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
			return m, func() tea.Msg { return screenChangeMsg{target: ScreenContainerList} }
		case "?":
			return m, func() tea.Msg { return screenChangeMsg{target: ScreenHelp} }
		case "enter":
			if strings.TrimSpace(m.filePath) == "" {
				m.errorMsg = "no build file selected"
				return m, nil
			}
			builder := services.BuildImageBuilder{
				Tag:         strings.TrimSpace(m.input.Value()),
				FilePath:    m.filePath,
				ContextPath: m.context,
			}
			cmd, err := builder.Build()
			if err != nil {
				m.errorMsg = err.Error()
				return m, nil
			}
			m.preview = &CommandPreviewModal{Title: "Build Image", Command: cmd}
			return m, nil
		}
	}

	updatedInput, cmd := m.input.Update(msg)
	m.input = updatedInput
	return m, cmd
}

// View renders the build screen.
func (m BuildScreen) View() string {
	builder := strings.Builder{}
	builder.WriteString(RenderTitle("Build Image") + "\n\n")
	if m.filePath != "" {
		builder.WriteString(RenderMuted("File: "+m.filePath) + "\n\n")
	}
	builder.WriteString(m.input.View() + "\n")
	if m.loading {
		builder.WriteString("\n" + RenderMuted("Building...") + "\n")
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
		builder.WriteString(m.viewport.View())
	}
	builder.WriteString("\n" + RenderMuted("Keys: enter=preview, ?=help, esc=back") + "\n")
	return builder.String()
}

func (m BuildScreen) executeCommandCmd(command models.Command) tea.Cmd {
	return func() tea.Msg {
		result, err := m.executor.Execute(command)
		return buildResultMsg{result: result, err: err}
	}
}
