package unit

import (
	"testing"

	"container-tui/src/services"
)

func TestDestructiveActionMetadata(t *testing.T) {
	metadata := services.DestructiveActionMetadata()
	if len(metadata) == 0 {
		t.Fatalf("expected destructive action metadata to be non-empty")
	}
	if _, ok := metadata[services.ActionDeleteContainer]; !ok {
		t.Fatalf("expected delete container metadata")
	}
	if _, ok := metadata[services.ActionStopDaemon]; !ok {
		t.Fatalf("expected stop daemon metadata")
	}
}
