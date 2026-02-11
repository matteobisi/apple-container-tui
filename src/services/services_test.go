package services

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"container-tui/src/models"
)

type stubExecutor struct {
	result models.Result
	err    error
}

func (s stubExecutor) Execute(cmd models.Command) (models.Result, error) {
	return s.result, s.err
}

func jsonMarshal(entry LogEntry) ([]byte, error) {
	data, err := json.Marshal(entry)
	if err != nil {
		return nil, err
	}
	return append(data, '\n'), nil
}

func TestNormalizeRequiredToken(t *testing.T) {
	if _, err := normalizeRequiredToken("", "name"); err == nil {
		t.Fatalf("expected error for empty value")
	}
	if _, err := normalizeRequiredToken("bad value", "name"); err == nil {
		t.Fatalf("expected error for whitespace")
	}
	if value, err := normalizeRequiredToken("ok", "name"); err != nil || value != "ok" {
		t.Fatalf("unexpected result: %v %q", err, value)
	}
}

func TestBuilders(t *testing.T) {
	start := StartContainerBuilder{ContainerID: "abc"}
	cmd, err := start.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cmd.Args[0] != "start" {
		t.Fatalf("unexpected args: %v", cmd.Args)
	}
	if _, err := start.Build(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := (StartContainerBuilder{ContainerID: ""}).Build(); err == nil {
		t.Fatalf("expected error for missing container id")
	}

	pull := PullImageBuilder{Reference: "nginx:latest"}
	cmd, err = pull.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cmd.Args[0] != "image" || cmd.Args[1] != "pull" {
		t.Fatalf("unexpected args: %v", cmd.Args)
	}

	stop := StopContainerBuilder{ContainerID: "abc"}
	cmd, err = stop.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cmd.Args[0] != "stop" {
		t.Fatalf("unexpected args: %v", cmd.Args)
	}
	if _, err := (StopContainerBuilder{ContainerID: "  "}).Build(); err == nil {
		t.Fatalf("expected error for blank container id")
	}

	deleteCmd, err := DeleteContainerBuilder{ContainerID: "abc"}.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if deleteCmd.Args[0] != "delete" {
		t.Fatalf("unexpected args: %v", deleteCmd.Args)
	}
	if _, err := (DeleteContainerBuilder{ContainerID: ""}).Build(); err == nil {
		t.Fatalf("expected error for missing container id")
	}

	listCmd, err := ListContainersBuilder{}.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if listCmd.Args[0] != "list" {
		t.Fatalf("unexpected args: %v", listCmd.Args)
	}
	if err := (ListContainersBuilder{}).Validate(); err != nil {
		t.Fatalf("unexpected validation error: %v", err)
	}

	startDaemonCmd, err := StartDaemonBuilder{}.Build()
	if err != nil || startDaemonCmd.Args[0] != "system" {
		t.Fatalf("unexpected daemon start: %v %v", startDaemonCmd, err)
	}
	stopDaemonCmd, err := StopDaemonBuilder{}.Build()
	if err != nil || stopDaemonCmd.Args[1] != "stop" {
		t.Fatalf("unexpected daemon stop: %v %v", stopDaemonCmd, err)
	}
	statusCmd, err := CheckDaemonStatusBuilder{}.Build()
	if err != nil || statusCmd.Args[1] != "status" {
		t.Fatalf("unexpected daemon status: %v %v", statusCmd, err)
	}
}

func TestBuildImageBuilderTrims(t *testing.T) {
	builder := BuildImageBuilder{Tag: " tag ", FilePath: " ./Containerfile ", ContextPath: " . "}
	cmd, err := builder.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cmd.Args[2] != "tag" || cmd.Args[4] != "./Containerfile" {
		t.Fatalf("unexpected args: %v", cmd.Args)
	}
}

func TestBuildImageBuilderErrors(t *testing.T) {
	builder := BuildImageBuilder{Tag: "", FilePath: "file", ContextPath: "."}
	if err := builder.Validate(); err == nil {
		t.Fatalf("expected error for missing tag")
	}
	builder = BuildImageBuilder{Tag: "tag", FilePath: "", ContextPath: "."}
	if err := builder.Validate(); err == nil {
		t.Fatalf("expected error for missing file path")
	}
	builder = BuildImageBuilder{Tag: "tag", FilePath: "file", ContextPath: ""}
	if err := builder.Validate(); err == nil {
		t.Fatalf("expected error for missing context")
	}
}

func TestPullImageBuilderErrors(t *testing.T) {
	builder := PullImageBuilder{Reference: ""}
	if err := builder.Validate(); err == nil {
		t.Fatalf("expected error for missing reference")
	}
}

func TestParseContainerList(t *testing.T) {
	output := "CONTAINER ID  IMAGE  COMMAND  CREATED  STATUS  PORTS\n" +
		"abc123        nginx  web     1m      running  0.0.0.0:8080->80/tcp\n"
	containers, err := ParseContainerList(output)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(containers) != 1 {
		t.Fatalf("expected one container")
	}
	if containers[0].ID != "abc123" {
		t.Fatalf("unexpected id: %s", containers[0].ID)
	}

	altOutput := "ID  IMAGE  OS  ARCH  STATE  ADDR  CPUS  MEMORY  STARTED\n" +
		"def456  docker.io/library/ubuntu:latest  linux  arm64  running  192.168.64.2/24  4  1024MB  2026-02-11T12:05:33Z\n"
	containers, err = ParseContainerList(altOutput)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(containers) != 1 {
		t.Fatalf("expected one container")
	}
	if containers[0].ID != "def456" {
		t.Fatalf("unexpected id: %s", containers[0].ID)
	}

	containers, err = ParseContainerList("")
	if err != nil || len(containers) != 0 {
		t.Fatalf("expected empty result")
	}
}

func TestParserHelpers(t *testing.T) {
	_, err := headerIndices("MISSING", []string{"CONTAINER ID"})
	if err == nil {
		t.Fatalf("expected error for missing header")
	}

	columns := sliceColumns("abc   def", []int{0, 6})
	if len(columns) != 2 || columns[0] != "abc" {
		t.Fatalf("unexpected columns: %v", columns)
	}

	ports := parsePortMappings("0.0.0.0:8080->80/tcp")
	if len(ports) != 1 || ports[0].HostPort != 8080 {
		t.Fatalf("unexpected ports: %v", ports)
	}
}

func TestParseDaemonStatus(t *testing.T) {
	status := ParseDaemonStatus("not running")
	if status.Running {
		t.Fatalf("expected not running")
	}
	status = ParseDaemonStatus("running")
	if !status.Running {
		t.Fatalf("expected running")
	}
}

func TestDetectBuildFile(t *testing.T) {
	dir := t.TempDir()
	containerfile := filepath.Join(dir, "Containerfile")
	if err := os.WriteFile(containerfile, []byte("FROM scratch"), 0o600); err != nil {
		t.Fatalf("write file: %v", err)
	}
	source, err := DetectBuildFile(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if source.FileType != models.BuildFileTypeContainerfile {
		t.Fatalf("unexpected file type: %s", source.FileType)
	}
}

func TestDetectBuildFileMissing(t *testing.T) {
	_, err := DetectBuildFile(t.TempDir())
	if err == nil {
		t.Fatalf("expected error for missing build file")
	}
}

func TestFormatError(t *testing.T) {
	message := FormatError(nil, "unauthorized")
	if !strings.Contains(message, "authentication") {
		t.Fatalf("unexpected message: %s", message)
	}
	message = FormatError(nil, "daemon connection failed")
	if !strings.Contains(message, "daemon") {
		t.Fatalf("unexpected message: %s", message)
	}
	message = FormatError(nil, "")
	if !strings.Contains(message, "unknown") {
		t.Fatalf("unexpected message: %s", message)
	}
}

func TestDestructiveActionMetadata(t *testing.T) {
	metadata := DestructiveActionMetadata()
	if metadata[ActionDeleteContainer].Label == "" {
		t.Fatalf("expected metadata")
	}
}

func TestIsExactMatch(t *testing.T) {
	if !IsExactMatch("a", "a") {
		t.Fatalf("expected exact match")
	}
}

func TestConfigManagerLoad(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	configPath := filepath.Join(home, ".config", "apple-tui", "config")
	if err := os.MkdirAll(filepath.Dir(configPath), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	content := []byte("theme_mode = \"dark\"\n")
	if err := os.WriteFile(configPath, content, 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}
	manager, err := NewConfigManager()
	if err != nil {
		t.Fatalf("manager: %v", err)
	}
	config, used, err := manager.Load()
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if used == "" || config.ThemeMode != "dark" {
		t.Fatalf("unexpected config: %s", used)
	}

	os.Remove(configPath)
	config, used, err = manager.Load()
	if err != nil {
		t.Fatalf("load defaults: %v", err)
	}
	if used != "" || config.DefaultBuildFile == "" {
		t.Fatalf("expected defaults")
	}
}

func TestLogWriterWrite(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	writer, err := NewLogWriter(0)
	if err != nil {
		t.Fatalf("log writer: %v", err)
	}
	entry := BuildLogEntry(models.Command{Executable: "container"}, models.Result{ExitCode: 0, Status: models.ResultSuccess, Duration: time.Second}, false)
	if err := writer.Write(entry); err != nil {
		t.Fatalf("write: %v", err)
	}
	if _, err := os.Stat(filepath.Join(home, "Library", "Application Support", "apple-tui", "command.log")); err != nil {
		t.Fatalf("expected log file: %v", err)
	}
}

func TestLogWriterRotate(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	writer, err := NewLogWriter(1)
	if err != nil {
		t.Fatalf("log writer: %v", err)
	}
	logPath := filepath.Join(home, "Library", "Application Support", "apple-tui", "command.log")
	if err := os.MkdirAll(filepath.Dir(logPath), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	oldEntry := LogEntry{Command: "container list", StartTime: time.Now().AddDate(0, 0, -10)}
	data, err := jsonMarshal(oldEntry)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if err := os.WriteFile(logPath, data, 0o600); err != nil {
		t.Fatalf("write log: %v", err)
	}
	if err := writer.rotateIfNeeded(); err != nil {
		t.Fatalf("rotate: %v", err)
	}
}

func TestDryRunExecutor(t *testing.T) {
	exec := DryRunExecutor{}
	result, err := exec.Execute(models.Command{Executable: "container"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Status != models.ResultSuccess {
		t.Fatalf("unexpected status: %s", result.Status)
	}
}

func TestRealExecutor(t *testing.T) {
	exec := RealExecutor{}
	result, err := exec.Execute(models.Command{Executable: "/bin/echo", Args: []string{"ok"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ExitCode != 0 {
		t.Fatalf("expected exit code 0")
	}

	result, err = exec.Execute(models.Command{Executable: "/usr/bin/false"})
	if err == nil {
		t.Fatalf("expected error for false")
	}
	if result.Status != models.ResultError {
		t.Fatalf("expected error status")
	}
}

func TestLoggingExecutor(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	writer, err := NewLogWriter(0)
	if err != nil {
		t.Fatalf("log writer: %v", err)
	}
	delegate := stubExecutor{result: models.Result{ExitCode: 0, Status: models.ResultSuccess}}
	executor := NewLoggingExecutor(delegate, writer, false)
	if _, err := executor.Execute(models.Command{Executable: "container"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	noWriter := NewLoggingExecutor(delegate, nil, false)
	if _, err := noWriter.Execute(models.Command{Executable: "container"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCheckCLIWithStub(t *testing.T) {
	dir := t.TempDir()
	binPath := filepath.Join(dir, "container")
	content := "#!/bin/sh\nif [ \"$1\" = \"system\" ] && [ \"$2\" = \"version\" ]; then\n  exit 0\nfi\nexit 1\n"
	if err := os.WriteFile(binPath, []byte(content), 0o700); err != nil {
		t.Fatalf("write stub: %v", err)
	}
	oldPath := os.Getenv("PATH")
	t.Setenv("PATH", dir+string(os.PathListSeparator)+oldPath)
	if err := CheckCLI(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCheckCLIFailure(t *testing.T) {
	oldPath := os.Getenv("PATH")
	t.Setenv("PATH", "")
	if err := CheckCLI(context.Background()); err == nil {
		t.Fatalf("expected error")
	}
	t.Setenv("PATH", oldPath)
}
