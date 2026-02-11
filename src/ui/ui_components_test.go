package ui

import (
	"strings"
	"testing"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"

	"container-tui/src/models"
)

func TestRenderStatusBar(t *testing.T) {
	result := RenderStatusBar(20, "left", "right")
	if !strings.Contains(result, "left") {
		t.Fatalf("unexpected status bar: %q", result)
	}
	result = RenderStatusBar(0, "left", "right")
	if result == "" {
		t.Fatalf("expected status bar output")
	}
}

func TestRenderResult(t *testing.T) {
	result := RenderResult(models.Result{Status: models.ResultSuccess, Stdout: "ok", Stderr: ""})
	if !strings.Contains(result, "Status") {
		t.Fatalf("unexpected result: %q", result)
	}
}

func TestCommandPreviewView(t *testing.T) {
	modal := CommandPreviewModal{Title: "Preview", Command: models.Command{Executable: "container"}}
	view := modal.View()
	if !strings.Contains(view, "container") {
		t.Fatalf("unexpected view: %q", view)
	}
}

func TestProgressModel(t *testing.T) {
	model := NewProgressModel()
	model.SetPercent(0.5)
	view := model.View(10)
	if view == "" {
		t.Fatalf("expected progress view")
	}
}

func TestSpinnerModel(t *testing.T) {
	model := NewSpinnerModel()
	model.SetActive(true)
	updated, _ := model.Update(spinner.TickMsg{})
	view := updated.View()
	if view == "" {
		t.Fatalf("expected spinner view")
	}
	model.SetActive(false)
	updated, _ = model.Update(spinner.TickMsg{})
	if updated.View() != "" {
		t.Fatalf("expected empty view")
	}
}

func TestThemeHelpers(t *testing.T) {
	ApplyTheme("dark")
	if RenderTitle("Title") == "" {
		t.Fatalf("expected render output")
	}
	if RenderError("err") == "" {
		t.Fatalf("expected render output")
	}
}

func TestYesNoConfirmModal(t *testing.T) {
	modal := YesNoConfirmModal{}
	confirmed, canceled := modal.Handle(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
	if !confirmed || canceled {
		t.Fatalf("expected confirmation")
	}
}

func TestTypeToConfirmModal(t *testing.T) {
	modal := NewTypeToConfirmModal("Delete", "match", models.Command{Executable: "container"})
	modal.input.SetValue("nope")
	_, confirmed, _ := modal.Handle(tea.KeyMsg{Type: tea.KeyEnter})
	if confirmed {
		t.Fatalf("expected not confirmed")
	}
	modal.input.SetValue("match")
	_, confirmed, _ = modal.Handle(tea.KeyMsg{Type: tea.KeyEnter})
	if !confirmed {
		t.Fatalf("expected confirmation")
	}
}

func TestHelpScreenView(t *testing.T) {
	help := HelpScreen{Version: "1.0.0"}
	view := help.View()
	if !strings.Contains(view, "1.0.0") {
		t.Fatalf("expected version in view")
	}
}

func TestDefaultKeyMap(t *testing.T) {
	keys := DefaultKeyMap()
	if len(keys.Quit.Keys()) == 0 {
		t.Fatalf("expected quit keys")
	}
}
