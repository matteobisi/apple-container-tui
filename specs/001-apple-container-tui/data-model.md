# Data Model: Apple Container TUI

**Phase 1 Output** | **Date**: 2026-02-11  
**Purpose**: Define entities, attributes, and relationships for the container TUI

## Overview

This data model describes the entities managed by the Apple Container TUI. All entities are value objects or aggregates that represent either user interface state or container system state retrieved from the Apple Container CLI.

## Core Entities

### Container

Represents a managed container instance with its current state and metadata.

**Attributes**:
- `ID` (string, required): Unique container identifier from Apple Container
- `Name` (string, required): Human-readable container name
- `Image` (string, required): Image reference the container was created from
- `Status` (enum, required): Current container state
  - Possible values: `running`, `stopped`, `paused`, `created`, `unknown`
- `Created` (timestamp, optional): When the container was created
- `Ports` (list of PortMapping, optional): Exposed port mappings

**Relationships**:
- One Container is created from one ImageReference
- One Container may have zero or more PortMappings

**Validation Rules**:
- `ID` must not be empty
- `Name` must not be empty
- `Status` must be one of the defined enum values

**Source**: Parsed from `container list --all` output per contracts/cli-commands.md

---

### ImageReference

A user-supplied or system-managed reference to a container image.

**Attributes**:
- `Registry` (string, optional): Registry host (e.g., `docker.io`, `ghcr.io`)
- `Repository` (string, required): Repository name (e.g., `library/alpine`)
- `Tag` (string, optional): Image tag (default: `latest`)
- `Digest` (string, optional): SHA256 digest for immutable reference

**Relationships**:
- One ImageReference may be used by zero or more Containers

**Validation Rules**:
- `Repository` must not be empty
- If `Tag` and `Digest` are both present, Digest takes precedence
- Full reference format: `[registry/]repository[:tag][@digest]`

**Source**: User input (pull operation) or parsed from Container metadata

---

### BuildSource

A local filesystem reference to a container build file.

**Attributes**:
- `FilePath` (string, required): Absolute or relative path to build file
- `FileType` (enum, required): Type of build file
  - Possible values: `Containerfile`, `Dockerfile`
- `WorkingDirectory` (string, required): Build context directory
- `Exists` (boolean, computed): Whether the file exists on disk

**Relationships**:
- One BuildSource may produce one or more ImageReferences (after build)

**Validation Rules**:
- `FilePath` must not be empty
- `WorkingDirectory` must be a valid directory
- File must exist before build can proceed
- File must be readable

**Source**: User-selected file via file picker or path input

**Note**: The TUI auto-detects `FileType` by checking for `Containerfile` first, then `Dockerfile` per clarification session 2026-02-11.

---

### DaemonStatus

The current operational state of the Apple Container daemon.

**Attributes**:
- `Running` (boolean, required): Whether daemon is currently running
- `Version` (string, optional): Daemon version string
- `LastChecked` (timestamp, required): When status was last queried

**Relationships**:
- None (system-level singleton state)

**Validation Rules**:
- `LastChecked` must not be in the future

**Source**: Retrieved from daemon health check (e.g., `acrun --version` or daemon ping)

**Refresh**: Only refreshed when user explicitly requests it (per FR-015: manual refresh only)

---

### CommandRun

A record of a command execution with its preview, outcome, and captured output.

**Attributes**:
- `Command` (string, required): Full command string as would be executed
- `DryRun` (boolean, required): Whether this was a dry-run (preview only)
- `ExitCode` (integer, optional): Process exit code (0 = success)
- `Stdout` (string, optional): Captured standard output
- `Stderr` (string, optional): Captured standard error
- `StartTime` (timestamp, required): When command execution started
- `Duration` (duration, optional): How long the command took
- `Status` (enum, required): Execution outcome
  - Possible values: `pending`, `running`, `success`, `error`, `cancelled`

**Relationships**:
- One CommandRun may relate to one Container operation
- One CommandRun may relate to one ImageReference operation (pull)
- One CommandRun may relate to one BuildSource operation (build)

**Validation Rules**:
- `Command` must not be empty
- `ExitCode` is required if Status is `success` or `error`
- `Duration` is required if Status is `success` or `error`

**Source**: Captured during command execution by executor service

**Lifecycle**: Created when user initiates action, updated during execution, finalized with result

---

### UserConfig

Persisted user preferences and application configuration.

**Attributes**:
- `DefaultBuildFile` (string, optional): Default build file name to search for
- `ConfirmDestructiveActions` (boolean, required, default: true): Whether to show confirmation for deletions
- `ThemeMode` (enum, optional): UI color scheme
  - Possible values: `auto`, `light`, `dark`
- `RefreshOnFocus` (boolean, required, default: false): Whether to auto-refresh when TUI gains focus
- `LogRetentionDays` (integer, optional): How many days to keep command logs

**Relationships**:
- None (singleton per user)

**Validation Rules**:
- `LogRetentionDays` must be >= 0 if set
- `DefaultBuildFile` must be valid filename if set

**Source**: Loaded from config file at startup, written on changes

