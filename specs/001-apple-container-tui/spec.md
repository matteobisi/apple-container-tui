# Feature Specification: Apple Container TUI

**Feature Branch**: `001-apple-container-tui`  
**Created**: 2026-02-11  
**Status**: Draft  
**Input**: User description: "i want to build this TUI with a modern language that permits with easy dependancy to have some facilitator like resize if terminal windows is resized,  decent graphics. The program could be shipped as script that must be transformed in executable or as binary self contained ( if needd we could create conifg file in standard location like ~/.config/apple-tui/config)"

## Clarifications

### Session 2026-02-11

- Q: Which container build file names should the TUI support? → A: Support Containerfile and Dockerfile (auto-detect).
- Q: Where should user configuration be stored? → A: Support both ~/.config/apple-tui/config and ~/Library/Application Support/apple-tui/config (read both, write to one).
- Q: Should the TUI support keyboard-only navigation, mouse interaction, or both? → A: Keyboard-only navigation.
- Q: What should happen if the Apple Container CLI is not installed or not found in PATH when the TUI starts? → A: Show error and exit.
- Q: Should the container list auto-refresh periodically, or only refresh manually when the user requests it? → A: Manual refresh only.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Manage Containers From a Menu (Priority: P1)

As a developer, I want to list containers and start or stop a selected container
from a TUI menu so I do not need to remember CLI commands.

**Why this priority**: Container lifecycle control is the core day-to-day task
and delivers immediate value without additional setup.

**Independent Test**: Can be fully tested by listing containers, selecting one
from the menu, and starting or stopping it with visible status feedback.

**Acceptance Scenarios**:

1. **Given** the daemon is running and at least one container exists, **When**
   I open the TUI, **Then** I see a list of containers with their status.
2. **Given** a stopped container is selected, **When** I choose Start, **Then**
   the command preview is shown and the container transitions to running.
3. **Given** a running container is selected, **When** I choose Stop, **Then**
   the command preview is shown and the container transitions to stopped.

---

### User Story 2 - Pull and Build Images (Priority: P2)

As a developer, I want to pull images and build from a selected container file
so I can prepare containers without leaving the TUI.

**Why this priority**: Image and build actions are frequent follow-ups to
starting or stopping containers and enable end-to-end workflows.

**Independent Test**: Can be fully tested by pulling a known image and building
from a chosen file, with progress and final status shown in the UI.

**Acceptance Scenarios**:

1. **Given** I enter a valid image reference, **When** I select Pull, **Then**
   I see progress and a success or failure result.
2. **Given** I select a valid Containerfile or Dockerfile, **When** I select
  Build, **Then** I see progress and a success or failure result.

---

### User Story 3 - Safe Destructive and Daemon Actions (Priority: P3)

As a developer, I want to delete stopped containers and start or stop the
daemon with safety checks so I can manage system state confidently.

**Why this priority**: These actions are less frequent but higher risk, so they
must be safe and explicit once core workflows exist.

**Independent Test**: Can be fully tested by deleting a stopped container and
starting or stopping the daemon with confirmation prompts and clear results.

**Acceptance Scenarios**:

1. **Given** a stopped container is selected, **When** I choose Delete, **Then**
   I must confirm and see the command preview before deletion occurs.
2. **Given** the daemon is stopped, **When** I choose Start Daemon, **Then** I
   see confirmation, command preview, and a running status afterward.
3. **Given** the daemon is running, **When** I choose Stop Daemon, **Then** I
   must confirm and see the daemon status change to stopped.

---

### Edge Cases

- No containers exist and the list view should show an empty state with
  guidance.
- The daemon is stopped and container actions should surface a clear error.
- A pull or build fails and the error output is surfaced to the user.
- The terminal is resized to a small window and the UI remains usable.
- The user lacks permission to start or stop the daemon and receives guidance.
- The Apple Container CLI is not installed or not in PATH and the TUI shows an
  error and exits immediately.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The system MUST list containers with name/ID and status.
- **FR-002**: Users MUST be able to start a selected container.
- **FR-003**: Users MUST be able to stop a selected container.
- **FR-004**: Users MUST be able to delete a selected stopped container with
  explicit confirmation.
- **FR-005**: The system MUST support pulling an image by user-provided
  reference.
- **FR-006**: Users MUST be able to select a local Containerfile or Dockerfile
  and start a build.
- **FR-007**: Users MUST be able to start and stop the container daemon.
- **FR-008**: Every action MUST display a command preview before execution and
  support a dry-run mode (activated via --dry-run CLI flag at startup) that does not execute the command.
- **FR-009**: Action results MUST show success or failure and include stdout
  and stderr output.
- **FR-010**: The UI MUST remain usable and consistent after terminal resize.
- **FR-011**: User preferences MUST be stored locally and support both
  ~/.config/apple-tui/config and ~/Library/Application Support/apple-tui/config
  (read both paths with first match used, write to ~/Library/Application Support/apple-tui/config per macOS standard). The write path MUST be documented in the
  app.
- **FR-012**: The application MUST run on macOS 26.x on Apple Silicon (M-series)
  CPUs.
- **FR-013**: The UI MUST support keyboard-only navigation without requiring
  mouse interaction.
- **FR-014**: The application MUST check for the Apple Container CLI at startup
  and exit with a clear error message if not found in PATH.
- **FR-015**: The container list MUST refresh only when explicitly requested by
  the user, not automatically or periodically.

### Key Entities *(include if feature involves data)*

- **Container**: A managed container with identity, status, and image reference.
- **ImageReference**: A user-supplied reference used to pull an image.
- **BuildSource**: A local file reference used to start a build.
- **DaemonStatus**: The current running state of the container daemon.
- **CommandRun**: A record of an action with command preview, outcome, and
  captured output.
- **UserConfig**: Saved user preferences such as default paths and UI options.

## Assumptions

- The app is used on a local developer machine and is not a remote management
  tool.
- Bulk deletion of all stopped containers is out of scope for the MVP.
- The Apple Container CLI is installed and available to the user.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 90% of first-time users can start or stop a container within
  2 minutes of launching the TUI.
- **SC-002**: 100% of actions show a command preview before execution.
- **SC-003**: 95% of users can complete a pull or build action without leaving
  the TUI.
- **SC-004**: The UI remains readable and navigable after resizing the terminal
  to 80x24 characters.
- **SC-005**: 0 destructive actions occur without an explicit confirmation
  step.
