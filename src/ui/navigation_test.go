package ui

import (
	"testing"

	"container-tui/src/models"
)

type navExecutor struct{}

func (navExecutor) Execute(cmd models.Command) (models.Result, error) {
	if len(cmd.Args) > 0 && cmd.Args[0] == "list" {
		return models.Result{Stdout: "CONTAINER ID  IMAGE  COMMAND  CREATED  STATUS  PORTS\nabc nginx web 1m running", Status: models.ResultSuccess}, nil
	}
	return models.Result{Status: models.ResultSuccess, ExitCode: 0}, nil
}

func TestNavigationStackPushPop(t *testing.T) {
	app := NewAppModel(navExecutor{}, "1.0.0")
	app.pushView(ScreenContainerList)
	app.pushView(ScreenImageList)

	got := app.popView(ScreenHelp)
	if got != ScreenImageList {
		t.Fatalf("expected ScreenImageList, got %s", got)
	}
	got = app.popView(ScreenHelp)
	if got != ScreenContainerList {
		t.Fatalf("expected ScreenContainerList, got %s", got)
	}
	got = app.popView(ScreenHelp)
	if got != ScreenHelp {
		t.Fatalf("expected fallback ScreenHelp, got %s", got)
	}
}
