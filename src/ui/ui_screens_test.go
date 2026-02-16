package ui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"container-tui/src/models"
	"container-tui/src/services"
)

type fakeExecutor struct {
	listOutput string
}

func (f fakeExecutor) Execute(cmd models.Command) (models.Result, error) {
	if len(cmd.Args) > 0 && cmd.Args[0] == "list" {
		return models.Result{Stdout: f.listOutput, Status: models.ResultSuccess}, nil
	}
	return models.Result{Status: models.ResultSuccess}, nil
}

func TestFilePickerEsc(t *testing.T) {
	screen := NewFilePickerScreen(fakeExecutor{})
	_, cmd := screen.Update(tea.KeyMsg{Type: tea.KeyEscape})
	if cmd == nil {
		t.Fatalf("expected command")
	}
	msg := cmd()
	if _, ok := msg.(BackToListMsg); !ok {
		t.Fatalf("unexpected message: %#v", msg)
	}
}

func TestImagePullPreview(t *testing.T) {
	screen := NewImagePullScreen(fakeExecutor{})
	screen.input.SetValue("nginx:latest")
	updated, _ := screen.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if updated.preview == nil {
		t.Fatalf("expected preview")
	}
}

func TestBuildScreenPreview(t *testing.T) {
	screen := NewBuildScreen(fakeExecutor{}, "./Containerfile")
	screen.input.SetValue("my-image:latest")
	updated, _ := screen.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if updated.preview == nil {
		t.Fatalf("expected preview")
	}
}

func TestDaemonControlConfirm(t *testing.T) {
	screen := NewDaemonControlScreen(fakeExecutor{})
	updated, _ := screen.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
	if updated.confirm == nil {
		t.Fatalf("expected confirm modal")
	}
}

func TestContainerListCache(t *testing.T) {
	screen := NewContainerListScreen(fakeExecutor{})
	screen.hasLoaded = true
	if cmd := screen.fetchContainersCmd(false); cmd != nil {
		t.Fatalf("expected nil cmd when cached")
	}
}

func TestContainerListFetch(t *testing.T) {
	output := "CONTAINER ID  IMAGE  COMMAND  CREATED  STATUS  PORTS\n" +
		"abc123        nginx  web     1m      running  \n"
	screen := NewContainerListScreen(fakeExecutor{listOutput: output})
	cmd := screen.fetchContainersCmd(true)
	if cmd == nil {
		t.Fatalf("expected fetch cmd")
	}
	msg := cmd()
	loaded, ok := msg.(containerListLoadedMsg)
	if !ok || loaded.err != nil || len(loaded.containers) != 1 {
		t.Fatalf("unexpected load: %#v", msg)
	}
}

func TestContainerListSelectedContainerEmpty(t *testing.T) {
	screen := NewContainerListScreen(fakeExecutor{})
	if _, ok := screen.selectedContainer(); ok {
		t.Fatalf("expected no selection")
	}
}

func TestServicesImports(t *testing.T) {
	_ = services.ListContainersBuilder{}
}
