package contract

import (
	"reflect"
	"testing"

	"container-tui/src/services"
)

func TestDeleteContainerBuilderBuildsCommand(t *testing.T) {
	builder := services.DeleteContainerBuilder{ContainerID: "abc123"}

	if err := builder.Validate(); err != nil {
		t.Fatalf("expected no validation error, got %v", err)
	}

	cmd, err := builder.Build()
	if err != nil {
		t.Fatalf("expected no build error, got %v", err)
	}

	if cmd.Executable != "container" {
		t.Fatalf("expected executable 'container', got %q", cmd.Executable)
	}
	if !reflect.DeepEqual(cmd.Args, []string{"delete", "abc123"}) {
		t.Fatalf("expected args [delete abc123], got %v", cmd.Args)
	}
}
