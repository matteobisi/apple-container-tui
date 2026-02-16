# Contract: Container Submenu

**Feature**: 002-refactor-menu-images  
**Component**: Container context menu for lifecycle and diagnostic operations  
**Date**: 2026-02-16

## Purpose

Define the user interactions and command contracts for the container submenu, which provides context-sensitive actions for a selected container (start/stop, view logs, enter shell).

---

## User Interface Contract

### Entry Point
- **Trigger**: User presses `Enter` on a container in the main container list
- **Precondition**: At least one container exists (list is not empty)
- **Effect**: Displays container submenu with context-appropriate options

### Menu Display

**For Stopped Container**:
```
Container: [container-name] (stopped)

> Start container
  Tail container log
  Back

Keys: up/down=navigate, enter=select, esc=back
```

**For Running Container**:
```
Container: [container-name] (running)

> Stop container
  Tail container log
  Enter container
  Back

Keys: up/down=navigate, enter=select, esc=back
```

### Navigation
- **Arrow keys (up/down)**: Move highlight between menu options
- **Enter**: Execute highlighted action
- **Esc**: Return to container list (equivalent to selecting "Back")

### Menu Options

| Option | Visible When | Action | Maps To (Spec) |
|--------|-------------|--------|----------------|
| Start container | Container stopped | Initiate start workflow | Existing start flow |
| Stop container | Container running | Initiate stop workflow | Existing stop flow |
| Tail container log | Always | Stream live logs | FR-009, FR-010 |
| Enter container | Container running | Open interactive shell | FR-011, FR-012, FR-013 |
| Back | Always | Return to container list | FR-008 |

---

## Command Contracts

### 1. Tail Container Log

**User Action**: Select "Tail container log"

**Command Generated**:
```bash
container logs -f <containerName>
```

**Parameters**:
- `<containerName>`: Container name or ID from selected container

**Expected Output**:
- Live stream of log lines to stdout
- Continues until container stops or user presses Esc

**Error Scenarios**:
- Container not found: Display error, return to submenu
- Container stopped during streaming: Show "Container stopped" message, wait for Esc (FR-034)
- Permission denied: Display error, return to submenu

**Success Criteria**:
- Log lines appear in <100ms of container output (SC-003)
- Esc key returns to container submenu (FR-010)
- Scrollable viewport for historical logs

---

### 2. Enter Container (Shell Detection)

**User Action**: Select "Enter container" (only visible for running containers)

**Phase 1: Shell Detection** (before interactive session)

**Probe Sequence**:
```bash
# Attempt 1
container exec <containerName> which bash
# If exit code 0: use bash

# Attempt 2 (if 1 failed)
container exec <containerName> which sh
# If exit code 0: use sh

# Attempt 3 (if 2 failed)
container exec <containerName> test -x /bin/sh
# If exit code 0: use /bin/sh

# Attempt 4 (if 3 failed)
container exec <containerName> test -x /bin/bash
# If exit code 0: use /bin/bash

# Attempt 5 (if 4 failed)
container exec <containerName> which ash
# If exit code 0: use ash

# If all failed: error (FR-033)
```

**Phase 2: Interactive Session** (after shell detected)

**Command Generated**:
```bash
container exec -it <containerName> <detectedShell>
```

**Parameters**:
- `<containerName>`: Container name or ID
- `<detectedShell>`: Shell command from detection phase (e.g., "bash", "/bin/sh")

**Expected Behavior**:
- TUI suspends rendering
- Terminal switches to raw mode for interactive shell
- User can execute commands in container
- `exit` command or Ctrl+D returns to TUI
- TUI resumes, navigates back to container submenu (FR-013)

**Error Scenarios**:
- No shell found (all probes failed): Display error "No supported shell found in container", stay in submenu (FR-033 + Clarification Q3)
- Container stopped during detection: Display error, return to submenu
- Shell session crashes: Display error, return to submenu
- Exec permission denied: Display error, return to submenu

