# CLI Command Contracts: Apple Container TUI

**Phase 1 Output** | **Date**: 2026-02-11  
**Purpose**: Define mappings between TUI actions and Apple Container CLI commands

## Overview

This document specifies the exact Apple Container CLI commands that each TUI action must generate. All commands follow the `container <subcommand> [options] [arguments]` pattern and are executed via Go's `os/exec.Command` for safe argument handling.

## Command Mapping Table

| TUI Action | CLI Command | User Story | Required Confirmation |
|------------|-------------|------------|----------------------|
| List Containers | `container list --all` | US1 | No |
| Start Container | `container start <id>` | US1 | No |
| Stop Container | `container stop <id>` | US1 | No |
| Delete Container | `container delete <id>` | US3 | Yes (type-to-confirm) |
| Pull Image | `container image pull <reference>` | US2 | No |
| Build from File | `container build -t <tag> -f <file> <context>` | US2 | No |
| Start Daemon | `container system start` | US3 | Yes (y/n confirm) |
| Stop Daemon | `container system stop` | US3 | Yes (y/n confirm) |
| Check Daemon Status | `container system status` | All | No |
| Refresh Container List | `container list --all` | All | No |
| Check CLI Availability | `container system version` | Startup | No |

## Command Specifications

### List Containers

**Command**: `container list --all`

**Purpose**: Retrieve all containers (running and stopped) for the container list view

**Arguments**: None

**Flags**:
- `--all`: Include stopped containers (default only shows running)

**Expected Output** (table format):
```
CONTAINER ID  IMAGE          COMMAND  CREATED       STATUS     PORTS
abc123        nginx:latest   nginx    2 hours ago   running    0.0.0.0:8080->80/tcp
def456        alpine:latest  sh       1 day ago     stopped
```

**Parsing Strategy**:
- Table format with headers
- Parse columns: CONTAINER ID, IMAGE, COMMAND, CREATED, STATUS, PORTS
- Map to Container entity attributes

**Error Cases**:
- Daemon not running → Exit code 1, stderr contains "daemon" or "connection"
- No containers exist → Empty table (no error)

---

### Start Container

**Command**: `container start <container-id>`

**Purpose**: Start a stopped container

**Arguments**:
- `<container-id>` (required): Container ID or name from Container.ID

**Flags**: None (for basic start; -a/-i available but not used in TUI)

**Expected Output**:
```
abc123
```
(Container ID echoed on success)

**Exit Codes**:
- 0: Success
- Non-zero: Error (container not found, already running, etc.)

**Error Cases**:
- Container not found → stderr contains "not found" or "no such container"
- Container already running → stderr contains "already" or "running"
- Daemon not running → stderr contains "daemon" or "connection"

**Command Builder Validation**:
- `container-id` must not be empty
- `container-id` must match a known Container.ID

---

### Stop Container

**Command**: `container stop <container-id>`

**Purpose**: Stop a running container gracefully (SIGTERM then SIGKILL after timeout)

**Arguments**:
- `<container-id>` (required): Container ID or name

**Flags**:
- `--time 5` (default, not specified): Wait 5 seconds before SIGKILL

**Expected Output**:
```
abc123
```

**Exit Codes**:
- 0: Success
- Non-zero: Error

**Error Cases**:
- Container not found → stderr contains "not found"
- Container already stopped → stderr contains "already stopped" or "not running"

**Command Builder Validation**:
- `container-id` must not be empty
- `container-id` must match a known Container.ID

---

### Delete Container

**Command**: `container delete <container-id>`

**Purpose**: Permanently delete a stopped container

**Arguments**:
- `<container-id>` (required): Container ID or name

**Flags**: None (no `--force` - only delete stopped containers per constitution)

**Expected Output**:
```
abc123
```

**Exit Codes**:
- 0: Success
- Non-zero: Error

