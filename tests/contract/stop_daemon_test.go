package contract

import (
	"reflect"
	"testing"

	"container-tui/src/services"
)

func TestStopDaemonBuilderBuildsCommand(t *testing.T) {
	builder := services.StopDaemonBuilder{}

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
	if !reflect.DeepEqual(cmd.Args, []string{"system", "stop"}) {
		t.Fatalf("expected args [system stop], got %v", cmd.Args)
	}
}
