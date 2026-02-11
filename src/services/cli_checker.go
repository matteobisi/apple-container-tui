package services

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
)

// CheckCLI verifies the Apple Container CLI is available and responding.
func CheckCLI(ctx context.Context) error {
	if _, err := exec.LookPath("container"); err != nil {
		return fmt.Errorf("apple container CLI not found in PATH. Please install from https://github.com/apple/container")
	}

	command := exec.CommandContext(ctx, "container", "system", "version")
	var stderr bytes.Buffer
	command.Stderr = &stderr

	if err := command.Run(); err != nil {
		message := strings.TrimSpace(stderr.String())
		if message == "" {
			message = err.Error()
		}
		return fmt.Errorf("failed to verify Apple Container CLI: %s", message)
	}

	return nil
}
