package contract

import (
	"reflect"
	"testing"

	"container-tui/src/services"
)

func TestImageInspectBuilderBuildsCommand(t *testing.T) {
	builder := services.ImageInspectBuilder{ImageReference: "ubuntu:latest"}
	cmd, err := builder.Build()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if cmd.Executable != "container" {
		t.Fatalf("expected executable container, got %q", cmd.Executable)
	}
	if !reflect.DeepEqual(cmd.Args, []string{"image", "inspect", "ubuntu:latest"}) {
		t.Fatalf("unexpected args: %v", cmd.Args)
	}
}
