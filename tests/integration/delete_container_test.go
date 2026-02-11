//go:build integration

package integration

import (
	"os"
	"os/exec"
	"testing"
)

func TestDeleteContainerWorkflow(t *testing.T) {
	if _, err := exec.LookPath("container"); err != nil {
		t.Skip("container CLI not available")
	}
	if err := exec.Command("container", "system", "version").Run(); err != nil {
		t.Skip("container CLI not responding")
	}

	containerID := os.Getenv("APPLE_TUI_DELETE_CONTAINER_ID")
	if containerID == "" {
		t.Skip("set APPLE_TUI_DELETE_CONTAINER_ID to test delete")
	}

	if err := exec.Command("container", "delete", containerID).Run(); err != nil {
		t.Fatalf("container delete failed: %v", err)
	}
}
