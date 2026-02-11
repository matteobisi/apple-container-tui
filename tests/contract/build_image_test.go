package contract

import (
	"reflect"
	"testing"

	"container-tui/src/services"
)

func TestBuildImageBuilderRequiresInputs(t *testing.T) {
	builder := services.BuildImageBuilder{}

	if err := builder.Validate(); err == nil {
		t.Fatalf("expected validation error for missing inputs")
	}
}

func TestBuildImageBuilderBuildsCommand(t *testing.T) {
	builder := services.BuildImageBuilder{
		Tag:         "my-image:latest",
		FilePath:    "./Containerfile",
		ContextPath: ".",
	}

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
	if !reflect.DeepEqual(cmd.Args, []string{"build", "-t", "my-image:latest", "-f", "./Containerfile", "."}) {
		t.Fatalf("expected args [build -t my-image:latest -f ./Containerfile .], got %v", cmd.Args)
	}
}
