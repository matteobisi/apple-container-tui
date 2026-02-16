package ui

import (
	"log"
	"os"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"container-tui/src/models"
	"container-tui/src/services"
)

// AppModel holds global UI state.
type AppModel struct {
	width  int
	height int
	keys   KeyMap
	active ActiveScreen
	stack  []ActiveScreen

	selectedContainer *models.Container
	selectedImage     *models.Image
	navDebugEnabled   bool

	containerList  ContainerListScreen
	containerSub   ContainerSubmenuScreen
	containerLogs  ContainerLogsScreen
	containerShell ContainerShellScreen
	imageList      ImageListScreen
	imageSub       ImageSubmenuScreen
	imageInspect   ImageInspectScreen
	imagePull      ImagePullScreen
	filePicker     FilePickerScreen
	buildScreen    BuildScreen
	daemonControl  DaemonControlScreen
	help           HelpScreen
	spinner        SpinnerModel
}

// NewAppModel creates the initial app model.
func NewAppModel(executor services.CommandExecutor, version string) AppModel {
	return AppModel{
		keys:            DefaultKeyMap(),
		active:          ScreenContainerList,
		stack:           []ActiveScreen{},
		navDebugEnabled: os.Getenv("ACTUI_DEBUG_NAV") == "1",
		containerList:   NewContainerListScreen(executor),
		containerSub:    NewContainerSubmenuScreen(executor),
		containerLogs:   NewContainerLogsScreen(executor),
		containerShell:  NewContainerShellScreen(executor),
		imageList:       NewImageListScreen(executor),
		imageSub:        NewImageSubmenuScreen(executor),
		imageInspect:    NewImageInspectScreen(executor),
		imagePull:       NewImagePullScreen(executor),
		filePicker:      NewFilePickerScreen(executor),
		buildScreen:     NewBuildScreen(executor, ""),
		daemonControl:   NewDaemonControlScreen(executor),
		help:            HelpScreen{Version: version},
		spinner:         NewSpinnerModel(),
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
		m.containerSub, _ = m.containerSub.Update(message)
		m.containerLogs, _ = m.containerLogs.Update(message)
		m.containerShell, _ = m.containerShell.Update(message)
		m.imageList, _ = m.imageList.Update(message)
		m.imageSub, _ = m.imageSub.Update(message)
		m.imageInspect, _ = m.imageInspect.Update(message)
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
		origin := m.active
		if message.push {
			m.pushView(m.active)
		}
		if message.container != nil {
			containerCopy := *message.container
			m.selectedContainer = &containerCopy
			m.containerSub = m.containerSub.SetContainer(containerCopy)
			m.containerLogs = m.containerLogs.SetContainer(containerCopy)
			m.containerShell = m.containerShell.SetContainer(containerCopy)
		}
		if message.image != nil {
			imageCopy := *message.image
			m.selectedImage = &imageCopy
			m.imageSub = m.imageSub.SetImage(imageCopy)
			m.imageInspect = m.imageInspect.SetImage(imageCopy)
		}
		if message.target == ScreenImagePull {
			if origin == ScreenImageList {
				m.imagePull = m.imagePull.SetReturnTarget(ScreenImageList)
			} else {
				m.imagePull = m.imagePull.SetReturnTarget(ScreenContainerList)
			}
		}
		if message.target == ScreenFilePicker {
			if origin == ScreenImageList {
				m.filePicker = m.filePicker.SetReturnTarget(ScreenImageList)
			} else {
				m.filePicker = m.filePicker.SetReturnTarget(ScreenContainerList)
			}
		}
		m.active = message.target
		m.logNavigation("screen-change", origin, m.active)
		switch m.active {
		case ScreenContainerList:
			m.stack = []ActiveScreen{}
			cmd = m.containerList.Init()
		case ScreenContainerSubmenu:
			cmd = m.containerSub.Init()
		case ScreenContainerLogs:
			cmd = m.containerLogs.Init()
		case ScreenContainerShell:
			cmd = m.containerShell.Init()
		case ScreenImageList:
			cmd = m.imageList.Init()
		case ScreenImageSubmenu:
			cmd = m.imageSub.Init()
		case ScreenImageInspect:
			cmd = m.imageInspect.Init()
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
		m.buildScreen = m.buildScreen.SetReturnTarget(message.returnTarget)
		if m.width > 0 && m.height > 0 {
			m.buildScreen, _ = m.buildScreen.Update(tea.WindowSizeMsg{Width: m.width, Height: m.height})
		}
		m.active = ScreenBuild
		cmd = m.buildScreen.Init()
		skipScreenUpdate = true
	case BackToListMsg:
		origin := m.active
		m.active = m.popView(ScreenContainerList)
		m.logNavigation("back-to-list", origin, m.active)
		cmd = m.initForActive()
		skipScreenUpdate = true
	case BackToSubmenuMsg:
		origin := m.active
		if m.selectedImage != nil {
			m.active = ScreenImageSubmenu
			m.imageSub = m.imageSub.SetImage(*m.selectedImage)
			cmd = m.imageSub.Init()
		} else if m.selectedContainer != nil {
			m.active = ScreenContainerSubmenu
			m.containerSub = m.containerSub.SetContainer(*m.selectedContainer)
			cmd = m.containerSub.Init()
		}
		m.logNavigation("back-to-submenu", origin, m.active)
		skipScreenUpdate = true
	}

	if !skipScreenUpdate {
		switch m.active {
		case ScreenImagePull:
			updated, updateCmd := m.imagePull.Update(msg)
			m.imagePull = updated
			cmd = tea.Batch(cmd, updateCmd)
		case ScreenContainerSubmenu:
			updated, updateCmd := m.containerSub.Update(msg)
			m.containerSub = updated
			cmd = tea.Batch(cmd, updateCmd)
		case ScreenContainerLogs:
			updated, updateCmd := m.containerLogs.Update(msg)
			m.containerLogs = updated
			cmd = tea.Batch(cmd, updateCmd)
		case ScreenContainerShell:
			updated, updateCmd := m.containerShell.Update(msg)
			m.containerShell = updated
			cmd = tea.Batch(cmd, updateCmd)
		case ScreenImageList:
			updated, updateCmd := m.imageList.Update(msg)
			m.imageList = updated
			cmd = tea.Batch(cmd, updateCmd)
		case ScreenImageSubmenu:
			updated, updateCmd := m.imageSub.Update(msg)
			m.imageSub = updated
			cmd = tea.Batch(cmd, updateCmd)
		case ScreenImageInspect:
			updated, updateCmd := m.imageInspect.Update(msg)
			m.imageInspect = updated
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
	case ScreenContainerSubmenu:
		return m.containerSub.View() + "\n" + status
	case ScreenContainerLogs:
		return m.containerLogs.View() + "\n" + status
	case ScreenContainerShell:
		return m.containerShell.View() + "\n" + status
	case ScreenImageList:
		return m.imageList.View() + "\n" + status
	case ScreenImageSubmenu:
		return m.imageSub.View() + "\n" + status
	case ScreenImageInspect:
		return m.imageInspect.View() + "\n" + status
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
	case ScreenContainerSubmenu:
		label = "Container Actions"
		if m.containerSub.preview != nil {
			preview = m.containerSub.preview.Command.String()
		}
	case ScreenContainerLogs:
		label = "Container Logs"
	case ScreenContainerShell:
		label = "Container Shell"
	case ScreenImageList:
		label = "Images"
	case ScreenImageSubmenu:
		label = "Image Actions"
	case ScreenImageInspect:
		label = "Image Inspect"
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
	case ScreenContainerLogs:
		return m.containerLogs.loading
	case ScreenContainerShell:
		return m.containerShell.loading
	case ScreenImageList:
		return m.imageList.loading
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

func (m *AppModel) pushView(screen ActiveScreen) {
	m.stack = append(m.stack, screen)
	m.logNavigation("push", m.active, screen)
}

func (m *AppModel) popView(fallback ActiveScreen) ActiveScreen {
	if len(m.stack) == 0 {
		return fallback
	}
	idx := len(m.stack) - 1
	view := m.stack[idx]
	m.stack = m.stack[:idx]
	m.logNavigation("pop", m.active, view)
	return view
}

func (m AppModel) logNavigation(event string, from ActiveScreen, to ActiveScreen) {
	if !m.navDebugEnabled {
		return
	}
	log.Printf("[nav] event=%s from=%s to=%s stack_depth=%d", event, from, to, len(m.stack))
}

func (m AppModel) initForActive() tea.Cmd {
	switch m.active {
	case ScreenContainerSubmenu:
		return m.containerSub.Init()
	case ScreenContainerLogs:
		return m.containerLogs.Init()
	case ScreenContainerShell:
		return m.containerShell.Init()
	case ScreenImageList:
		return m.imageList.Init()
	case ScreenImageSubmenu:
		return m.imageSub.Init()
	case ScreenImageInspect:
		return m.imageInspect.Init()
	case ScreenImagePull:
		return m.imagePull.Init()
	case ScreenFilePicker:
		return m.filePicker.Init()
	case ScreenBuild:
		return m.buildScreen.Init()
	case ScreenDaemonControl:
		return m.daemonControl.Init()
	case ScreenHelp:
		return m.help.Init()
	default:
		return m.containerList.Init()
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
