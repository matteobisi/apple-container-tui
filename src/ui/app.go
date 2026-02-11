package ui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"container-tui/src/services"
)

// AppModel holds global UI state.
type AppModel struct {
	width  int
	height int
	keys   KeyMap
	active ActiveScreen

	containerList ContainerListScreen
	imagePull     ImagePullScreen
	filePicker    FilePickerScreen
	buildScreen   BuildScreen
	daemonControl DaemonControlScreen
	help          HelpScreen
	spinner       SpinnerModel
}

// NewAppModel creates the initial app model.
func NewAppModel(executor services.CommandExecutor, version string) AppModel {
	return AppModel{
		keys:          DefaultKeyMap(),
		active:        ScreenContainerList,
		containerList: NewContainerListScreen(executor),
		imagePull:     NewImagePullScreen(executor),
		filePicker:    NewFilePickerScreen(executor),
		buildScreen:   NewBuildScreen(executor, ""),
		daemonControl: NewDaemonControlScreen(executor),
		help:          HelpScreen{Version: version},
		spinner:       NewSpinnerModel(),
	}
}

// Init starts any initial commands.
func (m AppModel) Init() tea.Cmd {
	return m.containerList.Init()
}

// Update handles incoming messages.
func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	skipScreenUpdate := false

	switch message := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = message.Width
		m.height = message.Height
		m.containerList, _ = m.containerList.Update(message)
		m.imagePull, _ = m.imagePull.Update(message)
		m.filePicker, _ = m.filePicker.Update(message)
		m.buildScreen, _ = m.buildScreen.Update(message)
		m.daemonControl, _ = m.daemonControl.Update(message)
		m.help, _ = m.help.Update(message)
	case tea.KeyMsg:
		if keyMatches(message, m.keys.Quit) {
			return m, tea.Quit
		}
	case screenChangeMsg:
		m.active = message.target
		switch m.active {
		case ScreenContainerList:
			cmd = m.containerList.Init()
		case ScreenImagePull:
			cmd = m.imagePull.Init()
		case ScreenFilePicker:
			cmd = m.filePicker.Init()
		case ScreenBuild:
			cmd = m.buildScreen.Init()
		case ScreenDaemonControl:
			cmd = m.daemonControl.Init()
		case ScreenHelp:
			cmd = m.help.Init()
		}
		skipScreenUpdate = true
	case buildFileSelectedMsg:
		m.buildScreen = NewBuildScreen(m.filePicker.executor, message.path)
		if m.width > 0 && m.height > 0 {
			m.buildScreen, _ = m.buildScreen.Update(tea.WindowSizeMsg{Width: m.width, Height: m.height})
		}
		m.active = ScreenBuild
		cmd = m.buildScreen.Init()
		skipScreenUpdate = true
	}

	if !skipScreenUpdate {
		switch m.active {
		case ScreenImagePull:
			updated, updateCmd := m.imagePull.Update(msg)
			m.imagePull = updated
			cmd = tea.Batch(cmd, updateCmd)
		case ScreenFilePicker:
			updated, updateCmd := m.filePicker.Update(msg)
			m.filePicker = updated
			cmd = tea.Batch(cmd, updateCmd)
		case ScreenBuild:
			updated, updateCmd := m.buildScreen.Update(msg)
			m.buildScreen = updated
			cmd = tea.Batch(cmd, updateCmd)
		case ScreenDaemonControl:
			updated, updateCmd := m.daemonControl.Update(msg)
			m.daemonControl = updated
			cmd = tea.Batch(cmd, updateCmd)
		case ScreenHelp:
			updated, updateCmd := m.help.Update(msg)
			m.help = updated
			cmd = tea.Batch(cmd, updateCmd)
		default:
			updated, updateCmd := m.containerList.Update(msg)
			m.containerList = updated
			cmd = tea.Batch(cmd, updateCmd)
		}
	}

	m.spinner.SetActive(m.isLoading())
	updatedSpinner, spinnerCmd := m.spinner.Update(msg)
	m.spinner = updatedSpinner
	cmd = tea.Batch(cmd, spinnerCmd)

	return m, cmd
}

// View renders the UI.
func (m AppModel) View() string {
	left, right := m.statusBarInfo()
	status := RenderStatusBar(m.width, left, right)

	switch m.active {
	case ScreenImagePull:
		return m.imagePull.View() + "\n" + status
	case ScreenFilePicker:
		return m.filePicker.View() + "\n" + status
	case ScreenBuild:
		return m.buildScreen.View() + "\n" + status
	case ScreenDaemonControl:
		return m.daemonControl.View() + "\n" + status
	case ScreenHelp:
		return m.help.View() + "\n" + status
	default:
		return m.containerList.View() + "\n" + status
	}
}

func (m AppModel) statusBarInfo() (string, string) {
	spinner := m.spinner.View()
	label := ""
	preview := ""

	switch m.active {
	case ScreenImagePull:
		label = "Pull Image"
		if m.imagePull.preview != nil {
			preview = m.imagePull.preview.Command.String()
		}
	case ScreenFilePicker:
		label = "File Picker"
	case ScreenBuild:
		label = "Build"
		if m.buildScreen.preview != nil {
			preview = m.buildScreen.preview.Command.String()
		}
	case ScreenDaemonControl:
		label = "Daemon"
		if m.daemonControl.confirm != nil {
			preview = m.daemonControl.confirm.Command.String()
		}
	case ScreenHelp:
		label = "Help"
	default:
		label = "Containers"
		if m.containerList.preview != nil {
			preview = m.containerList.preview.Command.String()
		} else if m.containerList.confirm != nil {
			preview = m.containerList.confirm.Command.String()
		}
	}

	left := label
	if spinner != "" {
		left = spinner + " " + label
	}
	left = RenderMuted(left)
	if preview != "" {
		preview = RenderMuted("Preview: " + preview)
	}

	return left, preview
}

func (m AppModel) isLoading() bool {
	switch m.active {
	case ScreenImagePull:
		return m.imagePull.loading
	case ScreenBuild:
		return m.buildScreen.loading
	case ScreenDaemonControl:
		return m.daemonControl.loading
	default:
		return m.containerList.loading
	}
}

func keyMatches(msg tea.KeyMsg, binding key.Binding) bool {
	for _, bindingKey := range binding.Keys() {
		if msg.String() == bindingKey {
			return true
		}
	}
	return false
}
