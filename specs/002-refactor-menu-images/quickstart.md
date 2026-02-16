# Quickstart: Enhanced Menu Navigation and Image Management

**Feature**: 002-refactor-menu-images  
**Audience**: Developers implementing this feature  
**Date**: 2026-02-16

## Overview

This guide walks through implementing the menu refactoring and image management features. It covers the three main components:

1. **Container submenu**: Context-sensitive actions for containers
2. **Image list view**: Browse and manage local images  
3. **Image submenu**: Detailed operations on selected images

Read this guide before starting implementation to understand the architecture, integration points, and testing strategy.

---

## Architecture Overview

### Navigation Model

The feature uses a **navigation stack pattern** where each view is an independent Bubbletea model:

```
                  ┌──────────────┐
                  │   AppModel   │ (Main coordinator)
                  └──────┬───────┘
                         │
         ┌───────────────┼───────────────┐
         │               │               │
    ┌────▼────┐    ┌─────▼─────┐   ┌────▼─────┐
    │Container│    │  Image    │   │ Daemon   │
    │  List   │    │   List    │   │ Control  │
    └────┬────┘    └─────┬─────┘   └──────────┘
         │               │
    ┌────▼────┐    ┌─────▼─────┐
    │Container│    │  Image    │
    │ Submenu │    │ Submenu   │
    └────┬────┘    └─────┬─────┘
         │               │
  ┌──────┼──────┐   ┌────▼─────┐
  │      │      │   │  Image   │
┌─▼──┐ ┌▼───┐ ┌▼──┐│ Inspect  │
│Logs│ │Shell│ │.. ││          │
└────┘ └────┘ └───┘└──────────┘
```

### Key Principles

1. **Each view is a separate model**: Container list, container submenu, image list, image submenu are independent models
2. **Navigation stack for back-navigation**: Push/pop pattern for Esc key behavior
3. **Message passing**: Views communicate via Bubbletea messages (Update → Msg → Update cycle)
4. **Reuse existing patterns**: Type-to-confirm, command builders, viewport components already exist

---

## Implementation Roadmap

### Phase 1: Navigation Infrastructure (Foundation)

**Goal**: Set up navigation state management before adding views

**Files to Modify**:
- `src/ui/app.go`: Add navigation stack and view states
- `src/ui/keys.go`: Add navigation keys (Esc, arrow keys)

**Tasks**:
1. Define `ViewType` enum in `app.go`
2. Add navigation stack field to `AppModel`
3. Implement `pushView()` and `popView()` methods
4. Update main `Update()` to route messages to active view
5. Update main `View()` to render active view

**Testing**:
- Unit test navigation stack push/pop logic
- Verify Esc key pops from stack correctly

**Completion Criteria**: Can push/pop between existing views (container list, help, daemon control)

---

### Phase 2: Container Submenu (First Submenu Implementation)

**Goal**: Implement submenu pattern for containers (establishes pattern for image submenu)

**New Files**:
- `src/ui/container_submenu.go`: Submenu view model
- `src/ui/container_logs.go`: Log streaming view
- `src/ui/container_shell.go`: Interactive shell wrapper
- `src/services/container_logs_builder.go`: Command builder for `container logs -f`
- `src/services/container_exec_builder.go`: Command builder for `container exec -it`
- `src/services/shell_detector.go`: Shell detection logic

**Files to Modify**:
- `src/ui/container_list.go`: Change Enter key from toggle to open submenu
- `src/ui/app.go`: Add container submenu routing

**Implementation Order**:
1. **Start simple**: Container submenu with just "Start", "Stop", "Back" (reuse existing start/stop)
2. **Add logs streaming**: ContainerLogsModel with viewport, async log reading
3. **Add shell detection**: Probe sequence, caching
4. **Add interactive shell**: Suspend TUI, exec container, resume TUI

**Testing**:
- Contract tests for logs/exec command builders
- Unit tests for shell detection logic
- Integration test: Enter submenu → Select logs → View logs → Esc → Back to submenu

**Completion Criteria**: Can navigate into container submenu, select action, return to list

---

### Phase 3: Image List View

**Goal**: Display local images and support pull/build/prune operations

**New Files**:
- `src/ui/image_list.go`: Image list view model
- `src/models/image.go`: Image entity
- `src/services/image_list_builder.go`: Command builder for `container image list`
- `src/services/image_prune_builder.go`: Command builder for `container image prune`

**Files to Modify**:
- `src/ui/container_list.go`: Add 'i' key handler to navigate to image list
- `src/ui/keys.go`: Add 'i' key constant
- `src/ui/app.go`: Add image list routing

