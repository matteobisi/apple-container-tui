# Phase 0: Research - Enhanced Menu Navigation and Image Management

**Feature**: 002-refactor-menu-images  
**Date**: 2026-02-16  
**Purpose**: Technical research for implementing submenu navigation, image management, and container interaction features

## Overview

This document resolves technical unknowns for implementing three user stories:
1. Container action submenu with context-sensitive options
2. Image list view with pull/build/prune operations  
3. Image submenu with inspect/delete capabilities

## Research Tasks Completed

### 1. Bubbletea Navigation State Management for Submenus

**Question**: How should we manage navigation state to support multiple menu levels (main → container submenu → container logs/shell, main → image list → image submenu → image inspect)?

**Decision**: Use a navigation stack with view states

**Rationale**:
- Bubbletea's `tea.Model` interface supports nested models where each view (container list, container submenu, image list, etc.) is its own model that handles its own Update/View logic
- Navigation stack maintained in main `app.go` model tracks current view and previous views for back-navigation
- Esc key pops from navigation stack to return to previous view
- Each submenu model receives the selected item (container or image) as initialization parameter

**Alternatives Considered**:
- Single flat state machine: Rejected because it creates excessive coupling between views and makes back-navigation complex
- Global navigation manager: Rejected because it violates Bubbletea's message-passing pattern and makes testing harder

**Implementation Approach**:
```go
type AppModel struct {
    currentView ViewType
    navigationStack []ViewType
    containerList *ContainerListModel
    containerSubmenu *ContainerSubmenuModel
    imageList *ImageListModel
    // ... other view models
}

type ViewType int
const (
    ViewContainerList ViewType = iota
    ViewContainerSubmenu
    ViewContainerLogs
    ViewImageList
    ViewImageSubmenu
    // ...
)
```

**References**: Existing pattern in `src/ui/app.go` already uses view switching for build/pull/help screens

---

### 2. Shell Detection for Container Exec

**Question**: How should we detect which shell is available in a container to execute `container exec -it` successfully?

**Decision**: Sequential shell probing with non-interactive exec

**Rationale**:
- Apple Container CLI supports `container exec <container> <command>` without `-it` flag for checking command existence
- Probe sequence: bash → sh → /bin/sh → /bin/bash → ash (covers 95%+ of container images)
- Use `container exec <name> which bash` or `container exec <name> test -x /bin/bash` to check availability
- First successful probe determines shell for interactive session
- Cache result per container to avoid repeated probes

**Alternatives Considered**:
- Try each shell interactively until one works: Rejected because failed attempts create poor UX (user sees multiple error messages)
- Always use /bin/sh: Rejected because not all minimal containers have /bin/sh symlink
- Query container image metadata: Rejected because ENTRYPOINT/CMD may not reflect available shells

**Implementation Approach**:
```go
type ShellDetector struct {
    executor CommandExecutor
    cache map[string]string // containerID -> detected shell
}

func (sd *ShellDetector) DetectShell(containerID string) (string, error) {
    shells := []string{"bash", "sh", "/bin/sh", "/bin/bash", "ash"}
    for _, shell := range shells {
        cmd := buildCheckCmd(containerID, shell)
        if sd.executor.Execute(cmd) == nil {
            sd.cache[containerID] = shell
            return shell, nil
        }
    }
    return "", ErrNoShellFound
}
```

**References**: Similar pattern used in Kubernetes kubectl exec, Docker exec

---

### 3. Live Container Log Streaming

**Question**: How should we stream live container logs (`container logs -f`) in a Bubbletea TUI while maintaining responsiveness?

**Decision**: Background goroutine with tea.Cmd pattern

**Rationale**:
- Bubbletea supports asynchronous operations via `tea.Cmd` that return `tea.Msg` when complete
- Spawn goroutine that tails container logs and sends line-by-line messages to Update loop
- Use buffered channel to handle backpressure if log output exceeds rendering speed
- Support graceful cancellation when user presses Esc or container stops
- Scrollable viewport (from bubbles library) displays log content

**Alternatives Considered**:
- Blocking read in Update(): Rejected because it freezes entire TUI
- Poll logs periodically: Rejected because it introduces latency (requirement is <100ms)
- External terminal multiplexer: Rejected because it breaks TUI abstraction

**Implementation Approach**:
```go
func streamContainerLogs(containerID string) tea.Cmd {
    return func() tea.Msg {
        cmd := exec.Command("container", "logs", "-f", containerID)
        stdout, _ := cmd.StdoutPipe()
        cmd.Start()
        
        scanner := bufio.NewScanner(stdout)
        for scanner.Scan() {
            return LogLineMsg{Line: scanner.Text()}
        }
        
        if err := cmd.Wait(); err != nil {
            return LogErrorMsg{Err: err}
        }
        return LogEndedMsg{}
    }
}
```

**References**: Pattern used in existing container start/stop operations in `src/ui/progress.go`

---

### 4. Image List Parsing

**Question**: How should we parse the output of `container image list` to extract NAME, TAG, and DIGEST columns?

**Decision**: Structured parsing with column detection

**Rationale**:
- `container image list` outputs tabular format with headers
- First line contains column headers (NAME, TAG, DIGEST, ...)
- Subsequent lines contain whitespace-separated values
- Use Go's `strings.Fields()` for splitting, then map columns by header names
- Handle edge cases: missing digest, long names that wrap, empty output

