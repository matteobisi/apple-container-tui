package contract

import (
	"strings"
	"testing"

	"container-tui/src/models"
	"container-tui/src/services"
)

func TestRealExecutorRunsCommand(t *testing.T) {
	executor := services.RealExecutor{}
	cmd := models.Command{Executable: "echo", Args: []string{"hello"}}

	result, err := executor.Execute(cmd)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.ExitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", result.ExitCode)
	}
	if result.Status != models.ResultSuccess {
		t.Fatalf("expected success status, got %s", result.Status)
	}
	if !strings.Contains(result.Stdout, "hello") {
		t.Fatalf("expected stdout to contain output")
	}
}