**Persistence**: Dual-path support per FR-011:
- Read from: `~/.config/apple-tui/config` OR `~/Library/Application Support/apple-tui/config` (check both, first found wins)
- Write to: `~/Library/Application Support/apple-tui/config` (macOS standard)

**Format**: TOML (using Viper library)

---

### PortMapping

A network port binding for a container.

**Attributes**:
- `HostPort` (integer, required): Port on host machine
- `ContainerPort` (integer, required): Port inside container
- `Protocol` (enum, required): Network protocol
  - Possible values: `tcp`, `udp`

**Relationships**:
- Zero or more PortMappings belong to one Container

**Validation Rules**:
- `HostPort` must be in valid port range (1-65535)
- `ContainerPort` must be in valid port range (1-65535)

**Source**: Parsed from Container metadata

---

## UI State Entities

These entities are specific to the TUI application state and do not represent container system state.

### Screen

The current active screen/view in the TUI.

**Possible Values** (enum):
- `ContainerList`: Main container listing view
- `ContainerDetail`: Detailed view of a single container
- `ImagePull`: Image pull workflow screen
- `BuildSelect`: Build file selection screen
- `BuildProgress`: Build execution and progress screen  
- `DaemonControl`: Daemon start/stop control screen
- `Settings`: User preferences configuration
- `CommandPreview`: Command preview and confirmation modal
- `Help`: Help and keyboard shortcuts screen

---

### SelectionState

Tracks which items are currently selected in list views.

**Attributes**:
- `SelectedContainerID` (string, optional): ID of currently selected container
- `SelectedImageRef` (string, optional): Currently selected image reference
- `CursorPosition` (integer, required): Current cursor position in active list

**Validation Rules**:
- `CursorPosition` must be >= 0

---

## Entity Relationships Diagram

```
UserConfig (singleton)
    
DaemonStatus (singleton)
    
ImageReference
    ├─ (1:N) → Container
    └─ (built from) BuildSource
    
Container
    ├─ (1:N) → PortMapping
    └─ (related to) CommandRun
    
BuildSource
    └─ (produces) ImageReference

CommandRun
    ├─ (relates to) Container
    ├─ (relates to) ImageReference
    └─ (relates to) BuildSource

Screen (current UI state)
    
SelectionState (current UI selection)
    └─ (references) Container | ImageReference
```

## Data Flow

### Startup
1. Load `UserConfig` from filesystem (dual-path check)
2. Check `DaemonStatus`
  3. If Apple Container CLI not found, ERROR and exit (per FR-014)
4. Initialize `Screen` to `ContainerList`
5. Fetch initial Container list (once, no auto-refresh per FR-015)

### Container List View
1. User selects Container → Update `SelectionState.SelectedContainerID` 
2. User chooses action → Transition to `CommandPreview` Screen
3. `CommandPreview` builds `CommandRun` with `DryRun=true`
4. User confirms → Execute `CommandRun` with `DryRun=false`
5. Update Container state from CLI output

### Image Pull
1. User enters `ImageReference` details
2. Validate `ImageReference`
3. Build `CommandRun` for pull operation
4. Show `CommandPreview`, user confirms
5. Execute pull, capture `Stdout`/`Stderr` in `CommandRun`
6. On success, refresh container list

### Build
1. User selects `BuildSource` (file picker)
2. Validate `BuildSource.Exists`
3. Auto-detect `FileType` (Containerfile > Dockerfile)
4. Build `CommandRun` for build operation
5. Show `CommandPreview`, user confirms
6. Execute build with streaming output to `BuildProgress` screen
7. On success, new `ImageReference` available

### Daemon Control
1. Fetch current `DaemonStatus`
2. User chooses start/stop action
3. Build `CommandRun` for daemon operation
4. Show confirmation (destructive action for stop)
5. Execute, update `DaemonStatus`

### Manual Refresh (per FR-015)
1. User presses refresh key (e.g., `r`)
2. Re-fetch Container list from CLI
3. Update SelectionState if selected container no longer exists
4. Refresh DaemonStatus

## Persistence and Caching

**Persisted**:
- `UserConfig`: Written to disk on change

**Cached in Memory** (until manual refresh):
- Container list
- DaemonStatus

**Not Persisted**:
- `CommandRun` history: Logged to file but not loaded at startup
- All UI state entities (`Screen`, `SelectionState`)

**Command Logs**:
- Each `CommandRun` appended to log file: `~/Library/Application Support/apple-tui/command.log`
- Format: JSON lines (one CommandRun per line)
- Rotation: Based on `UserConfig.LogRetentionDays`

## Validation Strategy

All validation rules above should be enforced in:
1. **Input forms**: Real-time validation as user types
2. **Command building**: Validation before creating CommandRun
3. **CLI parsing**: Validation when parsing Apple Container CLI output

Validation errors should be surfaced in the UI with:
- Clear error messages
- Suggested corrections where possible
- Highlighted fields for form inputs

## Notes

- All timestamps are stored in UTC, displayed in local timezone
- All file paths are normalized to absolute paths for consistency
- Container state is authoritative from CLI, never modified by TUI
- The TUI is read-only except for destructive actions (delete)
