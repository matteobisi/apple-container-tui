package contract

import (
	"reflect"
	"testing"

	"container-tui/src/services"
)

func TestPullImageBuilderRequiresReference(t *testing.T) {
	builder := services.PullImageBuilder{}

	if err := builder.Validate(); err == nil {
		t.Fatalf("expected validation error for empty reference")
	}
}

func TestPullImageBuilderBuildsCommand(t *testing.T) {
	builder := services.PullImageBuilder{Reference: "nginx:latest"}

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
	if !reflect.DeepEqual(cmd.Args, []string{"image", "pull", "nginx:latest"}) {
		t.Fatalf("expected args [image pull nginx:latest], got %v", cmd.Args)
	}
}
