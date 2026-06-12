package ui

import "container-tui/src/models"

// ActiveScreen identifies the currently displayed screen.
type ActiveScreen string

const (
	// ScreenContainerList shows the container list.
	ScreenContainerList ActiveScreen = "containers"
	// ScreenContainerSubmenu shows actions for selected container.
	ScreenContainerSubmenu ActiveScreen = "container-submenu"
	// ScreenContainerLogs shows tail logs for selected container.
	ScreenContainerLogs ActiveScreen = "container-logs"
	// ScreenContainerShell runs/represents interactive shell for selected container.
	ScreenContainerShell ActiveScreen = "container-shell"
	// ScreenImageList shows local image list.
	ScreenImageList ActiveScreen = "image-list"
	// ScreenImageSubmenu shows actions for selected image.
	ScreenImageSubmenu ActiveScreen = "image-submenu"
	// ScreenImageInspect shows image inspect output.
	ScreenImageInspect ActiveScreen = "image-inspect"
	// ScreenImagePull shows the image pull workflow.
	ScreenImagePull ActiveScreen = "image-pull"
	// ScreenRegistries shows runtime-managed registry entries.
	ScreenRegistries ActiveScreen = "registries"
	// ScreenMachineList shows container machines.
	ScreenMachineList ActiveScreen = "machine-list"
	// ScreenMachineSubmenu shows actions for selected container machine.
	ScreenMachineSubmenu ActiveScreen = "machine-submenu"
	// ScreenMachineInspect shows container machine inspect output.
	ScreenMachineInspect ActiveScreen = "machine-inspect"
	// ScreenMachineLogs shows container machine logs.
	ScreenMachineLogs ActiveScreen = "machine-logs"
	// ScreenMachineEditResources edits container machine resources.
	ScreenMachineEditResources ActiveScreen = "machine-edit-resources"
	// ScreenMachineCreate creates a container machine.
	ScreenMachineCreate ActiveScreen = "machine-create"
	// ScreenFilePicker shows the build file picker.
	ScreenFilePicker ActiveScreen = "file-picker"
	// ScreenBuild shows the build workflow.
	ScreenBuild ActiveScreen = "build"
	// ScreenContainerExport shows the container export workflow.
	ScreenContainerExport ActiveScreen = "container-export"
	// ScreenDaemonControl shows daemon start/stop controls.
	ScreenDaemonControl ActiveScreen = "daemon-control"
	// ScreenHelp shows the help screen.
	ScreenHelp ActiveScreen = "help"
)

type screenChangeMsg struct {
	target    ActiveScreen
	container *models.Container
	image     *models.Image
	machine   *models.ContainerMachine
	push      bool
}

type buildFileSelectedMsg struct {
	path         string
	returnTarget ActiveScreen
}

// BackToListMsg is emitted when child views want to return to their parent list.
type BackToListMsg struct{}

// BackToSubmenuMsg is emitted when nested detail views return to submenu.
type BackToSubmenuMsg struct{}
