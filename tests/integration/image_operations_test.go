//go:build integration

package integration

import (
	"os/exec"
	"testing"
)

func TestImageOperations(t *testing.T) {
	if _, err := exec.LookPath("container"); err != nil {
		t.Skip("container CLI not available")
	}

	if err := exec.Command("container", "system", "version").Run(); err != nil {
		t.Skip("container CLI not responding")
	}
}