**Error Cases**:
- Container not found → stderr contains "not found"
- Container still running → stderr contains "running" or "cannot remove"

**Destructive Action**: YES

**Confirmation Required**: Type-to-confirm (user must type container name/ID)

**Command Builder Validation**:
- `container-id` must not be empty
- Container.Status must be `stopped` (pre-validation before command)
- User must type exact container ID or name to confirm

**Pre-flight Check**: Verify container is stopped before allowing deletion

---

### Pull Image

**Command**: `container image pull <reference>`

**Purpose**: Pull a container image from a registry

**Arguments**:
- `<reference>` (required): Full image reference (e.g., `nginx:latest`, `docker.io/library/alpine:3.18`)

**Flags**:
- `--progress ansi` (default): Show ANSI progress bars

**Expected Output** (progress with ANSI codes):
```
Pulling nginx:latest
[====================>] 100%
sha256:abc123...
```

**Exit Codes**:
- 0: Success
- Non-zero: Error (network error, not found, auth required, etc.)

**Error Cases**:
- Image not found → stderr contains "not found" or "manifest unknown"
- Network error → stderr contains "connection" or "timeout"
- Authentication required → stderr contains "unauthorized" or "authentication"

**Command Builder Validation**:
- `reference` must not be empty
- `reference` should match ImageReference format

**Streaming**: Progress output should be displayed in real-time (buffered or streamed to UI)

---

### Build from File

**Command**: `container build -t <tag> -f <file> <context-dir>`

**Purpose**: Build an OCI image from a Containerfile or Dockerfile

**Arguments**:
- `<context-dir>` (required): Build context directory (usually directory containing build file)

**Flags**:
- `-t, --tag <tag>` (required): Tag for the built image
- `-f, --file <file>` (required): Path to Containerfile or Dockerfile

**Expected Output** (build progress):
```
Step 1/5 : FROM alpine:latest
 ---> abc123def456
Step 2/5 : RUN apk add curl
 ---> Running in xyz789...
...
Successfully built abc123
Successfully tagged my-image:latest
```

**Exit Codes**:
- 0: Success
- Non-zero: Build error

**Error Cases**:
- File not found → stderr contains "no such file" or "not found"
- Build error → stderr contains build step errors
- Builder not running → stderr contains "builder" or "buildkit"

**Command Builder Validation**:
- `tag` must not be empty
- `file` must exist (BuildSource.Exists must be true)
- `context-dir` must exist and be a directory

**File Type Auto-detection**:
Per clarification 2026-02-11, check for `Containerfile` first, then `Dockerfile` if not found.

**Streaming**: Build output should be streamed to UI in real-time

---

### Start Daemon

**Command**: `container system start`

**Purpose**: Start the Apple Container system services

**Arguments**: None

**Flags**: 
- May prompt for kernel install (use `--disable-kernel-install` to suppress if needed)

**Expected Output**:
```
Starting container services...
Services started successfully
```

**Exit Codes**:
- 0: Success
- Non-zero: Error (already running, permission denied, etc.)

**Error Cases**:
- Already running → stderr contains "already running" or "already started"
- Permission denied → stderr contains "permission" or "sudo"

**Destructive Action**: NO, but impactful

**Confirmation Required**: Yes (y/n confirmation)

**Command Builder Validation**: None (no arguments)

---

### Stop Daemon

**Command**: `container system stop`

**Purpose**: Stop the Apple Container system services

**Arguments**: None

**Flags**: None

**Expected Output**:
```
Stopping container services...
Services stopped successfully
```

**Exit Codes**:
- 0: Success
- Non-zero: Error (not running, permission error, etc.)

**Error Cases**:
- Not running → stderr contains "not running" or "not started"
- Permission denied → stderr may contain "permission"

**Destructive Action**: YES (stops all containers)

**Confirmation Required**: Yes (y/n confirmation with warning about stopping all containers)

**Command Builder Validation**: None (no arguments)

---

