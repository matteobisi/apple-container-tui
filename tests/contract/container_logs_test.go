package contract

import (
	"reflect"
	"testing"

	"container-tui/src/services"
)

func TestContainerLogsBuilderBuildsCommand(t *testing.T) {
	builder := services.ContainerLogsBuilder{ContainerName: "abc123"}
	cmd, err := builder.Build()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if cmd.Executable != "container" {
		t.Fatalf("expected executable container, got %q", cmd.Executable)
	}
	if !reflect.DeepEqual(cmd.Args, []string{"logs", "-f", "abc123"}) {
		t.Fatalf("unexpected args: %v", cmd.Args)
	}
}
