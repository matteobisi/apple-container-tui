package contract

import (
	"strings"
	"testing"

	"container-tui/src/models"
	"container-tui/src/services"
)

func TestDryRunExecutorDoesNotExecute(t *testing.T) {
	executor := services.DryRunExecutor{}
	cmd := models.Command{Executable: "container", Args: []string{"list", "--all"}}

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
	if !strings.Contains(result.Stdout, cmd.String()) {
		t.Fatalf("expected stdout to include command string")
	}
}
