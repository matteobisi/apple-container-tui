package contract

import (
	"reflect"
	"testing"

	"container-tui/src/services"
)

func TestContainerExecBuilderBuildsCommand(t *testing.T) {
	builder := services.ContainerExecBuilder{ContainerName: "abc123", Shell: "/bin/sh"}
	cmd, err := builder.Build()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if cmd.Executable != "container" {
		t.Fatalf("expected executable container, got %q", cmd.Executable)
	}
	if !reflect.DeepEqual(cmd.Args, []string{"exec", "-it", "abc123", "/bin/sh"}) {
		t.Fatalf("unexpected args: %v", cmd.Args)
	}
}
