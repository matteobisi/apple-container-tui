package contract

import (
	"reflect"
	"testing"

	"container-tui/src/services"
)

func TestImageListBuilderBuildsCommand(t *testing.T) {
	builder := services.ImageListBuilder{}
	cmd, err := builder.Build()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if cmd.Executable != "container" {
		t.Fatalf("expected executable container, got %q", cmd.Executable)
	}
	if !reflect.DeepEqual(cmd.Args, []string{"image", "list"}) {
		t.Fatalf("unexpected args: %v", cmd.Args)
	}
}
