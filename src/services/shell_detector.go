package services

import (
	"fmt"

	"container-tui/src/models"
)

// ShellDetector detects an available shell for a given container.
type ShellDetector struct {
	executor CommandExecutor
	cache    map[string]string
}

// NewShellDetector creates a ShellDetector.
func NewShellDetector(executor CommandExecutor) *ShellDetector {
	return &ShellDetector{executor: executor, cache: map[string]string{}}
}

// DetectShell probes shells in order: bash, sh, /bin/sh, /bin/bash, ash.
func (d *ShellDetector) DetectShell(containerName string) (string, error) {
	if cached, ok := d.cache[containerName]; ok && cached != "" {
		return cached, nil
	}

	type probe struct {
		shell string
		args  []string
	}
	probes := []probe{
		{shell: "bash", args: []string{"exec", containerName, "which", "bash"}},
		{shell: "sh", args: []string{"exec", containerName, "which", "sh"}},
		{shell: "/bin/sh", args: []string{"exec", containerName, "test", "-x", "/bin/sh"}},
		{shell: "/bin/bash", args: []string{"exec", containerName, "test", "-x", "/bin/bash"}},
		{shell: "ash", args: []string{"exec", containerName, "which", "ash"}},
	}

	for _, probe := range probes {
		cmd := models.Command{Executable: "container", Args: probe.args}
		result, err := d.executor.Execute(cmd)
		if err == nil && result.ExitCode == 0 {
			d.cache[containerName] = probe.shell
			return probe.shell, nil
		}
	}

	return "", fmt.Errorf("no supported shell found in container %q", containerName)
}
