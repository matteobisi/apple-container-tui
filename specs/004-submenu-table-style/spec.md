# Feature Specification: Consistent Table Styling for Submenus

**Feature Branch**: `004-submenu-table-style`  
**Created**: February 17, 2026  
**Status**: Draft  
**Input**: User description: "Apply table formatting style to container and image submenus. Extend visual improvements from 003-tui-table-format (bold headers, horizontal separators, inverse video selection) to submenu screens for consistent TUI styling across the entire application."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Container Submenu Visual Consistency (Priority: P1)

When users navigate to a container submenu (by pressing Enter on a container in the list), they want the same professional visual styling they experience in the main container list: clear section separation, bold headers, and inverse video selection highlighting. This provides a consistent interaction pattern throughout the application.

**Why this priority**: Container operations are the primary use case. Users interact with container submenus frequently (start/stop/logs/shell), so visual inconsistency here is most noticeable and impacts user confidence in the application's polish.

**Independent Test**: Navigate to any container submenu, verify the container info section has a horizontal separator, menu items use inverse video selection (no cursor prefix), and section headers are bold. Can be tested without touching image submenus.

**Acceptance Scenarios**:

1. **Given** user is on container list, **When** user presses Enter on a container, **Then** submenu displays container info with bold "Container Details" header followed by horizontal separator line
2. **Given** container submenu is displayed, **When** user navigates menu items with arrow keys, **Then** selected item is highlighted with inverse video (full-row background/foreground swap) instead of colored text with cursor prefix
3. **Given** container submenu has multiple action items, **When** viewing the menu, **Then** a horizontal separator line appears between the container info section and the action menu section
4. **Given** user is in container submenu, **When** viewing available actions, **Then** action section has bold "Available Actions" header for clear visual hierarchy

---

### User Story 2 - Image Submenu Visual Consistency (Priority: P2)

When users navigate to an image submenu (by pressing Enter on an image in the list), they want the same professional visual styling as container submenus and list views: clear section separation, bold headers, and inverse video selection highlighting.

**Why this priority**: While less frequently used than container submenus, image management is still a core workflow. Visual consistency across all screens reinforces the application's professional appearance and reduces cognitive load when switching between views.

**Independent Test**: Navigate to any image submenu, verify the image info section has a horizontal separator, menu items use inverse video selection, and section headers are bold. Can be tested independently of container submenu changes.

**Acceptance Scenarios**:

1. **Given** user is on image list, **When** user presses Enter on an image, **Then** submenu displays image info with bold "Image Details" header followed by horizontal separator line
2. **Given** image submenu is displayed, **When** user navigates menu items with arrow keys, **Then** selected item is highlighted with inverse video instead of colored text with cursor prefix
3. **Given** image submenu has multiple action items, **When** viewing the menu, **Then** a horizontal separator line appears between the image info section and the action menu section
4. **Given** user is in image submenu, **When** viewing available actions, **Then** action section has bold "Available Actions" header

---

### User Story 3 - Help Screen Visual Consistency (Priority: P3)

The help screen should use the same visual conventions: bold section headers and horizontal separators to organize keyboard shortcuts and command explanations. This ensures every screen in the application follows the established design language.

**Why this priority**: Help screen is accessed less frequently (only when users need reference) but should still match the visual style. Lower priority because it doesn't impact primary workflows, but important for completeness.

**Independent Test**: Press '?' to open help screen, verify section headers are bold, horizontal separators divide sections. Can be tested after submenu work is complete.

**Acceptance Scenarios**:

1. **Given** user presses '?' from any screen, **When** help screen displays, **Then** section headers (Navigation, Actions, etc.) are rendered in bold
2. **Given** help screen is displayed, **When** viewing content, **Then** horizontal separator lines divide major sections for improved scannability

---

### Edge Cases

- What happens when submenu has only one action? (Still show headers and separators for consistency, even with single item)
- What happens if container/image info is very long? (Truncate with ellipsis to maintain single-line display per field)
- What happens when terminal is too narrow to display full menu text? (Menu items wrap gracefully or truncate with ellipsis, maintaining minimum 80-char width requirement)
- What happens if terminal doesn't support inverse video? (Fall back to existing color accent highlighting, ensuring selection is still visible)
- What happens in daemon control screen? (Apply same styling: bold headers, separators, inverse video selection)

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Container submenu MUST display container information section with bold "Container Details" header
- **FR-002**: Container submenu MUST display horizontal separator line (─ characters) between container info and action menu sections
- **FR-003**: Container submenu action items MUST use inverse video selection highlighting instead of cursor prefix when navigated with arrow keys
- **FR-004**: Container submenu MUST display bold "Available Actions" header above action list
- **FR-005**: Image submenu MUST display image information section with bold "Image Details" header
- **FR-006**: Image submenu MUST display horizontal separator line between image info and action menu sections
- **FR-007**: Image submenu action items MUST use inverse video selection highlighting instead of cursor prefix
- **FR-008**: Image submenu MUST display bold "Available Actions" header above action list
- **FR-009**: Help screen MUST display section headers (Navigation, Actions, etc.) in bold text
- **FR-010**: Help screen MUST use horizontal separator lines to divide major content sections
- **FR-011**: All submenus and screens MUST maintain consistent spacing matching the table layout style from feature 003-tui-table-format
- **FR-012**: Inverse video selection MUST use Lipgloss .Reverse(true) style to match container/image list behavior

### Key Entities

- **Submenu Layout**: Consists of three visual sections: (1) Details header with bold styling, (2) Information display area (container/image details), (3) Horizontal separator, (4) Actions header with bold styling, (5) Action menu items with inverse video selection. Must maintain visual hierarchy through consistent use of bold headers and separators.

- **Menu Item Row**: Represents one actionable item in a submenu. Displays action text (e.g., "Stop container", "Inspect image") with full-row inverse video highlighting when selected, replacing the previous cursor-prefix approach. Must support keyboard navigation (arrow keys) with instant visual feedback.

- **Section Header**: Bold text label that identifies content sections (e.g., "Container Details", "Available Actions", "Navigation Keys"). Must be rendered using Lipgloss bold style to match table headers from 003-tui-table-format.

- **Visual Separator**: Horizontal line composed of box-drawing characters (─) spanning full terminal width. Must be used to divide sections (info from actions, help sections from each other) for improved visual organization.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Users can instantly identify selected menu items in submenus through inverse video highlighting without relying on cursor prefix recognition
- **SC-002**: Visual hierarchy in submenus is clear within 1 second of screen display through bold headers and horizontal separators
- **SC-003**: Container and image submenus maintain visual consistency with list views, reducing cognitive load when switching between screens
- **SC-004**: Help screen content scannability improves through bold section headers and separators, reducing time to find specific command reference
- **SC-005**: All screens adhere to the same design language established in feature 003-tui-table-format
- **SC-006**: Terminal width constraints (minimum 80 characters) are respected with appropriate text wrapping or truncation in all submenu screens

## Assumptions

- Lipgloss library already supports .Reverse(true) and .Bold(true) styles (confirmed in 003-tui-table-format implementation)
- Terminal emulator supports ANSI inverse video escape sequences (standard in modern terminals)
- Unicode box-drawing characters (─) are available in terminal font (same assumption as 003-tui-table-format)
- Existing submenu navigation logic (arrow keys, Enter, Esc) remains unchanged - only visual presentation is modified
- Container and image submenu View() methods can be refactored without breaking existing functionality
- Help screen has clearly defined sections that can be separated with horizontal lines