### Check Daemon Status

**Command**: `container system status`

**Purpose**: Check if the Apple Container daemon services are running

**Arguments**: None

**Flags**: None

**Expected Output**:
```
Container services are running
API Server: Reachable
```

**Exit Codes**:
- 0: Running
- Non-zero: Not running or unreachable

**Parsing Strategy**:
- Parse for "running" vs "not running"
- Map to DaemonStatus.Running boolean

**Command Builder Validation**: None (no arguments)

---

### Refresh Container List

**Command**: `container list --all`

**Purpose**: Same as List Containers, but explicitly called when user requests manual refresh

**See**: List Containers above

**Manual Trigger**: User presses refresh key (e.g., `r` key)

---

### Check CLI Availability

**Command**: `container system version`

**Purpose**: Verify Apple Container CLI is installed and accessible in PATH

**Arguments**: None

**Flags**: None

**Expected Output** (table format):
```
COMPONENT   VERSION  BUILD   COMMIT
CLI         1.2.3    release abc123
```

**Exit Codes**:
- 0: CLI available
- 127: Command not found (not in PATH)

**Startup Check** (per FR-014):
- Run at application startup
- If exit code 127 or command not found, show error and EXIT
- Error message: "Apple Container CLI not found in PATH. Please install from https://github.com/apple/container"

**Command Builder Validation**: None (no arguments)

---

## Command Builder Interface (Go)

### Interface Definition

```go
package commands

type Command struct {
    Executable string
    Args       []string
}

func (c Command) String() string {
    // Returns formatted command string for preview
    // Example: "container start abc123"
}

type CommandBuilder interface {
    Validate() error
    Build() (Command, error)
}
```

### Example Implementations

#### Start Container Builder

```go
type StartContainerBuilder struct {
    containerID string
}

func NewStartContainer(containerID string) *StartContainerBuilder {
    return &StartContainerBuilder{containerID: containerID}
}

func (b *StartContainerBuilder) Validate() error {
    if b.containerID == "" {
        return errors.New("container ID is required")
    }
    return nil
}

func (b *StartContainerBuilder) Build() (Command, error) {
    if err := b.Validate(); err != nil {
        return Command{}, err
    }
    return Command{
        Executable: "container",
        Args:       []string{"start", b.containerID},
    }, nil
}
```

#### Pull Image Builder

```go
type PullImageBuilder struct {
    reference string
}

func NewPullImage(reference string) *PullImageBuilder {
    return &PullImageBuilder{reference: reference}
}

func (b *PullImageBuilder) Validate() error {
    if b.reference == "" {
        return errors.New("image reference is required")
    }
    // Optional: validate reference format
    return nil
}

func (b *PullImageBuilder) Build() (Command, error) {
    if err := b.Validate(); err != nil {
        return Command{}, err
    }
    return Command{
        Executable: "container",
        Args:       []string{"image", "pull", "--progress", "ansi", b.reference},
    }, nil
}
```

## Execution Pattern

### Standard Execution (with os/exec)

```go
func (e *RealExecutor) Execute(cmd Command) (*Result, error) {
    execCmd := exec.Command(cmd.Executable, cmd.Args...)
    
    var stdout, stderr bytes.Buffer
    execCmd.Stdout = &stdout
    execCmd.Stderr = &stderr
    
    startTime := time.Now()
    err := execCmd.Run()
    duration := time.Since(startTime)
    
    exitCode := 0
    if exitErr, ok := err.(*exec.ExitError); ok {
        exitCode = exitErr.ExitCode()
    }
    
    return &Result{
        Command:   cmd.String(),
        DryRun:    false,
        ExitCode:  exitCode,
        Stdout:    stdout.String(),
        Stderr:    stderr.String(),
        StartTime: startTime,
        Duration:  duration,
        Status:    determineStatus(exitCode),
    }, err
}
```

### Dry-Run Execution

