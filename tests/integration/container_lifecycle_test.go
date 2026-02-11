//go:build integration

package integration

import (
	"os"
	"os/exec"
	"testing"
)

func TestContainerLifecycle(t *testing.T) {
	if _, err := exec.LookPath("container"); err != nil {
		t.Skip("container CLI not available")
	}

	if err := exec.Command("container", "system", "version").Run(); err != nil {
		t.Skip("container CLI not responding")
	}

	if err := exec.Command("container", "list", "--all").Run(); err != nil {
		t.Fatalf("container list failed: %v", err)
	}

	containerID := os.Getenv("APPLE_TUI_TEST_CONTAINER_ID")
	if containerID == "" {
		t.Skip("set APPLE_TUI_TEST_CONTAINER_ID to test start/stop")
	}

	if err := exec.Command("container", "start", containerID).Run(); err != nil {
		t.Fatalf("container start failed: %v", err)
	}
	if err := exec.Command("container", "stop", containerID).Run(); err != nil {
		t.Fatalf("container stop failed: %v", err)
	}
}