**Success Criteria**:
- Shell detection succeeds in 95% of standard containers (SC-004)
- Detection cached per container (avoid repeated probes)
- Interactive session fully functional (can run commands, see output, use Ctrl+C, etc.)
- Clean return to TUI after exit

---

### 3. Start Container

**User Action**: Select "Start container" (only visible for stopped containers)

**Command Generated**:
```bash
container start <containerName>
```

**Parameters**:
- `<containerName>`: Container name or ID

**Expected Behavior**:
- Delegates to existing start container workflow (already implemented)
- Shows progress/spinner during start
- Returns to container list on completion (not to submenu)

**Error Scenarios**:
- Container already running: Display error
- Daemon not running: Display error, offer to start daemon

**Success Criteria**:
- Consistent with existing start flow
- Command logged for observability (Principle IV)

---

### 4. Stop Container

**User Action**: Select "Stop container" (only visible for running containers)

**Command Generated**:
```bash
container stop <containerName>
```

**Parameters**:
- `<containerName>`: Container name or ID

**Expected Behavior**:
- Delegates to existing stop container workflow (already implemented)
- Shows progress/spinner during stop
- Returns to container list on completion (not to submenu)

**Error Scenarios**:
- Container already stopped: Display error
- Stop timeout: Display timeout warning, offer force-stop

**Success Criteria**:
- Consistent with existing stop flow
- Command logged for observability (Principle IV)

---

## State Management

**Navigation State Changes**:
```
Initial: ContainerList view
User presses Enter on container
→ NavigationState.selectedContainer = selected container
→ NavigationState.currentView = ContainerSubmenu
→ Push ContainerList to navigationStack

User selects "Tail container log"
→ NavigationState.currentView = ContainerLogs
→ Push ContainerSubmenu to navigationStack

User presses Esc in ContainerLogs
→ Pop from navigationStack (return to ContainerSubmenu)
→ NavigationState.currentView = ContainerSubmenu

User presses Esc in ContainerSubmenu
→ Pop from navigationStack (return to ContainerList)
→ NavigationState.currentView = ContainerList
→ NavigationState.selectedContainer = null
```

---

## Testing Requirements

### Unit Tests
- [ ] Menu option generation based on container state (stopped vs running)
- [ ] Shell detection probe sequence logic
- [ ] Shell detection caching

### Contract Tests
- [ ] `container logs -f` command generated correctly with container name
- [ ] Shell detection probe commands generated correctly per probe sequence
- [ ] `container exec -it <name> <shell>` command generated with detected shell
- [ ] Start/stop commands delegate to existing builders

### Integration Tests
- [ ] Full flow: Enter submenu → Select "Tail logs" → View logs → Esc → Back to submenu
- [ ] Full flow: Enter submenu → Select "Enter container" → Shell detected → Interactive session → Exit → Back to submenu
- [ ] Shell detection with no available shells → Error message displayed → Stay in submenu
- [ ] Log streaming interrupted (container stopped) → Message displayed → Wait for Esc → Return to submenu

---

## Acceptance Criteria (from Spec)

Links to User Story 1 acceptance scenarios:

1. ✅ **Given** container list with stopped container, **When** press Enter, **Then** submenu shows "Start container", "Tail container log", "Back"
2. ✅ **Given** container list with running container, **When** press Enter, **Then** submenu shows "Stop container", "Tail container log", "Enter container", "Back"
3. ✅ **Given** container submenu, **When** select "Tail container log", **Then** live streaming logs with Esc to exit
4. ✅ **Given** container submenu, **When** select "Enter container", **Then** shell auto-detected and interactive session started
5. ✅ **Given** interactive shell session, **When** exit shell or press designated key, **Then** return to container submenu
6. ✅ **Given** container submenu, **When** select "Back" or press Esc, **Then** return to main container list

**Contract Complete**: 2026-02-16
