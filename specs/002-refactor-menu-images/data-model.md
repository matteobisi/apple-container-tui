# Phase 1: Data Model - Enhanced Menu Navigation and Image Management

**Feature**: 002-refactor-menu-images  
**Date**: 2026-02-16  
**Purpose**: Define data structures and entities for submenu navigation and image management

## Overview

This data model supports three main capabilities:
1. **Navigation State Management**: Track current view, navigation stack, and view transitions
2. **Image Entities**: Represent container images with name, tag, digest
3. **Shell Detection**: Track shell availability and detection results

All entities are technology-agnostic descriptions. Implementation details (Go structs, JSON serialization) are deferred to implementation phase.

---

## Core Entities

### 1. Image

**Description**: Represents a container image in the local repository

**Attributes**:
- `name` (string, required): Image repository name (e.g., "ubuntu", "ghcr.io/apple/container-builder-shim/builder")
- `tag` (string, required): Image tag (e.g., "latest", "0.7.0")
- `digest` (string, optional): Image content digest/hash (e.g., "sha256:32c70b3752ac28fd6f47c019...")

**Relationships**:
- None (images are independent entities)

**Validation Rules**:
- `name` must not be empty
- `tag` must not be empty
- `digest` format: sha256:[hexadecimal string] or empty

**State Transitions**: N/A (images are immutable once pulled)

**Usage Context**:
- Populated from `container image list` command output
- Displayed in image list view (FR-016)
- Used as identifier for inspect/delete operations

---

### 2. Navigation State

**Description**: Tracks the current view and navigation history within the TUI

**Attributes**:
- `currentView` (ViewType enum, required): Active view being displayed
- `navigationStack` (array of ViewType, required): History of views for back-navigation
- `selectedContainer` (Container reference, optional): Container selected when entering container submenu
- `selectedImage` (Image reference, optional): Image selected when entering image submenu

**ViewType Enumeration**:
- `ContainerList`: Main container list view
- `ContainerSubmenu`: Context menu for a specific container
- `ContainerLogs`: Live log streaming view
- `ContainerShell`: Interactive shell session view
- `ImageList`: Image list view
- `ImageSubmenu`: Context menu for a specific image
- `ImageInspect`: Image inspection JSON view
- `ImagePull`: Image pull workflow (existing)
- `ImageBuild`: Image build workflow (existing)
- `DaemonControl`: Daemon management view (existing)
- `Help`: Help screen (existing)

**Relationships**:
- References `Container` when `selectedContainer` is set
- References `Image` when `selectedImage` is set

**Validation Rules**:
- `currentView` must be a valid ViewType
- `navigationStack` maintains chronological order (newest at end)
- `selectedContainer` must be non-null when currentView is ContainerSubmenu/ContainerLogs/ContainerShell
- `selectedImage` must be non-null when currentView is ImageSubmenu/ImageInspect

**State Transitions**:
```
ContainerList --[Enter on container]--> ContainerSubmenu
ContainerSubmenu --[Select "Tail logs"]--> ContainerLogs
ContainerSubmenu --[Select "Enter container"]--> ContainerShell
ContainerList --[Press 'i']--> ImageList
ImageList --[Enter on image]--> ImageSubmenu
ImageSubmenu --[Select "Inspect"]--> ImageInspect
Any view --[Press Esc]--> {pop from navigationStack}
```

**Usage Context**:
- Maintained by main app model
- Drives view rendering and message routing
- Enables Esc-based back-navigation (FR-004)

---

### 3. Shell Detection Result

**Description**: Outcome of attempting to find an available shell in a container

**Attributes**:
- `containerID` (string, required): Container identifier
- `detectedShell` (string, optional): Shell command that was successfully detected (e.g., "bash", "/bin/sh")
- `probeAttempts` (array of ShellProbe, required): History of probe attempts
- `error` (string, optional): Error message if no shell found

**ShellProbe Structure**:
- `shellCommand` (string): Shell being tested (e.g., "bash")
- `success` (boolean): Whether shell was found
- `probedAt` (timestamp): When probe occurred

**Relationships**:
- References `Container` via containerID

**Validation Rules**:
- `containerID` must not be empty
- If `detectedShell` is set, at least one probeAttempt must have success=true
- If `error` is set, all probeAttempts must have success=false

**State Transitions**:
```
Initial: No detection attempted
Probing: Attempting detection (probeAttempts growing)
Success: detectedShell set, error null
Failure: detectedShell null, error set
```

**Usage Context**:
- Used by shell detection service before executing interactive shell
- Cached to avoid repeated probes for same container
- Displayed in error message if detection fails (FR-033)

---

### 4. Container Submenu Option

**Description**: A selectable action in the container context menu

**Attributes**:
- `label` (string, required): Display text (e.g., "Start container", "Tail container log")
- `action` (ActionType enum, required): Operation to perform
- `visible` (boolean, required): Whether option appears in menu
- `command` (string, optional): CLI command that will be executed (for preview)

**ActionType Enumeration**:
- `StartContainer`: Start a stopped container
- `StopContainer`: Stop a running container
- `TailLogs`: Stream container logs
- `EnterShell`: Open interactive shell
- `Back`: Return to container list

**Relationships**:
- Associated with a specific `Container` (from NavigationState.selectedContainer)