**Implementation Order**:
1. **Display empty list**: Image list view with no data, shows "No images" message
2. **Parse and display**: Execute `container image list`, parse output, display in table
3. **Add pull/build**: Wire 'p' and 'b' keys to existing workflows, ensure they return to image list (not container list)
4. **Add prune**: Type-to-confirm, execute prune, refresh list

**Testing**:
- Unit tests for image list parser (well-formed, empty, malformed output)
- Contract tests for image list/prune command builders
- Integration test: Press 'i' → View images → Press 'p' → Pull → Return to image list

**Completion Criteria**: Can browse images, pull/build/prune from image list view

---

### Phase 4: Image Submenu and Inspect

**Goal**: Add detailed image actions (inspect, delete)

**New Files**:
- `src/ui/image_submenu.go`: Image submenu view model
- `src/ui/image_inspect.go`: Inspection viewport view
- `src/services/image_inspect_builder.go`: Command builder for `container image inspect | jq`
- `src/services/image_delete_builder.go`: Command builder for `container image rm`

**Files to Modify**:
- `src/ui/image_list.go`: Add Enter key handler to open image submenu
- `src/ui/app.go`: Add image submenu and inspect routing

**Implementation Order**:
1. **Image submenu skeleton**: Just "Inspect", "Delete", "Back" options with arrow navigation
2. **Inspect view**: Execute inspect | jq, display in viewport with scrolling
3. **Delete with confirmation**: Type-to-confirm, execute delete, handle "in use" error

**Testing**:
- Contract tests for inspect/delete command builders
- Integration test: Select image → Press Enter → Select "Inspect" → View JSON → Esc → Back
- Integration test: Select image → Press Enter → Select "Delete" → Confirm → Image removed

**Completion Criteria**: Can inspect and delete images from submenu

---

## Code Patterns to Follow

### 1. Bubbletea Model Template

Every new view follows this pattern:

```go
package ui

import (
    tea "github.com/charmbracelet/bubbletea"
)

// Model state
type MyViewModel struct {
    width  int
    height int
    cursor int
    items  []string
    err    error
}

// Constructor
func NewMyViewModel() MyViewModel {
    return MyViewModel{
        cursor: 0,
        items:  []string{},
    }
}

// Initialize (optional async operations)
func (m MyViewModel) Init() tea.Cmd {
    return nil  // or return tea.Cmd for async work
}

// Handle messages (business logic)
func (m MyViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "up":
            if m.cursor > 0 {
                m.cursor--
            }
        case "down":
            if m.cursor < len(m.items)-1 {
                m.cursor++
            }
        case "esc":
            return m, func() tea.Msg { return BackToListMsg{} }
        }
    case tea.WindowSizeMsg:
        m.width = msg.Width
        m.height = msg.Height
    }
    return m, nil
}

// Render view
func (m MyViewModel) View() string {
    s := "My View\n\n"
    for i, item := range m.items {
        cursor := " "
        if i == m.cursor {
            cursor = ">"
        }
        s += fmt.Sprintf("%s %s\n", cursor, item)
    }
    return s
}
```

### 2. Command Builder Pattern

Follow existing pattern in `src/services/`:

```go
package services

type ImageListBuilder struct{}

func NewImageListBuilder() *ImageListBuilder {
    return &ImageListBuilder{}
}

func (b *ImageListBuilder) Build() *Command {
    return &Command{
        executable: "container",
        args:       []string{"image", "list"},
        logType:    "image_list",
    }
}
```

### 3. Async Operations (Logs, Exec)

Use tea.Cmd for non-blocking operations:

```go
func streamLogs(containerID string) tea.Cmd {
    return func() tea.Msg {
        // Spawn goroutine or use exec.Command
        // Send messages back for each log line
        return LogLineMsg{Line: "container output"}
    }
}

// In Update():
case StartStreamingMsg:
    return m, streamLogs(m.containerID)
case LogLineMsg:
    m.lines = append(m.lines, msg.Line)
    return m, nil
```

### 4. Type-to-Confirm Pattern

Reuse existing component:

```go
import "container-tui/src/ui/typetoconfirm"

// In model:
confirmModel := typetoconfirm.New(
    "Type 'prune' to remove all unused images",
    "prune",
)

// In Update():
if confirmModel.Confirmed() {
    return m, executePrune()
}
```

---

## Integration Points

### With Existing Features

| Existing Feature | Integration Point | Notes |
|------------------|-------------------|-------|
| Container start/stop | Container submenu calls existing workflows | Just call existing functions |
| Image pull | Image list 'p' key calls existing workflow | Modify return destination to image list |
| Image build | Image list 'b' key calls existing workflow | Modify return destination to image list |
| Type-to-confirm | Reuse for prune and delete | Same pattern as container delete |
| Command logging | All new commands must log | Use existing CommandExecutor |
| Viewport (help) | Reuse for image inspect, container logs | Already in bubbles library |

