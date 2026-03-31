package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"container-tui/src/models"
	"container-tui/src/services"
)

type containerExportResultMsg struct {
	result services.ExportWorkflowResult
	err    error
}

type containerExportCleanupMsg struct {
	result models.Result
	err    error
}

// ContainerExportScreen collects a destination directory and runs the export workflow.
type ContainerExportScreen struct {
	executor  services.CommandExecutor
	container models.Container
	input     textinput.Model
	plan      *services.ContainerExportPlan
	preview   *CommandPreviewModal
	confirm   *YesNoConfirmModal
	loading   bool
	errorMsg  string
	result    *models.Result
	progress  ProgressModel
	width     int
}

func NewContainerExportScreen(executor services.CommandExecutor) ContainerExportScreen {
	input := textinput.New()
	input.Placeholder = "."
	input.Prompt = "Destination directory: "
	input.Focus()
	return ContainerExportScreen{executor: executor, input: input, progress: NewProgressModel()}
}

func (m ContainerExportScreen) SetContainer(container models.Container) ContainerExportScreen {
	m.container = container
	m.errorMsg = ""
	m.result = nil
	m.preview = nil
	m.confirm = nil
	m.plan = nil
	m.loading = false
	m.progress.SetPercent(0)
	return m
}

func (m ContainerExportScreen) Init() tea.Cmd {
	return textinput.Blink
}

func (m ContainerExportScreen) Update(msg tea.Msg) (ContainerExportScreen, tea.Cmd) {
	switch message := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = message.Width
		return m, nil
	case containerExportResultMsg:
		m.loading = false
		result := message.result.Result
		m.result = &result
		m.progress.SetPercent(1)
		if message.err != nil {
			m.errorMsg = services.FormatError(message.err, result.Stderr)
			return m, nil
		}
		m.errorMsg = ""
		if m.plan != nil && m.plan.CleanupCommand.Executable != "" {
			m.confirm = &YesNoConfirmModal{
				Title:   "Delete Temporary Export Image",
				Body:    fmt.Sprintf("Archive saved to %s\n\nDelete the temporary image %s now?", message.result.ArchivePath, m.plan.GeneratedImageRef),
				Command: m.plan.CleanupCommand,
				Warning: true,
			}
		}
		return m, nil
	case containerExportCleanupMsg:
		m.loading = false
		if m.result == nil {
			result := models.Result{Status: models.ResultSuccess}
			m.result = &result
		}
		if message.err != nil {
			m.result.Status = models.ResultError
			m.result.Stderr = strings.TrimSpace(strings.Join([]string{m.result.Stderr, services.FormatError(message.err, message.result.Stderr)}, "\n\n"))
			m.errorMsg = services.FormatError(message.err, message.result.Stderr)
			return m, nil
		}
		m.errorMsg = ""
		cleanupMessage := fmt.Sprintf("Temporary export image deleted: %s", m.plan.GeneratedImageRef)
		cleanupStdout := strings.TrimSpace(message.result.Stdout)
		if cleanupStdout != "" {
			cleanupMessage = cleanupMessage + "\n\n[cleanup]\n" + cleanupStdout
		}
		m.result.Stdout = strings.TrimSpace(strings.Join([]string{m.result.Stdout, cleanupMessage}, "\n\n"))
		return m, nil
	case tea.KeyMsg:
		if m.confirm != nil {
			confirmed, canceled := m.confirm.Handle(message)
			if confirmed {
				command := m.confirm.Command
				m.confirm = nil
				m.loading = true
				return m, m.executeCleanupCmd(command)
			}
			if canceled {
				if m.result != nil && m.plan != nil {
					retained := fmt.Sprintf("Temporary export image retained: %s", m.plan.GeneratedImageRef)
					m.result.Stdout = strings.TrimSpace(strings.Join([]string{m.result.Stdout, retained}, "\n\n"))
				}
				m.confirm = nil
				return m, nil
			}
			return m, nil
		}

		if m.preview != nil {
			switch strings.ToLower(message.String()) {
			case "y", "enter":
				m.preview = nil
				m.loading = true
				m.progress.SetPercent(0)
				if m.plan == nil {
					m.errorMsg = "export plan is missing"
					m.loading = false
					return m, nil
				}
				return m, m.executeExportCmd(*m.plan)
			case "n", "esc":
				m.preview = nil
				return m, nil
			}
			return m, nil
		}

		switch message.String() {
		case "esc":
			return m, func() tea.Msg { return BackToSubmenuMsg{} }
		case "?":
			return m, func() tea.Msg { return screenChangeMsg{target: ScreenHelp} }
		case "enter":
			workflow := services.NewExportWorkflowService(m.executor)
			plan, err := workflow.Plan(m.container, strings.TrimSpace(m.input.Value()))
			if err != nil {
				m.errorMsg = err.Error()
				return m, nil
			}
			m.plan = &plan
			m.preview = &CommandPreviewModal{Title: "Export Container", Commands: plan.Commands}
			return m, nil
		}
	}

	updatedInput, cmd := m.input.Update(msg)
	m.input = updatedInput
	return m, cmd
}

func (m ContainerExportScreen) View() string {
	builder := strings.Builder{}
	builder.WriteString(RenderTitle("Export Container") + "\n\n")
	builder.WriteString(RenderMuted("Container: "+m.container.Name+" ("+m.container.ID+")") + "\n\n")
	builder.WriteString(m.input.View() + "\n")
	if m.plan != nil {
		builder.WriteString(RenderMuted("Archive: "+m.plan.ArchivePath) + "\n")
	}
	if m.loading {
		builder.WriteString("\n" + RenderMuted("Exporting container...") + "\n")
	}
	if m.width > 0 {
		builder.WriteString(m.progress.View(m.width-4) + "\n")
	}
	if m.errorMsg != "" {
		builder.WriteString("\n" + RenderError("Error: "+m.errorMsg) + "\n")
	}
	if m.preview != nil {
		builder.WriteString("\n" + m.preview.View())
	}
	if m.confirm != nil {
		builder.WriteString("\n" + m.confirm.View())
	}
	if m.result != nil {
		builder.WriteString("\n\n" + RenderResult(*m.result))
	}
	builder.WriteString("\n" + RenderMuted("Keys: enter=preview, ?=help, esc=back") + "\n")
	return builder.String()
}

func (m ContainerExportScreen) executeExportCmd(plan services.ContainerExportPlan) tea.Cmd {
	return func() tea.Msg {
		workflow := services.NewExportWorkflowService(m.executor)
		result, err := workflow.Execute(plan)
		return containerExportResultMsg{result: result, err: err}
	}
}

func (m ContainerExportScreen) executeCleanupCmd(command models.Command) tea.Cmd {
	return func() tea.Msg {
		result, err := m.executor.Execute(command)
		return containerExportCleanupMsg{result: result, err: err}
	}
}