**Validation Rules**:
- `label` must not be empty
- `visible` determined by container state:
  - StartContainer: visible only if container stopped
  - StopContainer: visible only if container running
  - EnterShell: visible only if container running (FR-007)
  - TailLogs, Back: always visible

**State Transitions**: N/A (ephemeral, recreated on each submenu display)

**Usage Context**:
- Generated when entering container submenu
- Rendered as menu items with arrow-key navigation
- Executed when user presses Enter on highlighted option

---

### 5. Image Submenu Option

**Description**: A selectable action in the image context menu

**Attributes**:
- `label` (string, required): Display text (e.g., "Inspect image", "Delete image")
- `action` (ActionType enum, required): Operation to perform
- `command` (string, optional): CLI command that will be executed (for preview)

**ActionType Enumeration**:
- `InspectImage`: Display formatted JSON inspection
- `DeleteImage`: Remove image from local registry
- `Back`: Return to image list

**Relationships**:
- Associated with a specific `Image` (from NavigationState.selectedImage)

**Validation Rules**:
- `label` must not be empty
- All options always visible (no conditional visibility)

**State Transitions**: N/A (ephemeral, recreated on each submenu display)

**Usage Context**:
- Generated when entering image submenu
- Rendered as menu items with arrow-key navigation
- Executed when user presses Enter on highlighted option

---

### 6. Log Stream State

**Description**: State of an active container log streaming session

**Attributes**:
- `containerID` (string, required): Container being tailed
- `lines` (array of string, required): Log lines received so far
- `scrollOffset` (integer, required): Current scroll position in viewport
- `streaming` (boolean, required): Whether stream is active
- `error` (string, optional): Error message if stream failed/interrupted

**Relationships**:
- References `Container` via containerID

**Validation Rules**:
- `containerID` must not be empty
- `scrollOffset` must be >= 0 and <= len(lines)
- If `streaming` is false and error is null, stream ended normally
- If `error` is set, streaming must be false

**State Transitions**:
```
Initial: streaming=true, lines=[], error=null
Streaming: lines growing, scrollOffset may update
Interrupted: streaming=false, error set (container stopped/removed)
Ended: streaming=false, error=null (user pressed Esc)
```

**Usage Context**:
- Maintained by container logs view
- Updated asynchronously via tea.Cmd messages (FR-009)
- Displays error message on interruption (FR-034)

---

## Entity Relationships

```
NavigationState
├── currentView: ViewType
├── navigationStack: [ViewType]
├── selectedContainer ──> Container (existing entity)
└── selectedImage ──> Image

Image (independent)

ShellDetectionResult
└── containerID ──> Container

ContainerSubmenuOption
└── (associated with Container from NavigationState)

ImageSubmenuOption
└── (associated with Image from NavigationState)

LogStreamState
└── containerID ──> Container
```

---

## Data Flow Patterns

### Container Submenu Flow
```
1. User presses Enter on container in ContainerList
2. NavigationState.selectedContainer = selected container
3. NavigationState.currentView = ContainerSubmenu
4. Generate ContainerSubmenuOptions based on container state
5. User navigates options with arrow keys
6. User presses Enter on option → execute action
7. If action is "Tail logs":
   - NavigationState.currentView = ContainerLogs
   - Initialize LogStreamState
   - Start streaming
```

### Image List Flow
```
1. User presses 'i' in ContainerList
2. NavigationState.currentView = ImageList
3. Execute `container image list` command
4. Parse output into Image entities
5. Display images in list view
6. User presses Enter on image:
   - NavigationState.selectedImage = selected image
   - NavigationState.currentView = ImageSubmenu
```

### Shell Detection Flow
```
1. User selects "Enter container" from submenu
2. Check ShellDetectionResult cache for containerID
3. If not cached:
   - Create ShellDetectionResult
   - Probe shells in sequence (bash, sh, /bin/sh, /bin/bash, ash)
   - Set detectedShell on first success OR error on all failures
   - Cache result
4. If detectedShell found:
   - NavigationState.currentView = ContainerShell
   - Execute `container exec -it <id> <detectedShell>`
5. If no shell found:
   - Display error (FR-033)
   - Stay in ContainerSubmenu
```

---

## Validation Summary

| Entity | Key Validation | Enforcement Point |
|--------|---------------|-------------------|
| Image | name/tag not empty | Parse time (image list command) |
| NavigationState | selectedContainer set when in container views | State transition validation |
| ShellDetectionResult | detectedShell XOR error | Shell detection service |
| ContainerSubmenuOption | visible based on container state | Menu generation |
| ImageSubmenuOption | All always visible | Menu generation |
| LogStreamState | scrollOffset within bounds | Viewport update |

---

## Testing Strategy

**Unit Tests**:
- Image parsing from `container image list` output
- Navigation state transitions and stack management
- Shell detection probe logic
- Submenu option generation based on container state

**Contract Tests**:
- Image list command generation (FR-015)
- Image inspect command generation (FR-027)
- Image delete command generation (FR-031)
- Container logs command generation (FR-009)
- Container exec command generation (FR-012)

**Integration Tests**:
- Full navigation flow: ContainerList → ContainerSubmenu → ContainerLogs → back
- Full navigation flow: ContainerList → ImageList → ImageSubmenu → ImageInspect → back
- Shell detection with mock container responses
- Log streaming with simulated interruption

---

**Data Model Complete**: 2026-02-16  
**Next**: Phase 1 Contracts (container-submenu.md, image-list.md, image-submenu.md)
