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
	if len(cmd.Args) > 1 && cmd.Args[0] == "machine" && cmd.Args[1] == "list" {
		return models.Result{Stdout: `[{"id":"dev","image":"alpine:latest","state":"running","default":true,"cpus":4,"memory":"8G","homeMount":"rw"}]`, Status: models.ResultSuccess}, nil
	}
	if len(cmd.Args) > 0 && cmd.Args[0] == "registry" {
		return models.Result{Stdout: `[{"creationDate":"2025-11-04T01:17:34Z","id":"registry.example.com","labels":{},"modificationDate":"2025-11-04T01:17:34Z","name":"registry.example.com","username":"user"}]`, Status: models.ResultSuccess}, nil
	}
	if len(cmd.Args) > 1 && cmd.Args[0] == "system" && cmd.Args[1] == "status" {
		return models.Result{Stdout: `{"status":"running"}`, Status: models.ResultSuccess}, nil
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
	if _, ok := msg.(BackToListMsg); !ok {
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
	msg := daemonStatusMsg{status: models.DaemonStatus{State: models.DaemonStateRunning, Running: true}}
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

func TestRegistriesScreenFlow(t *testing.T) {
	screen := NewRegistriesScreen(flowExecutor{})
	updated, _ := screen.Update(registriesLoadedMsg{entries: []models.RegistryLogin{{Hostname: "registry.example.com", Username: "user"}}})
	if len(updated.entries) != 1 {
		t.Fatalf("expected one registry entry")
	}
	_, cmd := updated.Update(tea.KeyMsg{Type: tea.KeyEscape})
	if cmd == nil {
		t.Fatalf("expected back cmd")
	}
}

func TestMachineListAndSubmenuFlow(t *testing.T) {
	machine := models.ContainerMachine{ID: "dev", Image: "alpine:latest", State: models.MachineStateRunning, IsDefault: true, CPUs: 4, Memory: "8G", HomeMount: "rw"}
	screen := NewMachineListScreen(flowExecutor{})
	updated, _ := screen.Update(machineListLoadedMsg{machines: []models.ContainerMachine{machine}})
	if len(updated.machines) != 1 || !updated.hasLoaded {
		t.Fatalf("expected one loaded machine")
	}
	_, cmd := updated.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatalf("expected submenu navigation cmd")
	}
	msg := cmd()
	change, ok := msg.(screenChangeMsg)
	if !ok || change.target != ScreenMachineSubmenu || change.machine == nil || change.machine.ID != "dev" {
		t.Fatalf("unexpected navigation message: %#v", msg)
	}

	submenu := NewMachineSubmenuScreen(flowExecutor{}).SetMachine(machine)
	view := submenu.View()
	if !strings.Contains(view, "Stop machine") || strings.Contains(view, "Start machine") {
		t.Fatalf("unexpected running machine options: %s", view)
	}
	submenu.cursor = 2
	updatedSubmenu, _ := submenu.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if updatedSubmenu.preview == nil || !strings.Contains(updatedSubmenu.preview.Command.String(), "machine stop dev") {
		t.Fatalf("expected stop preview")
	}

	stopped := machine
	stopped.State = models.MachineStateStopped
	stoppedSubmenu := NewMachineSubmenuScreen(flowExecutor{}).SetMachine(stopped)
	if !strings.Contains(stoppedSubmenu.View(), "Start machine") {
		t.Fatalf("expected start action for stopped machine")
	}
	_, cmd = updated.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}})
	if cmd == nil {
		t.Fatalf("expected create navigation cmd")
	}
	if change, ok := cmd().(screenChangeMsg); !ok || change.target != ScreenMachineCreate {
		t.Fatalf("unexpected create navigation message")
	}
}

func TestMachineDetailAndEditBackFlows(t *testing.T) {
	machine := models.ContainerMachine{ID: "dev", Image: "alpine:latest", State: models.MachineStateRunning, CPUs: 4, Memory: "8G", HomeMount: "rw"}

	inspect := NewMachineInspectScreen(flowExecutor{}).SetMachine(machine)
	_, cmd := inspect.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	if cmd == nil {
		t.Fatalf("expected inspect back cmd")
	}
	if change, ok := cmd().(screenChangeMsg); !ok || change.target != ScreenMachineSubmenu {
		t.Fatalf("unexpected inspect back message")
	}

	logs := NewMachineLogsScreen(flowExecutor{}).SetMachine(machine)
	_, cmd = logs.Update(tea.KeyMsg{Type: tea.KeyEscape})
	if cmd == nil {
		t.Fatalf("expected logs back cmd")
	}

	edit := NewMachineEditResourcesScreen(flowExecutor{}).SetMachine(machine)
	updated, _ := edit.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if updated.preview == nil || !strings.Contains(updated.preview.Command.String(), "machine set -n dev cpus=4 memory=8G home-mount=rw") {
		t.Fatalf("expected resource edit preview")
	}
}

func TestMachineCreateFlow(t *testing.T) {
	screen := NewMachineCreateScreen(flowExecutor{})
	screen.inputs[0].SetValue("alpine:latest")
	screen.inputs[1].SetValue("dev")
	updated, _ := screen.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if updated.preview == nil || !strings.Contains(updated.preview.Command.String(), "machine create alpine:latest --name dev") {
		t.Fatalf("expected create preview")
	}
	updated, cmd := updated.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
	if cmd == nil || updated.loading != true {
		t.Fatalf("expected create execution")
	}
	updated, cmd = updated.Update(machineCreateResultMsg{result: models.Result{Status: models.ResultSuccess}})
	if cmd == nil || updated.result == nil {
		t.Fatalf("expected create result and return command")
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
