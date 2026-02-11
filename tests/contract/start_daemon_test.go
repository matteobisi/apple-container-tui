package contract

import (
	"reflect"
	"testing"

	"container-tui/src/services"
)

func TestStartDaemonBuilderBuildsCommand(t *testing.T) {
	builder := services.StartDaemonBuilder{}

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
	if !reflect.DeepEqual(cmd.Args, []string{"system", "start"}) {
		t.Fatalf("expected args [system start], got %v", cmd.Args)
	}
}
