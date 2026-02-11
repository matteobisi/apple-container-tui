package ui

import (
	"errors"
	"strings"
	"testing"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"

	"container-tui/src/models"
)

type statusExecutor struct {
	stdout string
	result models.Result
}

func (s statusExecutor) Execute(cmd models.Command) (models.Result, error) {
	if len(cmd.Args) > 0 && cmd.Args[0] == "system" && cmd.Args[1] == "status" {
		return models.Result{Stdout: s.stdout, Status: models.ResultSuccess}, nil
	}
	return s.result, nil
}

func TestAppModelWindowSizeAndView(t *testing.T) {
	app := NewAppModel(flowExecutor{}, "1.0.0")
	model, _ := app.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	updated := model.(AppModel)
	if updated.width != 80 || updated.height != 24 {
		t.Fatalf("unexpected size: %dx%d", updated.width, updated.height)
	}
	updated.containerList.containers = []models.Container{{Name: "web", Status: models.ContainerStatusRunning, Image: "nginx"}}
	view := updated.View()
	if !strings.Contains(view, "Containers") {
		t.Fatalf("unexpected view: %q", view)
	}
}

func TestAppModelScreenChangeAndBuildSelection(t *testing.T) {
	app := NewAppModel(flowExecutor{}, "1.0.0")
	model, _ := app.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
	appModel := model.(AppModel)
	model, _ = appModel.Update(screenChangeMsg{target: ScreenImagePull})
	appModel = model.(AppModel)
	if appModel.active != ScreenImagePull {
		t.Fatalf("expected image pull screen")
	}
	model, _ = appModel.Update(buildFileSelectedMsg{path: "./Containerfile"})
	appModel = model.(AppModel)
	if appModel.active != ScreenBuild {
		t.Fatalf("expected build screen")
	}
	if appModel.buildScreen.filePath != "./Containerfile" {
		t.Fatalf("unexpected build path: %s", appModel.buildScreen.filePath)
	}
	if appModel.buildScreen.viewport.Width == 0 || appModel.buildScreen.viewport.Height == 0 {
		t.Fatalf("expected build viewport sized")
	}
}

func TestAppModelSpinnerActive(t *testing.T) {
	app := NewAppModel(flowExecutor{}, "1.0.0")
	app.active = ScreenImagePull
	app.imagePull.loading = true
	model, _ := app.Update(spinner.TickMsg{})
	updated := model.(AppModel)
	if updated.spinner.View() == "" {
		t.Fatalf("expected spinner output")
	}
}

func TestContainerListDeleteConfirmation(t *testing.T) {
	screen := NewContainerListScreen(flowExecutor{})
	screen.containers = []models.Container{{ID: "abc", Name: "web", Status: models.ContainerStatusRunning}}
	updated, _ := screen.buildAndConfirmDelete()
	if updated.errorMsg == "" {
		t.Fatalf("expected error for running container")
	}

	screen.containers = []models.Container{{ID: "abc", Name: "web", Status: models.ContainerStatusStopped}}
	updated, _ = screen.buildAndConfirmDelete()
	if updated.confirm == nil {
		t.Fatalf("expected confirm modal")
	}
}

func TestContainerListExecuteCommandCmd(t *testing.T) {
	screen := NewContainerListScreen(flowExecutor{result: models.Result{Status: models.ResultSuccess}})
	cmd := screen.executeCommandCmd(models.Command{Executable: "container", Args: []string{"list"}})
	msg := cmd()
	if _, ok := msg.(commandExecutedMsg); !ok {
		t.Fatalf("unexpected message: %#v", msg)
	}
}

func TestDaemonControlCommands(t *testing.T) {
	screen := NewDaemonControlScreen(statusExecutor{stdout: "running", result: models.Result{Status: models.ResultSuccess}})
	cmd := screen.fetchStatusCmd()
	msg := cmd()
	if _, ok := msg.(daemonStatusMsg); !ok {
		t.Fatalf("unexpected status msg: %#v", msg)
	}
	action := screen.executeCommandCmd(models.Command{Executable: "container", Args: []string{"system", "start"}})
	if _, ok := action().(daemonActionMsg); !ok {
		t.Fatalf("unexpected action msg")
	}
}

func TestBuildScreenCancelAndError(t *testing.T) {
	screen := NewBuildScreen(flowExecutor{}, "./Containerfile")
	screen.input.SetValue("app:latest")
	updated, _ := screen.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if updated.preview == nil {
		t.Fatalf("expected preview")
	}
	updated, _ = updated.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
	if updated.preview != nil {
		t.Fatalf("expected preview cleared")
	}
	updated, _ = updated.Update(buildResultMsg{err: errors.New("boom"), result: models.Result{Status: models.ResultError, Stderr: "boom"}})
	if updated.errorMsg == "" {
		t.Fatalf("expected error message")
	}
}

func TestImagePullCancelAndError(t *testing.T) {
	screen := NewImagePullScreen(flowExecutor{})
	screen.input.SetValue("nginx:latest")
	updated, _ := screen.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if updated.preview == nil {
		t.Fatalf("expected preview")
	}
	updated, _ = updated.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
	if updated.preview != nil {
		t.Fatalf("expected preview cleared")
	}
	updated, _ = updated.Update(imagePullResultMsg{err: errors.New("oops"), result: models.Result{Status: models.ResultError, Stderr: "oops"}})
	if updated.errorMsg == "" {
		t.Fatalf("expected error message")
	}
}

func TestFilePickerWindowSize(t *testing.T) {
	screen := NewFilePickerScreen(flowExecutor{})
	updated, _ := screen.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	if updated.width != 120 || updated.height != 40 {
		t.Fatalf("unexpected dimensions")
	}
	if updated.picker.Height <= 0 {
		t.Fatalf("expected picker height set")
	}
}

func TestHelpScreenUpdateQuestion(t *testing.T) {
	help := HelpScreen{}
	_, cmd := help.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
	msg := cmd()
	if change, ok := msg.(screenChangeMsg); !ok || change.target != ScreenContainerList {
		t.Fatalf("unexpected message")
	}
}

func TestYesNoConfirmModalViewWarning(t *testing.T) {
	modal := YesNoConfirmModal{Title: "Stop", Body: "Stop services?", Command: models.Command{Executable: "container"}, Warning: true}
	view := modal.View()
	if !strings.Contains(view, "destructive") {
		t.Fatalf("expected warning in view")
	}
}

func TestTypeToConfirmCancel(t *testing.T) {
	modal := NewTypeToConfirmModal("Delete", "match", models.Command{Executable: "container"})
	_, confirmed, canceled := modal.Handle(tea.KeyMsg{Type: tea.KeyEscape})
	if confirmed || !canceled {
		t.Fatalf("expected cancel")
	}
}

func TestThemeResolveModes(t *testing.T) {
	ApplyTheme("light")
	_ = CurrentTheme()
	ApplyTheme("dark")
	_ = RenderSuccess("ok")
	ApplyTheme("auto")
	_ = RenderWarning("warn")
}
