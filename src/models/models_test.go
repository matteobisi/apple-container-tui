package models

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestCommandStringQuotes(t *testing.T) {
	command := Command{Executable: "container", Args: []string{"build", "-f", "Container file"}}
	result := command.String()
	if !strings.Contains(result, "\"Container file\"") {
		t.Fatalf("expected quoted argument, got %q", result)
	}
}

func TestCommandStringEmptyArg(t *testing.T) {
	command := Command{Executable: "container", Args: []string{"run", ""}}
	result := command.String()
	if !strings.Contains(result, "\"\"") {
		t.Fatalf("expected empty argument to be quoted, got %q", result)
	}
}

func TestContainerValidate(t *testing.T) {
	container := Container{ID: "", Name: "web", Status: ContainerStatusRunning}
	if err := container.Validate(); err == nil {
		t.Fatalf("expected error for missing id")
	}
	valid := Container{ID: "abc", Name: "web", Status: ContainerStatusStopped}
	if err := valid.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPortMappingValidate(t *testing.T) {
	mapping := PortMapping{HostPort: 0, ContainerPort: 80, Protocol: "tcp"}
	if err := mapping.Validate(); err == nil {
		t.Fatalf("expected error for host port")
	}
	mapping = PortMapping{HostPort: 8080, ContainerPort: 0, Protocol: "tcp"}
	if err := mapping.Validate(); err == nil {
		t.Fatalf("expected error for container port")
	}
	mapping = PortMapping{HostPort: 8080, ContainerPort: 80, Protocol: ""}
	if err := mapping.Validate(); err == nil {
		t.Fatalf("expected error for protocol")
	}
	valid := PortMapping{HostPort: 8080, ContainerPort: 80, Protocol: "tcp"}
	if err := valid.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestImageReferenceString(t *testing.T) {
	ref := ImageReference{Registry: "docker.io", Repository: "library/alpine", Tag: "latest"}
	if ref.String() != "docker.io/library/alpine:latest" {
		t.Fatalf("unexpected reference: %s", ref.String())
	}
	ref = ImageReference{Repository: "alpine", Digest: "sha256:abc"}
	if ref.String() != "alpine@sha256:abc" {
		t.Fatalf("unexpected digest reference: %s", ref.String())
	}
}

func TestImageReferenceValidate(t *testing.T) {
	ref := ImageReference{Repository: ""}
	if err := ref.Validate(); err == nil {
		t.Fatalf("expected error for missing repository")
	}
}

func TestBuildSourceValidateErrors(t *testing.T) {
	source := BuildSource{FilePath: "", FileType: BuildFileTypeContainerfile, WorkingDirectory: ""}
	if err := source.Validate(); err == nil {
		t.Fatalf("expected error for missing file path")
	}
}

func TestBuildSourceValidateAndExists(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "Containerfile")
	if err := os.WriteFile(filePath, []byte("FROM scratch"), 0o600); err != nil {
		t.Fatalf("write file: %v", err)
	}
	source := BuildSource{FilePath: filePath, FileType: BuildFileTypeContainerfile, WorkingDirectory: dir}
	if err := source.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	withExists := source.WithComputedExists()
	if !withExists.Exists {
		t.Fatalf("expected Exists true")
	}
}

func TestBuildSourceWithComputedExistsMissing(t *testing.T) {
	source := BuildSource{FilePath: filepath.Join(t.TempDir(), "Missingfile")}
	withExists := source.WithComputedExists()
	if withExists.Exists {
		t.Fatalf("expected Exists false")
	}
}

func TestDaemonStatus(t *testing.T) {
	status := DaemonStatus{Running: true, Version: "1.0", LastChecked: time.Now()}
	if !status.Running {
		t.Fatalf("expected running status")
	}
}

func TestDefaultUserConfig(t *testing.T) {
	config := DefaultUserConfig()
	if config.DefaultBuildFile == "" {
		t.Fatalf("expected default build file")
	}
}