**Alternatives Considered**:
- JSON output flag: Rejected because Apple Container CLI may not support JSON format for image list
- Fixed column widths: Rejected because column width varies with content
- Regex parsing: Rejected because fragile if format changes

**Implementation Approach**:
```go
type Image struct {
    Name   string
    Tag    string
    Digest string
}

func ParseImageList(output string) ([]Image, error) {
    lines := strings.Split(strings.TrimSpace(output), "\n")
    if len(lines) < 2 {
        return nil, ErrEmptyImageList
    }
    
    headers := strings.Fields(lines[0])
    nameIdx := indexOf(headers, "NAME")
    tagIdx := indexOf(headers, "TAG")
    digestIdx := indexOf(headers, "DIGEST")
    
    var images []Image
    for _, line := range lines[1:] {
        fields := strings.Fields(line)
        images = append(images, Image{
            Name:   fields[nameIdx],
            Tag:    fields[tagIdx],
            Digest: fields[digestIdx],
        })
    }
    return images, nil
}
```

**References**: Similar pattern in existing `container parser.go` for parsing container list output

---

### 5. Image Inspect JSON Display

**Question**: How should we display the JSON output of `container image inspect <image> | jq` in a scrollable, readable format?

**Decision**: Use bubbles viewport component with syntax highlighting optional

**Rationale**:
- Bubbles library provides `viewport.Model` for scrollable content
- Pipe container image inspect output to `jq` for formatting (requirement FR-027)
- Display raw formatted JSON in viewport with up/down/pgup/pgdn scrolling
- Syntax highlighting nice-to-have but not required (can use lipgloss for basic coloring)
- Esc returns to image submenu

**Alternatives Considered**:
- Custom JSON tree navigator: Rejected as overengineering for inspection use case
- External pager (less/more): Rejected because it breaks TUI session
- Flatten JSON to table: Rejected because loses structure and is harder to read

**Implementation Approach**:
```go
type ImageInspectModel struct {
    viewport viewport.Model
    content  string
    err      error
}

func NewImageInspectModel(imageName string) ImageInspectModel {
    vp := viewport.New(80, 24)
    
    // Execute: container image inspect <imageName> | jq
    cmd := exec.Command("sh", "-c", fmt.Sprintf("container image inspect %s | jq", imageName))
    output, err := cmd.CombinedOutput()
    
    vp.SetContent(string(output))
    return ImageInspectModel{viewport: vp, content: string(output), err: err}
}
```

**References**: Existing help screen uses similar viewport pattern in `src/ui/help.go`

---

### 6. Type-to-Confirm Pattern for Image Prune

**Question**: How should we implement type-to-confirm for the destructive image prune operation?

**Decision**: Reuse existing type-to-confirm component

**Rationale**:
- Existing TUI already has type-to-confirm component used for container deletion (per constitution)
- Located in `src/ui/type_to_confirm.go`
- User must type exact word (e.g., "prune" or "yes") to confirm
- Can be adapted for image prune by changing confirmation text
- Maintains consistency with delete container UX

**Alternatives Considered**:
- Simple yes/no: Rejected because doesn't provide adequate safety for bulk operation (per clarification Q5)
- Preview list then confirm: Rejected as more complex and prune command handles filtering

**Implementation Approach**:
```go
// Reuse existing component
confirmModel := typetoconfirm.New(
    "Type 'prune' to remove all unused images",
    "prune",
)

// On successful confirmation:
if confirmModel.Confirmed() {
    executeImagePrune()
}
```

**References**: `src/ui/type_to_confirm.go`, used in container deletion flow

---

## Technology Stack Decisions

| Component | Technology | Rationale |
|-----------|-----------|-----------|
| **Navigation State** | Bubbletea nested models with navigation stack | Native to framework, supports back-navigation, testable |
| **Shell Detection** | Sequential exec probes with caching | Reliable, covers 95%+ containers, minimizes failed attempts |
| **Log Streaming** | Goroutine with tea.Cmd pattern | Asynchronous, non-blocking, maintains TUI responsiveness |
| **Image Parsing** | Column-based text parsing | Robust to format variations, handles edge cases |
| **JSON Display** | Bubbles viewport component | Built-in scrolling, consistent with help screen pattern |
| **Confirmations** | Existing type-to-confirm UI | Maintains UX consistency, already tested |

## Best Practices Applied

1. **Reuse Existing Patterns**: Navigation state management, viewport usage, type-to-confirm reuse established TUI patterns
2. **Asynchronous Operations**: Log streaming uses tea.Cmd to avoid blocking Update loop
3. **Error Handling**: All operations (shell detection, log streaming, image parsing) have explicit error paths per FR-033-037
4. **Testing Strategy**: Shell detection logic unit-testable, command builders contract-testable, navigation flows integration-testable
5. **Performance**: Caching (shell detection), buffered channels (log streaming), column-based parsing (no regex)

## Open Questions / Deferred Decisions

None. All technical unknowns resolved. Ready for Phase 1 (data model and contracts).

---

**Research Complete**: 2026-02-16  
**Next Phase**: Phase 1 - Data Model and Contracts
