package unit

import (
	"testing"

	"container-tui/src/models"
	"container-tui/src/services"
)

type shellDetectorExecutor struct {
	calls int
}

func (e *shellDetectorExecutor) Execute(cmd models.Command) (models.Result, error) {
	e.calls++
	if len(cmd.Args) >= 4 && cmd.Args[0] == "exec" && cmd.Args[3] == "sh" {
		return models.Result{ExitCode: 0, Status: models.ResultSuccess}, nil
	}
	return models.Result{ExitCode: 1, Status: models.ResultError}, nil
}

func TestShellDetectorDetectsAndCachesShell(t *testing.T) {
	exec := &shellDetectorExecutor{}
	detector := services.NewShellDetector(exec)

	shell, err := detector.DetectShell("abc123")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if shell != "sh" {
		t.Fatalf("expected sh, got %q", shell)
	}
	firstCalls := exec.calls

	shell, err = detector.DetectShell("abc123")
	if err != nil {
		t.Fatalf("expected no error on cached call, got %v", err)
	}
	if shell != "sh" {
		t.Fatalf("expected sh from cache, got %q", shell)
	}
	if exec.calls != firstCalls {
		t.Fatalf("expected cached detection with no extra executor calls")
	}
}