### New Integration Points

1. **Shell Detection → Container Exec**: Shell detector must cache results per container
2. **Image List Parser → Image Entity**: Parser creates Image structs from CLI output
3. **Navigation Stack → All Views**: All views must support Esc key and navigation messages

---

## Testing Strategy

### Test Organization

```
tests/
├── contract/
│   ├── image_list_test.go          # Test container image list command
│   ├── image_inspect_test.go       # Test container image inspect | jq command
│   ├── image_prune_test.go         # Test container image prune command
│   ├── container_logs_test.go      # Test container logs -f command
│   └── container_exec_test.go      # Test container exec -it with shell
├── integration/
│   ├── navigation_flows_test.go    # Test full navigation paths
│   └── image_operations_test.go    # Update with new operations
└── unit/
    ├── shell_detector_test.go      # Test shell detection logic
    └── image_parser_test.go        # Test image list parsing
```

### Test Priority

1. **Contract tests first**: Ensure command generation is correct
2. **Unit tests for parsing**: Image list parser, shell detector
3. **Integration tests last**: Full navigation flows

### Mock Strategy

- Mock `CommandExecutor` for command builders
- Mock container/image CLI output for parsers
- Use test fixtures for well-formed and malformed CLI output

---

## Common Pitfalls

### 1. Forgetting to Update Navigation Stack

❌ **Wrong**:
```go
case tea.KeyMsg:
    if msg.String() == "i" {
        m.currentView = ViewImageList  // Stack not updated!
    }
```

✅ **Correct**:
```go
case tea.KeyMsg:
    if msg.String() == "i" {
        m.pushView(ViewImageList)
        m.currentView = ViewImageList
    }
```

### 2. Blocking Update() with Long Operations

❌ **Wrong**:
```go
func (m MyModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    output := exec.Command("container", "image", "list").Output()  // BLOCKS!
    m.images = parse(output)
    return m, nil
}
```

✅ **Correct**:
```go
func fetchImages() tea.Cmd {
    return func() tea.Msg {
        output := exec.Command("container", "image", "list").Output()
        return ImagesLoadedMsg{images: parse(output)}
    }
}

func (m MyModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    case ImagesLoadedMsg:
        m.images = msg.images
        return m, nil
}
```

### 3. Not Handling Window Resize

❌ **Wrong**:
```go
// Hardcoded dimensions
viewport := viewport.New(80, 24)
```

✅ **Correct**:
```go
case tea.WindowSizeMsg:
    m.viewport.Width = msg.Width - 4
    m.viewport.Height = msg.Height - 6
```

---

## Debugging Tips

### 1. Enable Debug Logging

Bubbletea supports debug logging to file:

```go
f, _ := tea.LogToFile("debug.log", "debug")
defer f.Close()
```

### 2. Print Messages

Add logging in Update():

```go
func (m MyModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    log.Printf("Received message: %T %+v\n", msg, msg)
    // ... rest of Update
}
```

### 3. Test Views in Isolation

Create standalone test programs for individual views:

```go
// test_container_submenu.go
func main() {
    p := tea.NewProgram(NewContainerSubmenuModel(testContainer))
    p.Run()
}
```

---

## Performance Considerations

1. **Image list parsing**: Should handle 100+ images in <2 seconds
2. **Log streaming**: Buffer channel to avoid blocking on fast output
3. **Shell detection**: Cache results per container to avoid repeated probes
4. **Viewport rendering**: Only render visible portion of large JSON

---

## Pre-Implementation Checklist

Before starting implementation:

- [ ] Read spec.md (functional requirements)
- [ ] Read research.md (technical decisions)
- [ ] Read data-model.md (entities and relationships)
- [ ] Read all three contracts (container-submenu, image-list, image-submenu)
- [ ] Review existing code in src/ui/ (app.go, container_list.go, help.go)
- [ ] Understand Bubbletea model pattern (Init/Update/View)
- [ ] Understand navigation stack pattern from this quickstart
- [ ] Set up test fixtures for container/image CLI output

---

## Getting Help

- **Bubbletea docs**: https://github.com/charmbracelet/bubbletea
- **Bubbles components**: https://github.com/charmbracelet/bubbles (viewport, list, etc.)
- **Existing patterns**: Look at src/ui/help.go (viewport), src/ui/container_list.go (list navigation), src/ui/type_to_confirm.go (confirmation)

---

**Ready to implement!** Start with Phase 1 (navigation infrastructure) and work through each phase sequentially.
