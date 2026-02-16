# Feature Specification: Enhanced Menu Navigation and Image Management

**Feature Branch**: `002-refactor-menu-images`  
**Created**: February 16, 2026  
**Status**: Draft  
**Input**: User description: "Refactor menu structure and add image management features: change container enter behavior to open submenu with start/stop/tail logs/enter container options, add images menu (i key) with image list/pull/build/prune, and image submenu with inspect/delete operations"

## Clarifications

### Session 2026-02-16

- Q: How should users select options within submenus (container submenu and image submenu)? → A: Arrow keys (up/down) to highlight + Enter to confirm selection
- Q: What should happen when a user tries to select "Enter container" on a stopped container? → A: Hidden
- Q: When shell auto-detection fails (no bash, sh, /bin/sh, /bin/bash, or ash found), what should the system do? → A: Show error, stay in submenu
- Q: When viewing live container logs and the container stops or is removed, what should happen? → A: Show message, wait for Esc
- Q: What type of confirmation should "Image prune" use (removing all unused images)? → A: Type-to-confirm

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Container Action Submenu (Priority: P1)

Users need quick access to common container operations (start, stop, view logs, enter shell) through an intuitive submenu instead of remembering individual keyboard shortcuts for each action.

**Why this priority**: This changes the core container interaction workflow. Getting this right is essential for user experience, and all other features depend on the pattern established here.

**Independent Test**: Can be fully tested by navigating the container list, pressing Enter on any container, and verifying the submenu displays with context-appropriate options (e.g., only "Start" shows for stopped containers). Delivers immediate value by simplifying container interaction.

**Acceptance Scenarios**:

1. **Given** I am viewing the container list with a stopped container selected, **When** I press Enter, **Then** a submenu appears showing "Start container", "Tail container log", "Enter container", and "Back" options
2. **Given** I am viewing the container list with a running container selected, **When** I press Enter, **Then** a submenu appears showing "Stop container", "Tail container log", "Enter container", and "Back" options
3. **Given** I am in a container submenu, **When** I select "Tail container log", **Then** I see live streaming logs from the container with the ability to exit back to the submenu using Esc
4. **Given** I am in a container submenu, **When** I select "Enter container", **Then** the system automatically detects and executes an available shell (trying bash, sh, /bin/sh, /bin/bash, ash in sequence) and I am placed in an interactive shell session
5. **Given** I am in an interactive container shell session, **When** I exit the shell or press the designated key combination, **Then** I return to the container submenu
6. **Given** I am in a container submenu, **When** I select "Back" or press Esc, **Then** I return to the main container list view

---

### User Story 2 - Image List View and Quick Actions (Priority: P2)

Users need to view all local container images and perform common image operations (pull new images, build from Containerfiles, prune unused images) from a dedicated image management screen.

**Why this priority**: Image management is a distinct workflow from container management. Separating it into its own screen improves organization and makes room for more image-specific features in the future.

**Independent Test**: Can be fully tested by pressing 'i' from the main menu and verifying the image list displays with NAME, TAG, and DIGEST columns. The pull, build, and prune features already exist and just need to be accessible from this new screen. Delivers value by consolidating image operations.

**Acceptance Scenarios**:

1. **Given** I am viewing the main container list, **When** I press 'i', **Then** I see the image list screen showing all local images with NAME, TAG, and DIGEST columns
2. **Given** I am viewing the image list, **When** I press 'p', **Then** I am taken to the existing image pull workflow
3. **Given** I am viewing the image list, **When** I press 'b', **Then** I am taken to the existing image build workflow
4. **Given** I am viewing the image list, **When** I press 'n', **Then** I am prompted with type-to-confirm for image pruning, and upon typing the confirmation word, unused images are removed
5. **Given** I am viewing the image list, **When** I press Esc, **Then** I return to the main container list view
6. **Given** the image list is displayed, **When** images are pulled, built, or pruned, **Then** the list automatically refreshes to show the current state

---

### User Story 3 - Image Details Submenu (Priority: P3)

Users need to inspect image details and delete specific images through a submenu similar to the container submenu, providing consistent navigation patterns across the application.

**Why this priority**: This enhances image management with inspection and deletion capabilities, but is less critical than the core viewing and quick actions. It follows the submenu pattern established in P1.

**Independent Test**: Can be fully tested by selecting an image from the image list, pressing Enter, and verifying the submenu shows inspection and deletion options. Delivers value by providing detailed image information and cleanup capabilities.

**Acceptance Scenarios**:

1. **Given** I am viewing the image list with an image selected, **When** I press Enter, **Then** a submenu appears showing "Inspect image", "Delete image", and "Back" options
2. **Given** I am in an image submenu, **When** I select "Inspect image", **Then** I see the formatted JSON output of the image inspection with the ability to scroll and exit using Esc
3. **Given** I am in an image submenu viewing inspection details, **When** I press Esc, **Then** I return to the image submenu
4. **Given** I am in an image submenu, **When** I select "Delete image", **Then** I am prompted to confirm deletion (with type-to-confirm for safety), and upon confirmation, the image is removed and I return to the image list
5. **Given** I am in an image submenu, **When** I select "Back" or press Esc, **Then** I return to the image list view

---

### Edge Cases

- What happens when attempting to enter a container shell but none of the standard shells (bash, sh, /bin/sh, /bin/bash, ash) are available in the container?
- How does the system handle pressing Enter on a container when the daemon is not running?
- What happens when viewing live container logs and the container stops or is removed?
- How does the system handle image deletion when the image is currently in use by a running or stopped container?
- What happens when attempting to prune images but no unused images exist?
- How does the system handle very long image names or digests that exceed the terminal width?
- What happens when the image list is empty?
- How does the system handle navigation (up/down) in empty lists or single-item lists?
- What happens when a user is viewing image inspection details and the image is deleted by another process?

