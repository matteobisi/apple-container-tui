package ui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"container-tui/src/models"
)

type flowExecutor struct {
	listOutput string
	result     models.Result
}

func (f flowExecutor) Execute(cmd models.Command) (models.Result, error) {
	if len(cmd.Args) > 0 && cmd.Args[0] == "list" {
		return models.Result{Stdout: f.listOutput, Status: models.ResultSuccess}, nil
	}
	return f.result, nil
}

func TestContainerListUpdateFlows(t *testing.T) {
	output := "CONTAINER ID  IMAGE  COMMAND  CREATED  STATUS  PORTS\n" +
		"abc123        nginx  web     1m      running  \n"
	exec := flowExecutor{listOutput: output, result: models.Result{Status: models.ResultSuccess}}
	screen := NewContainerListScreen(exec)

	msg := containerListLoadedMsg{containers: []models.Container{{ID: "abc", Name: "web", Status: models.ContainerStatusRunning, Image: "nginx"}}}
	updated, _ := screen.Update(msg)
	if !updated.hasLoaded {
		t.Fatalf("expected hasLoaded true")
	}

	updated.preview = &CommandPreviewModal{Command: models.Command{Executable: "container"}}
	updated, cmd := updated.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
	if cmd == nil || updated.pendingCmd == nil {
		t.Fatalf("expected command execution")
	}

	updated.confirm = nil
	updated, cmd = updated.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})
	if cmd == nil {
		t.Fatalf("expected refresh cmd")
	}
}

func TestContainerListView(t *testing.T) {
	screen := NewContainerListScreen(flowExecutor{})
	screen.errorMsg = "oops"
	screen.loading = true
	screen.preview = &CommandPreviewModal{Command: models.Command{Executable: "container"}}
	screen.result = &models.Result{Status: models.ResultSuccess}
	screen.containers = []models.Container{{Name: "web", Status: models.ContainerStatusStopped, Image: "nginx"}}
	view := screen.View()
	if !strings.Contains(view, "Containers") {
		t.Fatalf("unexpected view")
	}
}

func TestImagePullConfirmFlow(t *testing.T) {
	exec := flowExecutor{result: models.Result{Status: models.ResultSuccess}}
	screen := NewImagePullScreen(exec)
	screen.input.SetValue("nginx:latest")
	updated, _ := screen.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if updated.preview == nil {
		t.Fatalf("expected preview")
	}
	updated, cmd := updated.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
	if cmd == nil {
		t.Fatalf("expected execute cmd")
	}
	updated, _ = updated.Update(imagePullResultMsg{result: models.Result{Status: models.ResultSuccess}})
	if updated.result == nil {
		t.Fatalf("expected result")
	}
}

func TestImagePullEscHelp(t *testing.T) {
	screen := NewImagePullScreen(flowExecutor{})
	_, cmd := screen.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
	msg := cmd()
	if change, ok := msg.(screenChangeMsg); !ok || change.target != ScreenHelp {
		t.Fatalf("unexpected message")
	}
	_, cmd = screen.Update(tea.KeyMsg{Type: tea.KeyEscape})
	msg = cmd()
	if change, ok := msg.(screenChangeMsg); !ok || change.target != ScreenContainerList {
		t.Fatalf("unexpected message")
	}
}

func TestBuildScreenConfirmFlow(t *testing.T) {
	exec := flowExecutor{result: models.Result{Status: models.ResultSuccess}}
	screen := NewBuildScreen(exec, "./Containerfile")
	screen.input.SetValue("my-app:latest")
	updated, _ := screen.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if updated.preview == nil {
		t.Fatalf("expected preview")
	}
	updated, cmd := updated.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
	if cmd == nil {
		t.Fatalf("expected execute cmd")
	}
	updated, _ = updated.Update(buildResultMsg{result: models.Result{Status: models.ResultSuccess}})
	if updated.result == nil {
		t.Fatalf("expected result")
	}
}

func TestBuildScreenMissingFile(t *testing.T) {
	screen := NewBuildScreen(flowExecutor{}, "")
	updated, _ := screen.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if updated.errorMsg == "" {
		t.Fatalf("expected error")
	}
	_, cmd := updated.Update(tea.KeyMsg{Type: tea.KeyEscape})
	if cmd == nil {
		t.Fatalf("expected back cmd")
	}
}

func TestDaemonControlStatusFlow(t *testing.T) {
	exec := flowExecutor{result: models.Result{Stdout: "running", Status: models.ResultSuccess}}
	screen := NewDaemonControlScreen(exec)
	msg := daemonStatusMsg{status: models.DaemonStatus{Running: true}}
	updated, _ := screen.Update(msg)
	view := updated.View()
	if !strings.Contains(view, "running") {
		t.Fatalf("expected running in view")
	}
	updated, _ = updated.Update(daemonActionMsg{result: models.Result{Status: models.ResultSuccess}})
	if updated.result == nil {
		t.Fatalf("expected result")
	}
}

func TestDaemonControlStopConfirm(t *testing.T) {
	screen := NewDaemonControlScreen(flowExecutor{})
	updated, _ := screen.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'t'}})
	if updated.confirm == nil {
		t.Fatalf("expected confirm modal")
	}
	refreshed, _ := screen.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})
	if !refreshed.loading {
		t.Fatalf("expected loading")
	}
}

func TestAppModelUpdate(t *testing.T) {
	exec := flowExecutor{}
	app := NewAppModel(exec, "1.0.0")
	app.active = ScreenHelp
	model, cmd := app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	if cmd == nil {
		t.Fatalf("expected quit cmd")
	}
	appModel := model.(AppModel)
	left, right := appModel.statusBarInfo()
	if left == "" || right != "" {
		t.Fatalf("unexpected status info")
	}
	appModel.active = ScreenImagePull
	appModel.imagePull.preview = &CommandPreviewModal{Command: models.Command{Executable: "container"}}
	_, right = appModel.statusBarInfo()
	if !strings.Contains(right, "Preview") {
		t.Fatalf("expected preview")
	}
}

func TestHelpScreenUpdate(t *testing.T) {
	help := HelpScreen{}
	_, cmd := help.Update(tea.KeyMsg{Type: tea.KeyEscape})
	msg := cmd()
	if change, ok := msg.(screenChangeMsg); !ok || change.target != ScreenContainerList {
		t.Fatalf("unexpected message")
	}
}
