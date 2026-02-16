package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"container-tui/src/models"
	"container-tui/src/services"
)

type imageListLoadedMsg struct {
	images []models.Image
	err    error
}

type imageListActionMsg struct {
	result models.Result
	err    error
}

// ImageListScreen shows local images and image actions.
type ImageListScreen struct {
	executor services.CommandExecutor
	images   []models.Image
	cursor   int
	loading  bool
	errorMsg string
	result   *models.Result
	confirm  *TypeToConfirmModal
}

func NewImageListScreen(executor services.CommandExecutor) ImageListScreen {
	return ImageListScreen{executor: executor, images: []models.Image{}}
}

func (m ImageListScreen) Init() tea.Cmd {
	m.loading = true
	return m.fetchImagesCmd()
}

func (m ImageListScreen) Update(msg tea.Msg) (ImageListScreen, tea.Cmd) {
	switch message := msg.(type) {
	case imageListLoadedMsg:
		m.loading = false
		if message.err != nil {
			m.errorMsg = message.err.Error()
			return m, nil
		}
		m.images = message.images
		if m.cursor >= len(m.images) {
			m.cursor = max(0, len(m.images)-1)
		}
		if len(m.images) == 0 {
			m.errorMsg = "No images found"
		} else {
			m.errorMsg = ""
		}
		return m, nil
	case imageListActionMsg:
		m.loading = false
		m.result = &message.result
		if message.err != nil {
			m.errorMsg = services.FormatError(message.err, message.result.Stderr)
		}
		return m, m.fetchImagesCmd()
	case tea.KeyMsg:
		if m.confirm != nil {
			updatedConfirm, confirmed, canceled := m.confirm.Handle(message)
			m.confirm = &updatedConfirm
			if confirmed {
				command := updatedConfirm.Command
				m.confirm = nil
				m.loading = true
				return m, m.executeCommandCmd(command)
			}
			if canceled {
				m.confirm = nil
				return m, nil
			}
			return m, nil
		}

		switch message.String() {
		case "up", "k":
			m.cursor = max(0, m.cursor-1)
		case "down", "j":
			m.cursor = min(len(m.images)-1, m.cursor+1)
		case "r":
			m.loading = true
			return m, m.fetchImagesCmd()
		case "p":
			return m, func() tea.Msg { return screenChangeMsg{target: ScreenImagePull, push: true} }
		case "b":
			return m, func() tea.Msg { return screenChangeMsg{target: ScreenFilePicker, push: true} }
		case "n":
			cmd, err := (services.ImagePruneBuilder{}).Build()
			if err != nil {
				m.errorMsg = err.Error()
				return m, nil
			}
			confirm := NewTypeToConfirmModal("Prune Images", "prune", cmd)
			m.confirm = &confirm
			return m, nil
		case "esc":
			return m, func() tea.Msg { return screenChangeMsg{target: ScreenContainerList} }
		case "enter":
			if len(m.images) == 0 || m.cursor < 0 || m.cursor >= len(m.images) {
				return m, nil
			}
			selected := m.images[m.cursor]
			return m, func() tea.Msg { return screenChangeMsg{target: ScreenImageSubmenu, image: &selected, push: true} }
		}
	}
	return m, nil
}

func (m ImageListScreen) View() string {
	builder := strings.Builder{}
	builder.WriteString(RenderTitle("Images") + "\n\n")
	if m.loading {
		builder.WriteString(RenderMuted("Loading images...") + "\n")
	}
	if len(m.images) == 0 {
		builder.WriteString(RenderMuted("No images found. Press 'p' to pull an image or 'b' to build from Containerfile.") + "\n")
	} else {
		builder.WriteString("NAME\tTAG\tDIGEST\n")
		for i, image := range m.images {
			cursor := " "
			if i == m.cursor {
				cursor = ">"
			}
			digest := image.Digest
			if len(digest) > 24 {
				digest = digest[:24] + "..."
			}
			line := cursor + " " + image.Name + "\t" + image.Tag + "\t" + digest
			if i == m.cursor {
				line = RenderAccent(line)
			}
			builder.WriteString(line + "\n")
		}
	}
	if m.errorMsg != "" {
		builder.WriteString("\n" + RenderError("Error: "+m.errorMsg) + "\n")
	}
	if m.result != nil {
		builder.WriteString("\n" + RenderResult(*m.result) + "\n")
	}
	if m.confirm != nil {
		builder.WriteString("\n" + m.confirm.View() + "\n")
	}
	builder.WriteString("\n" + RenderMuted("Keys: up/down=navigate, enter=submenu, p=pull, b=build, n=image-prune, r=refresh, esc=back") + "\n")
	return builder.String()
}

func (m ImageListScreen) fetchImagesCmd() tea.Cmd {
	return func() tea.Msg {
		cmd, err := (services.ImageListBuilder{}).Build()
		if err != nil {
			return imageListLoadedMsg{err: err}
		}
		result, err := m.executor.Execute(cmd)
		if err != nil {
			return imageListLoadedMsg{err: err}
		}
		images, parseErr := services.ParseImageList(result.Stdout)
		if parseErr != nil {
			return imageListLoadedMsg{err: parseErr}
		}
		return imageListLoadedMsg{images: images}
	}
}

func (m ImageListScreen) executeCommandCmd(command models.Command) tea.Cmd {
	return func() tea.Msg {
		result, err := m.executor.Execute(command)
		return imageListActionMsg{result: result, err: err}
	}
}