## Requirements *(mandatory)*

### Functional Requirements

**Menu Structure & Navigation**:
- **FR-001**: System MUST replace the current Enter-to-toggle behavior in the container list with Enter-to-open-submenu behavior
- **FR-002**: Main menu MUST remove 'p' (pull) and 'b' (build) shortcuts and replace them with 'i' (images) shortcut
- **FR-003**: Main menu help text MUST change from "enter=toggle" to "enter=submenu" and from "p=pull, b=build" to "i=images"
- **FR-004**: System MUST provide consistent Esc key behavior across all submenus to return to the previous screen
- **FR-004a**: All submenus MUST support arrow key (up/down) navigation to highlight options and Enter key to confirm selection

**Container Submenu**:
- **FR-005**: Container submenu MUST display only contextually relevant options based on container state (Start for stopped, Stop for running)
- **FR-006**: Container submenu MUST include "Tail container log" option for all containers regardless of state
- **FR-007**: Container submenu MUST include "Enter container" option only for running containers (hidden for stopped containers)
- **FR-008**: Container submenu MUST include "Back" option to return to container list
- **FR-009**: "Tail container log" MUST execute `container logs -f [containerName]` and stream live logs
- **FR-010**: Live log view MUST support Esc key to return to container submenu
- **FR-011**: "Enter container" MUST automatically attempt to execute interactive shell in this order: bash, sh, /bin/sh, /bin/bash, ash, using the first available option
- **FR-012**: "Enter container" MUST execute `container exec -it [containerName] [shell]` for the detected shell
- **FR-013**: Interactive shell session MUST allow users to exit back to the container submenu

**Image List View**:
- **FR-014**: System MUST display image list when user presses 'i' from main container list
- **FR-015**: Image list MUST execute `container image list` command to retrieve images
- **FR-016**: Image list MUST display three columns: NAME, TAG, and DIGEST
- **FR-017**: Image list MUST support keyboard navigation (up/down arrows) to select images
- **FR-018**: Image list menu MUST display: "Keys: up/down, enter=submenu, p=pull, b=build, n=image-prune, esc=back to main menu"
- **FR-019**: Image list MUST support 'p' key to access existing image pull functionality
- **FR-020**: Image list MUST support 'b' key to access existing image build functionality
- **FR-021**: Image list MUST support 'n' key to trigger image prune operation
- **FR-022**: Image prune MUST execute `container image prune` command after user confirmation via type-to-confirm pattern
- **FR-023**: Image list MUST support Esc key to return to main container list
- **FR-024**: Image list MUST automatically refresh after pull, build, or prune operations

**Image Submenu**:
- **FR-025**: System MUST display image submenu when user presses Enter on a selected image
- **FR-026**: Image submenu MUST include "Inspect image", "Delete image", and "Back" options
- **FR-027**: "Inspect image" MUST execute `container image inspect [imageReference] | jq` and display formatted output, where imageReference is the full image reference in format NAME:TAG (e.g., "ubuntu:latest") or NAME@DIGEST for images without tags
- **FR-028**: Image inspection view MUST support scrolling for long output using arrow keys, Page Up/Down, Home/End keys
- **FR-029**: Image inspection view MUST support Esc key to return to image submenu
- **FR-030**: "Delete image" MUST prompt for confirmation before deletion
- **FR-031**: "Delete image" MUST execute `container image rm [imageName]` upon confirmation
- **FR-032**: Image submenu MUST support Esc or "Back" to return to image list

**Error Handling**:
- **FR-033**: System MUST display error message if attempting to enter container shell but no supported shell is available, and keep user in the container submenu
- **FR-034**: System MUST handle gracefully when container logs stream is interrupted (container stopped/removed) by displaying an informational message and waiting for user to press Esc to return to submenu
- **FR-035**: System MUST prevent image deletion if image is in use and display appropriate error message
- **FR-036**: System MUST handle empty image lists with appropriate message
- **FR-037**: System MUST truncate or wrap long image names/digests that exceed terminal width

### Key Entities

- **Container Submenu**: Navigation menu displaying context-sensitive operations for a specific container, including lifecycle actions (start/stop), diagnostics (logs), and interaction (shell access)
- **Image List View**: Display screen showing all local container images with identifying information (name, tag, digest) and operations toolbar
- **Image Submenu**: Navigation menu displaying operations for a specific image, including inspection and deletion
- **Menu Navigation State**: System state tracking current view (main container list, container submenu, image list, image submenu, etc.) to enable proper navigation flow and back-navigation
- **Shell Detection Result**: Outcome of attempting to find an available shell in a container, including the detected shell command or failure indication

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Users can navigate from container list to container operations (start/stop/logs/shell) in exactly 2 keystrokes (Enter + selection) instead of having to remember multiple different keyboard shortcuts
- **SC-002**: Users can access any image management operation (view list, pull, build, prune, inspect, delete) within 3 keystrokes from the main menu
- **SC-003**: Container logs stream in real-time with less than 100ms latency from when container outputs to when user sees it
- **SC-004**: Shell auto-detection succeeds in finding an available shell in 95% of standard container images
- **SC-005**: Navigation pattern is consistent across all screens - Esc always returns to previous screen, Enter always opens submenu or executes action
- **SC-006**: Image list displays all images in under 2 seconds for repositories with up to 100 local images
- **SC-007**: Users can perform complete image lifecycle (pull, inspect, delete) without leaving the TUI or referring to external documentation