```go
func (e *DryRunExecutor) Execute(cmd Command) (*Result, error) {
    return &Result{
        Command:   cmd.String(),
        DryRun:    true,
        ExitCode:  0,
        Stdout:    fmt.Sprintf("Would execute: %s", cmd.String()),
        Stderr:    "",
        StartTime: time.Now(),
        Duration:  0,
        Status:    "pending",
    }, nil
}
```

## Testing Contracts

### Unit Tests (Command Building)

```go
func TestStartContainerBuilder(t *testing.T) {
    tests := []struct {
        name        string
        containerID string
        want        string
        wantErr     bool
    }{
        {
            name:        "valid container ID",
            containerID: "abc123",
            want:        "container start abc123",
            wantErr:     false,
        },
        {
            name:        "empty container ID",
            containerID: "",
            wantErr:     true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            builder := NewStartContainer(tt.containerID)
            cmd, err := builder.Build()
            
            if tt.wantErr {
                assert.Error(t, err)
                return
            }
            
            assert.NoError(t, err)
            assert.Equal(t, tt.want, cmd.String())
        })
    }
}
```

### Contract Tests (with Dry-Run)

```go
func TestDestructiveActionsRequireConfirmation(t *testing.T) {
    // Verify that delete commands are marked as requiring confirmation
    builder := NewDeleteContainer("test-container")
    cmd, _ := builder.Build()
    
    // In actual implementation, this would check metadata or annotations
    assert.True(t, cmd.RequiresConfirmation())
    assert.Equal(t, "type-to-confirm", cmd.ConfirmationType())
}
```

### Integration Tests (Optional, requires Apple Container)

```go
// +build integration

func TestStartContainerIntegration(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test")
    }
    
    // Prerequisites: container daemon running, test container exists
    executor := &RealExecutor{}
    builder := NewStartContainer("test-container")
    cmd, _ := builder.Build()
    
    result, err := executor.Execute(cmd)
    assert.NoError(t, err)
    assert.Equal(t, 0, result.ExitCode)
}
```

## Error Handling Patterns

### Common Error Patterns

| Error Type | stderr Pattern | Recommended User Message |
|------------|----------------|--------------------------|
| Daemon not running | `connection`, `daemon` | "Container daemon is not running. Start with Daemon → Start" |
| Not found | `not found`, `no such` | "Container/Image not found: <id>" |
| Permission denied | `permission`, `sudo` | "Permission denied. May require sudo or system configuration" |
| Already running | `already`, `running` | "Container is already running" |
| Network error | `connect`, `timeout` | "Network error. Check connection and try again" |
| Auth required | `unauthorized`, `authentication` | "Registry authentication required. Use 'container registry login'" |

### Error Display in TUI

```go
func FormatError(result *Result) string {
    if result.ExitCode == 0 {
        return "" // No error
    }
    
    errorMsg := result.Stderr
    if errorMsg == "" {
        errorMsg = result.Stdout
    }
    
    // Add context
    return fmt.Sprintf(
        "Command failed: %s\nExit code: %d\nError: %s",
        result.Command,
        result.ExitCode,
        errorMsg,
    )
}
```

## Constitution Compliance

### Command Safety (Principle I)

✅ All commands support dry-run via DryRunExecutor  
✅ Command preview displayed before execution (Command.String())  
✅ Destructive actions require explicit confirmation (delete, daemon stop)  

### Tested Command Contracts (Principle V)

✅ Unit tests for all command builders  
✅ Validation before command building  
✅ Contract tests for destructive action metadata  
✅ Optional integration tests for end-to-end verification  

## Notes

- All commands use `os/exec.Command` with separate args (no shell execution)
- Command preview uses `Command.String()` for safe display formatting
- Exit code 0 always indicates success; non-zero indicates error
- ANSI progress codes should be handled or stripped in TUI display
- Streaming output (build, pull) should be buffered or displayed incrementally
