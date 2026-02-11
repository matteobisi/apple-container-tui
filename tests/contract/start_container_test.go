package contract

import (
	"reflect"
	"testing"

	"container-tui/src/services"
)

func TestStartContainerBuilderRequiresID(t *testing.T) {
	builder := services.StartContainerBuilder{}

	if err := builder.Validate(); err == nil {
		t.Fatalf("expected validation error for empty container ID")
	}
}

func TestStartContainerBuilderBuildsCommand(t *testing.T) {
	builder := services.StartContainerBuilder{ContainerID: "abc123"}

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
	if !reflect.DeepEqual(cmd.Args, []string{"start", "abc123"}) {
		t.Fatalf("expected args [start abc123], got %v", cmd.Args)
	}
}
