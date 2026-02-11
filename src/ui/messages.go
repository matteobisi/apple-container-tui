package ui

// ActiveScreen identifies the currently displayed screen.
type ActiveScreen string

const (
	// ScreenContainerList shows the container list.
	ScreenContainerList ActiveScreen = "containers"
	// ScreenImagePull shows the image pull workflow.
	ScreenImagePull ActiveScreen = "image-pull"
	// ScreenFilePicker shows the build file picker.
	ScreenFilePicker ActiveScreen = "file-picker"
	// ScreenBuild shows the build workflow.
	ScreenBuild ActiveScreen = "build"
	// ScreenDaemonControl shows daemon start/stop controls.
	ScreenDaemonControl ActiveScreen = "daemon-control"
	// ScreenHelp shows the help screen.
	ScreenHelp ActiveScreen = "help"
)

type screenChangeMsg struct {
	target ActiveScreen
}

type buildFileSelectedMsg struct {
	path string
}
